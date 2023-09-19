package h

import (
	"encoding/json"
	"strconv"
	"strings"

	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
)

// Function for checking availability of element in an array
func Contains[T comparable](s []T, e T) bool {
	for _, v := range s {
		if v == e {
			return true
		}
	}
	return false
}
// Function which converts an array with integers into array with string elements
func IntArrayToStrArray(arr []int) []string {
	var res []string
	for _, el := range arr {
		res = append(res, strconv.Itoa(el))
	}
	return res
}
// Function which checks validation of name for file or directory
func IsValidName(name string) bool {
	const bannedSymbols = "<>:«/\\|?*"
	// Name cannot be empty or longer than 50
	if name == "" || len(name) > 50 {
		return false
	}
	// Name cannot include a banned symbols
	for _, char := range bannedSymbols {
		if strings.Contains(name, string(char)) {
			return false
		}
	}
	return true
}
// Function which parse string and convert its array of integers
func ParseIds(jsonBuffer string) ([]int, error) {
	ids := []int{}
	// Returning an empty array if string is empty
	if len(jsonBuffer) == 0 {
		return ids, nil
	}
	// Replacing curly braces to square 
	jsonBuffer = strings.Replace(jsonBuffer, "{", "[", -1)
	jsonBuffer = strings.Replace(jsonBuffer, "}", "]", -1)
	// Parsing JSON string
    err := json.Unmarshal([]byte(jsonBuffer), &ids)
    return ids, err
}
// Function which returns unique name (uuid) for something
func GenerateUniqueName() string {
	return uuid.New().String()
}
func HashData(data ...string) (string, error) {
	var s string
	for _, elem := range data {
		s += elem
	}
    bytes, err := bcrypt.GenerateFromPassword([]byte(s), 14)
    return string(bytes), err
}

func CheckHash(hash string, data ...string) bool {
	var s string
	for _, elem := range data {
		s += elem
	}
    err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(s))
    return err == nil
}