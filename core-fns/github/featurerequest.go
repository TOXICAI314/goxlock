package github

import "fmt"

// requests the feature to the gievn url
func FeatureRequest() (url string,err error) {
	return fmt.Sprintf(`%s?template=feature_request.md`,GithubReportURL),nil
}