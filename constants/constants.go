package constants

const (
	// Version is the tool's version.
	Version = "0.1.0"
	// ConfigVersion is the version of the configuration file.
	ConfigVersion uint8 = 1
)

var (
	// RequiredGitLabScopes is a list with the scopes required for a GitLab repository.
	RequiredGitLabScopes = []string{"read_api", "read_user", "read_repository"}
	// RequiredGitHubScopes is a list with the scopes required for a GitHub repository.
	RequiredGitHubScopes = []string{"repo", "read:user"}
)
