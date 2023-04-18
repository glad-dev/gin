package github

import "net/url"

func (hub Details) CheckTokenScope(_ *url.URL) (string, error) { // TODO
	return "GitHub token for account " + hub.Username, nil
}
