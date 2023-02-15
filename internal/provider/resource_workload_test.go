package cpln

import (
	"fmt"
	client "terraform-provider-cpln/internal/provider/client"
	"testing"

	"github.com/go-test/deep"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
)

func TestControlPlane_BuildContainersServerless(t *testing.T) {

	unitTestWorkload := client.Workload{}
	buildContainers(generateFlatTestContainer("serverless"), &unitTestWorkload)

	if diff := deep.Equal(unitTestWorkload.Spec.Containers, generateTestContainers("serverless")); diff != nil {
		t.Errorf("Container was not built correctly. Diff: %s", diff)
	}
}

func TestControlPlane_BuildContainersStandard(t *testing.T) {

	unitTestWorkload := client.Workload{}
	buildContainers(generateFlatTestContainer("standard"), &unitTestWorkload)

	if diff := deep.Equal(unitTestWorkload.Spec.Containers, generateTestContainers("standard")); diff != nil {
		t.Errorf("Container was not built correctly. Diff: %s", diff)
	}
}

func generateTestContainers(workloadType string) *[]client.ContainerSpec {

	newContainer := client.ContainerSpec{
		Name:             GetString("container-01"),
		Image:            GetString("gcr.io/knative-samples/helloworld-go"),
		Memory:           GetString("128Mi"),
		CPU:              GetString("50m"),
		Command:          GetString("override-command"),
		InheritEnv:       GetBool(false),
		WorkingDirectory: GetString("/usr"),
	}

	if workloadType == "serverless" {
		newContainer.Port = GetInt(8080)
	} else if workloadType == "standard" {
		newContainer.Ports = &[]client.PortSpec{
			{
				Protocol: GetString("http"),
				Number:   GetInt(80),
			},
			{
				Protocol: GetString("http2"),
				Number:   GetInt(8080),
			},
			{
				Protocol: GetString("grpc"),
				Number:   GetInt(3000),
			},
			{
				Protocol: GetString("tcp"),
				Number:   GetInt(3001),
			},
		}
	}

	newContainer.Args = &[]string{
		"arg-01",
		"arg-02",
	}

	newContainer.Env = &[]client.NameValue{
		{
			Name:  GetString("env-name-01"),
			Value: GetString("env-value-01"),
		},
		{
			Name:  GetString("env-name-02"),
			Value: GetString("env-value-02"),
		},
	}

	newContainer.Volumes = &[]client.VolumeSpec{
		{
			Uri:  GetString("s3://bucket"),
			Path: GetString("/testpath01"),
		},
		{
			Uri:  GetString("azureblob://storageAccount/container"),
			Path: GetString("/testpath02"),
		},
	}

	newContainer.Metrics = &client.Metrics{
		Path: GetString("/metrics"),
		Port: GetInt(8181),
	}

	newContainer.ReadinessProbe = &client.HealthCheckSpec{

		InitialDelaySeconds: GetInt(1),
		PeriodSeconds:       GetInt(11),
		TimeoutSeconds:      GetInt(2),
		SuccessThreshold:    GetInt(2),
		FailureThreshold:    GetInt(4),

		TCPSocket: &client.TCPSocket{
			Port: GetInt(8181),
		},
	}

	newContainer.LivenessProbe = &client.HealthCheckSpec{

		InitialDelaySeconds: GetInt(2),
		PeriodSeconds:       GetInt(10),
		TimeoutSeconds:      GetInt(3),
		SuccessThreshold:    GetInt(1),
		FailureThreshold:    GetInt(5),

		HTTPGet: &client.HTTPGet{
			Path:   GetString("/path"),
			Port:   GetInt(8282),
			Scheme: GetString("HTTPS"),
			HTTPHeaders: &[]client.NameValue{
				{
					Name:  GetString("header-name-01"),
					Value: GetString("header-value-01"),
				},
				{
					Name:  GetString("header-name-02"),
					Value: GetString("header-value-02"),
				},
			},
		},
	}

	newContainer.LifeCycle = &client.LifeCycleSpec{

		PreStop: &client.LifeCycleInner{
			Exec: &client.Exec{
				Command: &[]string{
					"lc_pre_1", "lc_pre_2", "lc_pre_3",
				},
			},
		},

		PostStart: &client.LifeCycleInner{
			Exec: &client.Exec{
				Command: &[]string{
					"lc_post_1", "lc_post_2", "lc_post_3",
				},
			},
		},
	}

	testContainers := make([]client.ContainerSpec, 1)
	testContainers[0] = newContainer

	return &testContainers
}

