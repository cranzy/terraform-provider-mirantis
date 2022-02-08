package client

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
)

// HealthResponse structure from MSR
type HealthResponse struct {
	Error   string `json:"error"`
	Healthy bool   `json:"healthy"`
}

// IsHealthy method - checking the MSR health endpoint
func (c *Client) IsHealthy(ctx context.Context) (bool, error) {
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, fmt.Sprintf("%s/health", c.MsrURL), nil)

	if err != nil {
		return false, err
	}

	body, err := c.doRequest(req)
	if err != nil {
		return false, err
	}

	hResponse := &HealthResponse{}

	if err := json.Unmarshal(body, hResponse); err != nil {
		return false, err
	}

	return hResponse.Healthy, nil
}
