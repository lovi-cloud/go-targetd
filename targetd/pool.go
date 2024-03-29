package targetd

import (
	"context"
	"fmt"
)

// Pool is volume pool
type Pool struct {
	Name     string `json:"name"`
	Size     int64  `json:"size"`
	FreeSize int64  `json:"free_size"`
	Type     string `json:"type"`
	UUID     int64  `json:"uuid"`
}

// GetPoolList retrieve list of volume pool
func (c *Client) GetPoolList(ctx context.Context) ([]Pool, error) {
	method := "pool_list"

	req, err := c.newRequest(ctx, method, nil)
	if err != nil {
		return nil, fmt.Errorf(ErrCreateRequest+": %w", err)
	}

	var pools []Pool
	if err := c.request(req, &pools); err != nil {
		return nil, fmt.Errorf(ErrRequest+": %w", err)
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
