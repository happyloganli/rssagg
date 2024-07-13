package auth

import (
	"errors"
	"net/http"
	"strings"
)

// GetAPIKey Expected header: Authorization: ApiKey {api_key}
func GetAPIKey(header http.Header) (string, error) {
	val := header.Get("Authorization")
	if val == "" {
		return "", errors.New("no Authorization header found")
	}

	vals := strings.Split(val, " ")
	if len(vals) != 2 {
		return "", errors.New("invalid Authorization header")
	}

	if vals[0] != "ApiKey" {
		return "", errors.New("invalid Authorization header first part")
	}

	return vals[1], nil
}
