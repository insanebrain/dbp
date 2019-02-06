package utils

import (
    "errors"
    "fmt"
    "github.com/insanebrain/dbp/assets"
    "github.com/insanebrain/dbp/config"
    "github.com/insanebrain/dbp/model"
    "github.com/logrusorgru/aurora"
    "github.com/sirupsen/logrus"
    "github.com/xlab/treeprint"
    "io/ioutil"
    "os"
    "path"
    "path/filepath"
    "strings"
)

func GetAllImagesData(dirPath string) []*model.ImageData {
    var listImages []*model.ImageData
    _ = filepath.Walk(dirPath, func(fp string, fi os.FileInfo, err error) error {

        if err != nil {
            logrus.Errorf("Cannot walk there : %s.", err)
            return nil
        }
        if fi.IsDir() {
            return nil // not a file.  ignore.
        }
        matched, err := filepath.Match(model.DataFilename, fi.Name())
        if err != nil {
            logrus.Errorf("Malformed pattern : %s.", err)
            return err
        }
        if matched {
            var imageData model.ImageData
            logrus.Debugf("Try to load dbp file : %s", path.Dir(fp))
            err = imageData.Load(path.Dir(fp))
            if err != nil {
                logrus.Errorf("Could not load : %s.", err)
            } else {
                listImages = append(listImages, &imageData)
            }
        }
        return nil
    })
    return listImages
}

func FindDockerFile(searchPath string, root string) (string, error) {
    var matched = true

    if !path.IsAbs(root) {
        root, _ = filepath.Abs(root)
    }

    if !path.IsAbs(searchPath) {
        searchPath, _ = filepath.Abs(searchPath)
    }

    if _, err := os.Stat(searchPath + string(filepath.Separator) + "Dockerfile"); os.IsNotExist(err) {
        matched = false
    }

    if matched {
        return searchPath + string(filepath.Separator), nil
    } else if searchPath != root {
        return FindDockerFile(path.Dir(searchPath), root)
    }

    return "", errors.New("dockerfile not find")
}

// exclude extension file + find a dockerfile path for multiple similar sub dir
func ExcludeExtFileAndMergePath(filesUpdated []string) []string {
    var imageChangedPaths []string
    for _, file := range filesUpdated {
        if !Contains(strings.Split(config.Get().Build.ExtensionExclude, ","), filepath.Ext(file)) {
            absolutePath, _ := filepath.Abs(config.Get().CurrentPath + string(filepath.Separator) + path.Dir(file))
            absolutePath, _ = FindDockerFile(absolutePath, config.Get().CurrentPath)
            if absolutePath != "" && !Contains(imageChangedPaths, absolutePath) {
                imageChangedPaths = append(imageChangedPaths, absolutePath)
            }
        }
    }
    return imageChangedPaths
}

// mark images to build
func MarkImagesToBuild(imagePathsToBuild []string) []*model.ImageData {
    allImages := GetAllImagesData(config.Get().CurrentPath)
    for _, dir := range imagePathsToBuild {
        for _, image := range allImages {
            if strings.Contains(dir, image.Dir) {
                image.HasToBuild = true
                break
            }
        }
    }
    return allImages
}

// Sort with parent
func SortImages(imagesToSort []*model.ImageData) []*model.ImageData {
    var imagesSorted []*model.ImageData
    for _, mainImageData := range imagesToSort {
        for _, imageData := range imagesToSort {
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

    for _, imageData := range imagesToSort {
        if !imageData.HasParentToBuild {
            imagesSorted = append(imagesSorted, imageData)
        }
    }

    return imagesSorted
}

func Contains(slice []string, element string) bool {
    for _, item := range slice {
        if item == element {
            return true
        }
    }
    return false
}

func DisplayChildren(imageToDisplay []*model.ImageData) {
    tree := treeprint.New()
    displayChildren(imageToDisplay, tree)
    fmt.Println(tree.String())
}

func displayChildren(children []*model.ImageData, tree treeprint.Tree) () {
    for _, child := range children {

        name := child.GetFullName()
        if child.HasToBuild {
            name = aurora.Green(name).String()
        }

        if len(child.Children) > 0 {
            node := tree.AddBranch(name)
            displayChildren(child.Children, node)
        } else {
            tree.AddNode(name)
        }
    }
}

func BuildImages(imagesToBuild []*model.ImageData) error {
    for _, imageData := range imagesToBuild {
        if imageData.HasToBuild {
            logrus.Debugf("started building %s", imageData.GetFullName())
            err := Build(imageData)
            logrus.Debugf("end building %s", imageData.GetFullName())
            if err != nil {
                return err
            }
        }

        if len(imageData.Children) > 0 {
            err := BuildImages(imageData.Children)
            if err != nil {
                return err
            }
        }
    }
    return nil
}

func PushImages(imagesToBuild []*model.ImageData) error {
    for _, imageData := range imagesToBuild {
        if imageData.HasToBuild {
            for _, tag := range imageData.GetTags() {
                logrus.Debugf("started pushing %s", tag)
                err := Push(tag)
                logrus.Debugf("end pushing %s", tag)
                if err != nil {
                    return err
                }
            }
        }

        if len(imageData.Children) > 0 {
            err := PushImages(imageData.Children)
            if err != nil {
                return err
            }
        }
    }
    return nil
}

func GetImageTemplate() ([]byte, error) {
    var tmpl []byte
    templateFile := config.Get().Template.ImagePath

    if templateFile != "" {
        templateFile, _ = filepath.Abs(config.Get().CurrentPath + string(filepath.Separator) + templateFile)
        file, err := ioutil.ReadFile(templateFile)
        if err != nil {
            return nil, err
        }
        tmpl = file
    } else {
        file, err := assets.Asset("tmpl/image-readme.tmpl")
        if err != nil {
            return nil, err
        }
        tmpl = file
    }
    return tmpl, nil
}

func GetIndexTemplate() ([]byte, error) {
    var tmpl []byte
    templateFile := config.Get().Template.IndexPath

    if templateFile != "" {
        templateFile, _ = filepath.Abs(config.Get().CurrentPath + string(filepath.Separator) + templateFile)
        file, err := ioutil.ReadFile(templateFile)
        if err != nil {
            return nil, err
        }
        tmpl = file
    } else {
        file, err := assets.Asset("tmpl/index-readme.tmpl")
        if err != nil {
            return nil, err
        }
        tmpl = file
    }
    return tmpl, nil
}
