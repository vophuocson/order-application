# Deployment Guide - User API

Hướng dẫn deployment ứng dụng User API lên AWS ECS sử dụng Terraform.

## 📋 Mục lục

- [Yêu cầu](#yêu-cầu)
- [Kiến trúc](#kiến-trúc)
- [Cấu hình](#cấu-hình)
- [Deployment](#deployment)
- [Quản lý](#quản-lý)
- [Monitoring](#monitoring)
- [Troubleshooting](#troubleshooting)

## 🔧 Yêu cầu

### Tools cần thiết

```bash
# AWS CLI
aws --version  # >= 2.0

# Terraform
terraform --version  # >= 1.0

# Docker
docker --version  # >= 20.0

# Make
make --version
```

### AWS Credentials

```bash
# Cấu hình AWS credentials
aws configure

# Hoặc export environment variables
export AWS_ACCESS_KEY_ID="your-access-key"
export AWS_SECRET_ACCESS_KEY="your-secret-key"
export AWS_REGION="us-east-1"
```

## 🏗️ Kiến trúc

```
┌─────────────────────────────────────────────────────────┐
│                        Internet                          │
└────────────────────┬────────────────────────────────────┘
                     │
            ┌────────▼────────┐
            │  Application    │
            │  Load Balancer  │
            └────────┬────────┘
                     │
    ┌────────────────┼────────────────┐
    │                │                │
┌───▼───┐       ┌───▼───┐       ┌───▼───┐
│  ECS  │       │  ECS  │       │  ECS  │
│ Task 1│       │ Task 2│       │ Task 3│
└───┬───┘       └───┬───┘       └───┬───┘
    │                │                │
    └────────────────┼────────────────┘
                     │
            ┌────────▼────────┐
            │   RDS Postgres  │
            │   Multi-AZ      │
            └─────────────────┘
```

### Infrastructure Components

- **VPC**: 3 Availability Zones
  - Public Subnets: ALB, NAT Gateways
  - Private Subnets: ECS Tasks
  - Database Subnets: RDS
- **ECS Fargate**: Serverless containers
- **Application Load Balancer**: HTTP/HTTPS traffic
- **RDS PostgreSQL**: Managed database
- **ECR**: Docker image registry
- **CloudWatch**: Logs and monitoring
- **Secrets Manager**: Database credentials

## ⚙️ Cấu hình

### 1. Environment Variables

```bash
# Copy và cấu hình file .env
cp .env.example .env
# Chỉnh sửa .env với thông tin của bạn
```

### 2. Terraform Variables

```bash
# Copy và cấu hình terraform variables
cp terraform/terraform.tfvars.example terraform/terraform.tfvars

# Hoặc sử dụng environment-specific configs
# Development
cp terraform/environments/dev.tfvars terraform/terraform.tfvars

# Production
cp terraform/environments/prod.tfvars terraform/terraform.tfvars
```

### 3. Cấu hình quan trọng

**terraform/terraform.tfvars:**
```hcl
# Required
project_name = "user-api"
environment  = "prod"
aws_region   = "us-east-1"

# Optional - SSL Certificate
certificate_arn = "arn:aws:acm:us-east-1:xxx:certificate/xxx"

# Optional - Scaling
ecs_min_capacity = 2
ecs_max_capacity = 10

# Optional - Database
rds_multi_az = true  # Production: true, Dev: false
```

## 🚀 Deployment

### Quick Start - Development

```bash
# 1. Build và test locally
make docker-up

# 2. Test API
curl http://localhost:8080/health

# 3. Stop local environment
make docker-down
```

### Full Deployment - AWS

#### Step 1: Deploy Infrastructure

```bash
# Initialize Terraform
make tf-init

# Review changes
make tf-plan ENV=prod

# Deploy infrastructure
make tf-apply ENV=prod

# Save outputs
make tf-output > infrastructure-outputs.txt
```

#### Step 2: Build và Push Docker Image

```bash
# Get ECR repository URL from terraform output
ECR_URL=$(cd terraform && terraform output -raw ecr_repository_url)

# Login to ECR
make aws-login

# Build and push image
make docker-push ENV=prod
```

#### Step 3: Deploy Application

```bash
# Deploy to ECS
make ecs-deploy ENV=prod

# Check deployment status
make ecs-status ENV=prod

# View logs
make ecs-logs ENV=prod
```

#### One-Command Deployment

```bash
# Deploy everything at once
make deploy-all ENV=prod
```

## 🔄 Continuous Deployment

### Update Application Code

```bash
# 1. Make code changes
# 2. Build and push new image
make docker-push ENV=prod VERSION=1.0.1

# 3. Update ECS service
make ecs-deploy ENV=prod
```

### Update Infrastructure

```bash
# 1. Modify terraform files
# 2. Review changes
make tf-plan ENV=prod

# 3. Apply changes
make tf-apply ENV=prod
```

## 📊 Quản lý

### Scaling

```bash
# Scale service manually
make ecs-scale COUNT=5 ENV=prod

# Auto-scaling is configured by default:
# - CPU > 70% → scale up
# - Memory > 80% → scale up
```

### Logs

```bash
# View real-time logs
make ecs-logs ENV=prod

# View logs in AWS Console
# CloudWatch → Log groups → /ecs/user-api-prod
```

### Database Access

```bash
# Get database endpoint
cd terraform && terraform output db_instance_endpoint

# Get database credentials from Secrets Manager
aws secretsmanager get-secret-value \
  --secret-id user-api-prod-db-password \
  --query SecretString \
  --output text | jq .
```

### Execute Commands in Container

```bash
# Access running container
make ecs-exec ENV=prod
```

## 📈 Monitoring

### CloudWatch Dashboards

Truy cập: AWS Console → CloudWatch → Dashboards

**Metrics quan trọng:**
- ECS CPU Utilization
- ECS Memory Utilization
- ALB Request Count
- ALB Target Response Time
- RDS CPU Utilization
- RDS Database Connections

### Alarms

Tạo CloudWatch Alarms cho:
- High CPU (>80%)
- High Memory (>85%)
- Error Rate (>5%)
- Database Connections (>80% max)

### Cost Monitoring

```bash
# Estimate monthly cost
aws ce get-cost-and-usage \
  --time-period Start=2024-01-01,End=2024-01-31 \
  --granularity MONTHLY \
  --metrics BlendedCost \
  --group-by Type=TAG,Key=Project
```

**Estimated Monthly Costs:**
- Development: $30-50
- Production: $100-200
  - NAT Gateway: ~$45
  - RDS t3.small: ~$30
  - ECS Fargate (2 tasks): ~$25
  - ALB: ~$20
  - Data Transfer: ~$10-50

## 🔍 Troubleshooting

### ECS Service không start

```bash
# Check service events
aws ecs describe-services \
  --cluster user-api-prod \
  --services user-api-prod-service \
  --query 'services[0].events[0:5]'

# Check task logs
make ecs-logs ENV=prod
```

### Health Check Failed

```bash
# Check ALB target health
aws elbv2 describe-target-health \
  --target-group-arn $(cd terraform && terraform output -raw target_group_arn)

# Verify health check endpoint
curl http://$(cd terraform && terraform output -raw alb_dns_name)/health
```

### Database Connection Issues

```bash
# Check security groups
# ECS tasks should have access to RDS on port 5432

# Verify connection from task
make ecs-exec ENV=prod
# In container:
# wget -O - https://your-rds-endpoint:5432
```

### High Costs

```bash
# Check NAT Gateway usage
# Consider using single NAT for dev environment

# Check RDS instance size
# Downsize for non-production environments

# Enable VPC Flow Logs only when needed
```

## 🗑️ Cleanup

### Destroy Resources

```bash
# Destroy application
aws ecs update-service \
  --cluster user-api-prod \
  --service user-api-prod-service \
  --desired-count 0

# Destroy infrastructure
make tf-destroy ENV=prod
```

**⚠️ Warning:** Ensure you have backups before destroying production resources!

## 📚 Additional Resources

- [Terraform AWS Provider](https://registry.terraform.io/providers/hashicorp/aws/latest/docs)
- [AWS ECS Best Practices](https://docs.aws.amazon.com/AmazonECS/latest/bestpracticesguide/intro.html)
- [AWS Well-Architected Framework](https://aws.amazon.com/architecture/well-architected/)

## 🆘 Support

For issues or questions:
1. Check CloudWatch Logs
2. Review ECS Service Events
3. Check Security Groups
4. Verify IAM Permissions

## 🔐 Security Best Practices

1. **Never commit credentials** to git
2. Use **AWS Secrets Manager** for sensitive data
3. Enable **VPC Flow Logs** in production
4. Configure **SSL/TLS** with ACM certificates
5. Regular **security updates** for dependencies
6. Enable **container scanning** in ECR
7. Use **IAM roles** instead of access keys
8. Enable **deletion protection** for production RDS
9. Regular **backup testing**
10. **Multi-AZ** deployment for production


