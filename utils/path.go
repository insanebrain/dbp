package utils

import (
    "errors"
    "github.com/insanebrain/dbp/model"
    "path"
    "strings"
)

const DockerHubUrl = "https://hub.docker.com/"

func GetUrl(currentImage model.ImageData, destImage model.ImageData, ) (string, error) {
    switch destImage.GetImageType() {
    case model.OFFICIAL:
        return DockerHubUrl + "_/" + destImage.Name, nil
    case model.UNOFFICIAL:
        return DockerHubUrl + "r/" + destImage.Name, nil
    case model.REGISTRY:
        if destImage.RelativeDir == "" {
            return "", nil
        }
        depth := getDirDepth(currentImage.RelativeDir)
        return path.Join(strings.Repeat("../", depth), destImage.RelativeDir), nil
    default:
        return "", errors.New("unable to define ReadMeUrl (undefined ImageData type)")
    }

}

func getDirDepth(dir string) int {
    trimmed := strings.Trim(dir, "/")
    if trimmed == "" {
        return 0
    }
    joined := path.Join(trimmed)
    return len(strings.Split(joined, "/"))
}