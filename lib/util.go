package lib

import (
	"io/ioutil"
	"path"

	"github.com/mitchellh/go-homedir"

	"github.com/pkg/errors"
	"google.golang.org/api/compute/v1"
)

func GenerateRDPFile(filePath, ip, userName string) error {
	contents := "full address:s:" + ip + "\n" + "username:s:" + userName
	return ioutil.WriteFile(filePath, []byte(contents), 0644)
}

func ExtractNatIpFromInstance(instance *compute.Instance) (string, error) {
	if len(instance.NetworkInterfaces) == 0 {
		return "", errors.New("no NetworkInterfaces found")
	}
	accessConfigs := instance.NetworkInterfaces[0].AccessConfigs
	if len(accessConfigs) == 0 {
		return "", errors.New("no AccessConfigs found")
	}
	return accessConfigs[0].NatIP, nil
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

func PanicIfErrorExist(err error) {
	if err != nil {
		panic(err)
	}
}
