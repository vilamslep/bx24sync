package app


import(
	"io"
	"strings"
	"regexp"
)

//checkinput
func HandlerCheckInput(reader io.Reader) (bool, error) {
	body, err := io.ReadAll(reader)

	if err != nil {
		return false, err
	}

	content := strings.ReplaceAll(string(body), "\n", "")

	regStr := `^{"#",+[[:xdigit:]]{8}(-[[:xdigit:]]{4}){3}-[[:xdigit:]]{12},[\d]{1,6}:[[:xdigit:]]{32}}$`

	if matched, err := regexp.MatchString(regStr, content); err != nil {
		return false, err
	} else {
		return matched, nil
	}
}
