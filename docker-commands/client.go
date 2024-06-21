package commands

import "github.com/docker/docker/client"

type Client struct {
	c *client.Client
}

func NewClient() (*Client, error) {
	c, err := client.NewClientWithOpts(client.FromEnv, client.WithAPIVersionNegotiation())
	if err != nil {
		return nil, err
	}

	return &Client{c: c}, nil
}

func (c *Client) Close() {
	c.c.Close()
}
