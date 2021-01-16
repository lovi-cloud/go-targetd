package targetd

import (
	"context"
	"fmt"
)

// Volume is volume
type Volume struct {
	Name string `json:"name"`
	Size int    `json:"size"`
	UUID string `json:"uuid"`
}

// GetVolumeList retrieve list of volume in pool
func (c *Client) GetVolumeList(ctx context.Context, poolName string) ([]Volume, error) {
	method := "vol_list"

	param := map[string]interface{}{
		"pool": poolName,
	}

	req, err := c.newRequest(ctx, method, param)
	if err != nil {
		return nil, fmt.Errorf(ErrCreateRequest+": %w", err)
	}

	var volumes []Volume
	if err := c.request(req, &volumes); err != nil {
		return nil, fmt.Errorf(ErrRequest+": %w", err)
	}

	return volumes, nil
}

// CreateVolume create a volume
func (c *Client) CreateVolume(ctx context.Context, poolName, volumeName string, sizeByte int) error {
	method := "vol_create"

	param := map[string]interface{}{
		"pool": poolName,
		"name": volumeName,
		"size": sizeByte,
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

// DestroyVolume delete a volume
func (c *Client) DestroyVolume(ctx context.Context, poolName, volumeName string) error {
	method := "vol_destroy"

	param := map[string]interface{}{
		"pool": poolName,
		"name": volumeName,
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

// CopyVolume copy a volume.
// you need to set zfs_enable_copy if use ZFS backend.
func (c *Client) CopyVolume(ctx context.Context, poolName, originalVolumeName, newVolumeName string, sizeByte int) error {
	method := "vol_copy"

	param := map[string]interface{}{
		"pool":     poolName,
		"vol_orig": originalVolumeName,
		"vol_new":  newVolumeName,
		"size":     sizeByte,
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

// GetVolume retrieve volume from volume name
func (c *Client) GetVolume(ctx context.Context, poolName, volumeName string) (*Volume, error) {
	volumes, err := c.GetVolumeList(ctx, poolName)
	if err != nil {
		return nil, fmt.Errorf("faeiled to retrieve list of volume: %w", err)
	}

	for _, vol := range volumes {
		if vol.Name == volumeName {
			return &vol, nil
		}
	}

	return nil, fmt.Errorf("%s is not found in vol_list", volumeName)
}
