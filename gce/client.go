package gce

import (
	"context"

	daisyCompute "github.com/GoogleCloudPlatform/compute-image-tools/daisy/compute"
	"github.com/mpppk/gce-auto-connect/password"
	compute "google.golang.org/api/compute/v1"
)

type Client struct {
	daisyClient daisyCompute.Client
	project     string
	zone        string
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

func NewClient(ctx context.Context, project, zone string) (*Client, error) {
	daisyClient, err := daisyCompute.NewClient(ctx)
	return &Client{
		daisyClient: daisyClient,
		project:     project,
		zone:        zone,
	}, err
}
