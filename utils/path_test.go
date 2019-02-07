package utils

import (
    "github.com/insanebrain/dbp/model"
    "github.com/stretchr/testify/assert"
    "testing"
)

type testArgs struct {
    currentImage model.ImageData
    destImage    model.ImageData
}

func getTestArgs(destName string, currentRelativeDir string, destRelativeDir string) testArgs {
    return testArgs{
        model.ImageData{
            RelativeDir: currentRelativeDir,
        },
        model.ImageData{
            Name:        destName,
            RelativeDir: destRelativeDir,
        },
    }
}

func Test_Utils_Path_GetReadMeUrl_ShouldReturnAValidUrl(t *testing.T) {

    var tests = []struct {
        msg      string
        param    testArgs
        expected string
    }{
        {
            msg:      "official url",
            param:    getTestArgs("image", "image", ""),
            expected: DockerHubUrl + "_/image",
        },
        {
            msg:      "unofficial url",
            param:    getTestArgs("user/image", "image", ""),
            expected: DockerHubUrl + "r/user/image",
        },
        {
            msg:      "registry url to external",
            param:    getTestArgs("example.fr/image", "image", ""),
            expected: "",
        },
        {
            msg:      "registry url to internal 1 up 1 down",
            param:    getTestArgs("example.fr/image", "image", "otherimage"),
            expected: "../otherimage",
        },
        {
            msg:      "registry url to internal 2 up 1 down",
            param:    getTestArgs("example.fr/image", "namespace/image", "otherimage"),
            expected: "../../otherimage",
        },
        {
            msg:      "registry url to internal 2 up 2 down",
            param:    getTestArgs("example.fr/image", "namespace/image", "namespace/otherimage"),
            expected: "../../namespace/otherimage",
        },
    }

    for _, test := range tests {
        actual, err := GetUrl(test.param.currentImage, test.param.destImage)
        assert.Equal(t, test.expected, actual, test.msg)
        assert.Nil(t, err)
    }

}

func Test_Utils_Path_GetReadMeUrl_ShouldReturnAnError(t *testing.T) {

    test := getTestArgs("", "", "")

    _, err := GetUrl(test.currentImage, test.destImage)
    assert.Error(t, err)

}

func Test_Utils_Path_getDirDepth_ShouldReturnAValidDepth(t *testing.T) {
    var tests = []struct {
        param    string
        expected int
    }{
        {"", 0},
        {"/", 0},
        {"//", 0},
        {"///", 0},
        {"path", 1},
        {"/path", 1},
        {"path/", 1},
        {"/path/", 1},
        {"/path/", 1},
        {"///path///", 1},
        {"path/path", 2},
        {"/path/path", 2},
        {"path/path/", 2},
        {"/path/path/", 2},
        {"///path///path///", 2},
        {"path/path/path", 3},
        {"/path/path/path", 3},
        {"path/path/path/", 3},
        {"/path/path/path/", 3},
        {"///path///path///path///", 3},
    }

    for _, test := range tests {
        assert.Equal(t, test.expected, getDirDepth(test.param))
    }
}

func Test_Utils_Path_GetImageType_ShouldReturnAValidType(t *testing.T) {

    var tests = []struct {
        param    string
        expected ImageType
    }{
        {"image", OFFICIAL},
        {"user/image", UNOFFICIAL},
        {"example.fr/image", REGISTRY},
        {"example.fr/namespace/image", REGISTRY},
        {"example.fr/namespace/subnamespace/image", REGISTRY},
        {"127.0.0.1:5000/docker/docker", REGISTRY},
        {"example.fr:5000/docker/docker", REGISTRY},
    }

    for _, test := range tests {
        imageData := model.ImageData{
            Name: test.param,
        }
        assert.Equal(t, test.expected, GetImageType(imageData), "param ["+test.param+"]")
    }

}

func Test_Utils_Path_GetImageType_ShouldReturnUndefined(t *testing.T) {
    imageData := model.ImageData{
        Name: "",
    }
    assert.Equal(t, UNDEFINED, GetImageType(imageData))
}
