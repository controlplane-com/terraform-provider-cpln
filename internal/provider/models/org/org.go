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

// Observability //

type ObservabilityModel struct {
	LogsRetentionDays    types.Int32 `tfsdk:"logs_retention_days"`
	MetricsRetentionDays types.Int32 `tfsdk:"metrics_retention_days"`
	TracesRetentionDays  types.Int32 `tfsdk:"traces_retention_days"`
	DefaultAlertEmails   types.Set   `tfsdk:"default_alert_emails"`
}

// Security //

type SecurityModel struct {
	ThreatDetection []SecurityThreatDetectionModel `tfsdk:"threat_detection"`
}

// Security -> Threat Detection //

type SecurityThreatDetectionModel struct {
	Enabled         types.Bool                           `tfsdk:"enabled"`
	MinimumSeverity types.String                         `tfsdk:"minimum_severity"`
	Syslog          []SecurityThreatDetectionSyslogModel `tfsdk:"syslog"`
}

// Security -> Threat Detection -> Syslog //

type SecurityThreatDetectionSyslogModel struct {
	Transport types.String `tfsdk:"transport"`
	Host      types.String `tfsdk:"host"`
	Port      types.Int32  `tfsdk:"port"`
}

// Status //

type StatusModel struct {
	AccountLink types.String `tfsdk:"account_link"`
	Active      types.Bool   `tfsdk:"active"`
}

func (s StatusModel) AttributeTypes() attr.Type {
	return types.ObjectType{
		AttrTypes: map[string]attr.Type{
			"account_link": types.StringType,
			"active":       types.BoolType,
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

// Coralogix Logging //

type CoralogixLoggingModel struct {
	Cluster     types.String `tfsdk:"cluster"`
	Credentials types.String `tfsdk:"credentials"`
	App         types.String `tfsdk:"app"`
	Subsystem   types.String `tfsdk:"subsystem"`
}

// Datadog Logging //

type DatadogLoggingModel struct {
	Host        types.String `tfsdk:"host"`
	Credentials types.String `tfsdk:"credentials"`
}

// Logzio Logging //

type LogzioLoggingModel struct {
	ListenerHost types.String `tfsdk:"listener_host"`
	Credentials  types.String `tfsdk:"credentials"`
}

// Elastic Logging //

type ElasticLoggingModel struct {
	AWS          []ElasticLoggingAwsModel          `tfsdk:"aws"`
	ElasticCloud []ElasticLoggingElasticCloudModel `tfsdk:"elastic_cloud"`
	Generic      []ElasticLoggingGenericModel      `tfsdk:"generic"`
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

// Elastic Logging -> Elastic Cloud //

type ElasticLoggingElasticCloudModel struct {
	Index       types.String `tfsdk:"index"`
	Type        types.String `tfsdk:"type"`
	Credentials types.String `tfsdk:"credentials"`
	CloudID     types.String `tfsdk:"cloud_id"`
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

// Cloud Watch Logging //

type CloudWatchModel struct {
	Region        types.String `tfsdk:"region"`
	Credentials   types.String `tfsdk:"credentials"`
	RetentionDays types.Int32  `tfsdk:"retention_days"`
	GroupName     types.String `tfsdk:"group_name"`
	StreamName    types.String `tfsdk:"stream_name"`
	ExtractFields types.Map    `tfsdk:"extract_fields"`
}

// Fluentd Logging //

type FluentdLoggingModel struct {
	Host types.String `tfsdk:"host"`
	Port types.Int32  `tfsdk:"port"`
}

// Stackdriver Logging //

type StackdriverLoggingModel struct {
	Credentials types.String `tfsdk:"credentials"`
	Location    types.String `tfsdk:"location"`
}

// Syslog Logging //

type SyslogLoggingModel struct {
	Host     types.String `tfsdk:"host"`
	Port     types.Int32  `tfsdk:"port"`
	Mode     types.String `tfsdk:"mode"`
	Format   types.String `tfsdk:"format"`
	Severity types.Int32  `tfsdk:"severity"`
}
