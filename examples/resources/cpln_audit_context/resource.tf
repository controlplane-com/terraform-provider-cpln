resource "cpln_audit_context" "example" {
    name = "audit-context-example"
    description = "audit context description" 
    
    tags = {
        terraform_generated = "true"
        example = "true"
    }
}