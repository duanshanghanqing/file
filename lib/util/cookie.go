package util

import (
	"strings"
)

func GetCookie(cookiesStr string) map[string]string {
	cookieMap := make(map[string]string)

	if cookiesStr == "" {
		return cookieMap
	}

	for _, val := range strings.Split(cookiesStr, ";") {
		split := strings.Split(val, "=")
		key := strings.Replace(split[0], " ", "", -1)
		value := split[1]
		cookieMap[key] = value
	}
	return cookieMap
}
