package cmd

import (
    "github.com/insanebrain/dbp/config"
    "github.com/insanebrain/dbp/utils"
    "gopkg.in/alecthomas/kingpin.v2"
)

type ListConfig struct {
    Name string
}

func (l *ListConfig) run(c *kingpin.ParseContext) error {

    allImages := utils.GetAllImagesData(config.Get().CurrentPath)
    utils.DisplayChildren(utils.SortImages(allImages))

    return nil
}

func ConfigureListCmd(app *kingpin.Application) {
    listConfig := &ListConfig{}

    app.Command("list", "List all images of directory").Action(listConfig.run)
    //listCmd.Arg("name", "Filter with name").StringVar(&listConfig.Name)
}
