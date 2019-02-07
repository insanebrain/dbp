package model

import (
    "github.com/stretchr/testify/assert"
    "testing"
)

func Test_Model_ImageData_GetImageType_ShouldReturnAValidType(t *testing.T) {

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
        imageData := ImageData{
            Name: test.param,
        }
        assert.Equal(t, test.expected, imageData.GetImageType(), "param ["+test.param+"]")
    }

}

func Test_Model_ImageData_GetImageType_ShouldReturnUndefined(t *testing.T) {
        imageData := ImageData{
            Name: "",
        }
        assert.Equal(t, UNDEFINED, imageData.GetImageType())
}
