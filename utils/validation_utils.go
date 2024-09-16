package utils

import "net/http"

func ValidateAndRespond(w http.ResponseWriter, checks ...func() (int, error)) bool {
	for _, check := range checks {
		status, err := check()
		if err != nil {
			WriteErrorResponse(w, err.Error(), status)
			return false
		}
	}
	return true
}
