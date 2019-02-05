package cmd

import (
	"gopkg.in/alecthomas/kingpin.v2"
	"insanebrain/dbp/cmd/build"
)

func ConfigureBuildCmd(app *kingpin.Application) {
	buildConfig := build.Config{}
	buildCmd := app.Command("build", "Build docker images")

	build.ConfigureDirtyCmd(buildCmd, &buildConfig)
	build.ConfigureCommitCmd(buildCmd, &buildConfig)
	buildCmd.Flag("push", "Push image").
		BoolVar(&buildConfig.PushNeeded)
}
