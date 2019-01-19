package gce

import (
	"context"

	"github.com/pkg/errors"
	"golang.org/x/oauth2/google"

	daisyCompute "github.com/GoogleCloudPlatform/compute-image-tools/daisy/compute"
	"github.com/mpppk/gce-auto-connect/password"
	compute "google.golang.org/api/compute/v1"
)

type Client struct {
	daisyClient    daisyCompute.Client
	computeService *compute.Service
	project        string
	zone           string
}

func (c *Client) GetInstance(instanceName string) (*compute.Instance, error) {
	return c.daisyClient.GetInstance(c.project, c.zone, instanceName)
}

func (c *Client) StartInstance(instanceName string) error {
	return c.daisyClient.StartInstance(c.project, c.zone, instanceName)
}

func (c *Client) ResetPassword(instanceName, userName string) (string, error) {
	return password.ResetPassword(c.daisyClient, instanceName, c.zone, c.project, userName)
}

func (c *Client) ListInstances() ([]*compute.Instance, error) {
	instanceList, err := c.computeService.Instances.List(c.project, c.zone).Do()
	if err != nil {
		return nil, errors.Wrap(err, "failed to fetch instances(project: "+
			c.project+", zone: "+c.zone+")")
	}
	return instanceList.Items, nil

}

func NewClient(ctx context.Context, project, zone string) (*Client, error) {
	daisyClient, err := daisyCompute.NewClient(ctx)
	if err != nil {
		return nil, errors.Wrap(err, "failed to create daisy client")
	}
	client, err := google.DefaultClient(ctx)
	if err != nil {
		return nil, errors.Wrap(err, "failed to create default client")
	}
	service, err := compute.New(client)
	if err != nil {
		return nil, errors.Wrap(err, "failed to create new client")
	}
	return &Client{
		daisyClient:    daisyClient,
		computeService: service,
		project:        project,
		zone:           zone,
	}, err
}
