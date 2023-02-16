package cpln

import (
	"fmt"
)

type AuditContext struct {
	Base
}

// GetAuditContext - Get Audit Context by name
func (c *Client) GetAuditContext(name string) (*AuditContext, int, error) {

	auditCtx, code, err := c.GetResource(fmt.Sprintf("auditctx/%s", name), new(AuditContext))
	if err != nil {
		return nil, code, err
	}

	return auditCtx.(*AuditContext), code, err
}

// CreateAuditContext - Create a new Audit Context
func (c *Client) CreateAuditContext(auditCtx AuditContext) (*AuditContext, int, error) {

	code, err := c.CreateResource("auditctx", *auditCtx.Name, auditCtx)
	if err != nil {
		return nil, code, err
	}

	return c.GetAuditContext(*auditCtx.Name)
}

// UpdateAuditContext - Update an existing Audit Context
func (c *Client) UpdateAuditContext(auditCtx AuditContext) (*AuditContext, int, error) {

	code, err := c.UpdateResource(fmt.Sprintf("auditctx/%s", *auditCtx.Name), auditCtx)
	if err != nil {
		return nil, code, err
	}

	return c.GetAuditContext(*auditCtx.Name)
}
