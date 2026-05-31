package proxy

import (
	"net/url"
)

type URL string

func NewURL(raw string) (URL, error) {
	u, err := url.ParseRequestURI(raw)
	if err != nil {
		return "", err
	}

	if u.Scheme != "http" && u.Scheme != "https" {
		return "", ErrUnsupportedScheme
	}

	if u.Host == "" {
		return "", ErrEmptyHost
	}

	return URL(u.String()), nil
}