func generateFlatTestContainer(workloadType string) []interface{} {

	c := map[string]interface{}{
		"name":  "container-01",
		"image": "gcr.io/knative-samples/helloworld-go",
		// "port":              8080,
		"memory":            "128Mi",
		"cpu":               "50m",
		"command":           "override-command",
		"working_directory": "/usr",
		"inherit_env":       false,
	}

	if workloadType == "serverless" {
		c["port"] = 8080
	} else if workloadType == "standard" {

		port_01 := make(map[string]interface{})
		port_01["protocol"] = "http"
		port_01["number"] = 80

		port_02 := make(map[string]interface{})
		port_02["protocol"] = "http2"
		port_02["number"] = 8080

		port_03 := make(map[string]interface{})
		port_03["protocol"] = "grpc"
		port_03["number"] = 3000

		port_04 := make(map[string]interface{})
		port_04["protocol"] = "tcp"
		port_04["number"] = 3001

		c["ports"] = []interface{}{
			port_01,
			port_02,
			port_03,
			port_04,
		}
	}

	c["args"] = []interface{}{
		"arg-01",
		"arg-02",
	}

	envs := map[string]interface{}{
		"env-name-01": "env-value-01",
		"env-name-02": "env-value-02",
	}

	c["env"] = envs

	volume_01 := make(map[string]interface{})
	volume_01["uri"] = "s3://bucket"
	volume_01["path"] = "/testpath01"

	volume_02 := make(map[string]interface{})
	volume_02["uri"] = "azureblob://storageAccount/container"
	volume_02["path"] = "/testpath02"

	c["volume"] = []interface{}{
		volume_01,
		volume_02,
	}

	metrics := make(map[string]interface{})
	metrics["path"] = "/metrics"
	metrics["port"] = 8181

	c["metrics"] = []interface{}{
		metrics,
	}

	readiness := make(map[string]interface{})

	readiness["initial_delay_seconds"] = 1
	readiness["period_seconds"] = 11
	readiness["timeout_seconds"] = 2
	readiness["success_threshold"] = 2
	readiness["failure_threshold"] = 4

	tcpSocket := make(map[string]interface{})
	tcpSocket["port"] = 8181

	tcpSocketAsInterface := []interface{}{tcpSocket}
	readiness["tcp_socket"] = tcpSocketAsInterface

	c["readiness_probe"] = []interface{}{readiness}

	liveness := make(map[string]interface{})

	liveness["initial_delay_seconds"] = 2
	liveness["period_seconds"] = 10
	liveness["timeout_seconds"] = 3
	liveness["success_threshold"] = 1
	liveness["failure_threshold"] = 5

	h := make(map[string]interface{})
	h["path"] = "/path"
	h["port"] = 8282
	h["scheme"] = "HTTPS"

	headers := make(map[string]interface{})
	headers["header-name-01"] = "header-value-01"
	headers["header-name-02"] = "header-value-02"

	h["http_headers"] = headers

	hi := []interface{}{h}

	liveness["http_get"] = hi

	c["liveness_probe"] = []interface{}{liveness}

	localContainers := []interface{}{
		c,
	}

	return localContainers
}

func TestControlPlane_BuildOptions(t *testing.T) {

	unitTestWorkload := client.Workload{}

	buildOptions(generateFlatTestOptions(), &unitTestWorkload, false, "")

	if diff := deep.Equal(unitTestWorkload.Spec.DefaultOptions, generateTestOptions("serverless")); diff != nil {
		t.Errorf("Options was not built correctly. Diff: %s", diff)
	}
}

