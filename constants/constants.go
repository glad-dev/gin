package constants

const (
	Version             = "1.0.0"
	ConfigVersion uint8 = 1
	ProgramName         = "gn" // TODO: Remove this once name is set
)

var RequiredScopes = []string{"read_api", "read_user", "read_repository"}
