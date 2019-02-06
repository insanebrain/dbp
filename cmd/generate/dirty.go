package generate

import (
    "github.com/insanebrain/dbp/config"
    "github.com/insanebrain/dbp/utils"
    "github.com/sirupsen/logrus"
    "gopkg.in/alecthomas/kingpin.v2"
)

type DirtyConfig struct {
}

func (d *DirtyConfig) run(c *kingpin.ParseContext) error {
    filesUpdated, err := utils.GetStatus(config.Get().CurrentPath)

    if err != nil {
        logrus.Error(err)
    }

    imageChangedPaths := utils.ExcludeExtFileAndMergePath(filesUpdated)

    allImages := utils.MarkImagesToBuild(imageChangedPaths)

    imageToBuild := utils.SortImages(allImages)

    utils.DisplayChildren(imageToBuild)

    templateFile, err := utils.GetImageTemplate()
    if err != nil {
        logrus.Fatalf("template could not be reach : %s", err)
    }
    err = utils.GenerateReadmeImages(imageToBuild, templateFile)

    if err != nil {
        logrus.Errorf("something went wrong when generate readme of image : %s", err)
    }

    return nil
}

func ConfigureDirtyCmd(generateCmd *kingpin.CmdClause) {
    dirtyConfig := &DirtyConfig{}

    generateCmd.Command("dirty", "Generate readme of images for dirty repo").Action(dirtyConfig.run)
}
