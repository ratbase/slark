package version

// Version is the current version of the application.
// This will be overwritten at build time by the -X linker flag.
var Version = "dev"

// GetVersion returns the current version of the application.
func GetVersion() string {
	return Version
}
