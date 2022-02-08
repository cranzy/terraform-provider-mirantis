package client

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
)

// User struct
type User struct {
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
func (c *Client) CreateUser(ctx context.Context, user User) (User, error) {
	body, err := json.Marshal(user)
	if err != nil {
		return User{}, fmt.Errorf("create user failed in MSR client. %w ", err)
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodPost, c.createEnziUrl("accounts"), bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")

	if err != nil {
		return User{}, fmt.Errorf("request creation failed in MSR client. %w ", err)
	}

	body, err = c.doRequest(req)

	if err != nil {
		return User{}, err
	}

	if err := json.Unmarshal(body, &user); err != nil {
		return User{}, fmt.Errorf("create user failed in MSR client. %w ", err)
	}

	return user, nil
}

// DeleteUser deletes a user from in Enzi
func (c *Client) DeleteUser(ctx context.Context, id string) error {
	url := fmt.Sprintf("%s/%s", c.createEnziUrl("accounts"), id)
	req, err := http.NewRequestWithContext(ctx, http.MethodDelete, url, nil)

	if err != nil {
		return fmt.Errorf("delete user failed in MSR client. %w ", err)
	}

	_, err = c.doRequest(req)

	if err != nil {
		return err
	}

	return nil
}

// ReadUser method retrieves a user from the enzi endpoint
func (c *Client) ReadUser(ctx context.Context, name string) (User, error) {
	url := fmt.Sprintf("%s/%s", c.createEnziUrl("accounts"), name)
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)

	if err != nil {
		return User{}, err
	}

	body, err := c.doRequest(req)

	if err != nil {
		return User{}, err
	}

	user := User{}
	if err := json.Unmarshal(body, &user); err != nil {
		return User{}, fmt.Errorf("get user failed in MSR client. %w ", err)
	}
	return user, nil
}

// UpdateUser updates a user in the enzi endpoint
func (c *Client) UpdateUser(ctx context.Context, user User) (User, error) {
	url := fmt.Sprintf("%s/%s", c.createEnziUrl("accounts"), user.ID)

	body, err := json.Marshal(user)
	if err != nil {
		return User{}, fmt.Errorf("update user failed in MSR client. %w ", err)
	}
	req, err := http.NewRequestWithContext(ctx, http.MethodPatch, url, bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")

	if err != nil {
		return User{}, err
	}

	body, err = c.doRequest(req)

	if err != nil {
		return User{}, err
	}

	if json.Unmarshal(body, &user) != nil {
		return User{}, fmt.Errorf("update user failed in MSR client. %w ", err)
	}
	return user, nil
}
