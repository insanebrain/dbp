package cmd

import (
    "fmt"
    "github.com/insanebrain/dbp/config"
    "gopkg.in/alecthomas/kingpin.v2"
)

type VersionCmd struct {
}

func (v *VersionCmd) run(c *kingpin.ParseContext) error {
    cnf := config.Get()
    fmt.Println(cnf.Version + " [" + cnf.Commit + "]")
    return nil
}

func ConfigureVersionCmd(app *kingpin.Application) {
    versionCmd := &VersionCmd{}
    app.Command("version", "Display version").Action(versionCmd.run)
}
