package buildinfo

import (
	"fmt"
	"os/user"
	"runtime"
	"strings"
)

// BuildInfo contains information about the build
type BuildInfo struct {
	Name      string `json:"name"`
	Version   string `json:"version"`
	BuildDate string `json:"buildDate"`
	GoVersion string `json:"goVersion"`
	OS        string `json:"os"`
	Platform  string `json:"platform"`
}

// set with ldFlags
var (
	name      string
	version   string
	buildDate string
)

var buildInfo BuildInfo

func localBuild() string {
	return fmt.Sprintf("%s-localbuild", getDeveloperName())
}
func getDeveloperName() string {
	usr, err := user.Current()

	if err != nil {
		return "unknown"
	}
	// On Windows the username may contain the domain, which we won't want.
	if strings.Contains(usr.Username, `\`) {
		parts := strings.Split(usr.Username, `\`)
		usr.Username = parts[len(parts)-1]
	}

	return usr.Username
}

// ProvideVersion provides the BuildInfo
func ProvideBuildInfo() BuildInfo {
	return buildInfo
}

func GenerateVersion(appName string) {
	GenerateVersionFromVersionYaml(nil, appName)
}
func GenerateVersionFromVersionYaml(versionyaml []byte, appName string) {
	v, n, b := generate(versionyaml, appName)
	if len(version) == 0 {
		version = v
	}
	if len(name) == 0 {
		name = n
	}
	if len(buildDate) == 0 {
		buildDate = b
	}
	if len(version) == 0 {
		version = localBuild()
	}
	buildInfo = BuildInfo{
		Name:      name,
		Version:   version,
		BuildDate: buildDate,
		GoVersion: runtime.Version(),
		OS:        runtime.GOOS,
		Platform:  runtime.GOARCH,
	}
}
