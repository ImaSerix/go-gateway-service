package proxy

import (
	"fmt"
	"net/url"
	"slices"
	"strings"
)

var allowedSchemes = []string{
	"http",
	"https",
	"ws",
}

func makeScheme(scheme string) (string, error) {

	trim := strings.TrimSpace(scheme)
	if trim == "" {
		return "", fmt.Errorf("%w: %s", ErrInvalidScheme, "empty scheme")
	}
	if !slices.Contains(allowedSchemes, trim) {
		return "", fmt.Errorf("%w: %s", ErrUnresolverScheme, trim)
	}

	return trim, nil
}

func makeHost(host string) (string, error) {

	trim := strings.TrimSpace(host)
	if trim == "" {
		return "", fmt.Errorf("%w: %s", ErrInvalidHost, "empty host")
	}

	return trim, nil
}

func MakeURL(scheme string, host string, path string) (*url.URL, error) {

	s, err := makeScheme(scheme)
	if err != nil {
		return nil, err
	}

	h, err := makeHost(host)
	if err != nil {
		return nil, err
	}

	u := &url.URL{
		Scheme: s,
		Host:   h,
	}

	trim := strings.TrimSpace(path)
	if trim == "" {
		return u, nil
	}

	for _, p := range strings.Split(trim, "/") {
		u.JoinPath(p)
	}

	return u, nil
}
