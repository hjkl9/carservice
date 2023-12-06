package sms

import (
	"regexp"
)

func CheckPhoneNumber(phoneNumber string) bool {
	// Check prefix of phone number.
	matched, _ := regexp.MatchString("0?(13|14|15|17|18)[0-9]{9}", phoneNumber)
	return matched
}
