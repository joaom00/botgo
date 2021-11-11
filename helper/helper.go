package helper

import "strings"

func Find(arr []string, el string) bool {
	for _, v := range arr {
		if v == strings.ToUpper(el) {
			return true
		}
	}

	return false
}
