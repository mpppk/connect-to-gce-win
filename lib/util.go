package lib

import (
	"io/ioutil"

	"github.com/AlecAivazis/survey"

	"github.com/pkg/errors"
	"google.golang.org/api/compute/v1"
)

func GenerateRDPFile(filePath, ip, userName string) error {
	contents := "full address:s:" + ip + "\n" + "username:s:" + userName
	return ioutil.WriteFile(filePath, []byte(contents), 0644)
}

func ChooseInstance(instances []*compute.Instance, project string, zone string) (*compute.Instance, error) {
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
	err := survey.AskOne(prompt, &chooseInstanceName, nil)
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

func PanicIfErrorExist(err error) {
	if err != nil {
		panic(err)
	}
}
