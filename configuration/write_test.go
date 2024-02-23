package configuration

import (
	"log"
	"net/url"
	"os"
	"reflect"
	"testing"

	"github.com/glad-dev/gin/constants"
	"github.com/glad-dev/gin/remote"
	rt "github.com/glad-dev/gin/remote/type"
)

func TestMain(m *testing.M) {
	dir, err := os.MkdirTemp("", "")
	if err != nil {
		log.Fatalf("Failed to create temp directory: %s", err)
	}

	err = os.Setenv("XDG_CONFIG_HOME", dir)
	if err != nil {
		log.Fatalf("Failed to set env: %s", err)
	}

	os.Exit(m.Run())
}

func TestWriteRead(t *testing.T) {
	conf := &Config{
		Colors: Colors{
			Blurred: "#AAAAAA",
			Border:  "#BBBBBB",
			Focused: "#CCCCCC",
		},
		Remotes: []Remote{
			{
				URL: url.URL{
					Scheme: "https",
					Host:   "github.com",
				},
				Type: rt.Github,
				Details: []remote.Details{
					{
						Token:     "Github token 1",
						TokenName: "Github token name 1",
						Username:  "Github username 1",
					},
					{
						Token:     "Github token 2",
						TokenName: "Github token name 2",
						Username:  "Github username 2",
					},
				},
			},
			{
				URL: url.URL{
					Scheme: "https",
					Host:   "gitlab.com",
				},
				Type: rt.Gitlab,
				Details: []remote.Details{
					{
						Token:     "Gitlab token",
						TokenName: "Gitlab token name",
						Username:  "Gitlab user name",
					},
				},
			},
		},

		Version: constants.ConfigVersion,
	}

	err := write(conf)
	if err != nil {
		t.Errorf("Failed to write config: %s", err)
	}

	loaded, err := Load()
	if err != nil {
		t.Errorf("Failed to read config: %s", err)
	}

	if !reflect.DeepEqual(*conf, *loaded) {
		t.Errorf("Written and loaded configs differ\nWritten: %+v\n Loaded: %+v", conf, loaded)
	}
}
