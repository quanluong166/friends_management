package helper

import "strings"

func GenerateRelationshipKey(userID1, userID2 string) string {
	return userID1 + "_" + userID2
}

// GetOtherEmailFromKey extracts the other email from the relationship key.
func GetOtherEmailFromKey(relationshipKey, email string) string {
	emails := strings.Split(relationshipKey, ":")
	if emails[0] == email {
		return emails[1]
	}
	return emails[0]
}
