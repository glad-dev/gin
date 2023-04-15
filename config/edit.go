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

	u, err := checkURLStr(url)
	if err != nil {
		return err
	}

	// Check if there are any changes
	old := wrapper.Configs[wrapperIndex]
	if old.URL == *u {
		for _, detail := range old.Details {
			if detail.Token == token {
				return ErrUpdateSameValues
			}
		}
	}

	rd := RemoteDetails{
		Token: token,
	}

	err = rd.Init(u)
	if err != nil {
		return err
	}

	wrapper.Configs[wrapperIndex].Details[detailsIndex] = rd

	return Write(wrapper)
}
