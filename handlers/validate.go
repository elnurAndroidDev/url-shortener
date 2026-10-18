package handlers

import (
	"errors"
	"net/url"
	"strings"
)

const maxURLLength = 2048

func validateURL(rawURL string) (string, error) {
	trimmedURL := strings.TrimSpace(rawURL)
	if len(trimmedURL) > maxURLLength {
		return "", errors.New("url is too long")
	}
	u, err := url.ParseRequestURI(trimmedURL)
	if err != nil {
		return "", errors.New("url must be absolute, like https://example.com")
	}
	if u.Scheme != "http" && u.Scheme != "https" {
		return "", errors.New("scheme should be http or https")
	}
	if u.Hostname() == "" {
		return "", errors.New("hostname should not be empty")
	}
	if u.User != nil {
		return "", errors.New("urls with username are not allowed")
	}
	return u.String(), nil
}
