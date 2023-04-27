resource "cpln_memcache" "example" {
    name 		= "memcache-example"
    description = "Memcache description for memcache-example" 
    
    tags = {
        terraform_generated = "true"
        acceptance_test 	= "true"
    }

    node_count = 1
    node_size  = 0.3
    version    = "1.5.22"

    options {
        eviction_disabled 	 = true
        idle_timeout_seconds = 600
        max_item_size 		 = 1024
        max_connections      = 1024
    }

    locations  = ["/org/{your-org}/location/aws-us-west-2"]
}