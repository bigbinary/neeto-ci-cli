package main

import (
	"fmt"
	"runtime"

	client "github.com/bigbinary/neeto-ci-cli/api/client"
	"github.com/bigbinary/neeto-ci-cli/cmd"
)

// injected as ldflags during building
var (
	version = "dev"
	commit  = "none"
	date    = "unknown"
)

func main() {
	// inject version information
	cmd.ReleaseVersion = version
	cmd.ReleaseCommit = commit
	cmd.ReleaseDate = date

	// Inject neetoCI User-Agent to identify the CLI in HTTP calls
	client.UserAgent = fmt.Sprintf("neetoCLI/%s (%s; %s; %s; %s; %s)", version, version, commit, date, runtime.GOOS, runtime.GOARCH)

	cmd.Execute()
}
