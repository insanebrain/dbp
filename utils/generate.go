package utils

import (
    "github.com/insanebrain/dbp/config"
    "github.com/insanebrain/dbp/model"
    "github.com/sirupsen/logrus"
    "os"
    "path/filepath"
    "text/template"
    "time"
)

func GenerateReadmeImages(imagesToGenerate []*model.ImageData, template []byte) error {
    for _, imageData := range imagesToGenerate {
        if imageData.HasToBuild {
            logrus.Debugf("started generating readme %s", imageData.GetFullName())
            err := GenerateReadmeImage(imageData, template)
            logrus.Debugf("end generating readme %s", imageData.GetFullName())
            if err != nil {
                return err
            }

            if imageData.HasLocalParent {
                logrus.Debugf("started generating readme parent of %s", imageData.GetFullName())
                err := GenerateReadmeImage(imageData.Parent, template)
                logrus.Debugf("end generating readme parent of %s", imageData.GetFullName())
                if err != nil {
                    return err
                }
            }
        }

        if len(imageData.Children) > 0 {
            err := GenerateReadmeImages(imageData.Children, template)
            if err != nil {
                return err
            }
        }
    }
    return nil
}

func GenerateReadmeImage(image *model.ImageData, data []byte) error {
    readmePath := image.Dir + string(filepath.Separator) + "README.md"
    additionalVars := template.FuncMap{
        "now": time.Now,
        "getUrl": GetUrl,
        "config": config.Get,
    }

    tmpl, err := template.New("image-readme").Funcs(additionalVars).Parse(string(data))
    if err != nil {
        return err
    }

    file, err := os.Create(readmePath)
    if err != nil {
        return err
    }
    defer file.Close()

    err = tmpl.Execute(file, image)
    if err != nil {
        return err
    }

    return nil
}

func GenerateReadmeIndex(imagesToGenerate []*model.ImageData, data []byte) error {
    readmePath := config.Get().CurrentPath + string(filepath.Separator) + "README.md"
    additionalVars := template.FuncMap{
        "now": time.Now,
        "config": config.Get,
    }

    tmpl, err := template.New("index-readme").Funcs(additionalVars).Parse(string(data))
    if err != nil {
        return err
    }

    file, err := os.Create(readmePath)
    if err != nil {
        return err
    }
    defer file.Close()

    err = tmpl.Execute(file, imagesToGenerate)
    if err != nil {
        return err
    }

    return nil
}
