package gce

import (
	"github.com/pkg/errors"
	"google.golang.org/api/compute/v1"
)

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
