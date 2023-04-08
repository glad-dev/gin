package config

import (
	"errors"
)

var ErrUpdateSameValues = errors.New("called update config with existing values")

func Update(wrapper *Wrapper, wrapperIndex int, detailsIndex int, url string, token string) error {
	if wrapperIndex < 0 || wrapperIndex >= len(wrapper.Configs) {
		return errors.New("update: invalid wrapper index")
	}

	if detailsIndex < 0 || detailsIndex >= len(wrapper.Configs[wrapperIndex].Details) {
		return errors.New("update: invalid details index")
	}

	// Check if there are any changes
	old := wrapper.Configs[wrapperIndex]
	if old.URL.String() == url {
		for _, detail := range old.Details {
			if detail.Token == token {
				return ErrUpdateSameValues
			}
		}
	}

	u, err := checkURLStr(url)
	if err != nil {
		return err
	}

	rd := RepoDetails{
		Token: token,
	}

	err = rd.Init(u)
	if err != nil {
		return err
	}

	wrapper.Configs[wrapperIndex].Details[detailsIndex] = rd

	return writeConfig(wrapper)
}
