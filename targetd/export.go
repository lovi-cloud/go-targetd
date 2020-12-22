package targetd

import (
	"context"
	"fmt"
)

// Export is volume export
type Export struct {
	InitiatorWwn string `json:"initiator_wwn"`
	LUN          int    `json:"lun"`
	VolName      string `json:"vol_name"`
	Pool         string `json:"pool"`
	VolUUID      string `json:"vol_uuid"`
	VolSize      int    `json:"vol_size"`
}

// ListExport retrieve list of export
func (c *Client) ListExport(ctx context.Context) ([]Export, error) {
	method := "export_list"

	req, err := c.newRequest(ctx, method, nil)
	if err != nil {
		return nil, fmt.Errorf(ErrCreateRequest+": %w", err)
	}

	var exports []Export
	if err := c.request(req, &exports); err != nil {
		return nil, fmt.Errorf(ErrRequest+": %w", err)
	}

	return exports, nil
}

// CreateExport create a export
func (c *Client) CreateExport(ctx context.Context, poolName, volumeName string, lunID int, initiatorWWN string) error {
	method := "export_create"

	param := map[string]interface{}{
		"pool":          poolName,
		"vol":           volumeName,
		"lun":           lunID,
		"initiator_wwn": initiatorWWN,
	}

	req, err := c.newRequest(ctx, method, param)
	if err != nil {
		return fmt.Errorf(ErrCreateRequest+": %w", err)
	}

	var i interface{}
	if err := c.request(req, &i); err != nil {
		return fmt.Errorf(ErrRequest+": %w", err)
	}

	return nil
}

// DestroyExport delete a export
func (c *Client) DestroyExport(ctx context.Context, poolName, volumeName, initiatorWWN string) error {
	method := "export_destroy"

	param := map[string]interface{}{
		"pool":          poolName,
		"vol":           volumeName,
		"initiator_wwn": initiatorWWN,
	}

	req, err := c.newRequest(ctx, method, param)
	if err != nil {
		return fmt.Errorf(ErrCreateRequest+": %w", err)
	}

	var i interface{}
	if err := c.request(req, &i); err != nil {
		return fmt.Errorf(ErrRequest+": %w", err)
	}

	return nil
}
