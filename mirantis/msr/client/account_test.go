package client_test

import (
	"context"
	"testing"

	"github.com/Mirantis/terraform-provider-mirantis/mirantis/msr/client"
)

const ORGNAME = "testorg"

// CreateOrgFixture creates a test org in MSR
func CreateOrgFixture(ctx context.Context, c client.Client) (client.Account, error) {
	org := client.Account{
		Name:  ORGNAME,
		IsOrg: true,
	}
	cOrg, err := c.CreateAccount(ctx, org)
	if err != nil {
		return client.Account{}, err
	}
	return cOrg, nil
}

func TestAddAndDeleteUser(t *testing.T) {
	c, err := CreateClientFixture()
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
	c, err := CreateClientFixture()
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
	c, err := CreateClientFixture()
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
	c, err := CreateClientFixture()
	if err != nil {
		t.Error(err)
		return
	}

	ctx := context.Background()

	o, err := CreateOrgFixture(ctx, c)

	if err != nil {
		t.Error(err)
	}
	if o.Name != ORGNAME {
		t.Error("Org created doesn't match expected org name")
	}
	if c.DeleteAccount(ctx, o.ID) != nil {
		t.Error(err)
	}
}
