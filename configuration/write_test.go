package configuration

import (
	"net/url"
	"reflect"
	"testing"

	"github.com/glad-dev/gin/constants"
	"github.com/glad-dev/gin/remote"
	rt "github.com/glad-dev/gin/remote/type"
)

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
				Details: []remote.Details{
					{
						Token:     "Github token 1",
						TokenName: "Github token name 1",
						Username:  "Github username 1",
						Type:      rt.Github,
					},
					{
						Token:     "Github token 2",
						TokenName: "Github token name 2",
						Username:  "Github username 2",
						Type:      rt.Github,
					},
				},
			},
			{
				URL: url.URL{
					Scheme: "https",
					Host:   "gitlab.com",
				},
				Details: []remote.Details{
					{
						Token:     "Gitlab token",
						TokenName: "Gitlab token name",
						Username:  "Gitlab user name",
						Type:      rt.Gitlab,
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
