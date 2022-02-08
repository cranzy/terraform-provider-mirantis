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

func createClientFixture() (client.Client, error) {
	if len(host) <= 0 && len(user) <= 0 && len(pass) <= 0 {
		return client.Client{}, errors.New("Missing required environment variables.")
	}

	return client.NewClient(host, user, pass)
}

func TestMSRClientCreation(t *testing.T) {
	c, err := createClientFixture()
	if err != nil {
		t.Error(err)
	}
	ctx := context.Background()
	healthy, err := c.IsHealthy(ctx)

	if !healthy {
		t.Errorf("MSR endpoint is not healthy%s", err)
	}
}

func TestAddAndDeleteUser(t *testing.T) {
	c, err := createClientFixture()
	if err != nil {
		t.Error(err)
		return
	}
	testUserName := "unittest3"
	user := client.User{
		Name:     testUserName,
		Password: client.GeneratePass(),
		FullName: "Unit Test",
		IsActive: true,
	}

	ctx := context.Background()
	rUser, err := c.CreateUser(ctx, user)

	if err != nil {
		t.Error(err)
	}
	if rUser.Name != testUserName {
		t.Error("User created doesn't match expected user")
	}
	err = c.DeleteUser(ctx, testUserName)
	if err != nil {
		t.Error(err)
	}
}

func TestUpdateUser(t *testing.T) {
	c, err := createClientFixture()
	if err != nil {
		t.Error(err)
	}
	testUserName := "unittest_updateuser2"

	ctx := context.Background()
	cUser := client.User{
		Name:     testUserName,
		Password: client.GeneratePass(),
		FullName: "Unit Test",
		IsActive: true,
	}

	testFullName := "Test Updateee"
	testIsActive := false
	cUser, err = c.CreateUser(ctx, cUser)
	if err != nil {
		t.Error("")
	}
	cUser.FullName = testFullName
	cUser.IsActive = testIsActive
	uUser, err := c.UpdateUser(ctx, cUser)
	if err != nil {
		t.Error(err)
	} else {
		if uUser.FullName != cUser.FullName {
			t.Error("Full Name didn't update as expected")
		} else if uUser.IsActive != cUser.IsActive {
			t.Error("IsActive doesn't match")
		}
	}
	err = c.DeleteUser(ctx, testUserName)
	if err != nil {
		t.Error(err)
	}
}

func TestReadUser(t *testing.T) {
	c, err := createClientFixture()
	if err != nil {
		t.Error(err)
	}
	testUserName := "unittest_readuser"
	cUser := client.User{
		Name:     testUserName,
		Password: client.GeneratePass(),
		FullName: "Unit Test",
		IsActive: true,
	}
	ctx := context.Background()
	rUser, err := c.CreateUser(ctx, cUser)
	if err != nil {
		t.Error(err)
	}
	rUser, err = c.ReadUser(ctx, rUser.Name)
	if err != nil {
		t.Error(err)
	} else {
		if rUser.Name != testUserName {
			t.Error("User created doesn't match expected")
		}
	}
	err = c.DeleteUser(ctx, testUserName)
	if err != nil {
		t.Error(err)
	}
}
