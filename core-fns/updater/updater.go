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
	"time"
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
	client := &http.Client{
		Timeout: 30 * time.Second,
	}
	url := `https://api.github.com/repos/TOXICAI314/goxlock/releases/latest`
	resp,err := client.Get(url)
	if err != nil {
		return &config.FunctionFailError{
			Cause: err.Error(),
			Message: fmt.Sprintf(`No such web data or history found on - %s`,url),
			ElapsedTime: time.Now(),
			Provider: `unlocker.CheckUpdates`,
		}
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		return &config.FunctionCancelError{
			Cause: `Expected an ok status code`,
			Message: fmt.Sprintf(`Expected a 200 status code , got - %d`,resp.StatusCode),
			ElapsedTime: time.Now(),
			Provider: `unlocker.CheckUpdates`,
		}
	}

	data,err := io.ReadAll(resp.Body)
	if err != nil {
		return &config.FunctionFailError{
			Cause: err.Error(),
			Message: `The given data is not in readable format`,
			ElapsedTime: time.Now(),
			Provider: `unlocker.CheckUpdates`,
		}
	}

	var release Release
	err = json.Unmarshal(data,&release)
	if err != nil {
		return &config.FunctionFailError{
			Cause: err.Error(),
			Message: `Expected a json like data but got otherwise`,
			ElapsedTime: time.Now(),
			Provider: `unlocker.CheckUpdates`,
		}
	}

	if release.TagName == `v`+config.Version {
		fmt.Printf(`Already the latest verison of the app installed - %s`,config.Version)
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
		return &config.FunctionFailError{
			Cause: `Access denied to the temp dir`,
			Message: `The access to the temporart directory is not given`,
			ElapsedTime: time.Now(),
			Provider: `unlocker.CheckUpdates`,
		}
	}

	var index int
	var found bool
	for i,data := range release.Assets {
		if data.Name != `goxlock_setup.exe` {
			continue
		}
		exedresp,err := http.Get(data.BrowserDownloadURL)
		if err != nil {
			return &config.FunctionFailError{
				Cause: err.Error(),
				Message:fmt.Sprintf(`Cannot fetch from the the url that is given for the latest update - %s`,release.Assets[0].BrowserDownloadURL),
				ElapsedTime: time.Now(),
				Provider: `unlocker.CheckUpdates`,
			}
		}
		defer exedresp.Body.Close()
		if exedresp.StatusCode != http.StatusOK {
			return &config.FunctionCancelError{
				Cause: `Expected an ok status code`,
				Message: fmt.Sprintf(`Expected a 200 status code , got - %d`,resp.StatusCode),
				ElapsedTime: time.Now(),
				Provider: `unlocker.CheckUpdates`,
			}
		}

		_,err = io.Copy(permittedTempFile,exedresp.Body)
		if err != nil {
			return &config.FunctionFailError{
				Cause: err.Error(),
				Message: `Cannot copy the data into the temp file to install the update`,
				ElapsedTime: time.Now(),
				Provider: `unlocker.CheckUpdates`,
			}
		}
		index = i
		found = true
	}

	if !found {
		return &config.FunctionCancelError{
			Cause: `No exe url found to download`,
			Message: `Expected exe url got none`,
			ElapsedTime: time.Now(),
			Provider: `unlocker.CheckUpdates`,
		}
	} 

	permittedTempFile.Seek(0,0)
	hash := sha256.New()
	_,err = io.Copy(hash,permittedTempFile)
	if err != nil {
		return &config.FunctionFailError{
			Cause: err.Error(),
			Message: `The given hash cannot be made from the downloaded file`,
			ElapsedTime: time.Now(),
			Provider: `unlocker.CheckUpdates`,
		}
	}
	actualDigest := fmt.Sprintf("%x", hash.Sum(nil))
	permittedTempFile.Close()

	givenDigest := release.Assets[index].ShaDigest
	expectedDigest := strings.TrimPrefix(givenDigest, "sha256:")

	if !strings.EqualFold(actualDigest, expectedDigest) {
		return &config.FunctionCancelError{
			Cause: `Unmatched Hashes`,
			Message: fmt.Sprintf("The hashes are unmatched from the provided hash\nProvided - %s\nExpected - %s",actualDigest,expectedDigest),
			ElapsedTime: time.Now(),
			Provider: `unlocker.CheckUpdates`,
		}
	}

	cmd := exec.Command(
				"cmd", 
				"/c", 
				"timeout /t 5 /nobreak > NUL && "+permittedTempFile.Name(),
				"/NORESTART",
			)
	err = cmd.Start()
	if err != nil {
		return &config.FunctionFailError{
			Cause: err.Error(),
			Message: `Cannot start the update of the executable file`,
			ElapsedTime: time.Now(),
			Provider: `unlocker.CheckUpdates`,
		}
	}
	fmt.Printf("Exitting the Program to redownload the application\nThe download will start after 5 sec .")
	os.Exit(0)
	return nil
}