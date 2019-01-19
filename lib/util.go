package lib

import (
	"context"
	"io/ioutil"
	"path"

	"github.com/AlecAivazis/survey"
	"github.com/mitchellh/go-homedir"
	"golang.org/x/oauth2/google"

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

func listInstances(project string, zone string) ([]*compute.Instance, error) {
	ctx := context.Background()
	client, err := google.DefaultClient(ctx)
	if err != nil {
		return nil, errors.Wrap(err, "failed to create default client")
	}
	services, err := compute.New(client)
	if err != nil {
		return nil, errors.Wrap(err, "failed to create new client")
	}

	instanceList, err := services.Instances.List(project, zone).Do()
	if err != nil {
		return nil, errors.Wrap(err, "failed to fetch instances")
	}
	return instanceList.Items, nil

}

func ChooseInstance(project string, zone string) (*compute.Instance, error) {
	instances, err := listInstances(project, zone)
	if err != nil {
		return nil, errors.Wrap(err, "failed to list instances")
	}

	if len(instances) <= 0 {
		return nil, errors.New("no instances does found")
	}

	if len(instances) == 1 {
		instance := instances[0]
		return instance, nil
	}

	var instanceNames []string
	for _, instance := range instances {
		instanceNames = append(instanceNames, instance.Name)
	}

	var chooseInstanceName string
	prompt := &survey.Select{
		Message: "Choose instance",
		Options: instanceNames,
	}
	err = survey.AskOne(prompt, &chooseInstanceName, nil)
	if err != nil {
		return nil, errors.Wrap(err, "failed to choose instance")
	}

	for _, instance := range instances {
		if instance.Name == chooseInstanceName {
			return instance, nil
		}
	}

	return nil, errors.New("unexpected error occurred in ChooseInstance")
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
