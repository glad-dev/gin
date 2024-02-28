package configuration

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"testing"
	"time"
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

func mockGitlabAPI() {
	http.HandleFunc("/api/v4/personal_access_tokens/self", func(w http.ResponseWriter, r *http.Request) {
		// Set yesterday as "created_at" date
		createdAt := time.Now().AddDate(0, 0, -1).Format("2006-01-02T15:04:05.000Z")
		expiresAt := time.Now().AddDate(0, 0, 1).Format("2006-01-02")

		_, _ = fmt.Fprintf(w, ""+
			"{"+
			"\"id\": 123,"+
			"\"name\": \"Test token\","+
			"\"revoked\": false,"+
			"\"created_at\": \"%s\","+
			"\"scopes\": [\"read_api\",\"read_user\",\"read_repository\"],"+
			"\"user_id\": 456,"+
			"\"active\": true,"+
			"\"expires_at\": \"%s\""+
			"}", createdAt, expiresAt,
		)
	})
	
	err := http.ListenAndServe(":3333", nil)
	if err != nil {
		log.Fatalf("Failed to run server: %s", err)
	}
}
