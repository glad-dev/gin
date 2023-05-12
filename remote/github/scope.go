package github

import "net/url"

// CheckTokenScope is not implemented for a GitHub host.
func (hub Details) CheckTokenScope(_ *url.URL) (string, error) { // TODO: Implement this
	return "GitHub token for account " + hub.Username, nil
}
