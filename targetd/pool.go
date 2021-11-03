package targetd

import (
	"context"
	"encoding/json"
	"fmt"
)

// Pool is volume pool
type Pool struct {
	Name     string `json:"name"`
	Size     int64  `json:"size"`
	FreeSize int64  `json:"free_size"`
	Type     string `json:"type"`
	UUID     string `json:"uuid"`
}

// GetPoolList retrieve list of volume pool
func (c *Client) GetPoolList(ctx context.Context) ([]Pool, error) {
	method := "pool_list"

	req, err := c.newRequest(ctx, method, nil)
	if err != nil {
		return nil, fmt.Errorf(ErrCreateRequest+": %w", err)
	}

	type jsonPool struct {
		Name     string `json:"name"`
		Size     int64  `json:"size"`
		FreeSize int64  `json:"free_size"`
		Type     string `json:"type"`
		// targetd response not quoted number, but `uuid` is string
		// ref: https://github.com/open-iscsi/targetd/blob/d694b77c0dd0cc00d72761c6584bbc302d621a04/API.md?plain=1#L25
		UUID json.Number `json:"uuid"`
	}

	var resp []jsonPool
	if err := c.request(req, &resp); err != nil {
		return nil, fmt.Errorf(ErrRequest+": %w", err)
	}

	var pools []Pool
	for _, r := range resp {
		pools = append(pools, Pool{
			Name:     r.Name,
			Size:     r.Size,
			FreeSize: r.FreeSize,
			Type:     r.Type,
			UUID:     r.UUID.String(),
		})
	}

	return pools, nil
}

// GetPool retrieve pool from name of pool
func (c *Client) GetPool(ctx context.Context, poolName string) (*Pool, error) {
	pools, err := c.GetPoolList(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to retrieve list of pool: %w", err)
	}

	for _, pool := range pools {
		if pool.Name == poolName {
			return &pool, nil
		}
	}

	return nil, fmt.Errorf("%s is not found in pool_list", poolName)
}
