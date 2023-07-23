package app

import (
	"fmt"
	"os"

	"gopkg.in/yaml.v2"
)

const (
	PATH      = "."
	FILE_PATH = ".trunk"
)

var App *AppConfig

func InitApp() {
	App = &AppConfig{}
	App.init()
}

type Config struct {
	ReleaseCanidateName string `yaml:"rc-name"`
	ReleaseName         string `yaml:"release-name"`
}

type AppConfig struct {
	Config Config
}

func (app *AppConfig) IsEmpty() bool {
	return len(app.Config.ReleaseName) == 0 && len(app.Config.ReleaseCanidateName) == 0
}

func (app *AppConfig) Save(releaseName string, rcName string) error {
	c := Config{ReleaseName: releaseName,
		ReleaseCanidateName: rcName}
	data, err := yaml.Marshal(&c)
	if err != nil {
		return err
	}

	fmt.Println(string(data[:]))
	err = os.WriteFile(FILE_PATH, data, 0)
	if err != nil {
		return err
	}
	app.Config = c
	return nil
}

func (app *AppConfig) init() {
	f, err := os.Open(FILE_PATH)
	if err != nil {
		f, err = os.Create(FILE_PATH)
		if err != nil {
			fmt.Println("")
		}
	}
	defer f.Close()

	var cfg Config
	decoder := yaml.NewDecoder(f)
	err = decoder.Decode(&cfg)
	if err != nil {
		fmt.Println("")
	}
	app.Config = cfg
}
