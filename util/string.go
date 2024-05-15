package utils

import (
	"github.com/google/uuid"
	"regexp"
	"strings"
)

func CapitalizeEachWord(input string) string {
	words := strings.Fields(input)
	for i, word := range words {
		if len(word) > 0 {
			words[i] = strings.ToUpper(string(word[0])) + strings.ToLower(word[1:])
		}
	}
	return strings.Join(words, " ")
}
func RemoveSpecialChars(input string) string {
	re := regexp.MustCompile(`[\n\w]`)
	result := re.ReplaceAllString(input, "")
	return result
}
func GenerateUUIDString() string {
	id := uuid.New()
	return id.String()
}
