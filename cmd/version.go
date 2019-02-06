package cmd

import (
    "fmt"
    "gopkg.in/alecthomas/kingpin.v2"
)

type VersionCmd struct {
    VersionNumber string
}

func (v *VersionCmd) run(c *kingpin.ParseContext) error {
    fmt.Println(v.VersionNumber)

    return nil
}

func ConfigureVersionCmd(app *kingpin.Application, versionNumber string) {
    versionCmd := &VersionCmd{}
    versionCmd.VersionNumber = versionNumber
    app.Command("version", "Display version").Action(versionCmd.run)
}
