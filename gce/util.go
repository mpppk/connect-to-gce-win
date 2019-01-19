package gce

import (
	"context"

	"github.com/pkg/errors"
	"golang.org/x/oauth2/google"
	compute "google.golang.org/api/compute/v1"
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
