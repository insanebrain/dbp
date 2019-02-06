package cmd

import (
    "github.com/insanebrain/dbp/cmd/generate"
    "gopkg.in/alecthomas/kingpin.v2"
)

func ConfigureGenerateCmd(app *kingpin.Application) {
    generateCmd := app.Command("generate", "Generate readme of images")

    generate.ConfigureDirtyCmd(generateCmd)
    generate.ConfigureCommitCmd(generateCmd)
    generate.ConfigureAllCmd(generateCmd)
    generate.ConfigureIndexCmd(generateCmd)
}
