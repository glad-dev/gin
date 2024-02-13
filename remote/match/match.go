package match

import (
	"context"
	"fmt"
	"net/http"
	"net/url"

	"github.com/glad-dev/gin/log"
	rt "github.com/glad-dev/gin/remote/type"

	"github.com/shurcooL/graphql"
	"golang.org/x/oauth2"
)

// Match contains the information needed to request data from a remote.
type Match struct {
	URL       url.URL
	Token     string
	Username  string
	TokenName string
	Type      rt.Type
}

func (m *Match) GraphqlClient() (*graphql.Client, error) {
	var tc *http.Client
	if len(m.Token) > 0 {
		ctx := context.Background()
		tc = oauth2.NewClient(ctx, oauth2.StaticTokenSource(
			&oauth2.Token{AccessToken: m.Token},
		))
	}

	apiURL, err := m.Type.GraphqlAPIURL(&m.URL)
	if err != nil {
		log.Error("Failed to get API URL", "error", err, "match-url", m.URL.String())

		return nil, fmt.Errorf("invalid API url: %w", err)
	}

	return graphql.NewClient(apiURL, tc), nil
}
