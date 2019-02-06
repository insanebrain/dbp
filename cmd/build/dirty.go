package build

import (
    "github.com/insanebrain/dbp/config"
    "github.com/insanebrain/dbp/utils"
    "github.com/sirupsen/logrus"
    "gopkg.in/alecthomas/kingpin.v2"
)

type DirtyConfig struct {
    BuildConfig *Config
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

    err = utils.BuildImages(imageToBuild)

    if err != nil {
        logrus.Errorf("something went wrong when building : %s", err)
    }

    if d.BuildConfig.PushNeeded {
        err = utils.PushImages(imageToBuild)
    }

    if err != nil {
        logrus.Errorf("something went wrong when pushing : %s", err)
    }

    return nil
}

func ConfigureDirtyCmd(buildCmd *kingpin.CmdClause, buildConfig *Config) {
    dirtyConfig := &DirtyConfig{BuildConfig: buildConfig}

    buildCmd.Command("dirty", "Build docker images for dirty repo").Action(dirtyConfig.run)
}
