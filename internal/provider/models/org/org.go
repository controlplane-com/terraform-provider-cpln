package org

import (
	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

/*** Main Models ***/

// Auth Config //

type AuthConfigModel struct {
	DomainAutoMembers types.Set  `tfsdk:"domain_auto_members"`
	SamlOnly          types.Bool `tfsdk:"saml_only"`
}

func (a AuthConfigModel) AttributeTypes() attr.Type {
	return types.ObjectType{
		AttrTypes: map[string]attr.Type{
			"domain_auto_members": types.SetType{ElemType: types.StringType},
			"saml_only":           types.BoolType,
		},
	}
}

// Observability //

type ObservabilityModel struct {
	LogsRetentionDays    types.Int32 `tfsdk:"logs_retention_days"`
	MetricsRetentionDays types.Int32 `tfsdk:"metrics_retention_days"`
	TracesRetentionDays  types.Int32 `tfsdk:"traces_retention_days"`
	DefaultAlertEmails   types.Set   `tfsdk:"default_alert_emails"`
}

func (o ObservabilityModel) AttributeTypes() attr.Type {
	return types.ObjectType{
		AttrTypes: map[string]attr.Type{
			"logs_retention_days":    types.Int32Type,
			"metrics_retention_days": types.Int32Type,
			"traces_retention_days":  types.Int32Type,
			"default_alert_emails":   types.SetType{ElemType: types.StringType},
		},
	}
}

// Security //

type SecurityModel struct {
	ThreatDetection types.List `tfsdk:"threat_detection"`
}

func (s SecurityModel) AttributeTypes() attr.Type {
	return types.ObjectType{
		AttrTypes: map[string]attr.Type{
			"threat_detection": types.ListType{ElemType: SecurityThreatDetectionModel{}.AttributeTypes()},
		},
	}
}

// Security -> Threat Detection //

type SecurityThreatDetectionModel struct {
	Enabled         types.Bool   `tfsdk:"enabled"`
	MinimumSeverity types.String `tfsdk:"minimum_severity"`
	Syslog          types.List   `tfsdk:"syslog"`
}

func (s SecurityThreatDetectionModel) AttributeTypes() attr.Type {
	return types.ObjectType{
		AttrTypes: map[string]attr.Type{
			"enabled":          types.BoolType,
			"minimum_severity": types.StringType,
			"syslog":           types.ListType{ElemType: SecurityThreatDetectionSyslogModel{}.AttributeTypes()},
		},
	}
}

// Security -> Threat Detection -> Syslog //

type SecurityThreatDetectionSyslogModel struct {
	Transport types.String `tfsdk:"transport"`
	Host      types.String `tfsdk:"host"`
	Port      types.Int32  `tfsdk:"port"`
}

func (s SecurityThreatDetectionSyslogModel) AttributeTypes() attr.Type {
	return types.ObjectType{
		AttrTypes: map[string]attr.Type{
			"transport": types.StringType,
			"host":      types.StringType,
			"port":      types.Int32Type,
		},
	}
}

// Status //

type StatusModel struct {
	AccountLink    types.String `tfsdk:"account_link"`
	Active         types.Bool   `tfsdk:"active"`
	EndpointPrefix types.String `tfsdk:"endpoint_prefix"`
}

func (s StatusModel) AttributeTypes() attr.Type {
	return types.ObjectType{
		AttrTypes: map[string]attr.Type{
			"account_link":    types.StringType,
			"active":          types.BoolType,
			"endpoint_prefix": types.StringType,
		},
	}
}

/*** Logging Models ***/

// S3 Logging //

type S3LoggingModel struct {
	Bucket      types.String `tfsdk:"bucket"`
	Region      types.String `tfsdk:"region"`
	Prefix      types.String `tfsdk:"prefix"`
	Credentials types.String `tfsdk:"credentials"`
}

func (s S3LoggingModel) AttributeTypes() attr.Type {
	return types.ObjectType{
		AttrTypes: map[string]attr.Type{
			"bucket":      types.StringType,
			"region":      types.StringType,
			"prefix":      types.StringType,
			"credentials": types.StringType,
		},
	}
}

// Coralogix Logging //

type CoralogixLoggingModel struct {
	Cluster     types.String `tfsdk:"cluster"`
	Credentials types.String `tfsdk:"credentials"`
	App         types.String `tfsdk:"app"`
	Subsystem   types.String `tfsdk:"subsystem"`
}

func (c CoralogixLoggingModel) AttributeTypes() attr.Type {
	return types.ObjectType{
		AttrTypes: map[string]attr.Type{
			"cluster":     types.StringType,
			"credentials": types.StringType,
			"app":         types.StringType,
			"subsystem":   types.StringType,
		},
	}
}

// Datadog Logging //

type DatadogLoggingModel struct {
	Host        types.String `tfsdk:"host"`
	Credentials types.String `tfsdk:"credentials"`
}

func (d DatadogLoggingModel) AttributeTypes() attr.Type {
	return types.ObjectType{
		AttrTypes: map[string]attr.Type{
			"host":        types.StringType,
			"credentials": types.StringType,
		},
	}
}

// Logzio Logging //

type LogzioLoggingModel struct {
	ListenerHost types.String `tfsdk:"listener_host"`
	Credentials  types.String `tfsdk:"credentials"`
}

func (l LogzioLoggingModel) AttributeTypes() attr.Type {
	return types.ObjectType{
		AttrTypes: map[string]attr.Type{
			"listener_host": types.StringType,
			"credentials":   types.StringType,
		},
	}
}

// Elastic Logging //

type ElasticLoggingModel struct {
	AWS          types.List `tfsdk:"aws"`
	ElasticCloud types.List `tfsdk:"elastic_cloud"`
	Generic      types.List `tfsdk:"generic"`
}

func (e ElasticLoggingModel) AttributeTypes() attr.Type {
	return types.ObjectType{
		AttrTypes: map[string]attr.Type{
			"aws":           types.ListType{ElemType: ElasticLoggingAwsModel{}.AttributeTypes()},
			"elastic_cloud": types.ListType{ElemType: ElasticLoggingElasticCloudModel{}.AttributeTypes()},
			"generic":       types.ListType{ElemType: ElasticLoggingGenericModel{}.AttributeTypes()},
		},
	}
}

// Elastic Logging -> AWS //

type ElasticLoggingAwsModel struct {
	Host        types.String `tfsdk:"host"`
	Port        types.Int32  `tfsdk:"port"`
	Index       types.String `tfsdk:"index"`
	Type        types.String `tfsdk:"type"`
	Credentials types.String `tfsdk:"credentials"`
	Region      types.String `tfsdk:"region"`
}

func (e ElasticLoggingAwsModel) AttributeTypes() attr.Type {
	return types.ObjectType{
		AttrTypes: map[string]attr.Type{
			"host":        types.StringType,
			"port":        types.Int32Type,
			"index":       types.StringType,
			"type":        types.StringType,
			"credentials": types.StringType,
			"region":      types.StringType,
		},
	}
}

// Elastic Logging -> Elastic Cloud //

type ElasticLoggingElasticCloudModel struct {
	Index       types.String `tfsdk:"index"`
	Type        types.String `tfsdk:"type"`
	Credentials types.String `tfsdk:"credentials"`
	CloudID     types.String `tfsdk:"cloud_id"`
}

func (e ElasticLoggingElasticCloudModel) AttributeTypes() attr.Type {
	return types.ObjectType{
		AttrTypes: map[string]attr.Type{
			"index":       types.StringType,
			"type":        types.StringType,
			"credentials": types.StringType,
			"cloud_id":    types.StringType,
		},
	}
}

// Elastic Logging -> Generic //

type ElasticLoggingGenericModel struct {
	Host        types.String `tfsdk:"host"`
	Port        types.Int32  `tfsdk:"port"`
	Path        types.String `tfsdk:"path"`
	Index       types.String `tfsdk:"index"`
	Type        types.String `tfsdk:"type"`
	Credentials types.String `tfsdk:"credentials"`
}

func (e ElasticLoggingGenericModel) AttributeTypes() attr.Type {
	return types.ObjectType{
		AttrTypes: map[string]attr.Type{
			"host":        types.StringType,
			"port":        types.Int32Type,
			"path":        types.StringType,
			"index":       types.StringType,
			"type":        types.StringType,
			"credentials": types.StringType,
		},
	}
}

// Cloud Watch Logging //

type CloudWatchModel struct {
	Region        types.String `tfsdk:"region"`
	Credentials   types.String `tfsdk:"credentials"`
	RetentionDays types.Int32  `tfsdk:"retention_days"`
	GroupName     types.String `tfsdk:"group_name"`
	StreamName    types.String `tfsdk:"stream_name"`
	ExtractFields types.Map    `tfsdk:"extract_fields"`
}

func (c CloudWatchModel) AttributeTypes() attr.Type {
	return types.ObjectType{
		AttrTypes: map[string]attr.Type{
			"region":         types.StringType,
			"credentials":    types.StringType,
			"retention_days": types.Int32Type,
			"group_name":     types.StringType,
			"stream_name":    types.StringType,
			"extract_fields": types.MapType{ElemType: types.StringType},
		},
	}
}

// Fluentd Logging //

type FluentdLoggingModel struct {
	Host types.String `tfsdk:"host"`
	Port types.Int32  `tfsdk:"port"`
}

func (f FluentdLoggingModel) AttributeTypes() attr.Type {
	return types.ObjectType{
		AttrTypes: map[string]attr.Type{
			"host": types.StringType,
			"port": types.Int32Type,
		},
	}
}

// Stackdriver Logging //

type StackdriverLoggingModel struct {
	Credentials types.String `tfsdk:"credentials"`
	Location    types.String `tfsdk:"location"`
}

func (s StackdriverLoggingModel) AttributeTypes() attr.Type {
	return types.ObjectType{
		AttrTypes: map[string]attr.Type{
			"credentials": types.StringType,
			"location":    types.StringType,
		},
	}
}

// Syslog Logging //

type SyslogLoggingModel struct {
	Host     types.String `tfsdk:"host"`
	Port     types.Int32  `tfsdk:"port"`
	Mode     types.String `tfsdk:"mode"`
	Format   types.String `tfsdk:"format"`
	Severity types.Int32  `tfsdk:"severity"`
}

func (s SyslogLoggingModel) AttributeTypes() attr.Type {
	return types.ObjectType{
		AttrTypes: map[string]attr.Type{
			"host":     types.StringType,
			"port":     types.Int32Type,
			"mode":     types.StringType,
			"format":   types.StringType,
			"severity": types.Int32Type,
		},
	}
}
