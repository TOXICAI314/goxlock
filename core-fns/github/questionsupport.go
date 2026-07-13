package github

import "fmt"

// reports the question to the given url page
func QuestionAndSupport() (url string,err error) {
	return fmt.Sprintf(`%s?template=question_support.md`,GithubReportURL),nil
}