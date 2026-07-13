package github

import (
	"fmt"
)

// reports the bug to the github page by retruning a url which can be opened later by the parent
func BugReport() (url string, err error) {
	return fmt.Sprintf(`%s?template=bug_report.md`, GithubReportURL), nil
}
