package h

import (
	"encoding/json"
	"strconv"
	"strings"
)

const bannedSymbols = "<>:«/\\|?*"

func Contains[T comparable](s []T, e T) bool {
	for _, v := range s {
		if v == e {
			return true
		}
	}
	return false
}

func IntArrayToStrArray(arr []int) []string {
	var res []string
	for _, el := range arr {
		res = append(res, strconv.Itoa(el))
	}
	return res
}

func IsValidName(name string) bool {
	if name == "" || len(name) > 50 {
		return false
	}
	for _, char := range bannedSymbols {
		if strings.Contains(name, string(char)) {
			return false
		}
	}
	return true
}

func ParseIds(jsonBuffer string) ([]int, error) {
	ids := []int{}
	if len(jsonBuffer) == 0 {
		return ids, nil
	}
	jsonBuffer = strings.Replace(jsonBuffer, "{", "[", -1)
	jsonBuffer = strings.Replace(jsonBuffer, "}", "]", -1)

    err := json.Unmarshal([]byte(jsonBuffer), &ids)
    return ids, err
}