package utils

import "regexp"

func FindEmails(text string) []string {
	emailPattern := regexp.MustCompile(`[a-zA-Z0-9._%+\-]+@[a-zA-Z0-9.\-]+\.[a-zA-Z]{2,}`)
	matches := emailPattern.FindAllString(text, -1)
	if matches == nil {
		return []string{}
	}
	return matches
}
