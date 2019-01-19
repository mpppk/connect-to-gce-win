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

func GetConfigDirName() string {
	return path.Join(".config", "hlb")
}

func GetConfigFileName() string {
	return ".hlb.yaml"
}

func GetConfigDirPath() (string, error) {
	dir, err := homedir.Dir()
	return path.Join(dir, GetConfigDirName()), errors.Wrap(err, "Error occurred in hlblib.GetConfigDirPath")
}

func GetConfigFilePath() (string, error) {
	configDirPath, err := GetConfigDirPath()
	return path.Join(configDirPath, GetConfigFileName()), errors.Wrap(err, "Error occurred in hlblib.GetConfigFilePath")
}
