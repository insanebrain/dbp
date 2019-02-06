package utils

import (
    "context"
    "encoding/base64"
    "encoding/json"
    "errors"
    "fmt"
    "github.com/docker/distribution/reference"
    "github.com/docker/docker/api/types"
    "github.com/docker/docker/client"
    "github.com/docker/docker/pkg/archive"
    "github.com/docker/docker/pkg/jsonmessage"
    "github.com/docker/docker/pkg/term"
    "github.com/insanebrain/dbp/model"
    "github.com/sirupsen/logrus"
    "io/ioutil"
    "os"
    "path"
    "strings"
)

const defaultDockerAPIVersion = "v1.39"

type AuthConfig struct {
    AuthConfigs map[string]types.AuthConfig `json:"auths,omitempty"`
    HttpHeaders struct {
        UserAgent string `json:"User-Agent,omitempty"`
    }
}

func (authConfig *AuthConfig) GetAuthConfigs() map[string]types.AuthConfig {
    authConfigs := map[string]types.AuthConfig{}

    for hostname, config := range authConfig.AuthConfigs {
        data, err := base64.StdEncoding.DecodeString(config.Auth)
        if err != nil {
            logrus.Debug("cannot decode base64 string from .docker/config.json")
        }

        usernamePassword := strings.SplitN(string(data), ":", 2)
        if len(usernamePassword) != 2 {
            logrus.Debug("base64 string length is more than 2")
        }

        authConfigs[hostname] = types.AuthConfig{
            Username:      usernamePassword[0],
            Password:      usernamePassword[1],
            Auth:          config.Auth,
            ServerAddress: hostname,
        }
    }

    return authConfigs
}

func getDockerClient() *client.Client {

    cli, _ := client.NewClientWithOpts(client.FromEnv, client.WithVersion(defaultDockerAPIVersion))

    return cli
}

func GetAuthConfig() (AuthConfig, error) {
    authConfig := AuthConfig{}
    configFile, err := ioutil.ReadFile(path.Join(os.Getenv("HOME"), ".docker", "config.json"))

    if err != nil {
        return authConfig, err
    }
    err = json.Unmarshal(configFile, &authConfig)

    if err != nil {
        return authConfig, err
    }

    return authConfig, nil
}

// Build the container using the native docker api
func Build(imageData *model.ImageData) error {
    imageDir := imageData.Dir
    tags := imageData.GetTags()

    dockerBuildContext, err := archive.TarWithOptions(imageDir, &archive.TarOptions{})

    defer dockerBuildContext.Close()

    if err != nil {
        return err
    }
    authConfig, _ := GetAuthConfig()

    cli := getDockerClient()

    args := map[string]*string{
    }
    options := types.ImageBuildOptions{
        SuppressOutput: false,
        Remove:         true,
        ForceRemove:    true,
        PullParent:     true,
        Tags:           tags,
        BuildArgs:      args,
        AuthConfigs:    authConfig.GetAuthConfigs(),
    }
    buildResponse, err := cli.ImageBuild(context.Background(), dockerBuildContext, options)
    if err != nil {
        return err
    }
    defer buildResponse.Body.Close()

    termFd, isTerm := term.GetFdInfo(os.Stderr)
    return jsonmessage.DisplayJSONMessagesStream(buildResponse.Body, os.Stderr, termFd, isTerm, nil)
}

func Push(tag string) error {
    cli := getDockerClient()
    authConfig, _ := GetAuthConfig()
    authConfigs := authConfig.GetAuthConfigs()
    ref, err := reference.ParseNormalizedNamed(tag)

    if _, ok := authConfigs[reference.Domain(ref)]; !ok {
        return errors.New(fmt.Sprintf("unable to find docker credential of %s.\n did you forget to docker login ?", reference.Domain(ref)))
    }

    buf, err := json.Marshal(authConfigs[reference.Domain(ref)])
    if err != nil {
        return err
    }
    options := types.ImagePushOptions{
        RegistryAuth: base64.URLEncoding.EncodeToString(buf),
        All:          false,
    }
    pushResponse, err := cli.ImagePush(context.Background(), tag, options)
    if err != nil {
        return err
    }
    defer pushResponse.Close()

    termFd, isTerm := term.GetFdInfo(os.Stderr)
    return jsonmessage.DisplayJSONMessagesStream(pushResponse, os.Stderr, termFd, isTerm, nil)
}
