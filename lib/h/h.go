package h

import (
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