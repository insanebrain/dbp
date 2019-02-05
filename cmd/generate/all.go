package generate

import (
	"github.com/sirupsen/logrus"
	"gopkg.in/alecthomas/kingpin.v2"
	"insanebrain/dbp/config"
	"insanebrain/dbp/utils"
)

type AllConfig struct {
}

func (d *AllConfig) run(c *kingpin.ParseContext) error {
	allImages := utils.GetAllImagesData(config.Get().CurrentPath)

	templateFile, err := utils.GetImageTemplate()
	if err != nil {
		logrus.Fatal("template could not be reach : %s", err)
	}
	for _, image := range allImages {
		image.HasToBuild = true
	}
	imagesSorted := utils.SortImages(allImages)
	utils.DisplayChildren(imagesSorted)

	for _, image := range allImages {
		logrus.Debugf("started generating readme %s", image.GetFullName())
		err := utils.GenerateReadmeImage(image, templateFile)
		logrus.Debugf("end generating readme %s", image.GetFullName())

		if err != nil {
			logrus.Errorf("something went wrong when generate readme of image : %s", err)
		}
	}

	return nil
}

func ConfigureAllCmd(generateCmd *kingpin.CmdClause) {
	dirtyConfig := &AllConfig{}

	generateCmd.Command("all", "Generate readme of images for all").Action(dirtyConfig.run)
}
