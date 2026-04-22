package main

import (
	"github.com/gongdaowen/apix/cmd"
)

// Version is set at build time via ldflags
var Version = "dev"
var BuildTime = "unknown"
var CommitHash = "unknown"

func main() {
	cmd.SetVersion(Version)
	cmd.SetBuildInfo(BuildTime, CommitHash)
	cmd.Execute()
}