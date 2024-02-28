package configuration

import (
	"errors"
	"fmt"
	"os"
	"reflect"
	"sort"
	"strings"
	"testing"

	"github.com/glad-dev/gin/configuration/location"
	rt "github.com/glad-dev/gin/remote/type"
)

func checkRemote(remote Remote, urlHost string, remoteType rt.Type, tokens []string) error {
	if strings.HasPrefix(urlHost, "https://") {
		urlHost = urlHost[len("https://"):]
	} else if strings.HasPrefix(urlHost, "http://") {
		urlHost = urlHost[len("http://"):]
	}

	if remote.Type != remoteType {
		return fmt.Errorf("remote type mismatch. Expected %s, got %s for tokens: %+v", remoteType.String(), remote.Type.String(), tokens)
	}

	if len(remote.Details) != len(tokens) {
		return fmt.Errorf("details mismatch. Expected '%+v', got %+v\n", tokens, remote.Details)
	}

	if remote.URL.Host != urlHost {
		return fmt.Errorf("URL hosts don't match. Expected '%s', got '%s'", urlHost, remote.URL.Host)
	}

	storedTokens := make([]string, 0, len(remote.Details))
	for _, detail := range remote.Details {
		storedTokens = append(storedTokens, detail.Token)
	}

	// Sort tokens to ensure that they are in the same order
	sort.Strings(tokens)
	sort.Strings(storedTokens)

	if !reflect.DeepEqual(tokens, storedTokens) {
		return fmt.Errorf("unexpected tokens. Expected %+v, got %+v", tokens, storedTokens)
	}

	return nil
}

func TestAppend(t *testing.T) {
	// Delete existing config if it exists
	path, err := location.Get()
	if err != nil {
		t.Fatalf("Failed to get config location: %s", err)
	}

	err = os.Remove(path)
	if err != nil && !errors.Is(err, os.ErrNotExist) {
		t.Fatalf("Failed to remove existing config: %s", err)
	}

	details := [][]struct {
		urlStr     string
		token      string
		remoteType rt.Type
	}{
		{
			{
				urlStr:     "http://127.0.0.1" + gitlabMockPort,
				token:      "Legitimate token gitlab 1",
				remoteType: rt.Gitlab,
			},
			{
				urlStr:     "http://127.0.0.1" + gitlabMockPort,
				token:      "Legitimate token gitlab 2",
				remoteType: rt.Gitlab,
			},
		},
		{
			{
				urlStr:     "http://localhost" + gitlabMockPort,
				token:      "Legitimate token gitlab 1",
				remoteType: rt.Gitlab,
			},
			{
				urlStr:     "http://localhost" + gitlabMockPort,
				token:      "Legitimate token gitlab 2",
				remoteType: rt.Gitlab,
			},
		},
	}

	for i, group := range details {
		tokens := make([]string, 0, 2)

		for _, d := range group {
			tokens = append(tokens, d.token)
			err = Append(d.urlStr, d.token, d.remoteType)
			if err != nil {
				t.Fatalf("Failed to append: %s", err)
			}

			// Check if append successful
			config, err := Load()
			if err != nil {
				t.Fatalf("Failed to load config: %s", err)
			}

			if len(config.Remotes) != (i + 1) {
				t.Fatalf("Config has invalid amount of remotes. Expected %d, got %d", i+1, len(config.Remotes))
			}

			err = checkRemote(config.Remotes[i], d.urlStr, d.remoteType, tokens)
			if err != nil {
				t.Fatalf("Verifying loaded remote: %s", err)
			}
		}
	}
}
