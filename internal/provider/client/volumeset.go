package cpln

import "fmt"

type VolumeSet struct {
	Base
	Spec        *VolumeSetSpec   `json:"spec,omitempty"`
	SpecReplace *VolumeSetSpec   `json:"$replace/spec,omitempty"`
	Status      *VolumeSetStatus `json:"status,omitempty"`
}

type VolumeSetSpec struct {
	InitialCapacity  *int                `json:"initialCapacity,omitempty"`
	PerformanceClass *string             `json:"performanceClass,omitempty"`
	FileSystemType   *string             `json:"fileSystemType,omitempty"`
	Snapshots        *VolumeSetSnapshots `json:"snapshots,omitempty"`
	AutoScaling      *VolumeSetScaling   `json:"autoscaling,omitempty"`
}

type VolumeSetStatus struct {
	ParentID       *string        `json:"parentID,omitempty"`
	UsedByWorkload *string        `json:"usedByWorkload,omitempty"`
	BindingID      *string        `json:"bindingId,omitempty"`
	Locations      *[]interface{} `json:"locations,omitempty"`
}

type VolumeSetSnapshots struct {
	CreateFinalSnapshot *bool   `json:"createFinalSnapshot,omitempty"`
	RetentionDuration   *string `json:"retentionDuration,omitempty"`
	Schedule            *string `json:"schedule,omitempty"`
}

type VolumeSetScaling struct {
	MaxCapacity       *int     `json:"maxCapacity,omitempty"`
	MinFreePercentage *int     `json:"minFreePercentage,omitempty"`
	ScalingFactor     *float64 `json:"scalingFactor,omitempty"`
}

func (c *Client) GetVolumeSet(name string, gvc string) (*VolumeSet, int, error) {

	volumeSet, code, err := c.GetResource(fmt.Sprintf("gvc/%s/volumeset/%s", gvc, name), new(VolumeSet))
	if err != nil {
		return nil, code, err
	}

	return volumeSet.(*VolumeSet), code, err
}

func (c *Client) CreateVolumeSet(volumeSet VolumeSet, gvc string) (*VolumeSet, int, error) {

	code, err := c.CreateResource(fmt.Sprintf("gvc/%s/volumeset", gvc), *volumeSet.Name, volumeSet)
	if err != nil {
		return nil, code, err
	}

	return c.GetVolumeSet(*volumeSet.Name, gvc)
}

func (c *Client) UpdateVolumeSet(volumeSet VolumeSet, gvc string) (*VolumeSet, int, error) {

	code, err := c.UpdateResource(fmt.Sprintf("gvc/%s/volumeset/%s", gvc, *volumeSet.Name), volumeSet)
	if err != nil {
		return nil, code, err
	}

	return c.GetVolumeSet(*volumeSet.Name, gvc)
}

func (c *Client) DeleteVolumeSet(name string, gvc string) error {
	return c.DeleteResource(fmt.Sprintf("gvc/%s/volumeset/%s", gvc, name))
}
