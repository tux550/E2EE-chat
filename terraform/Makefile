VAR_FILE = secret.tfvars

init:
	terraform init
plan:
	terraform plan -var-file=$(VAR_FILE) -out=plan.out
graph:
	terraform graph | dot -Tpng > graph.png
apply:
	terraform apply plan.out
destroy:
	terraform destroy -var-file=$(VAR_FILE)
fmt:
	terraform fmt