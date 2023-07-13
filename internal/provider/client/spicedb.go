package cpln

import "fmt"

type Spicedb struct {
	Base
	Alias  *string        `json:"alias,omitempty"`
	Spec   *ClusterSpec   `json:"spec,omitempty"`
	Status *ClusterStatus `json:"status,omitempty"`
}

type ClusterSpec struct {
	Version   *string   `json:"version,omitempty"`
	Locations *[]string `json:"locations,omitempty"`
}

type ClusterStatus struct {
	ExternalEndpoint *string `json:"externalEndpoint,omitempty"`
}

func (c *Client) GetSpicedb(name string) (*Spicedb, int, error) {

	spicedb, code, err := c.GetResource(fmt.Sprintf("spicedbcluster/%s", name), new(Spicedb))

	if err != nil {
		return nil, code, err
	}

	return spicedb.(*Spicedb), code, err
}

func (c *Client) CreateSpicedb(spicedb Spicedb) (*Spicedb, int, error) {

	code, err := c.CreateResource("spicedbcluster", *spicedb.Name, spicedb)
	if err != nil {
		return nil, code, err
	}

	return c.GetSpicedb(*spicedb.Name)
}

func (c *Client) UpdateSpicedb(spicedb Spicedb) (*Spicedb, int, error) {

	code, err := c.UpdateResource(fmt.Sprintf("spicedbcluster/%s", *spicedb.Name), spicedb)
	if err != nil {
		return nil, code, err
	}

	return c.GetSpicedb(*spicedb.Name)
}

func (c *Client) DeleteSpicedb(name string) error {
	return c.DeleteResource(fmt.Sprintf("spicedbcluster/%s", name))
}
