package version

var UndefinedVersion = "undefined"
var DevVersion = "dev" // default version set if building without the linker flags
// This will be set by the linker during build
var Version = UndefinedVersion

func IsReleaseVersion() bool {
	return Version != UndefinedVersion && Version != DevVersion
}
