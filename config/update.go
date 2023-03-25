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

	lab := GitLab{
		URL:   *u,
		Token: token,
	}

	err = lab.GetUsername()
	if err != nil {
		return err
	}

	wrapper.Configs[index] = lab

	return writeConfig(wrapper)
}
