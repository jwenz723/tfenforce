# tfenforce

A simple example written in go of how to enforce policies/standards on terraform resources.

## Examples

## [No Violations](examples/no_violation)

Generate the plan and run tfenforce:

```bash
cd examples/no_violation
terraform init
terraform plan -out tfplan
terraform show -json tfplan > tfplan.json
go run ../../main.go -tfplan ./tfplan.json
```

This will output nothing.

## [IAM Policy Violation](examples/iam_violation)

Generate the plan and run tfenforce:

```bash
cd examples/iam_violation
terraform init
terraform plan -out tfplan
terraform show -json tfplan > tfplan.json
go run ../../main.go -tfplan ./tfplan.json
```

This will output a single violation:

```bash
cannot declare '*' IAM permission (aws_iam_policy.service-policy-default)
```

## [VPC Policy Violation](examples/vpc_violation)

Generate the plan and run tfenforce:

```bash
cd examples/vpc_violation
terraform init
terraform plan -out tfplan
terraform show -json tfplan > tfplan.json
go run ../../main.go -tfplan ./tfplan.json
```

This will output a single violation:

```bash
cannot define aws_vpc resources (aws_vpc.test)
```