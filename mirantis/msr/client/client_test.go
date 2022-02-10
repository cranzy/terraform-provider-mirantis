package client_test

import (
	"context"
	"errors"
	"os"
	"testing"

	"github.com/Mirantis/terraform-provider-mirantis/mirantis/msr/client"
)

var (
	host = os.Getenv("MSR_HOST_URL")
	user = os.Getenv("MSR_ADMIN_USER")
	pass = os.Getenv("MSR_ADMIN_PASS")
)

func CreateClientFixture() (client.Client, error) {
	if len(host) <= 0 && len(user) <= 0 && len(pass) <= 0 {
		return client.Client{}, errors.New("Missing required environment variables.")
	}

	return client.NewClient(host, user, pass)
}

func TestMSRClientCreation(t *testing.T) {
	c, err := CreateClientFixture()
	if err != nil {
		t.Error(err)
		return
	}
	ctx := context.Background()
	healthy, err := c.IsHealthy(ctx)

	if !healthy {
		t.Errorf("MSR endpoint is not healthy%s", err)
	}
}
