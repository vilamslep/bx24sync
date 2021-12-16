package app


import(
	"io"
	"strings"
	"regexp"
)

//body would be like {"#",8f8a65b4-94c4-4794-b3e2-800d18d503ca,151:80d60cc47a5468a511e5f0049748a86c}
func DefaultCheckInput(reader io.Reader) (bool, error) {
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
