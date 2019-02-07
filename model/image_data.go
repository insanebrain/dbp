package model

import (
    "bufio"
    "fmt"
    "github.com/go-yaml/yaml"
    "github.com/insanebrain/dbp/config"
    "github.com/sirupsen/logrus"
    "io/ioutil"
    "os"
    "path/filepath"
    "strings"
)

var DataFilename = "dbp.yml"

type ImageData struct {
    Name             string      `yaml:"name"`
    Tag              string      `yaml:"tag"`
    IsLatest         bool        `yaml:"isLatest"`
    Alias            []ImageData `yaml:"alias"`
    Dir              string
    RelativeDir      string
    Parent           *ImageData
    HasLocalParent   bool
    HasToBuild       bool
    HasParentToBuild bool
    Children         []*ImageData
    EnvVariables     map[string]string `yaml:"envvars"`
    Packages         map[string]string `yaml:"packages"`
}

func (imageData *ImageData) Load(dir string) error {
    imageData.Dir = dir
    imageData.RelativeDir = dir[len(config.Get().CurrentPath)+1:]
    dbpFilePath := imageData.Dir + string(filepath.Separator) + DataFilename
    dbpFile, err := ioutil.ReadFile(dbpFilePath)
    if err != nil {
        return fmt.Errorf("could not load file %s", dbpFilePath)
    }

    err = yaml.Unmarshal(dbpFile, &imageData)
    if err != nil {
        return fmt.Errorf("could not parse %s", dbpFilePath)
    }

    imageData.FindParentData(imageData.Dir)

    logrus.Debugf("Image %s loaded.", imageData.GetFullName())
    return nil
}

func (imageData *ImageData) FindParentData(dir string) {
    imageData.Dir = dir
    file, err := os.Open(imageData.Dir + "/Dockerfile")
    if err != nil {
        logrus.Warn(err)
    }
    defer file.Close()
    scannerDockerFile := bufio.NewScanner(file)

    for scannerDockerFile.Scan() {
        line := scannerDockerFile.Text()
        if strings.Contains(line, "FROM") {
            parent := strings.Split(strings.TrimSpace(strings.Replace(line, "FROM", "", 1)), ":")
            var parentImage ImageData
            parentImage.Name = parent[0]
            if len(parent) > 1 {
                parentImage.Tag = parent[1]
            } else {
                parentImage.Tag = "latest"
            }
            imageData.Parent = &parentImage
            break
        }

        if err := scannerDockerFile.Err(); err != nil {
            logrus.Warn(err)
        }
    }
}

func (imageData ImageData) GetFullName() string {
    return imageData.Name + ":" + imageData.Tag
}

func (imageData ImageData) GetParents() []*ImageData {
    var parents []*ImageData
    if imageData.Parent != nil {
        parents = append(parents, imageData.Parent)
        parents = append(parents, imageData.Parent.GetParents()...)
    }

    return parents
}

func (imageData ImageData) GetAllEnvVar() map[string]string {
    envVars := make(map[string]string)
    var images []*ImageData
    images = append(images, &imageData)
    images = append(images, imageData.GetParents()...)
    for i := len(images) - 1; i >= 0; i-- {
        for k, v := range images[i].EnvVariables {
            envVars[k] = v
        }
    }

    return envVars
}

func (imageData ImageData) GetAllPackages() map[string]string {
    packagesVars := make(map[string]string)
    var images []*ImageData
    images = append(images, &imageData)
    images = append(images, imageData.GetParents()...)
    for i := len(images) - 1; i >= 0; i-- {
        for k, v := range images[i].Packages {
            packagesVars[k] = v
        }
    }

    return packagesVars
}

func (imageData ImageData) GetTags() []string {
    var tags []string

    tags = append(tags, imageData.GetFullName())
    if imageData.IsLatest {
        tags = append(tags, imageData.Name+":latest")
    }
    for _, alias := range imageData.Alias {
        tags = append(tags, alias.GetFullName())
    }
    return tags
}