package configuration

import (
	"errors"
	"net/url"

	"github.com/glad-dev/gin/log"
)

func checkURLStr(urlStr string) (*url.URL, error) {
	u, err := url.ParseRequestURI(urlStr)
	if err != nil {
		log.Error("URL is invalid.", "url", urlStr)

		return nil, err
	}

	if !u.IsAbs() {
		log.Error("URL is not absolute.", "url", u.String())

		return nil, errors.New("URL is not absolute")
	}

	if u.Scheme != "https" && u.Scheme != "http" {
		log.Error("URL has invalid scheme.", "url", u.String())

		return nil, errors.New("URL has invalid scheme")
	}

	return u, nil
}
