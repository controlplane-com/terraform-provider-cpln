package cpln

import (
	"context"

	client "terraform-provider-cpln/internal/provider/client"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

	"encoding/json"
)

func resourceAgent() *schema.Resource {

	return &schema.Resource{
		CreateContext: resourceAgentCreate,
		ReadContext:   resourceAgentRead,
		UpdateContext: resourceAgentUpdate,
		DeleteContext: resourceAgentDelete,
		Schema: map[string]*schema.Schema{
			"cpln_id": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"name": {
				Type:         schema.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: NameValidator,
			},
			"description": {
				Type:             schema.TypeString,
				Optional:         true,
				ValidateFunc:     DescriptionValidator,
				DiffSuppressFunc: DiffSuppressDescription,
			},
			"tags": {
				Type:     schema.TypeMap,
				Optional: true,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
				ValidateFunc: TagValidator,
			},
			"self_link": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"user_data": {
				Type:      schema.TypeString,
				Computed:  true,
				Sensitive: true,
			},
		},
		Importer: &schema.ResourceImporter{},
	}
}

func resourceAgentCreate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {

	// log.Printf("[INFO] Method: resourceAgentCreate")

	agent := client.Agent{}
	agent.Name = GetString(d.Get("name"))
	agent.Description = GetString(d.Get("description"))
	agent.Tags = GetStringMap(d.Get("tags"))

	c := m.(*client.Client)
	newAgent, code, err := c.CreateAgent(agent)

	if code == 409 {
		return ResourceExistsHelper()
	}

	if err != nil {
		return diag.FromErr(err)
	}

	a, err := json.Marshal(newAgent.Status.BootstrapConfig)
	if err != nil {
		return diag.FromErr(err)
	}

	// User data is only available on create
	if err := d.Set("user_data", string(a)); err != nil {
		return diag.FromErr(err)
	}

	return setAgent(d, newAgent)
}

func resourceAgentRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {

	// log.Printf("[INFO] Method: resourceAgentRead")

	c := m.(*client.Client)
	agent, code, err := c.GetAgent(d.Id())

	if code == 404 {
		d.SetId("")
		return nil
	}

	if err != nil {
		return diag.FromErr(err)
	}

	return setAgent(d, agent)
}

func setAgent(d *schema.ResourceData, agent *client.Agent) diag.Diagnostics {

	if agent == nil {
		d.SetId("")
		return nil
	}

	d.SetId(*agent.Name)

	if err := SetBase(d, agent.Base); err != nil {
		return diag.FromErr(err)
	}

	if err := SetSelfLink(agent.Links, d); err != nil {
		return diag.FromErr(err)
	}

	return nil
}

func resourceAgentUpdate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {

	// log.Printf("[INFO] Method: resourceAgentUpdate")

	if d.HasChanges("description", "tags") {

		agentToUpdate := client.Agent{}
		agentToUpdate.Name = GetString(d.Get("name"))

		if d.HasChange("description") {
			agentToUpdate.Description = GetDescriptionString(d.Get("description"), *agentToUpdate.Name)
		}

		if d.HasChange("tags") {
			agentToUpdate.Tags = GetTagChanges(d)
		}

		c := m.(*client.Client)
		updatedAgent, _, err := c.UpdateAgent(agentToUpdate)
		if err != nil {
			return diag.FromErr(err)
		}

		return setAgent(d, updatedAgent)
	}

	return nil
}

func resourceAgentDelete(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {

	// log.Printf("[INFO] Method: resourceAgentDelete")

	c := m.(*client.Client)
	err := c.DeleteAgent(d.Id())
	if err != nil {
		return diag.FromErr(err)
	}

	d.SetId("")

	return nil
}
