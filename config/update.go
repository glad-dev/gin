package config

import (
	"errors"
)

func Update(wrapper *Wrapper, index int, url string, token string) error {
	if index < 0 || index >= len(wrapper.Configs) {
		return errors.New("update: invalid index")
	}

	u, err := checkURLStr(url)
	if err != nil {
		return err
	}

	wrapper.Configs[index] = GitLab{
		URL:   *u,
		Token: token,
	}

	return writeConfig(wrapper)
}
