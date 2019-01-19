package lib

import (
	"path"

	"github.com/mitchellh/go-homedir"
	"github.com/pkg/errors"
)

type Config struct {
	UserName     string
	Project      string
	Zone         string
	InstanceName string
}

func (c *Config) Validate() error {
	if c.UserName == "" {
		return errors.New("UserName is not specified")
	}

	if c.Project == "" {
		return errors.New("Project is not specified")
	}

	if c.Zone == "" {
		return errors.New("Zone is not specified")
	}
	return nil
}

func GetConfigDirName() string {
	return path.Join(".config", "connect-to-gce-win")
}

func GetConfigFileName() string {
	return ".connect-to-gce-win.yaml"
}

func GetConfigDirPath() (string, error) {
	dir, err := homedir.Dir()
	return path.Join(dir, GetConfigDirName()), errors.Wrap(err, "Error occurred in GetConfigDirPath")
}

func GetConfigFilePath() (string, error) {
	configDirPath, err := GetConfigDirPath()
	return path.Join(configDirPath, GetConfigFileName()), errors.Wrap(err, "Error occurred in GetConfigFilePath")
}
