package cmd

import (
	"gopkg.in/alecthomas/kingpin.v2"
	"insanebrain/dbp/cmd/generate"
)

func ConfigureGenerateCmd(app *kingpin.Application) {
	generateCmd := app.Command("generate", "Generate readme of images")

	generate.ConfigureDirtyCmd(generateCmd)
	generate.ConfigureCommitCmd(generateCmd)
	generate.ConfigureAllCmd(generateCmd)
	generate.ConfigureIndexCmd(generateCmd)
}
