package targetd

import (
	"bytes"
	"context"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
)

func (c *Client) newRequest(ctx context.Context, method string, param map[string]interface{}) (*http.Request, error) {
	spath := "/targetrpc"

	u := c.URL
	u.Path = spath

	body, err := marshalRequestBody(method, param)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal body: %w", err)
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodPost, u.String(), body)
	if err != nil {
		return nil, fmt.Errorf("failed to create new HTTP Request: %w", err)
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("User-Agent", userAgent)

	c.setAuthHeader(req)

	return req, nil
}

func (c *Client) setAuthHeader(req *http.Request) {
	data := []byte(fmt.Sprintf("%s:%s", c.User, c.Password))
	enc := base64.StdEncoding.EncodeToString(data)

	req.Header.Set("Authorization", fmt.Sprintf("Basic %s", enc))
}

func marshalRequestBody(method string, param map[string]interface{}) (io.Reader, error) {
	input := map[string]interface{}{
		"id":      1,
		"method":  method,
		"params":  param,
		"jsonrpc": "2.0",
	}

	jb, err := json.Marshal(input)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal json bytes: %w", err)
	}

	return bytes.NewBuffer(jb), nil
}

func (c *Client) request(req *http.Request, out interface{}) error {
	c.Logger.Printf("do request: %+v", req)

	resp, err := c.HTTPClient.Do(req)
	if err != nil {
		return fmt.Errorf("failed to do request: %w", err)
	}

	if err := c.decodeBody(resp, out); err != nil {
		return fmt.Errorf("failed to decode body: %w", err)
	}

	return nil
}

// Response is http Response from targetd
type Response struct {
	Result  interface{} `json:"result,omitempty"`
	Error   ErrorResp   `json:"error,omitempty"`
	ID      int         `json:"id"`
	JSONRPC string      `json:"jsonrpc"`
}

// ErrorResp is error response structure from targetd
type ErrorResp struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
}

// Error is interface for fmt.Error
func (e ErrorResp) Error() error {
	switch e.Code {
	case 0:
		return nil
	}

	return fmt.Errorf("targetd Error: %s (code: %d)", e.Message, e.Code)
}

func (c *Client) decodeBody(resp *http.Response, out interface{}) error {
	defer resp.Body.Close()
	decoder := json.NewDecoder(resp.Body)

	r := &Response{
		Result: out,
	}

	if err := decoder.Decode(r); err != nil {
		return fmt.Errorf("failed to decode response JSON: %w", err)
	}

	if r.Error.Error() != nil {
		return r.Error.Error()
	}

	return nil
}
