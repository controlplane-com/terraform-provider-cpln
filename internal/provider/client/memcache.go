package cpln

import "fmt"

type Memcache struct {
	Base
	Spec *MemcacheClusterSpec `json:"spec,omitempty"`
}

type MemcacheClusterSpec struct {
	NodeCount   *int             `json:"nodeCount,omitempty"`
	NodeSizeGiB *float64         `json:"nodeSizeGiB,omitempty"`
	Version     *string          `json:"version,omitempty"`
	Options     *MemcacheOptions `json:"options,omitempty"`
	Locations   *[]string        `json:"locations,omitempty"`
}

type MemcacheOptions struct {
	EvictionsDisabled  *bool `json:"evictionsDisabled,omitempty"`
	IdleTimeoutSeconds *int  `json:"idleTimeoutSeconds,omitempty"`
	MaxItemSizeKiB     *int  `json:"maxItemSizeKiB,omitempty"`
	MaxConnections     *int  `json:"maxConnections,omitempty"`
}

func (c *Client) GetMemcache(name string) (*Memcache, int, error) {

	memcache, code, err := c.GetResource(fmt.Sprintf("memcachecluster/%s", name), new(Memcache))
	if err != nil {
		return nil, code, err
	}

	return memcache.(*Memcache), code, err
}

func (c *Client) CreateMemcache(memcache Memcache) (*Memcache, int, error) {

	code, err := c.CreateResource("memcachecluster", *memcache.Name, memcache)
	if err != nil {
		return nil, code, err
	}

	return c.GetMemcache(*memcache.Name)
}

func (c *Client) UpdateMemcache(memcache Memcache) (*Memcache, int, error) {

	code, err := c.UpdateResource(fmt.Sprintf("memcachecluster/%s", *memcache.Name), memcache)
	if err != nil {
		return nil, code, err
	}

	return c.GetMemcache(*memcache.Name)
}

func (c *Client) DeleteMemcache(name string) error {
	return c.DeleteResource(fmt.Sprintf("memcachecluster/%s", name))
}
