# Terraform Infrastructure

Infrastructure as Code cho User API sử dụng Terraform modules.

## 📁 Cấu trúc

```
terraform/
├── main.tf                 # Main configuration
├── variables.tf            # Input variables
├── outputs.tf              # Output values
├── terraform.tfvars.example # Example variables
├── environments/           # Environment-specific configs
│   ├── dev.tfvars
│   └── prod.tfvars
└── modules/               # Terraform modules
    ├── networking/        # VPC, Subnets, NAT Gateway
    ├── database/          # RDS PostgreSQL
    └── ecs/              # ECS Cluster, Service, ALB
```

## 🚀 Quick Start

### 1. Initialize

```bash
cd terraform
terraform init
```

### 2. Configure Variables

```bash
# Copy example file
cp terraform.tfvars.example terraform.tfvars

# Or use environment-specific config
cp environments/dev.tfvars terraform.tfvars
```

### 3. Plan

```bash
# Development
terraform plan -var-file="environments/dev.tfvars"

# Production
terraform plan -var-file="environments/prod.tfvars"
```

### 4. Apply

```bash
# Development
terraform apply -var-file="environments/dev.tfvars"

# Production
terraform apply -var-file="environments/prod.tfvars"
```

## 📦 Modules

### Networking Module

Tạo VPC với:
- 3 Availability Zones
- Public subnets (ALB, NAT Gateway)
- Private subnets (ECS tasks)
- Database subnets (RDS)
- NAT Gateways cho internet access
- VPC Flow Logs (optional)

**Outputs:**
- `vpc_id`
- `public_subnet_ids`
- `private_subnet_ids`
- `database_subnet_ids`
- `nat_gateway_ips`

### Database Module

Tạo RDS PostgreSQL với:
- Automated backups
- Multi-AZ deployment (optional)
- Enhanced monitoring
- Secrets Manager integration
- Security groups
- Parameter groups

**Outputs:**
- `db_instance_endpoint`
- `db_secret_arn`
- `db_security_group_id`

### ECS Module

Tạo ECS infrastructure với:
- ECS Cluster (Fargate)
- ECR repository
- Task definitions
- ECS Service
- Application Load Balancer
- Auto Scaling (CPU & Memory)
- CloudWatch Logs
- IAM roles

**Outputs:**
- `ecr_repository_url`
- `ecs_cluster_name`
- `alb_dns_name`
- `application_url`

## 🔧 Configuration

### Required Variables

```hcl
environment  = "prod"        # Environment name
aws_region   = "us-east-1"   # AWS region
```

### Optional Variables

```hcl
# ECS
ecs_task_cpu      = 512
ecs_task_memory   = 1024
ecs_desired_count = 2

# RDS
rds_instance_class = "db.t3.micro"
rds_multi_az       = false

# SSL
certificate_arn = "arn:aws:acm:..."
```

## 📊 Outputs

After deployment:

```bash
# View all outputs
terraform output

# View specific output
terraform output alb_dns_name
terraform output ecr_repository_url
```

## 🏗️ Resource Tagging

Tất cả resources được tag với:
- `Project`: project_name
- `Environment`: environment
- `ManagedBy`: Terraform

## 💰 Cost Estimation

### Development (~$30-50/month)
- NAT Gateway: $45
- RDS t3.micro: $15
- ECS (1 task): $10
- ALB: $20

### Production (~$100-200/month)
- NAT Gateway (3 AZs): $135
- RDS t3.small Multi-AZ: $60
- ECS (2+ tasks): $25+
- ALB: $20
- Data Transfer: Variable

## 🔐 Security

### Best Practices

1. **Remote State**: Configure S3 backend
2. **State Locking**: Use DynamoDB table
3. **Encryption**: Enable for RDS and S3
4. **Secrets**: Use Secrets Manager
5. **Network**: Private subnets for apps
6. **Access**: IAM roles, not keys

### State Management

```hcl
# main.tf
terraform {
  backend "s3" {
    bucket         = "your-terraform-state"
    key            = "user-api/prod/terraform.tfstate"
    region         = "us-east-1"
    encrypt        = true
    dynamodb_table = "terraform-lock"
  }
}
```

## 🔄 Workflow

### Development Workflow

```bash
# 1. Make changes to .tf files
# 2. Format
terraform fmt -recursive

# 3. Validate
terraform validate

# 4. Plan
terraform plan -out=tfplan

# 5. Apply
terraform apply tfplan
```

### Multi-Environment

```bash
# Dev
terraform workspace new dev
terraform apply -var-file="environments/dev.tfvars"

# Prod
terraform workspace new prod
terraform apply -var-file="environments/prod.tfvars"
```

## 🧹 Cleanup

```bash
# Destroy resources
terraform destroy -var-file="environments/dev.tfvars"
```

⚠️ **Warning**: Always backup production data before destroying!

## 📚 Documentation

- [Networking Module](./modules/networking/)
- [Database Module](./modules/database/)
- [ECS Module](./modules/ecs/)

## 🐛 Troubleshooting

### Issue: Terraform State Lock

```bash
# Remove lock (if safe)
terraform force-unlock <LOCK_ID>
```

### Issue: Resource Already Exists

```bash
# Import existing resource
terraform import module.networking.aws_vpc.main vpc-xxxxx
```

### Issue: Plan Too Large

```bash
# Target specific resource
terraform plan -target=module.ecs
terraform apply -target=module.ecs
```


