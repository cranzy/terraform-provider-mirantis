package client

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
)

type Team struct {
	Description  string `json:"description"`
	ID           string `json:"id"`
	MembersCount int    `json:"membersCount"`
	Name         string `json:"name"`
	OrgID        string `json:"orgID"`
}

// CreateTeam creates a team in Enzin
func (c *Client) CreateTeam(ctx context.Context, orgID string, team Team) (Team, error) {
	body, err := json.Marshal(team)
	if err != nil {
		return Team{}, fmt.Errorf("create team failed in MSR client. %w ", err)
	}
	url := fmt.Sprintf("%s/%s/teams", c.createEnziUrl("accounts"), orgID)
	req, err := http.NewRequestWithContext(ctx, http.MethodPost, url, bytes.NewBuffer(body))
	if err != nil {
		return Team{}, fmt.Errorf("request creation failed in MSR client. %w ", err)
	}
	req.Header.Set("Content-Type", "application/json")

	resBody, err := c.doRequest(req)
	if err != nil {
		return Team{}, err
	}

	if err := json.Unmarshal(resBody, &team); err != nil {
		return Team{}, fmt.Errorf("create a team failed in MSR client. %w ", err)
	}

	return team, nil
}

// ReadTeam method retrieves a team from the enzi endpoint
func (c *Client) ReadTeam(ctx context.Context, orgID string, teamID string) (Team, error) {
	url := fmt.Sprintf("%s/%s/teams/%s", c.createEnziUrl("accounts"), orgID, teamID)
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)

	if err != nil {
		return Team{}, err
	}

	body, err := c.doRequest(req)

	if err != nil {
		return Team{}, err
	}

	team := Team{}
	if err := json.Unmarshal(body, &team); err != nil {
		return Team{}, fmt.Errorf("read team failed in MSR client. %w ", err)
	}
	return team, nil
}

// UpdateTeam updates a team in the enzi endpoint
func (c *Client) UpdateTeam(ctx context.Context, orgID string, team Team) (Team, error) {
	url := fmt.Sprintf("%s/%s/teams/%s", c.createEnziUrl("accounts"), orgID, team.ID)

	body, err := json.Marshal(team)
	if err != nil {
		return Team{}, fmt.Errorf("update team failed in MSR client. %w ", err)
	}
	req, err := http.NewRequestWithContext(ctx, http.MethodPatch, url, bytes.NewBuffer(body))

	if err != nil {
		return Team{}, err
	}

	req.Header.Set("Content-Type", "application/json")

	resBody, err := c.doRequest(req)

	if err != nil {
		return Team{}, err
	}

	if json.Unmarshal(resBody, &team) != nil {
		return Team{}, fmt.Errorf("update team failed in MSR client. %w ", err)
	}
	return team, nil
}

// DeleteTeam deletes a team from Enzi
func (c *Client) DeleteTeam(ctx context.Context, orgID string, teamID string) error {
	url := fmt.Sprintf("%s/%s/teams/%s", c.createEnziUrl("accounts"), orgID, teamID)
	req, err := http.NewRequestWithContext(ctx, http.MethodDelete, url, nil)

	if err != nil {
		return fmt.Errorf("delete team failed in MSR client. %w ", err)
	}

	_, err = c.doRequest(req)

	return err
}
