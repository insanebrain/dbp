package build

import (
    "github.com/insanebrain/dbp/config"
    "github.com/insanebrain/dbp/utils"
    "github.com/sirupsen/logrus"
    "gopkg.in/alecthomas/kingpin.v2"
)

type CommitConfig struct {
    BuildConfig *Config
    CommitId    string
}

func (cc *CommitConfig) run(c *kingpin.ParseContext) error {
    filesUpdated, err := utils.GetCommitFiles(config.Get().CurrentPath, cc.CommitId)

    if err != nil {
        logrus.Fatal(err)
    }

    imageChangedPaths := utils.ExcludeExtFileAndMergePath(filesUpdated)

    allImages := utils.MarkImagesToBuild(imageChangedPaths)

    imageToBuild := utils.SortImages(allImages)

    utils.DisplayChildren(imageToBuild)

    err = utils.BuildImages(imageToBuild)

    if err != nil {
        logrus.Fatalf("something went wrong when building : %s", err)
    }

    if cc.BuildConfig.PushNeeded {
        err = utils.PushImages(imageToBuild)
        if err != nil {
            logrus.Fatalf("something went wrong when pushing : %s", err)
        }
    }


    return nil
}

func ConfigureCommitCmd(buildCmd *kingpin.CmdClause, buildConfig *Config) {
    commitConfig := &CommitConfig{BuildConfig: buildConfig}

    commitCmd := buildCmd.Command("commit", "Build docker images for specific commit").Action(commitConfig.run)
    commitCmd.Arg("commit", "Commit sha").Required().StringVar(&commitConfig.CommitId)
}
