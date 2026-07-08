package updater

import (
	"crypto/sha256"
	"encoding/json"
	"fmt"
	"goxlock/config"
	"goxlock/utils"
	"io"
	"net/http"
	"os"
	"os/exec"
	"strings"
)

// - Release
// Gets the release wanted by the application (focused on the latest)
type Release struct {
    TagName string `json:"tag_name"`
    Assets  []Asset `json:"assets"`
}

// - Asset
// The data required to get the download data
type Asset struct {
    Name 				string `json:"name"`
    BrowserDownloadURL 	string `json:"browser_download_url"`
	ShaDigest			string `json:"digest"`
}



// - CheckForUpdate
// This will check for update from trusted sources
func CheckForUpdate() error {
	// Info : The safe url for goxlock latest release
	// url -> download url -> download
	url := `https://api.github.com/repos/TOXICAI314/goxlock/releases/latest`
	resp,err := http.Get(url)
	if err != nil {
		return &config.UserSafetyError{
			Cause: err.Error(),
			Message: fmt.Sprintf(`No such web data or history found on - %s`,url),
		}
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		return &config.UserSafetyError{
			Cause: `Expected an ok status code`,
			Message: fmt.Sprintf(`Expected a 200 status code , got - %d`,resp.StatusCode),
		}
	}

	data,err := io.ReadAll(resp.Body)
	if err != nil {
		return &config.UserSafetyError{
			Cause: err.Error(),
			Message: `The given data is not in readable format`,
		}
	}

	var release Release
	err = json.Unmarshal(data,&release)
	if err != nil {
		return &config.UserSafetyError{
			Cause: err.Error(),
			Message: `Expected a json like data but got otherwise`,
		}
	}

	if release.TagName == fmt.Sprintf(`v%.1f`,config.Version) {
		fmt.Printf(`Already the latest verison of the app installed - %.1f`,config.Version)
		return nil
	}

	ans := utils.GetYesORNo(fmt.Sprintf(`
		A new version of 'goxlock' is found - %s
		Want to update? (y/n) : `,release.TagName))
	
	if !ans {
		return nil
	}

	permittedTempFile,err := os.CreateTemp(``,`goxlock-updated-installer-*.exe`)
	if err != nil {
		return &config.UserSafetyError{
			Cause: `Access denied to the temp dir`,
			Message: `The access to the temporart directory is not given`,
		}
	}
	defer permittedTempFile.Close()

	var index int
	var found bool
	for i,data := range release.Assets {
		if data.Name != `goxlock_setup.exe` {
			continue
		}
		exedresp,err := http.Get(data.BrowserDownloadURL)
		if err != nil {
			return &config.UserSafetyError{
				Cause: err.Error(),
				Message:fmt.Sprintf(`Cannot fetch from the the url that is given for the latest update - %s`,release.Assets[0].BrowserDownloadURL),
			}
		}
		defer exedresp.Body.Close()
		if resp.StatusCode != http.StatusOK {
			return &config.UserSafetyError{
				Cause: `Expected an ok status code`,
				Message: fmt.Sprintf(`Expected a 200 status code , got - %d`,resp.StatusCode),
			}
		}

		_,err = io.Copy(permittedTempFile,exedresp.Body)
		if err != nil {
			return &config.UserSafetyError{
				Cause: err.Error(),
				Message: `Cannot copy the data into the temp file to install the update`,
			}
		}
		index = i
		found = true
	}

	if !found {
		return &config.UserSafetyError{
			Cause: `No exe url found to download`,
			Message: `Expected exe url got none`,
		}
	} 

	permittedTempFile.Seek(0,0)
	hash := sha256.New()
	_,err = io.Copy(hash,permittedTempFile)
	if err != nil {
		return &config.UserSafetyError{
			Cause: err.Error(),
			Message: `The given hash cannot be made from the downloaded file`,
		}
	}
	actualDigest := fmt.Sprintf("%x", hash.Sum(nil))

	givenDigest := release.Assets[index].ShaDigest
	expectedDigest := strings.TrimPrefix(givenDigest, "sha256:")

	if !strings.EqualFold(actualDigest, expectedDigest) {
		return &config.UserSafetyError{
			Cause: `Unmatched Hashes`,
			Message: fmt.Sprintf("The hashes are unmatched from the provided hash\nProvided - %s\nExpected - %s",actualDigest,expectedDigest),
		}
	}

	cmd := exec.Command(
				permittedTempFile.Name(),
				"/VERYSILENT",
				"/SUPPRESSMSGBOXES",
				"/NORESTART",
			)
	err = cmd.Start()
	if err != nil {
		return &config.UserSafetyError{
			Cause: err.Error(),
			Message: `Cannot start the update of the executable file`,
		}
	}
	fmt.Println(`Exitting the Program to redownload the application`)
	return nil
}