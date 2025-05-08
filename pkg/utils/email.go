package utils

import "regexp"

// IsValidEmail support to check whether input string is valid email format
func IsValidEmail(email string) bool {
	regex := `^[a-zA-Z0-9._%+\-]+@[a-zA-Z0-9.\-]+\.[a-zA-Z]{2,}$`
	re := regexp.MustCompile(regex)
	return re.MatchString(email)
}

// FindEmails find valid email string format in the input param
func FindEmails(text string) []string {
	emailPattern := regexp.MustCompile(`[a-zA-Z0-9._%+\-]+@[a-zA-Z0-9.\-]+\.[a-zA-Z]{2,}`)
	matches := emailPattern.FindAllString(text, -1)
	if matches == nil {
		return []string{}
	}
	return matches
}
