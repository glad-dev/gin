package config

import (
	"errors"
)

var ErrUpdateSameValues = errors.New("called update config with existing values")

func Update(wrapper *Wrapper, index int, url string, token string) error {
	if index < 0 || index >= len(wrapper.Configs) {
		return errors.New("update: invalid index")
	}

	// Check if there are any changes
	old := wrapper.Configs[index]
	if old.URL.String() == url && old.Token == token {
		return ErrUpdateSameValues
	}

	u, err := checkURLStr(url)
	if err != nil {
		return err
	}

	lab := GitLab{
		URL:   *u,
		Token: token,
	}

	err = lab.Init()
	if err != nil {
		return err
	}

	wrapper.Configs[index] = lab

	return writeConfig(wrapper)
}
