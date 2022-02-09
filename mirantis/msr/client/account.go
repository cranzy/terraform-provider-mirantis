package client

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
)

// Account struct
type Account struct {
	Name       string `json:"name"`
	ID         string `json:"id"`
	Password   string `json:"password"`
	FullName   string `json:"fullName"`
	IsActive   bool   `json:"isActive"`
	IsAdmin    bool   `json:"isAdmin"`
	IsOrg      bool   `json:"isOrg"`
	SearchLDAP bool   `json:"searchLDAP"`
}

// Create method - checking the MSR health endpoint
func (c *Client) CreateAccount(ctx context.Context, acc Account) (Account, error) {
	body, err := json.Marshal(acc)
	if err != nil {
		return Account{}, fmt.Errorf("create user failed in MSR client. %w ", err)
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodPost, c.createEnziUrl("accounts"), bytes.NewBuffer(body))
	if err != nil {
		return Account{}, fmt.Errorf("request creation failed in MSR client. %w ", err)
	}
	req.Header.Set("Content-Type", "application/json")

	resBody, err := c.doRequest(req)
	if err != nil {
		return Account{}, err
	}

	if err := json.Unmarshal(resBody, &acc); err != nil {
		return Account{}, fmt.Errorf("create account failed in MSR client. %w ", err)
	}

	return acc, nil
}

// DeleteAccount deletes a user from in Enzi
func (c *Client) DeleteAccount(ctx context.Context, id string) error {
	url := fmt.Sprintf("%s/%s", c.createEnziUrl("accounts"), id)
	req, err := http.NewRequestWithContext(ctx, http.MethodDelete, url, nil)

	if err != nil {
		return fmt.Errorf("delete account failed in MSR client. %w ", err)
	}

	_, err = c.doRequest(req)

	return err
}

// ReadAccount method retrieves a user from the enzi endpoint
func (c *Client) ReadAccount(ctx context.Context, id string) (Account, error) {
	url := fmt.Sprintf("%s/%s", c.createEnziUrl("accounts"), id)
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)

	if err != nil {
		return Account{}, err
	}

	body, err := c.doRequest(req)

	if err != nil {
		return Account{}, err
	}

	acc := Account{}
	if err := json.Unmarshal(body, &acc); err != nil {
		return Account{}, fmt.Errorf("read account failed in MSR client. %w ", err)
	}
	return acc, nil
}

// UpdateAccount updates a user in the enzi endpoint
func (c *Client) UpdateAccount(ctx context.Context, acc Account) (Account, error) {
	url := fmt.Sprintf("%s/%s", c.createEnziUrl("accounts"), acc.ID)

	body, err := json.Marshal(acc)
	if err != nil {
		return Account{}, fmt.Errorf("update user failed in MSR client. %w ", err)
	}
	req, err := http.NewRequestWithContext(ctx, http.MethodPatch, url, bytes.NewBuffer(body))

	if err != nil {
		return Account{}, err
	}

	req.Header.Set("Content-Type", "application/json")

	resBody, err := c.doRequest(req)

	if err != nil {
		return Account{}, err
	}

	if json.Unmarshal(resBody, &acc) != nil {
		return Account{}, fmt.Errorf("update user failed in MSR client. %w ", err)
	}
	return acc, nil
}
