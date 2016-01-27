package api

import (
	"encoding/json"
	"fmt"
	"github.com/hashicorp/nomad/client/allocdir"
	"io"
	"net/http"
	"net/url"
	"strconv"
)

// AllocFS is used to introspect an allocation directory on a Nomad client
type AllocFS struct {
	client *Client
}

// AllocFS returns an handle to the AllocFS endpoints
func (c *Client) AllocFS() *AllocFS {
	return &AllocFS{client: c}
}

// List is used to list the files at a given path of an allocation directory
func (a *AllocFS) List(alloc *Allocation, path string, q *QueryOptions) ([]*allocdir.AllocFileInfo, *QueryMeta, error) {
	node, _, err := a.client.Nodes().Info(alloc.NodeID, &QueryOptions{})
	if err != nil {
		return nil, nil, err
	}

	if node.HTTPAddr == "" {
		return nil, nil, fmt.Errorf("http addr of the node where alloc %q is running is not advertised", alloc.ID)
	}
	u := &url.URL{
		Scheme: "http",
		Host:   node.HTTPAddr,
		Path:   fmt.Sprintf("/v1/client/fs/ls/%s", alloc.ID),
	}
	v := url.Values{}
	v.Set("path", path)
	u.RawQuery = v.Encode()
	req := &http.Request{
		Method: "GET",
		URL:    u,
	}
	c := http.Client{}
	resp, err := c.Do(req)
	if err != nil {
		return nil, nil, err
	}
	decoder := json.NewDecoder(resp.Body)
	var files []*allocdir.AllocFileInfo
	if err := decoder.Decode(&files); err != nil {
		return nil, nil, err
	}
	return files, nil, nil
}

// Stat is used to stat a file at a given path of an allocation directory
func (a *AllocFS) Stat(alloc *Allocation, path string, q *QueryOptions) (*allocdir.AllocFileInfo, *QueryMeta, error) {
	node, _, err := a.client.Nodes().Info(alloc.NodeID, &QueryOptions{})
	if err != nil {
		return nil, nil, err
	}

	if node.HTTPAddr == "" {
		return nil, nil, fmt.Errorf("http addr of the node where alloc %q is running is not advertised", alloc.ID)
	}
	u := &url.URL{
		Scheme: "http",
		Host:   node.HTTPAddr,
		Path:   fmt.Sprintf("/v1/client/fs/stat/%s", alloc.ID),
	}
	v := url.Values{}
	v.Set("path", path)
	u.RawQuery = v.Encode()
	req := &http.Request{
		Method: "GET",
		URL:    u,
	}
	c := http.Client{}
	resp, err := c.Do(req)
	if err != nil {
		return nil, nil, err
	}
	decoder := json.NewDecoder(resp.Body)
	var file *allocdir.AllocFileInfo
	if err := decoder.Decode(&file); err != nil {
		return nil, nil, err
	}
	return file, nil, nil
}

// ReadAt is used to read bytes at a given offset until limit at the given path
// in an allocation directory
func (a *AllocFS) ReadAt(alloc *Allocation, path string, offset int64, limit int64, w io.Writer, q *QueryOptions) (*QueryMeta, error) {
	node, _, err := a.client.Nodes().Info(alloc.NodeID, &QueryOptions{})
	if err != nil {
		return nil, err
	}

	if node.HTTPAddr == "" {
		return nil, fmt.Errorf("http addr of the node where alloc %q is running is not advertised", alloc.ID)
	}
	u := &url.URL{
		Scheme: "http",
		Host:   node.HTTPAddr,
		Path:   fmt.Sprintf("/v1/client/fs/readat/%s", alloc.ID),
	}
	v := url.Values{}
	v.Set("path", path)
	v.Set("offset", strconv.FormatInt(offset, 10))
	v.Set("limit", strconv.FormatInt(limit, 10))
	u.RawQuery = v.Encode()
	req := &http.Request{
		Method: "GET",
		URL:    u,
	}
	c := http.Client{}
	resp, err := c.Do(req)
	if err != nil {
		return nil, err
	}
	io.Copy(w, resp.Body)
	return nil, nil
}
