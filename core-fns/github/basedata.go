package github

const (
	// api.github.com -> Rest api for the github (for automation and code execution for json)

	// Base API for the connection
	GithubAPIURL                 string = `https://api.github.com/repos/TOXICAI314/goxlock`
	GithubLatestReleaseAPIURL string = GithubAPIURL + `/releases/latest`


	// github.com -> The interactive browser form
	GithubRepoURL   string = "https://github.com/TOXICAI314/goxlock"
	GithubReportURL string = GithubRepoURL + `/issues/new`
)
