package generate

import (
    "github.com/insanebrain/dbp/config"
    "github.com/insanebrain/dbp/utils"
    "github.com/sirupsen/logrus"
    "gopkg.in/alecthomas/kingpin.v2"
)

type IndexConfig struct {
}

func (i *IndexConfig) run(c *kingpin.ParseContext) error {
    allImages := utils.GetAllImagesData(config.Get().CurrentPath)

    templateFile, err := utils.GetIndexTemplate()
    if err != nil {
        logrus.Fatalf("template could not be reach : %s", err)
    }

    for _, mainImageData := range allImages {
        for _, imageData := range allImages {
            if imageData.Parent != nil {
                isParent := mainImageData.Name == imageData.Parent.Name && mainImageData.Tag == imageData.Parent.Tag
                if isParent {
                    imageData.HasParentToBuild = true
                    imageData.HasLocalParent = true
                    imageData.Parent = mainImageData
                    mainImageData.Children = append(mainImageData.Children, imageData)
                }
            }
        }
        if len(mainImageData.Children) > 0 && mainImageData.HasToBuild {
            for _, child := range mainImageData.Children {
                child.HasToBuild = true
            }
        }
    }

    err = utils.GenerateReadmeIndex(allImages, templateFile)

    if err != nil {
        logrus.Fatalf("something went wrong when generate readme index : %s", err)
    }

    return nil
}

func ConfigureIndexCmd(generateCmd *kingpin.CmdClause) {
    indexConfig := &IndexConfig{}

    generateCmd.Command("index", "Generate a readme index").Action(indexConfig.run)
}
