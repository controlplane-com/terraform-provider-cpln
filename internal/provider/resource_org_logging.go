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
	"s3_logging", "coralogix_logging", "datadog_logging", "logzio_logging", "elastic_logging", "cloud_watch_logging", "fluentd_logging", "stackdriver_logging", "syslog_logging",
}

func resourceOrgLogging() *schema.Resource {

	return &schema.Resource{
		CreateContext: resourceOrgLoggingCreate,
		ReadContext:   resourceOrgLoggingRead,
		UpdateContext: resourceOrgLoggingUpdate,
		DeleteContext: resourceOrgLoggingDelete,
		Schema: map[string]*schema.Schema{
			"cpln_id": {
				Type:        schema.TypeString,
				Description: "The ID, in GUID format, of the org.",
				Computed:    true,
			},
			"name": {
				Type:        schema.TypeString,
				Description: "The name of the org.",
				Computed:    true,
			},
			"description": {
				Type:        schema.TypeString,
				Description: "The description of org.",
				Computed:    true,
			},
			"tags": {
				Type:        schema.TypeMap,
				Description: "Key-value map of the org's tags.",
				Computed:    true,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
			},
			"s3_logging": {
				Type:        schema.TypeList,
				Description: "[Documentation Reference](https://docs.controlplane.com/external-logging/s3)",
				Optional:    true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"bucket": {
							Type:        schema.TypeString,
							Description: "Name of S3 bucket.",
							Required:    true,
						},
						"region": {
							Type:        schema.TypeString,
							Description: "AWS region where bucket is located.",
							Required:    true,
						},
						"prefix": {
							Type:        schema.TypeString,
							Description: "Bucket path prefix. Default: \"/\".",
							Optional:    true,
							Default:     "/",
						},
						"credentials": {
							Type:        schema.TypeString,
							Description: "Full link to referenced AWS Secret.",
							Required:    true,
						},
					},
				},
			},
			"coralogix_logging": {
				Type:        schema.TypeList,
				Description: "[Documentation Reference](https://docs.controlplane.com/external-logging/coralogix)",
				Optional:    true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"cluster": {
							Type:        schema.TypeString,
							Description: "Coralogix cluster URI.",
							Required:    true,
						},
						"credentials": {
							Type:        schema.TypeString,
							Description: "Full link to referenced Opaque Secret.",
							Required:    true,
						},
						"app": {
							Type:        schema.TypeString,
							Description: "App name to be displayed in Coralogix dashboard.",
							Required:    true,
						},
						"subsystem": {
							Type:        schema.TypeString,
							Description: "Subsystem name to be displayed in Coralogix dashboard.",
							Required:    true,
						},
					},
				},
			},
			"datadog_logging": {
				Type:        schema.TypeList,
				Description: "[Documentation Reference](https://docs.controlplane.com/external-logging/datadog)",
				Optional:    true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"host": {
							Type:        schema.TypeString,
							Description: "Datadog host URI.",
							Required:    true,
						},
						"credentials": {
							Type:        schema.TypeString,
							Description: "Full link to referenced Opaque Secret.",
							Required:    true,
						},
					},
				},
			},
			"logzio_logging": {
				Type:        schema.TypeList,
				Description: "[Documentation Reference](https://docs.controlplane.com/external-logging/logz-io)",
				Optional:    true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"listener_host": {
							Type:        schema.TypeString,
							Description: "Logzio listener host URI.",
							Required:    true,
						},
						"credentials": {
							Type:        schema.TypeString,
							Description: "Full link to referenced Opaque Secret.",
							Required:    true,
						},
					},
				},
			},
			"elastic_logging": {
				Type:        schema.TypeList,
				Description: "For logging and analyzing data within an org using Elastic Logging.",
				Optional:    true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"aws": {
							Type:        schema.TypeList,
							Description: "For targeting Amazon Web Services (AWS) ElasticSearch.",
							Optional:    true,
							MaxItems:    1,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"host": {
										Type:        schema.TypeString,
										Description: "A valid AWS ElasticSearch hostname (must end with es.amazonaws.com).",
										Required:    true,
									},
									"port": {
										Type:        schema.TypeInt,
										Description: "Port. Default: 443",
										Required:    true,
									},
									"index": {
										Type:        schema.TypeString,
										Description: "Logging Index.",
										Required:    true,
									},
									"type": {
										Type:        schema.TypeString,
										Description: "Logging Type.",
										Required:    true,
									},
									"credentials": {
										Type:        schema.TypeString,
										Description: "Full Link to a secret of type `aws`.",
										Required:    true,
									},
									"region": {
										Type:        schema.TypeString,
										Description: "Valid AWS region.",
										Required:    true,
									},
								},
							},
						},
						"elastic_cloud": {
							Type:        schema.TypeList,
							Description: "For targeting Elastic Cloud.",
							Optional:    true,
							MaxItems:    1,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"index": {
										Type:        schema.TypeString,
										Description: "Logging Index.",
										Required:    true,
									},
									"type": {
										Type:        schema.TypeString,
										Description: "Logging Type.",
										Required:    true,
									},
									"credentials": {
										Type:        schema.TypeString,
										Description: "Full Link to a secret of type `userpass`.",
										Required:    true,
									},
									"cloud_id": {
										Type:        schema.TypeString,
										Description: "[Cloud ID](https://www.elastic.co/guide/en/cloud/current/ec-cloud-id.html)",
										Required:    true,
									},
								},
							},
						},
						"generic": {
							Type:        schema.TypeList,
							Description: "For targeting generic Elastic Search providers.",
							Optional:    true,
							MaxItems:    1,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"host": {
										Type:        schema.TypeString,
										Description: "A valid Elastic Search provider hostname.",
										Required:    true,
									},
									"port": {
										Type:        schema.TypeInt,
										Description: "Port. Default: 443",
										Required:    true,
									},
									"path": {
										Type:        schema.TypeString,
										Description: "Logging path.",
										Required:    true,
									},
									"index": {
										Type:        schema.TypeString,
										Description: "Logging Index.",
										Required:    true,
									},
									"type": {
										Type:        schema.TypeString,
										Description: "Logging Type.",
										Required:    true,
									},
									"credentials": {
										Type:        schema.TypeString,
										Description: "Full Link to a secret of type `userpass`.",
										Required:    true,
									},
								},
							},
						},
					},
				},
			},
			"cloud_watch_logging": {
				Type:     schema.TypeList,
				Optional: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"region": {
							Type:        schema.TypeString,
							Description: "Valid AWS region.",
							Required:    true,
						},
						"credentials": {
							Type:        schema.TypeString,
							Description: "Full Link to a secret of type `opaque`.",
							Required:    true,
						},
						"retention_days": {
							Type:        schema.TypeInt,
							Description: "Length, in days, for how log data is kept before it is automatically deleted.",
							Optional:    true,
						},
						"group_name": {
							Type:        schema.TypeString,
							Description: "A container for log streams with common settings like retention. Used to categorize logs by application or service type.",
							Required:    true,
						},
						"stream_name": {
							Type:        schema.TypeString,
							Description: "A sequence of log events from the same source within a log group. Typically represents individual instances of services or applications.",
							Required:    true,
						},
					},
				},
			},
			"fluentd_logging": {
				Type:     schema.TypeList,
				Optional: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"host": {
							Type:        schema.TypeString,
							Description: "The hostname or IP address of a remote log storage system.",
							Required:    true,
						},
						"port": {
							Type:        schema.TypeInt,
							Description: "Port. Default: 24224",
							Optional:    true,
							Default:     24224,
						},
					},
				},
			},
			"stackdriver_logging": {
				Type:     schema.TypeList,
				Optional: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"credentials": {
							Type:        schema.TypeString,
							Description: "Full Link to a secret of type `opaque`.",
							Required:    true,
						},
						"location": {
							Type:        schema.TypeString,
							Description: "A Google Cloud Provider region.",
							Required:    true,
						},
					},
				},
			},
			"syslog_logging": {
				Type:     schema.TypeList,
				Optional: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"host": {
							Type:        schema.TypeString,
							Description: "Hostname of Syslog Endpoint.",
							Required:    true,
						},
						"port": {
							Type:        schema.TypeInt,
							Description: "Port of Syslog Endpoint.",
							Required:    true,
						},
						"mode": {
							Type:        schema.TypeString,
							Description: "Log Mode. Valid values: TCP, TLS, or UDP.",
							Optional:    true,
							Default:     "tcp",
						},
						"format": {
							Type:        schema.TypeString,
							Description: "Log Format. Valid values: RFC3164 or RFC5424.",
							Optional:    true,
							Default:     "rfc5424",
						},
						"severity": {
							Type:        schema.TypeInt,
							Description: "Severity Level. See documentation for details. Valid values: 0 to 7.",
							Optional:    true,
							Default:     6,
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

	if e := orgLoggingValidate(loggings); e != nil {
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

	if d.HasChanges("s3_logging", "coralogix_logging", "datadog_logging", "logzio_logging", "elastic_logging", "cloud_watch_logging", "fluentd_logging", "stackdriver_logging", "syslog_logging") {

		c := m.(*client.Client)

		// Build regardless of changes
		loggings := buildMultipleLoggings(d, loggingNames...)

		if e := orgLoggingValidate(loggings); e != nil {
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
		var cloudWatchArray []client.CloudWatchLogging
		var fluentdArray []client.FluentdLogging
		var stackdriverArray []client.StackdriverLogging
		var syslogArray []client.SyslogLogging

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

			if logging.CloudWatch != nil {
				cloudWatchArray = append(cloudWatchArray, *logging.CloudWatch)
			}

			if logging.Fluentd != nil {
				fluentdArray = append(fluentdArray, *logging.Fluentd)
			}

			if logging.Stackdriver != nil {
				stackdriverArray = append(stackdriverArray, *logging.Stackdriver)
			}

			if logging.Syslog != nil {
				syslogArray = append(syslogArray, *logging.Syslog)
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

		if err := d.Set("cloud_watch_logging", flattenCloudWatchLogging(cloudWatchArray)); err != nil {
			return diag.FromErr(err)
		}

		if err := d.Set("fluentd_logging", flattenFluentdLogging(fluentdArray)); err != nil {
			return diag.FromErr(err)
		}

		if err := d.Set("stackdriver_logging", flattenStackdriverLogging(stackdriverArray)); err != nil {
			return diag.FromErr(err)
		}

		if err := d.Set("syslog_logging", flattenSyslogLogging(syslogArray)); err != nil {
			return diag.FromErr(err)
		}
	}

	return nil
}

/*** Build ***/

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

func buildCloudWatchLogging(logging []interface{}) []client.Logging {

	if len(logging) == 0 {
		return nil
	}

	var output []client.Logging

	for _, logs := range logging {

		if logs == nil {
			continue
		}

		log := logs.(map[string]interface{})

		tempLogging := client.Logging{
			CloudWatch: &client.CloudWatchLogging{
				Region:        GetString(log["region"].(string)),
				Credentials:   GetString(log["credentials"].(string)),
				RetentionDays: GetInt(log["retention_days"]),
				GroupName:     GetString(log["group_name"].(string)),
				StreamName:    GetString(log["stream_name"].(string)),
			},
		}

		output = append(output, tempLogging)
	}

	return output
}

func buildFluentdLogging(logging []interface{}) []client.Logging {

	if len(logging) == 0 {
		return nil
	}

	var output []client.Logging

	for _, logs := range logging {

		if logs == nil {
			continue
		}

		log := logs.(map[string]interface{})

		tempLogging := client.Logging{
			Fluentd: &client.FluentdLogging{
				Host: GetString(log["host"].(string)),
				Port: GetInt(log["port"].(int)),
			},
		}

		output = append(output, tempLogging)
	}

	return output
}

func buildStackdriverLogging(logging []interface{}) []client.Logging {

	if len(logging) == 0 {
		return nil
	}

	var output []client.Logging

	for _, logs := range logging {

		if logs == nil {
			continue
		}

		log := logs.(map[string]interface{})

		tempLogging := client.Logging{
			Stackdriver: &client.StackdriverLogging{
				Credentials: GetString(log["credentials"].(string)),
				Location:    GetString(log["location"].(string)),
			},
		}

		output = append(output, tempLogging)
	}

	return output
}

func buildSyslogLogging(logging []interface{}) []client.Logging {

	if len(logging) == 0 {
		return nil
	}

	var output []client.Logging

	for _, logs := range logging {

		if logs == nil {
			continue
		}

		log := logs.(map[string]interface{})

		tempLogging := client.Logging{
			Syslog: &client.SyslogLogging{
				Host:     GetString(log["host"].(string)),
				Port:     GetInt(log["port"].(int)),
				Mode:     GetString(log["mode"].(string)),
				Format:   GetString(log["format"].(string)),
				Severity: GetInt(log["severity"].(int)),
			},
		}

		output = append(output, tempLogging)
	}

	return output
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

/*** Flatten ***/

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

func flattenCloudWatchLogging(logs []client.CloudWatchLogging) []interface{} {

	if len(logs) == 0 {
		return nil
	}

	output := make([]interface{}, len(logs))

	for i, log := range logs {

		outputMap := make(map[string]interface{})

		outputMap["region"] = *log.Region
		outputMap["credentials"] = *log.Credentials

		if log.RetentionDays != nil {
			outputMap["retention_days"] = *log.RetentionDays
		}

		outputMap["group_name"] = *log.GroupName
		outputMap["stream_name"] = *log.StreamName

		output[i] = outputMap
	}

	return output
}

func flattenFluentdLogging(logs []client.FluentdLogging) []interface{} {

	if len(logs) == 0 {
		return nil
	}

	output := make([]interface{}, len(logs))

	for i, log := range logs {

		outputMap := make(map[string]interface{})

		outputMap["host"] = *log.Host
		outputMap["port"] = *log.Port

		output[i] = outputMap
	}

	return output
}

func flattenStackdriverLogging(logs []client.StackdriverLogging) []interface{} {

	if len(logs) == 0 {
		return nil
	}

	output := make([]interface{}, len(logs))

	for i, log := range logs {

		outputMap := make(map[string]interface{})

		outputMap["credentials"] = *log.Credentials
		outputMap["location"] = *log.Location

		output[i] = outputMap
	}

	return output
}

func flattenSyslogLogging(logs []client.SyslogLogging) []interface{} {

	if len(logs) == 0 {
		return nil
	}

	output := make([]interface{}, len(logs))

	for i, log := range logs {

		outputMap := make(map[string]interface{})

		outputMap["host"] = *log.Host
		outputMap["port"] = *log.Port
		outputMap["mode"] = *log.Mode
		outputMap["format"] = *log.Format
		outputMap["severity"] = *log.Severity

		output[i] = outputMap
	}

	return output
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

func orgLoggingValidate(loggings []client.Logging) diag.Diagnostics {

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
		case "cloud_watch_logging":
			loggingToAdd = buildCloudWatchLogging(logArray)
		case "fluentd_logging":
			loggingToAdd = buildFluentdLogging(logArray)
		case "stackdriver_logging":
			loggingToAdd = buildStackdriverLogging(logArray)
		case "syslog_logging":
			loggingToAdd = buildSyslogLogging(logArray)
		default:
			continue
		}

		if loggingToAdd != nil {
			loggings = append(loggings, loggingToAdd...)
		}
	}

	return loggings
}
