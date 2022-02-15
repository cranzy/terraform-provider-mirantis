package client_test

import (
	"context"
	"testing"

	"github.com/Mirantis/terraform-provider-mirantis/mirantis/msr/client"
)

func TestCRUDTeam(t *testing.T) {
	c, err := CreateClientFixture()
	if err != nil {
		t.Error(err)
		return
	}
	ctx := context.Background()

	o, err := CreateOrgFixture(ctx, c)
	if err != nil {
		t.Error(err)
		return
	}

	team := client.Team{
		Name:        "ReadOnly",
		Description: "Test team",
	}
	cTeam, err := c.CreateTeam(ctx, o.ID, team)
	if err != nil {
		t.Error(err)
	}
	team.Description = "Updated Description"
	team.ID = cTeam.ID
	uTeam, err := c.UpdateTeam(ctx, o.ID, team)
	if err != nil {
		t.Error(err)
	}
	rTeam, err := c.ReadTeam(ctx, o.ID, cTeam.ID)
	if err != nil {
		t.Error(err)
	}
	if rTeam.Name != uTeam.Name {
		t.Error("Team name doesn't match expected!")
	}
	if rTeam.Description != uTeam.Description {
		t.Error("Description doesn't match expected!")
	}
	if err := c.DeleteTeam(ctx, cTeam.OrgID, cTeam.ID); err != nil {
		t.Error(err)
	}
	if err := c.DeleteAccount(ctx, cTeam.OrgID); err != nil {
		t.Error(err)
	}
}

func TestAddUserToTeam(t *testing.T) {
	c, err := CreateClientFixture()
	if err != nil {
		t.Error(err)
		return
	}
	ctx := context.Background()

	o, err := CreateOrgFixture(ctx, c)
	if err != nil {
		t.Error(err)
		return
	}

	team := client.Team{
		Name:        "Blah",
		Description: "Test team",
	}
	cTeam, err := c.CreateTeam(ctx, o.ID, team)
	if err != nil {
		t.Error(err)
	}

	user := client.Account{
		Name:     "abc",
		Password: client.GeneratePass(),
		IsActive: true,
	}

	cUser, err := c.CreateAccount(ctx, user)
	if err != nil {
		t.Error(err)
	}
	if err := c.AddUserToTeam(ctx, o.ID, cTeam.ID, cUser); err != nil {
		t.Error("Couldn't add user to team: %w", err)
	}

	c.DeleteUserFromTeam(ctx, o.ID, cTeam.ID, cUser.ID)
	_, err = c.GetTeamUsers(ctx, o.ID, cTeam.ID)
	if err != nil {
		t.Error(err)
	}

	if err := c.DeleteAccount(ctx, cUser.ID); err != nil {
		t.Error(err)
	}
	if err := c.DeleteTeam(ctx, cTeam.OrgID, cTeam.ID); err != nil {
		t.Error(err)
	}
	if err := c.DeleteAccount(ctx, cTeam.OrgID); err != nil {
		t.Error(err)
	}
}