func generateTestOptions(workloadType string) *client.Options {

	if workloadType == "standard" {
		return &client.Options{
			CapacityAI:     GetBool(false),
			TimeoutSeconds: GetInt(30),
			Debug:          GetBool(false),
			Suspend:        GetBool(false),

			AutoScaling: &client.AutoScaling{
				Metric:           GetString("cpu"),
				Target:           GetInt(60),
				MaxScale:         GetInt(3),
				MinScale:         GetInt(2),
				MaxConcurrency:   GetInt(500),
				ScaleToZeroDelay: GetInt(400),
			},
		}
	}

	return &client.Options{
		CapacityAI:     GetBool(true),
		TimeoutSeconds: GetInt(30),
		Debug:          GetBool(false),
		Suspend:        GetBool(false),

		AutoScaling: &client.AutoScaling{
			Metric:           GetString("concurrency"),
			Target:           GetInt(100),
			MaxScale:         GetInt(3),
			MinScale:         GetInt(2),
			MaxConcurrency:   GetInt(500),
			ScaleToZeroDelay: GetInt(400),
		},
	}
}

func generateFlatTestOptions() []interface{} {

	as := map[string]interface{}{
		"metric":              "concurrency",
		"target":              100,
		"max_scale":           3,
		"min_scale":           2,
		"max_concurrency":     500,
		"scale_to_zero_delay": 400,
	}

	asi := []interface{}{
		as,
	}

	o := map[string]interface{}{
		"capacity_ai":     true,
		"timeout_seconds": 30,
		"autoscaling":     asi,
		"debug":           false,
		"suspend":         false,
	}

	return []interface{}{
		o,
	}
}

func TestControlPlane_BuildFirewallSpec(t *testing.T) {

	unitTestWorkload := client.Workload{}

	buildFirewallSpec(generateFlatTestFirewallSpec(true), &unitTestWorkload, false)

	if diff := deep.Equal(unitTestWorkload.Spec.FirewallConfig, generateTestFirewallSpec()); diff != nil {
		t.Errorf("FirewallSpec was not built correctly. Diff: %s", diff)
	}
}

func generateTestFirewallSpec() *client.FirewallSpec {

	return &client.FirewallSpec{
		External: &client.FirewallSpecExternal{
			InboundAllowCIDR: &[]string{"0.0.0.0/0"},
			// OutboundAllowCIDR:     &[]string{},
			OutboundAllowHostname: &[]string{"*.cpln.io", "*.controlplane.com"},
		},
		Internal: &client.FirewallSpecInternal{
			InboundAllowType: GetString("none"),
			// InboundAllowWorkload: &[]string{},
		},
	}
}

func generateFlatTestFirewallSpec(useSet bool) []interface{} {

	stringFunc := schema.HashSchema(StringSchema())
	e := make(map[string]interface{})

	if useSet {
		e["inbound_allow_cidr"] = schema.NewSet(stringFunc, []interface{}{"0.0.0.0/0"})
	} else {
		e["inbound_allow_cidr"] = []string{"0.0.0.0/0"}
	}

	// e["outbound_allow_cidr"] = []interface{}{}

	if useSet {
		e["outbound_allow_hostname"] = schema.NewSet(stringFunc, []interface{}{"*.cpln.io", "*.controlplane.com"})
	} else {
		e["outbound_allow_hostname"] = []interface{}{"*.cpln.io", "*.controlplane.com"}
	}

	i := make(map[string]interface{})
	i["inbound_allow_type"] = "none"
	// i["inbound_allow_workload"] = []interface{}{}

	fc := map[string]interface{}{
		"external": []interface{}{
			e,
		},
		"internal": []interface{}{
			i,
		},
	}

	return []interface{}{
		fc,
	}
}

func TestControlPlane_FlattenWorkloadStatus(t *testing.T) {

	endpoint := "endpoint"
	parent_id := "parent_id"
	canonical := "canonical"

	status := &client.WorkloadStatus{
		Endpoint:          GetString(endpoint),
		ParentID:          GetString(parent_id),
		CanonicalEndpoint: GetString(canonical),
	}

	flatStatus := map[string]interface{}{
		"endpoint":           "endpoint",
		"parent_id":          "parent_id",
		"canonical_endpoint": "canonical",
	}

	flatStatusArray := []interface{}{
		flatStatus,
	}

	flattenedStatus := flattenWorkloadStatus(status)

	if diff := deep.Equal(flattenedStatus, flatStatusArray); diff != nil {
		t.Errorf("Workload Status was not flattened correctly. Diff: %s", diff)
	}
}

