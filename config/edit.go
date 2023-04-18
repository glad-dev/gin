package config

import (
	"errors"

	"gn/config/remote"
	"gn/logger"
)

var ErrUpdateSameValues = errors.New("called update config with existing values")

func Update(wrapper *Wrapper, wrapperIndex int, detailsIndex int, url string, token string) error {
	if wrapperIndex < 0 || wrapperIndex >= len(wrapper.Remotes) {
		logger.Log.Error("Wrapper index is invalid.", "index", wrapperIndex, "len(remotes)", len(wrapper.Remotes))

		return errors.New("update: invalid wrapper index")
	}

	if detailsIndex < 0 || detailsIndex >= len(wrapper.Remotes[wrapperIndex].Details) {
		logger.Log.Error("Details index is invalid.", "index", detailsIndex, "len(details)", len(wrapper.Remotes[wrapperIndex].Details))

		return errors.New("update: invalid details index")
	}

	u, err := checkURLStr(url)
	if err != nil {
		return err
	}

	// Check if there are any changes
	old := wrapper.Remotes[wrapperIndex]
	if old.URL == *u {
		for _, detail := range old.Details {
			if detail.GetToken() == token {
				logger.Log.Warn("Attempted to update the remote '%s' with the same token", u.String())

				return ErrUpdateSameValues
			}
		}
	}

	var rd remote.Details
	if u.Host == "github.com" {
		rd = remote.GitHubDetails{
			Token: token,
		}
	} else {
		rd = remote.GitLabDetails{
			Token: token,
		}
	}

	rd, err = rd.Init(u)
	if err != nil {
		logger.Log.Errorf("Failed to initialize token: %s", err)

		return err
	}

	wrapper.Remotes[wrapperIndex].Details[detailsIndex] = rd

	return Write(wrapper)
}
