package config

import (
	"errors"
	"net/url"

	"github.com/glad-dev/gin/logger"
)

func checkURLStr(urlStr string) (*url.URL, error) {
	u, err := url.ParseRequestURI(urlStr)
	if err != nil {
		logger.Log.Error("URL is invalid.", "url", urlStr)

		return nil, err
	}

	if !u.IsAbs() {
		logger.Log.Error("URL is not absolute.", "url", u.String())

		return nil, errors.New("URL is not absolute")
	}

	if u.Scheme != "https" && u.Scheme != "http" {
		logger.Log.Error("URL has invalid scheme.", "url", u.String())

		return nil, errors.New("URL has invalid scheme")
	}

	return u, nil
}