func TestControlPlane_FlattenContainerServerless(t *testing.T) {

	containers := generateTestContainers("serverless")
	flattenedContainer := flattenContainer(containers)

	flatContainer := generateFlatTestContainer("serverless")

	if diff := deep.Equal(flatContainer, flattenedContainer); diff != nil {
		t.Errorf("Container was not flattened correctly. Diff: %s", diff)
	}
}

func TestControlPlane_FlattenContainerStandard(t *testing.T) {

	containers := generateTestContainers("standard")
	flattenedContainer := flattenContainer(containers)

	flatContainer := generateFlatTestContainer("standard")

	if diff := deep.Equal(flatContainer, flattenedContainer); diff != nil {
		t.Errorf("Container was not flattened correctly. Diff: %s", diff)
	}
}

func TestControlPlane_FlattenOptions(t *testing.T) {

	options := generateTestOptions("serverless")
	flatOptions := generateFlatTestOptions()
	flattenedOptions := flattenOptions([]client.Options{*options}, false, "")

	if diff := deep.Equal(flatOptions, flattenedOptions); diff != nil {
		t.Errorf("Options not flattened correctly. Diff: %s", diff)
	}
}

func TestControlPlane_FlattenFirewallSpec(t *testing.T) {

	spec := generateTestFirewallSpec()
	flattenedFirewallSpec := flattenFirewallSpec(spec)

	flatSpec := generateFlatTestFirewallSpec(false)

	if diff := deep.Equal(flatSpec, flattenedFirewallSpec); diff != nil {
		t.Errorf("FirewallSpec not flattened correctly. Diff: %s", diff)
	}
}

func TestAccControlPlaneWorkload_basic(t *testing.T) {

	var testWorkload client.Workload

	gName := "gvc-" + acctest.RandStringFromCharSet(10, acctest.CharSetAlphaNum)
	wName := "workload-" + acctest.RandStringFromCharSet(10, acctest.CharSetAlphaNum)
	randomName := acctest.RandStringFromCharSet(10, acctest.CharSetAlphaNum)

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t, "WORKLOAD") },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckControlPlaneWorkloadCheckDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccControlPlaneWorkload(randomName, gName, "GVC created using terraform for acceptance tests", wName, "Workload created using terraform for acceptance tests"),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckControlPlaneWorkloadExists("cpln_workload.new", wName, gName, &testWorkload),
					testAccCheckControlPlaneWorkloadAttributes(&testWorkload, "serverless"),
					resource.TestCheckResourceAttr("cpln_gvc.new", "description", "GVC created using terraform for acceptance tests"),
					resource.TestCheckResourceAttr("cpln_workload.new", "description", "Workload created using terraform for acceptance tests"),
				),
			},
			{
				Config: testAccControlPlaneWorkload(randomName, gName, "GVC created using terraform for acceptance tests", wName+"renamed", "Renamed Workload created using terraform for acceptance tests"),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckControlPlaneWorkloadExists("cpln_workload.new", wName+"renamed", gName, &testWorkload),
					testAccCheckControlPlaneWorkloadAttributes(&testWorkload, "serverless"),
					resource.TestCheckResourceAttr("cpln_workload.new", "description", "Renamed Workload created using terraform for acceptance tests"),
				),
			},
			{
				Config: testAccControlPlaneWorkload(randomName, gName, "GVC created using terraform for acceptance tests", wName+"renamed", "Updated Workload description created using terraform for acceptance tests"),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckControlPlaneWorkloadExists("cpln_workload.new", wName+"renamed", gName, &testWorkload),
					testAccCheckControlPlaneWorkloadAttributes(&testWorkload, "serverless"),
					resource.TestCheckResourceAttr("cpln_workload.new", "description", "Updated Workload description created using terraform for acceptance tests"),
				),
			},
			{
				Config: testAccControlPlaneStandardWorkload(randomName, gName, "GVC created using terraform for acceptance tests", wName+"standard", "Standard Workload description created using terraform for acceptance tests"),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckControlPlaneWorkloadExists("cpln_workload.new", wName+"standard", gName, &testWorkload),
					testAccCheckControlPlaneWorkloadAttributes(&testWorkload, "standard"),
					resource.TestCheckResourceAttr("cpln_workload.new", "description", "Standard Workload description created using terraform for acceptance tests"),
				),
			},
		},
	})
}

