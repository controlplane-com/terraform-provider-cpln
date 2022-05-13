package cpln

import (
	"fmt"
)

// Agent - Agent
type Agent struct {
	Base
	Status *AgentStatus `json:"status,omitempty"`
}

// AgentStatus - AgentStatus
type AgentStatus struct {
	BootstrapConfig *AgentBootstrapConfig `json:"bootstrapConfig,omitempty"`
}

type AgentBootstrapConfig struct {
	RegistrationToken *string `json:"registrationToken,omitempty"`
	AgentId           *string `json:"agentId,omitempty"`
	AgentLink         *string `json:"agentLink,omitempty"`
	HubEndpoint       *string `json:"hubEndpoint,omitempty"`
}

// GetAgent - Get Agent by name
func (c *Client) GetAgent(name string) (*Agent, int, error) {

	agent, code, err := c.GetResource("agent/"+name, new(Agent))

	if err != nil {
		return nil, code, err
	}

	return agent.(*Agent), code, err
}

// CreateAgent - Create an Agent
func (c *Client) CreateAgent(agent Agent) (*Agent, int, error) {
	return c.CreateResourceAgent(agent)
}

// UpdateAgent - Update an Agent
func (c *Client) UpdateAgent(agent Agent) (*Agent, int, error) {

	code, err := c.UpdateResource(fmt.Sprintf("agent/%s", *agent.Name), agent)
	if err != nil {
		return nil, code, err
	}

	return c.GetAgent(*agent.Name)
}

// DeleteAgent - Delete Agent by name
func (c *Client) DeleteAgent(name string) error {
	return c.DeleteResource(fmt.Sprintf("agent/%s", name))
}
