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
		return
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
	user := client.Account{
		Name:     testUserName,
		Password: client.GeneratePass(),
		FullName: "Unit Test",
		IsActive: true,
	}

	ctx := context.Background()

	rUser, err := c.CreateAccount(ctx, user)

	if err != nil {
		t.Error(err)
	}
	if rUser.Name != testUserName {
		t.Error("User created doesn't match expected user")
	}
	if c.DeleteAccount(ctx, testUserName) != nil {
		t.Error(err)
	}
}

func TestUpdateUser(t *testing.T) {
	c, err := createClientFixture()
	if err != nil {
		t.Error(err)
		return
	}
	testUserName := "unittest_updateuser2"

	ctx := context.Background()
	nUser := client.Account{
		Name:     testUserName,
		Password: client.GeneratePass(),
		FullName: "Unit Test",
		IsActive: true,
	}

	testFullName := "Test Updateee"
	testIsActive := false
	cUser, err := c.CreateAccount(ctx, nUser)
	if err != nil {
		t.Error("")
	}
	cUser.FullName = testFullName
	cUser.IsActive = testIsActive
	uUser, err := c.UpdateAccount(ctx, cUser)
	if err != nil {
		t.Error(err)
	} else {
		if uUser.FullName != cUser.FullName {
			t.Error("Full Name didn't update as expected")
		} else if uUser.IsActive != cUser.IsActive {
			t.Error("IsActive doesn't match")
		}
	}
	if c.DeleteAccount(ctx, testUserName) != nil {
		t.Error(err)
	}
}

func TestReadUser(t *testing.T) {
	c, err := createClientFixture()
	if err != nil {
		t.Error(err)
		return
	}
	testUserName := "unittest_readuser"
	nUser := client.Account{
		Name:     testUserName,
		Password: client.GeneratePass(),
		FullName: "Unit Test",
		IsActive: true,
	}
	ctx := context.Background()
	cUser, err := c.CreateAccount(ctx, nUser)
	if err != nil {
		t.Error(err)
	}
	rUser, err := c.ReadAccount(ctx, cUser.Name)
	if err != nil {
		t.Error(err)
	} else {
		if rUser.Name != testUserName {
			t.Error("User created doesn't match expected")
		}
	}
	if c.DeleteAccount(ctx, testUserName) != nil {
		t.Error(err)
	}
}

func TestAddAndDeleteOrg(t *testing.T) {
	c, err := createClientFixture()
	if err != nil {
		t.Error(err)
		return
	}
	testUserName := "unittest3"
	account := client.Account{
		Name:     testUserName,
		IsActive: true,
		IsOrg:    true,
	}

	ctx := context.Background()
	rAccount, err := c.CreateAccount(ctx, account)

	if err != nil {
		t.Error(err)
	}
	if rAccount.Name != testUserName {
		t.Error("User created doesn't match expected user")
	}
	if c.DeleteAccount(ctx, testUserName) != nil {
		t.Error(err)
	}
}