func testAccControlPlaneWorkload(randomName, gvcName, gvcDescription, workloadName, workloadDescription string) string {

	TestLogger.Printf("Inside testAccControlPlaneWorkload")

	return fmt.Sprintf(`

	variable "random-name" {
		type = string
		default = "%s"
	}

	resource "cpln_gvc" "new" {
		name        = "%s"	
		description = "%s"

		locations = ["aws-eu-central-1", "aws-us-west-2"]
	  
		tags = {
		  terraform_generated = "true"
		  acceptance_test = "true"
		}
	}

	resource "cpln_identity" "new" {

		gvc = cpln_gvc.new.name
	  
		name        = "terraform-identity-${var.random-name}"
		description = "Identity created using terraform"
	  
		tags = {
		  terraform_generated = "true"
		  acceptance_test     = "true"
		}
	}
	  
	resource "cpln_workload" "new" {

		gvc = cpln_gvc.new.name
	  
		name        = "%s"
		description = "%s"
	  
		tags = {
		  terraform_generated = "true"
		  acceptance_test = "true"
		}

		identity_link = cpln_identity.new.self_link
	  
		container {
		  name  = "container-01"
		  image = "gcr.io/knative-samples/helloworld-go"
		  port = 8080
		  memory = "128Mi"
		  cpu = "50m"

		  command = "override-command"
		  working_directory = "/usr"
	  
		  env = {
			env-name-01 = "env-value-01",
			env-name-02 = "env-value-02",
		  }
	  
		  args = ["arg-01", "arg-02"]

		  volume {
			uri  = "s3://bucket"
			path = "/testpath01"
		  }

		  volume {
			uri  = "azureblob://storageAccount/container"
			path = "/testpath02"
		  }

		  metrics {
			path = "/metrics"
			port = 8181
		  }

		  readiness_probe {

				tcp_socket {
					port = 8181
				}

			// exec {
			// 	command = ["test1", "test2"]
			// }
	  
				period_seconds       = 11
				timeout_seconds      = 2
				failure_threshold    = 4
				success_threshold    = 2
				initial_delay_seconds = 1
		  }

		  liveness_probe {

				http_get {
					path = "/path"
					port = 8282
					scheme = "HTTPS"
					http_headers = {
						header-name-01 = "header-value-01"
						header-name-02 = "header-value-02"
					}
				}
	  
				period_seconds       = 10
				timeout_seconds      = 3
				failure_threshold    = 5
				success_threshold    = 1
				initial_delay_seconds = 2
		  }

			lifecycle {
				pre_stop {
					exec {
						command = ["lc_pre_1", "lc_pre_2", "lc_pre_3"]
					}
				}
	
				post_start {
					exec {
						command = ["lc_post_1", "lc_post_2", "lc_post_3"]
					}
				}
			}
		}

		// container {
		// 	name  = "container-02"
		// 	image = "gcr.io/knative-samples/helloworld-go"
		// 	memory = "128Mi"
		// 	cpu = "50m"
		
		// 	env = {
		// 	  env-name-01 = "env-value-01",
		// 	  env-name-02 = "env-value-02",
		// 	}
		
		// 	args = ["arg-01", "arg-02"]
		// }
		 	  	  
		options {
		  capacity_ai = true
		  timeout_seconds = 30
		  suspend = false
	  
		  autoscaling {
			metric = "concurrency"
			target = 100
			max_scale = 3
			min_scale = 2
			max_concurrency = 500
			scale_to_zero_delay = 400
		  }
		}

		// locations = ["aws-eu-central-1", "aws-us-west-2", "azure-eastus2", "azure-eastus2"]

		local_options {
			location = "aws-eu-central-1"
			capacity_ai = true
			timeout_seconds = 30
			suspend = false
		
			autoscaling {
			  metric = "concurrency"
			  target = 100
			  max_scale = 3
			  min_scale = 2
			  max_concurrency = 500
			  scale_to_zero_delay = 400
			}
		}
	  
		firewall_spec {
		  external {
			inbound_allow_cidr =  ["0.0.0.0/0"]
			// outbound_allow_cidr =  []
			outbound_allow_hostname =  ["*.controlplane.com", "*.cpln.io"]
		  }
		  internal { 
			# Allowed Types: "none", "same-gvc", "same-org", "workload-list"
			inbound_allow_type = "none"
			inbound_allow_workload = []
		  }
		}
	  }
	  `, randomName, gvcName, gvcDescription, workloadName, workloadDescription)
}

