//go:generate go-bindata -pkg assets -o assets/tmpl.go tmpl/
package main

import (
    "fmt"
    "github.com/insanebrain/dbp/cmd"
    "github.com/insanebrain/dbp/config"
    "github.com/sirupsen/logrus"
    "gopkg.in/alecthomas/kingpin.v2"
    "os"
)

var (
    app     = kingpin.New("dbp", "A command-line dbp helper.")
    version = "master"
)

func main() {
    kingpin.Version(version)

    // Config flag to override config file path
    var configFile string
    configFileFlag := app.Flag("config", "Define specific configuration file")
    configFileFlag.Short('c')
    configFileFlag.StringVar(&configFile)

    // Config flag to override config file path
    var currentPath string
    currentPathFlag := app.Flag("path", "Define execution path")
    currentPathFlag.Short('p')
    currentPathFlag.StringVar(&currentPath)
    currentPathFlag.Default(".")

    // Add commands
    cmd.ConfigureBuildCmd(app)
    cmd.ConfigureListCmd(app)
    cmd.ConfigureGenerateCmd(app)
    cmd.ConfigureVersionCmd(app, version)
    app.PreAction(func(context *kingpin.ParseContext) error {
        // init config
        config.Load(currentPath, string(configFile))
        config.Get()

        // init logger
        levelStdoutParsed, err := logrus.ParseLevel(config.Get().Log.LevelStdout)

        if err != nil {
            fmt.Println("Could not parse log level stdout")
        }

        logrus.SetOutput(os.Stdout)
        logrus.SetLevel(levelStdoutParsed)

        //if config.Get().Log.PathFile != "" {
        //	levelFileParsed, err := log.ParseLevel(config.Get().Log.LevelFile)
        //	log.AddHook()
        //	if err != nil {
        //		fmt.Println("Could not parse log level stdout")
        //	}
        //}
        logrus.Debugf("Configuration loaded.")

        return nil
    })

    kingpin.MustParse(app.Parse(os.Args[1:]))
}
