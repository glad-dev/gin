package gitlab

import (
	"context"
	"errors"
	"fmt"
	"net/url"

	"github.com/glad-dev/gin/log"
	"github.com/glad-dev/gin/remote/match"
	remotetype "github.com/glad-dev/gin/remote/type"

	"github.com/shurcooL/graphql"
)

type queryUsername struct {
	CurrentUser struct {
		Username graphql.String
	}
}

func getUsername(token string, u *url.URL) (string, error) {
	m := match.Match{
		Token: token,
		URL:   *u,
		Type:  remotetype.Gitlab,
	}

	client, err := m.GraphqlClient()
	if err != nil {
		log.Error("Creating gitlab client",
			"error", err,
			"API-URL", u,
		)

		return "", fmt.Errorf("creating gitlab client: %w", err)
	}

	q := &queryUsername{}
	err = client.Query(context.Background(), q, nil)
	if err != nil {
		log.Error("Querying username",
			"error", err,
		)

		return "", fmt.Errorf("querying username: %w", err)
	}

	if len(q.CurrentUser.Username) == 0 {
		return "", errors.New("empty username")
	}

	return string(q.CurrentUser.Username), nil
}
