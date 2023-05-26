package config

import (
	"errors"
	"fmt"
	"net/url"

	"github.com/glad-dev/gin/logger"
	"github.com/glad-dev/gin/remote"
)

// Remote contains the remote's URL and a list of Details, containing the token, username and token name.
type Remote struct {
	URL     url.URL
	Details []remote.Details
}

// ToMatch casts the remote to a remote.Match if the remote contains one Details. An error is returned if there are none
// or more than one Details.
func (r *Remote) ToMatch() (*remote.Match, error) {
	if len(r.Details) == 0 {
		logger.Log.Error("Remote contains no details.")

		return nil, errors.New("failed to convert Remote to Match since Remote.Details contains no elements")
	}

	if len(r.Details) == 1 {
		return &remote.Match{
			URL:       r.URL,
			Token:     r.Details[0].GetToken(),
			Username:  r.Details[0].GetUsername(),
			TokenName: r.Details[0].GetTokenName(),
		}, nil
	}

	logger.Log.Info("Got a remote with multiple details.")

	return &remote.Match{
		URL: r.URL,
	}, errors.New("config contains multiple matching configs")
}

// ToMatchAtIndex casts the remote at index to a remote.Match.
func (r *Remote) ToMatchAtIndex(index int) (*remote.Match, error) {
	if index < 0 || index >= len(r.Details) {
		logger.Log.Error("Invalid index.", "index", index, "len(details)", len(r.Details))

		return nil, errors.New("ToMatchAtIndex: invalid index")
	}

	return &remote.Match{
		URL:       r.URL,
		Token:     r.Details[index].GetToken(),
		Username:  r.Details[index].GetUsername(),
		TokenName: r.Details[index].GetTokenName(),
	}, nil
}

func (r *Remote) checkSemantics() error {
	// Check URL
	_, err := checkURLStr(r.URL.String())
	if err != nil {
		return err
	}

	if len(r.Details) == 0 {
		return fmt.Errorf("config contains no details")
	}

	for _, details := range r.Details {
		if len(details.GetUsername()) == 0 {
			return fmt.Errorf("config contains empty username")
		}

		// Check if token is semantically correct. The tokens validity is not checked.
		if len(details.GetToken()) == 0 { // TODO: Get actual sizes
			return fmt.Errorf("config contains empty token")
		}

		if len(details.GetTokenName()) == 0 {
			return fmt.Errorf("config contains empty token name")
		}
	}

	return nil
}
