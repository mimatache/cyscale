package info

import (
	"time"
)

var version string
var commitHash string
var buildDate string
var appName string

// App contains information about the running application such as name and version
type App struct {
	Name      string `json:"name"`
	Version   string `json:"version"`
	Hash      string `json:"hash"`
	BuildDate string `json:"buildDate"`
}

// Instance contains information about the current instance of the running appplication
type Instance struct {
	Hostname string `json:"Hostname"`
	OS       string `json:"os"`
	Arch     string `json:"arch"`
}

// AppInfo transforms application information into a struct
func AppInfo() App {
	app := App{
		Name:      getAppName(),
		Version:   getVersion(),
		Hash:      getCommitHash(),
		BuildDate: getBuildDate(),
	}
	return app
}

func getVersion() string {
	if version == "" {
		return "dev"
	}
	return version
}

func getCommitHash() string {
	if commitHash == "" {
		return "dev"
	}
	return commitHash
}

func getBuildDate() string {
	if buildDate == "" {
		return time.Now().Format(time.RFC3339)
	}
	return buildDate
}

func getAppName() string {
	if appName == "" {
		return "cyscale-cli"
	}
	return appName
}
