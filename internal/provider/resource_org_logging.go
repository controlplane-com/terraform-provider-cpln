package cpln

import (
	"context"
	"sync"

	client "github.com/controlplane-com/terraform-provider-cpln/internal/provider/client"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

var resourceLock = &sync.Mutex{}

var loggingNames = []string{
	"s3_logging", "coralogix_logging", "datadog_logging", "logzio_logging", "elastic_logging",
}

func resourceOrgLogging() *schema.Resource {

	return &schema.Resource{
		CreateContext: resourceOrgLoggingCreate,
		ReadContext:   resourceOrgLoggingRead,
		UpdateContext: resourceOrgLoggingUpdate,
		DeleteContext: resourceOrgLoggingDelete,
		Schema: map[string]*schema.Schema{
			"cpln_id": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"name": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"description": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"tags": {
				Type:     schema.TypeMap,
				Computed: true,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
			},
			"s3_logging": {
				Type:     schema.TypeList,
				Optional: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"bucket": {
							Type:     schema.TypeString,
							Required: true,
						},
						"region": {
							Type:     schema.TypeString,
							Required: true,
						},
						"prefix": {
							Type:     schema.TypeString,
							Optional: true,
							Default:  "/",
						},
						"credentials": {
							Type:     schema.TypeString,
							Required: true,
						},
					},
				},
			},
			"coralogix_logging": {
				Type:     schema.TypeList,
				Optional: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"cluster": {
							Type:     schema.TypeString,
							Required: true,
						},
						"credentials": {
							Type:     schema.TypeString,
							Required: true,
						},
						"app": {
							Type:     schema.TypeString,
							Required: true,
						},
						"subsystem": {
							Type:     schema.TypeString,
							Required: true,
						},
					},
				},
			},
			"datadog_logging": {
				Type:     schema.TypeList,
				Optional: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"host": {
							Type:     schema.TypeString,
							Required: true,
						},
						"credentials": {
							Type:     schema.TypeString,
							Required: true,
						},
					},
				},
			},
			"logzio_logging": {
				Type:     schema.TypeList,
				Optional: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"listener_host": {
							Type:     schema.TypeString,
							Required: true,
						},
						"credentials": {
							Type:     schema.TypeString,
							Required: true,
						},
					},
				},
			},
			"elastic_logging": {
				Type:     schema.TypeList,
				Optional: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"aws": {
							Type:     schema.TypeList,
							Optional: true,
							MaxItems: 1,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"host": {
										Type:     schema.TypeString,
										Required: true,
									},
									"port": {
										Type:     schema.TypeInt,
										Required: true,
									},
									"index": {
										Type:     schema.TypeString,
										Required: true,
									},
									"type": {
										Type:     schema.TypeString,
										Required: true,
									},
									"credentials": {
										Type:     schema.TypeString,
										Required: true,
									},
									"region": {
										Type:     schema.TypeString,
										Required: true,
									},
								},
							},
						},
						"elastic_cloud": {
							Type:     schema.TypeList,
							Optional: true,
							MaxItems: 1,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"index": {
										Type:     schema.TypeString,
										Required: true,
									},
									"type": {
										Type:     schema.TypeString,
										Required: true,
									},
									"credentials": {
										Type:     schema.TypeString,
										Required: true,
									},
									"cloud_id": {
										Type:     schema.TypeString,
										Required: true,
									},
								},
							},
						},
						"generic": {
							Type:     schema.TypeList,
							Optional: true,
							MaxItems: 1,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"host": {
										Type:     schema.TypeString,
										Required: true,
									},
									"port": {
										Type:     schema.TypeInt,
										Required: true,
									},
									"path": {
										Type:     schema.TypeString,
										Required: true,
									},
									"index": {
										Type:     schema.TypeString,
										Required: true,
									},
									"type": {
										Type:     schema.TypeString,
										Required: true,
									},
									"credentials": {
										Type:     schema.TypeString,
										Required: true,
									},
								},
							},
						},
					},
				},
			},
		},
		Importer: &schema.ResourceImporter{},
	}
}

func resourceOrgLoggingCreate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {

	resourceLock.Lock()
	defer resourceLock.Unlock()

	// log.Printf("[INFO] Method: resourceOrgCreate")

	c := m.(*client.Client)

	currentOrg, _, err := c.GetOrg()

	if err != nil {
		return diag.FromErr(err)
	}

	if currentOrg.Spec != nil && currentOrg.Spec.Logging != nil {
		return diag.Errorf("only one 'cpln_org_logging' resource can be declared")
	}

	loggings := buildMultipleLoggings(d, loggingNames...)

	if e := orgLoggingValidate(d, loggings); e != nil {
		return e
	}

	org, _, err := c.UpdateOrgLogging(&loggings)
	if err != nil {
		return diag.FromErr(err)
	}

	return setOrgLogging(d, org)
}

func resourceOrgLoggingRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {

	// log.Printf("[INFO] Method: resourceOrgtRead")

	c := m.(*client.Client)
	org, _, err := c.GetOrg()

	if err != nil {
		return diag.FromErr(err)
	}

	return setOrgLogging(d, org)
}

func resourceOrgLoggingUpdate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {

	// log.Printf("[INFO] Method: resourceOrgUpdate")

	if d.HasChanges("s3_logging", "coralogix_logging", "datadog_logging", "logzio_logging", "elastic_logging") {

		c := m.(*client.Client)

		// Build regardless of changes
		loggings := buildMultipleLoggings(d, loggingNames...)

		if e := orgLoggingValidate(d, loggings); e != nil {
			return e
		}

		org, _, err := c.UpdateOrgLogging(&loggings)

		if err != nil {
			return diag.FromErr(err)
		}

		return setOrgLogging(d, org)
	}

	return nil
}

func resourceOrgLoggingDelete(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {

	// log.Printf("[INFO] Method: resourceOrgDelete")

	c := m.(*client.Client)

	_, _, err := c.UpdateOrgLogging(nil)

	if err != nil {
		return diag.FromErr(err)
	}

	d.SetId("")

	return nil
}

/*** Build Functions ***/
func buildS3Logging(logging []interface{}) []client.Logging {

	if len(logging) > 0 {

		var output []client.Logging

		for _, logs := range logging {

			if logs != nil {

				log := logs.(map[string]interface{})

				iLog := &client.S3Logging{}
				iLog.Bucket = GetString(log["bucket"])
				iLog.Region = GetString(log["region"])
				iLog.Prefix = GetString(log["prefix"])
				iLog.Credentials = GetString(log["credentials"])

				tempLogging := client.Logging{
					S3: iLog,
				}

				output = append(output, tempLogging)
			}
		}

		return output
	}

	return nil
}

func buildCoralogixLogging(logging []interface{}) []client.Logging {

	if len(logging) > 0 {

		var output []client.Logging

		for _, logs := range logging {

			if logs != nil {
				log := logs.(map[string]interface{})

				iLog := &client.CoralogixLogging{}
				iLog.Cluster = GetString(log["cluster"])
				iLog.Credentials = GetString(log["credentials"])
				iLog.App = GetString(log["app"])
				iLog.Subsystem = GetString(log["subsystem"])

				tempLogging := client.Logging{
					Coralogix: iLog,
				}

				output = append(output, tempLogging)
			}
		}

		return output
	}

	return nil
}

func buildDatadogLogging(logging []interface{}) []client.Logging {

	if len(logging) > 0 {

		var output []client.Logging

		for _, logs := range logging {

			if logs != nil {

				log := logs.(map[string]interface{})

				iLog := &client.DatadogLogging{}
				iLog.Host = GetString(log["host"])
				iLog.Credentials = GetString(log["credentials"])

				tempLogging := client.Logging{
					Datadog: iLog,
				}

				output = append(output, tempLogging)
			}
		}

		return output
	}

	return nil
}

func buildLogzioLogging(logging []interface{}) []client.Logging {

	if len(logging) > 0 {

		var output []client.Logging

		for _, logs := range logging {

			if logs != nil {
				log := logs.(map[string]interface{})

				iLog := &client.LogzioLogging{}
				iLog.ListenerHost = GetString(log["listener_host"])
				iLog.Credentials = GetString(log["credentials"])

				tempLogging := client.Logging{
					Logzio: iLog,
				}

				output = append(output, tempLogging)
			}
		}

		return output
	}

	return nil
}

func buildElasticLogging(logging []interface{}) []client.Logging {

	if len(logging) > 0 {

		var output []client.Logging

		for _, logs := range logging {

			if logs != nil {

				log := logs.(map[string]interface{})
				result := &client.ElasticLogging{}

				if log["aws"] != nil {
					result.AWS = buildAWSLogging(log["aws"].([]interface{}))
				}

				if log["elastic_cloud"] != nil {
					result.ElasticCloud = buildElasticCloudLogging(log["elastic_cloud"].([]interface{}))
				}

				if log["generic"] != nil {
					result.Generic = buildGenericLogging(log["generic"].([]interface{}))
				}

				tempLogging := client.Logging{
					Elastic: result,
				}

				output = append(output, tempLogging)
			}
		}

		return output
	}

	return nil
}

func buildAWSLogging(logging []interface{}) *client.AWSLogging {

	if len(logging) == 0 || logging[0] == nil {
		return nil
	}

	log := logging[0].(map[string]interface{})

	// Note: No need to check for nil because all of these fields are required.
	result := &client.AWSLogging{
		Host:        GetString(log["host"].(string)),
		Port:        GetInt(log["port"].(int)),
		Index:       GetString(log["index"].(string)),
		Type:        GetString(log["type"].(string)),
		Credentials: GetString(log["credentials"].(string)),
		Region:      GetString(log["region"].(string)),
	}

	return result
}

func buildElasticCloudLogging(logging []interface{}) *client.ElasticCloudLogging {

	if len(logging) == 0 || logging[0] == nil {
		return nil
	}

	log := logging[0].(map[string]interface{})

	// Note: No need to check for nil because all of these fields are required.
	result := &client.ElasticCloudLogging{
		Index:       GetString(log["index"].(string)),
		Type:        GetString(log["type"].(string)),
		Credentials: GetString(log["credentials"].(string)),
		CloudID:     GetString(log["cloud_id"].(string)),
	}

	return result
}

func buildGenericLogging(specs []interface{}) *client.GenericLogging {
	if len(specs) == 0 || specs[0] == nil {
		return nil
	}

	spec := specs[0].(map[string]interface{})

	return &client.GenericLogging{
		Host:        GetString(spec["host"].(string)),
		Port:        GetInt(spec["port"].(int)),
		Path:        GetString(spec["path"].(string)),
		Index:       GetString(spec["index"].(string)),
		Type:        GetString(spec["type"].(string)),
		Credentials: GetString(spec["credentials"].(string)),
	}
}

/*** Flatten Functions ***/
func flattenS3Logging(logs []client.S3Logging) []interface{} {

	if len(logs) > 0 {

		output := make([]interface{}, len(logs))

		for l, log := range logs {

			outputMap := make(map[string]interface{})

			outputMap["bucket"] = *log.Bucket
			outputMap["region"] = *log.Region
			outputMap["prefix"] = *log.Prefix
			outputMap["credentials"] = *log.Credentials

			output[l] = outputMap
		}

		return output
	}

	return nil
}

func flattenCoralogixLogging(logs []client.CoralogixLogging) []interface{} {

	if len(logs) > 0 {

		output := make([]interface{}, len(logs))

		for l, log := range logs {

			outputMap := make(map[string]interface{})

			outputMap["cluster"] = *log.Cluster
			outputMap["credentials"] = *log.Credentials
			outputMap["app"] = *log.App
			outputMap["subsystem"] = *log.Subsystem

			output[l] = outputMap
		}

		return output
	}

	return nil
}

func flattenDatadogLogging(logs []client.DatadogLogging) []interface{} {

	if len(logs) > 0 {

		output := make([]interface{}, len(logs))

		for l, log := range logs {

			outputMap := make(map[string]interface{})

			outputMap["host"] = *log.Host
			outputMap["credentials"] = *log.Credentials

			output[l] = outputMap
		}

		return output
	}

	return nil
}

func flattenLogzioLogging(logs []client.LogzioLogging) []interface{} {

	if len(logs) > 0 {

		output := make([]interface{}, len(logs))

		for l, log := range logs {

			outputMap := make(map[string]interface{})

			outputMap["listener_host"] = *log.ListenerHost
			outputMap["credentials"] = *log.Credentials

			output[l] = outputMap
		}

		return output
	}

	return nil
}

func flattenElasticLogging(logs []client.ElasticLogging) []interface{} {

	if len(logs) > 0 {

		output := make([]interface{}, len(logs))

		for l, log := range logs {

			result := make(map[string]interface{})

			if log.AWS != nil {
				result["aws"] = flattenAWSLogging(log.AWS)
			}

			if log.ElasticCloud != nil {
				result["elastic_cloud"] = flattenElasticCloudLogging(log.ElasticCloud)
			}

			if log.Generic != nil {
				result["generic"] = flattenGenericLogging(log.Generic)
			}

			output[l] = result
		}

		return output
	}

	return nil
}

func flattenAWSLogging(log *client.AWSLogging) []interface{} {

	if log == nil {
		return nil
	}

	result := make(map[string]interface{})

	result["host"] = *log.Host
	result["port"] = *log.Port
	result["index"] = *log.Index
	result["type"] = *log.Type
	result["credentials"] = *log.Credentials
	result["region"] = *log.Region

	return []interface{}{
		result,
	}
}

func flattenElasticCloudLogging(log *client.ElasticCloudLogging) []interface{} {

	if log == nil {
		return nil
	}

	result := make(map[string]interface{})

	result["index"] = *log.Index
	result["type"] = *log.Type
	result["credentials"] = *log.Credentials
	result["cloud_id"] = *log.CloudID

	return []interface{}{
		result,
	}
}

func flattenGenericLogging(logging *client.GenericLogging) []interface{} {
	if logging == nil {
		return nil
	}

	output := make(map[string]interface{})

	output["host"] = *logging.Host
	output["port"] = *logging.Port
	output["path"] = *logging.Path
	output["index"] = *logging.Index
	output["type"] = *logging.Type
	output["credentials"] = *logging.Credentials

	return []interface{}{
		output,
	}
}

/*** Helper Functions ***/
func setOrgLogging(d *schema.ResourceData, org *client.Org) diag.Diagnostics {

	if org == nil {
		d.SetId("")
		return nil
	}

	d.SetId(*org.Name)

	if err := SetBase(d, org.Base); err != nil {
		return diag.FromErr(err)
	}

	loggings := []client.Logging{}

	if org.Spec != nil {

		if org.Spec.Logging != nil {
			loggings = append(loggings, *org.Spec.Logging)
		}

		if org.Spec.ExtraLogging != nil && len(*org.Spec.ExtraLogging) > 0 {
			loggings = append(loggings, *org.Spec.ExtraLogging...)
		}

		var s3Array []client.S3Logging
		var coralogixArray []client.CoralogixLogging
		var dataDogArray []client.DatadogLogging
		var logzioArray []client.LogzioLogging
		var elasticArray []client.ElasticLogging

		for _, logging := range loggings {

			if logging.S3 != nil {
				s3Array = append(s3Array, *logging.S3)
			}

			if logging.Coralogix != nil {
				coralogixArray = append(coralogixArray, *logging.Coralogix)
			}

			if logging.Datadog != nil {
				dataDogArray = append(dataDogArray, *logging.Datadog)
			}

			if logging.Logzio != nil {
				logzioArray = append(logzioArray, *logging.Logzio)
			}

			if logging.Elastic != nil {
				elasticArray = append(elasticArray, *logging.Elastic)
			}
		}

		if err := d.Set("s3_logging", flattenS3Logging(s3Array)); err != nil {
			return diag.FromErr(err)
		}

		if err := d.Set("coralogix_logging", flattenCoralogixLogging(coralogixArray)); err != nil {
			return diag.FromErr(err)
		}

		if err := d.Set("datadog_logging", flattenDatadogLogging(dataDogArray)); err != nil {
			return diag.FromErr(err)
		}

		if err := d.Set("logzio_logging", flattenLogzioLogging(logzioArray)); err != nil {
			return diag.FromErr(err)
		}

		if err := d.Set("elastic_logging", flattenElasticLogging(elasticArray)); err != nil {
			return diag.FromErr(err)
		}
	}

	return nil
}

func orgLoggingValidate(d *schema.ResourceData, loggings []client.Logging) diag.Diagnostics {

	// Max of 4 external logging providers
	if len(loggings) > 4 {
		return diag.Errorf("max of 4 external logging providers allowed")
	}

	if len(loggings) == 0 {
		return diag.Errorf("at least one external logging providers must be defined")
	}

	return nil
}

func buildMultipleLoggings(d *schema.ResourceData, loggingTypes ...string) []client.Logging {

	loggings := []client.Logging{}

	for _, loggingType := range loggingTypes {

		logArray := d.Get(loggingType).([]interface{})

		if logArray == nil {
			continue
		}

		var loggingToAdd []client.Logging

		switch loggingType {
		case "s3_logging":
			loggingToAdd = buildS3Logging(logArray)
		case "coralogix_logging":
			loggingToAdd = buildCoralogixLogging(logArray)
		case "datadog_logging":
			loggingToAdd = buildDatadogLogging(logArray)
		case "logzio_logging":
			loggingToAdd = buildLogzioLogging(logArray)
		case "elastic_logging":
			loggingToAdd = buildElasticLogging(logArray)
		default:
			continue
		}

		if loggingToAdd != nil {
			loggings = append(loggings, loggingToAdd...)
		}
	}

	return loggings
}
