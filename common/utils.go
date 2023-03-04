package common

import log "github.com/sirupsen/logrus"

func AppRecover() {
	if err := recover(); err != nil {
		log.Println("Recovery error:", err)
	}
}

func HasString(arr []string, item string) bool {
	for i := range arr {
		if arr[i] == item {
			return true
		}
	}

	return false
}

func IsAdmin(requester Requester) bool {
	return requester.GetRole() == "admin" || requester.GetRole() == "mod"
}
