# Production Environment - Terraform Configuration

This directory contains the Terraform configuration for the **production** environment of the order application infrastructure on AWS.

## 📁 Structure

```
production/
├── main.tf           # Main resource definitions
├── variables.tf      # Variable definitions
├── outputs.tf        # Output definitions
├── versions.tf       # Terraform and provider version constraints
├── provider.tf       # AWS provider configuration
├── terraform.tfvars  # Production-specific values
└── README.md         # This file
```

## 🏗️ Infrastructure Components

This configuration deploys:

1. **Networking** - VPC, subnets (public, private, database), security groups
2. **Database** - RDS PostgreSQL with Multi-AZ, automated backups
3. **Compute** - ECS Fargate cluster and services
4. **Load Balancer** - Application Load Balancer (ALB)
5. **Container Registry** - Amazon ECR
6. **Logging** - CloudWatch Log Groups
7. **Authentication** - GitHub OIDC for CI/CD

## 🚀 Getting Started

### Prerequisites

1. **AWS CLI** configured with appropriate credentials
2. **Terraform** >= 1.5.0 installed
3. **S3 bucket** for state storage (already configured: `production-terraform-up-and-running-state`)
4. **DynamoDB table** for state locking (`terraform-state-lock`)

### Setup State Backend

Before running Terraform, ensure the S3 bucket and DynamoDB table exist:

```bash
# Create S3 bucket for state (if not exists)
aws s3 mb s3://production-terraform-up-and-running-state --region ap-southeast-1

# Enable versioning
aws s3api put-bucket-versioning \
  --bucket production-terraform-up-and-running-state \
  --versioning-configuration Status=Enabled

# Enable encryption
aws s3api put-bucket-encryption \
  --bucket production-terraform-up-and-running-state \
  --server-side-encryption-configuration '{
    "Rules": [{
      "ApplyServerSideEncryptionByDefault": {
        "SSEAlgorithm": "AES256"
      }
    }]
  }'

# Create DynamoDB table for state locking
aws dynamodb create-table \
  --table-name terraform-state-lock \
  --attribute-definitions AttributeName=LockID,AttributeType=S \
  --key-schema AttributeName=LockID,KeyType=HASH \
  --billing-mode PAY_PER_REQUEST \
  --region ap-southeast-1
```

### Initialize Terraform

```bash
cd terraform/production
terraform init
```

### Review Configuration

Edit `terraform.tfvars` to match your production requirements:

```bash
vim terraform.tfvars
```

Key values to update:
- `owner` - Your team name
- `certificate_arn` - Your ACM certificate for HTTPS
- `allowed_repos_branches` - Your GitHub repository details
- Resource sizes (CPU, memory, storage)

### Plan Changes

```bash
terraform plan -out=tfplan
```

Review the plan carefully before applying in production!

### Apply Configuration

```bash
terraform apply tfplan
```

### Get Outputs

```bash
# View all outputs
terraform output

# View specific output
terraform output application_url
terraform output alb_dns_name
```

## 🔒 Security Best Practices

This configuration implements several security best practices:

1. **State Encryption** - S3 backend with encryption enabled
2. **State Locking** - DynamoDB table prevents concurrent modifications
3. **Multi-AZ Deployment** - Database and application across multiple AZs
4. **Private Subnets** - Application and database in private subnets
5. **Security Groups** - Least privilege network access
6. **Secrets Management** - Database credentials in AWS Secrets Manager
7. **Automated Backups** - 30-day retention for RDS
8. **Final Snapshots** - Enabled to prevent data loss
9. **Default Tags** - All resources tagged for governance

## 📊 Resource Tagging

All resources are automatically tagged with:

```hcl
Environment = "production"
Project     = "user-api"
ManagedBy   = "Terraform"
Owner       = "DevOps Team"
```

## 🔄 Workflow

### Standard Changes

```bash
# 1. Make changes to .tf files or terraform.tfvars
# 2. Format code
terraform fmt -recursive

# 3. Validate configuration
terraform validate

# 4. Plan changes
terraform plan -out=tfplan

# 5. Review and apply
terraform apply tfplan
```

### Destroy Resources (⚠️ Use with extreme caution!)

```bash
terraform destroy
```

**Note**: With `rds_skip_final_snapshot = false`, a final snapshot will be created before destroying the database.

## 🌐 Accessing Your Application

After deployment, access your application via:

```bash
# Get the URL
terraform output application_url

# Or directly via ALB DNS
terraform output alb_dns_name
```

## 📝 Common Tasks

### Update Container Image

```bash
# Edit terraform.tfvars
container_image = "123456789012.dkr.ecr.ap-southeast-1.amazonaws.com/user-api:v1.0.1"

# Apply changes
terraform apply
```

### Scale ECS Services

```bash
# Edit terraform.tfvars
ecs_desired_count = 4

# Apply changes
terraform apply
```

### Update Database Size

```bash
# Edit terraform.tfvars
rds_instance_class = "db.t3.medium"

# Apply changes (will cause downtime for single-AZ or brief interruption for Multi-AZ)
terraform apply
```

## 🐛 Troubleshooting

### State Lock Issues

If Terraform is stuck with a state lock:

```bash
# Force unlock (only if you're sure no other operation is running)
terraform force-unlock <LOCK_ID>
```

### Refresh State

```bash
terraform refresh
```

### View Current State

```bash
terraform state list
terraform state show <RESOURCE_NAME>
```

## 📚 Additional Resources

- [Terraform AWS Provider Documentation](https://registry.terraform.io/providers/hashicorp/aws/latest/docs)
- [AWS ECS Best Practices](https://docs.aws.amazon.com/AmazonECS/latest/bestpracticesguide/)
- [Terraform Best Practices](https://www.terraform-best-practices.com/)

## 📞 Support

For issues or questions, contact: DevOps Team

---

**⚠️ Important**: This is the PRODUCTION environment. Always:
- Review plans carefully before applying
- Test changes in development/staging first
- Have a rollback plan
- Communicate changes to the team
- Follow the change management process

