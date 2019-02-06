package config

import (
    "context"
    "fmt"
    "github.com/heetch/confita"
    "github.com/heetch/confita/backend"
    "github.com/heetch/confita/backend/env"
    "github.com/heetch/confita/backend/file"
    "os"
    "path/filepath"
)

var defaultPath = "config.yml"

type Config struct {
    Log struct {
        LevelStdout string `config:"LOG_LEVEL_STDOUT" yaml:"levelStdout"`
        PathFile    string `config:"LOG_PATH_FILE" yaml:"pathFile"`
        LevelFile   string `config:"LOG_LEVEL_FILE" yaml:"levelFile"`
    } `yaml:"log"`
    Build struct {
        ExtensionExclude string `config:"BUILD_EXTENSION_EXCLUDE" yaml:"extensionExclude"`
    } `yaml:"build"`
    Template struct {
        ImagePath string `config:"TEMPLATE_IMAGE_PATH" yaml:"imagePath"`
        IndexPath string `config:"TEMPLATE_INDEX_PATH" yaml:"indexPath"`
    } `yaml:"template"`
    CurrentPath string
}

var instance Config
var configLoaded bool

func Get() Config {
    if !configLoaded {
        Load(".", "")
    }
    return instance
}

func Load(currentPath string, configPath string) {
    var backends []backend.Backend
    if configPath == "" {
        configPath = currentPath + string(filepath.Separator) + defaultPath
    }

    if _, err := os.Stat(configPath); !os.IsNotExist(err) {
        backends = append(backends, file.NewBackend(configPath))
    }
    backends = append(backends, env.NewBackend())

    loader := confita.NewLoader(backends...)

    // init default value
    instance = Config{}
    instance.Log.LevelStdout = "info"
    instance.Build.ExtensionExclude = ".md,.txt"
    instance.CurrentPath, _ = filepath.Abs(currentPath)

    err := loader.Load(context.Background(), &instance)
    if err != nil {
        fmt.Println("Cannot load configuration : ", err)
    }

    configLoaded = true
}