func testAccCheckControlPlaneWorkloadExists(resourceName, workloadName, gvcName string, workload *client.Workload) resource.TestCheckFunc {
	return func(s *terraform.State) error {

		TestLogger.Printf("Inside testAccCheckControlPlaneWorkloadExists. Resources Length: %d", len(s.RootModule().Resources))

		rs, ok := s.RootModule().Resources[resourceName]

		if !ok {
			return fmt.Errorf("Not found: %s", s)
		}

		if rs.Primary.ID != workloadName {
			return fmt.Errorf("Workload name does not match")
		}

		client := testAccProvider.Meta().(*client.Client)

		wl, _, err := client.GetWorkload(workloadName, gvcName)

		if err != nil {
			return err
		}

		if *wl.Name != workloadName {
			return fmt.Errorf("Workload name does not match")
		}

		*workload = *wl

		return nil
	}
}

func testAccCheckControlPlaneWorkloadAttributes(workload *client.Workload, workloadType string) resource.TestCheckFunc {
	return func(s *terraform.State) error {

		tags := *workload.Tags

		if tags["terraform_generated"] != "true" {
			return fmt.Errorf("Tags - workload terraform_generated attribute does not match")
		}

		if tags["acceptance_test"] != "true" {
			return fmt.Errorf("Tags - workload acceptance_test attribute does not match")
		}

		containers := generateTestContainers(workloadType)

		if diff := deep.Equal(containers, workload.Spec.Containers); diff != nil {
			return fmt.Errorf("Containers attributes does not match. Diff: %s", diff)
		}

		options := generateTestOptions(workloadType)

		if diff := deep.Equal(options, workload.Spec.DefaultOptions); diff != nil {
			return fmt.Errorf("Options attributes does not match. Diff: %s", diff)
		}

		firewallSpec := generateTestFirewallSpec()

		if diff := deep.Equal(firewallSpec, workload.Spec.FirewallConfig); diff != nil {
			return fmt.Errorf("FirewallSpec attributes does not match. Diff: %s", diff)
		}

		return nil
	}
}

