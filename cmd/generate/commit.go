package generate

import (
    "github.com/insanebrain/dbp/config"
    "github.com/insanebrain/dbp/utils"
    "github.com/sirupsen/logrus"
    "gopkg.in/alecthomas/kingpin.v2"
)

type CommitConfig struct {
    CommitId string
}

func (cc *CommitConfig) run(c *kingpin.ParseContext) error {
    filesUpdated, err := utils.GetCommitFiles(config.Get().CurrentPath, cc.CommitId)

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

func ConfigureCommitCmd(generateCmd *kingpin.CmdClause) {
    commitConfig := &CommitConfig{}

    commitCmd := generateCmd.Command("commit", "Generate readme of images for specific commit").Action(commitConfig.run)
    commitCmd.Arg("commit", "Commit sha").Required().StringVar(&commitConfig.CommitId)
}