func testAccControlPlaneStandardWorkload(randomName, gvcName, gvcDescription, workloadName, workloadDescription string) string {

	TestLogger.Printf("Inside testAccControlPlaneWorkload")

	return fmt.Sprintf(`

	variable "random-name" {
		type = string
		default = "%s"
	}

	resource "cpln_gvc" "new" {
		name        = "%s"	
		description = "%s"

		locations = ["aws-eu-central-1", "aws-us-west-2"]
	  
		tags = {
		  terraform_generated = "true"
		  acceptance_test = "true"
		}
	}

	resource "cpln_identity" "new" {

		gvc = cpln_gvc.new.name
	  
		name        = "terraform-identity-${var.random-name}"
		description = "Identity created using terraform"
	  
		tags = {
		  terraform_generated = "true"
		  acceptance_test     = "true"
		}
	}
	  
	resource "cpln_workload" "new" {

		gvc = cpln_gvc.new.name
	  
		name        = "%s"
		description = "%s"
	  
		tags = {
		  terraform_generated = "true"
		  acceptance_test = "true"
		}

		identity_link = cpln_identity.new.self_link

		type = "standard" 
	  
		container {
		  name  = "container-01"
		  image = "gcr.io/knative-samples/helloworld-go"
		  memory = "128Mi"
		  cpu = "50m"	  

		  ports {
		    protocol = "http"
			number   = "80" 
		  }

		  ports {
			protocol = "http2"
			number   = "8080" 
	      }

		  ports {
			protocol = "grpc"
			number   = "3000" 
	      }

		  ports {
			protocol = "tcp"
			number   = "3001" 
	      }


		  command = "override-command"
		  working_directory = "/usr"
	  
		  env = {
			env-name-01 = "env-value-01",
			env-name-02 = "env-value-02",
		  }
	  
		  args = ["arg-01", "arg-02"]

		  lifecycle {

			pre_stop {
				exec {
					command = ["lc_pre_1", "lc_pre_2", "lc_pre_3"]
				}
			}

			post_start {
				exec {
					command = ["lc_post_1", "lc_post_2", "lc_post_3"]
				}
			}
		  }

		  volume {
			uri  = "s3://bucket"
			path = "/testpath01"
		  }

		  volume {
			uri  = "azureblob://storageAccount/container"
			path = "/testpath02"
		  }

		  metrics {
			path = "/metrics"
			port = 8181
		  }

		  readiness_probe {

			tcp_socket {
			  port = 8181
			}

			// exec {
			// 	command = ["test1", "test2"]
			// }
	  
			period_seconds       = 11
			timeout_seconds      = 2
			failure_threshold    = 4
			success_threshold    = 2
			initial_delay_seconds = 1
		  }

		  liveness_probe {

			http_get {
				path = "/path"
				port = 8282
				scheme = "HTTPS"
				http_headers = {
					header-name-01 = "header-value-01"
					header-name-02 = "header-value-02"
				}
			}
	  
			period_seconds       = 10
			timeout_seconds      = 3
			failure_threshold    = 5
			success_threshold    = 1
			initial_delay_seconds = 2
		  }
		}
	 	  	  
		options {
		  capacity_ai = false
		  timeout_seconds = 30
		  suspend = false
	  
		  autoscaling {
			metric = "cpu"
			target = 60
			max_scale = 3
			min_scale = 2
			max_concurrency = 500
			scale_to_zero_delay = 400
		  }
		}

		// locations = ["aws-eu-central-1", "aws-us-west-2", "azure-eastus2", "azure-eastus2"]

		// local_options {
		// 	location = "aws-eu-central-1"
		// 	capacity_ai = true
		// 	timeout_seconds = 30
		
		// 	autoscaling {
		// 	  metric = "concurrency"
		// 	  target = 100
		// 	  max_scale = 3
		// 	  min_scale = 2
		// 	  max_concurrency = 500
		// 	  scale_to_zero_delay = 400
		// 	}
		// }
	  
		firewall_spec {
		  external {
			inbound_allow_cidr =  ["0.0.0.0/0"]
			// outbound_allow_cidr =  []
			outbound_allow_hostname =  ["*.controlplane.com", "*.cpln.io"]
		  }
		  internal { 
			# Allowed Types: "none", "same-gvc", "same-org", "workload-list"
			inbound_allow_type = "none"
			inbound_allow_workload = []
		  }
		}
	  }
	  `, randomName, gvcName, gvcDescription, workloadName, workloadDescription)
}

func testAccCheckControlPlaneWorkloadCheckDestroy(s *terraform.State) error {

	TestLogger.Printf("Inside testAccCheckControlPlaneWorkloadCheckDestroy. Resources Length: %d", len(s.RootModule().Resources))

	if len(s.RootModule().Resources) == 0 {
		return fmt.Errorf("Error In CheckDestroy. No Resources To Verify")
	}

	c := testAccProvider.Meta().(*client.Client)

	for _, rs := range s.RootModule().Resources {

		TestLogger.Printf("Inside testAccCheckControlPlaneWorkloadDestroy: rs.Type: %s", rs.Type)

		if rs.Type != "cpln_gvc" {
			continue
		}

		gvcName := rs.Primary.ID

		TestLogger.Printf("Inside testAccCheckControlPlaneWorkloadDestroy: gvcName: %s", gvcName)

		gvc, _, _ := c.GetGvc(gvcName)
		if gvc != nil {
			return fmt.Errorf("GVC still exists. Name: %s. Associated Workloads might still exist", *gvc.Name)
		}
	}

	return nil
}
