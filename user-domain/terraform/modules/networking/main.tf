locals {
  name = "${var.project_name}-${var.environment}"
  public_subnet_cidr = [for i, az in var.availability_zones :cidrsubnet(var.vpc_cidr,8,i)]
  private_subnet_cidr = [for i, az in var.availability_zones : cidrsubnet(var.vpc_cidr,8, i + 10)]
  private_database_cidr = [for i, az in var.availability_zones : cidrsubnet(var.vpc_cidr, 8, i +20)]
}

resource "aws_vpc" "main" {
  cidr_block = var.vpc_cidr
  enable_dns_hostnames = true
  enable_dns_support = true
  tags = merge(var.tags, {
    Name="${local.name}-vcp"
  })
}

resource "aws_internet_gateway" "main" {
  vpc_id = aws_vpc.main.id
  tags = merge(var.tags, {
    Name="${local.name}-igw"
  })
}

# for load balancer, NAT getaway
resource "aws_subnet" "public" {
  count = length(var.availability_zones)
  vpc_id = aws_vpc.main.id
  cidr_block = local.public_subnet_cidr[count.index]
  availability_zone = var.availability_zones[count.index]

   tags = merge(var.tags, {
    Name = "${local.name}-public-${var.availability_zones[count.index]}"
    Type = "public"
  })
}
# for ECS
resource "aws_subnet" "private" {
  count = length(var.availability_zones)
  cidr_block = local.private_subnet_cidr[count.index]
  availability_zone = var.availability_zones[count.index]

  vpc_id = aws_vpc.main.id
   tags = merge(var.tags, {
    Name = "${local.name}-private-${var.availability_zones[count.index]}"
    Type = "private"
  })
}

# for database 
resource "aws_subnet" "database" {
  vpc_id = aws_vpc.main.id
  count = length(var.availability_zones)
  cidr_block = local.private_database_cidr[count.index]
  availability_zone = var.availability_zones[count.index]

  tags = merge(var.tags, {
    Name = "${local.name}-database-${var.availability_zones[count.index]}"
    Type = "database"
  })
}

resource "aws_db_subnet_group" "main" {
  name = "${local.name}-db-subnet-group"
  subnet_ids = aws_subnet.database[*].id

  tags = merge(var.tags, {
    Name = "${local.name}-db-subnet-group"
  })
}

resource "aws_eip" "nat" {
  count = length(var.availability_zones)
  domain = "vpc"

   tags = merge(var.tags, {
    Name = "${local.name}-nat-eip-${count.index + 1}"
  })
  depends_on = [ aws_internet_gateway.main ]
}

resource "aws_nat_gateway" "main" {
  count = length(var.availability_zones)
  subnet_id = local.public_subnet_cidr[count.index]

   tags = merge(var.tags, {
    Name = "${local.name}-nat-${count.index + 1}"
  })

  depends_on = [ aws_internet_gateway.main ]
}

resource "aws_route_table" "public" {
  vpc_id = aws_vpc.main.id
  route {
    cidr_block = "0.0.0.0/0"
    gateway_id = aws_internet_gateway.main.id
  }

  tags = merge(var.tags, {
    Name = "${local.name}-public-rt"
  })
}

resource "aws_route_table" "private" {
  count = length(var.availability_zones)
  vpc_id = aws_vpc.main.id
  route {
    cidr_block = "0.0.0.0/0"
    nat_gateway_id = aws_nat_gateway[count.index].main.id
  }

    tags = merge(var.tags, {
    Name = "${local.name}-private-rt-${count.index + 1}"
  })
}

resource "aws_route_table" "database" {
  vpc_id = aws_vpc.main.id
  
   tags = merge(var.tags, {
    Name = "${local.name}-database-rt"
  })
}

resource "aws_route_table_association" "public" {
  count = length(var.availability_zones)
  route_table_id = aws_route_table.public.id
  subnet_id = aws_subnet.public[count.index].id
}

resource "aws_route_table_association" "private" {
  count = length(var.availability_zones)
  route_table_id = aws_route_table.private[i].private.id
  subnet_id = aws_subnet.private[count.index].id
}

resource "aws_route_table_association" "database" {
  count = length(var.availability_zones)
  route_table_id = aws_route_table.database.id
  subnet_id = aws_subnet.database[count.index].id
}


# VPC Flow Logs (optional)
resource "aws_flow_log" "main" {
  count = var.enable_flow_logs ? 1 : 0

  iam_role_arn    = aws_iam_role.flow_logs[0].arn
  log_destination = aws_cloudwatch_log_group.flow_logs[0].arn
  traffic_type    = "ALL"
  vpc_id          = aws_vpc.main.id

  tags = var.tags
}

resource "aws_cloudwatch_log_group" "flow_logs" {
  count = var.enable_flow_logs ? 1 : 0

  name              = "/aws/vpc/${local.name}"
  retention_in_days = 7

  tags = var.tags
}

resource "aws_iam_role" "flow_logs" {
  count = var.enable_flow_logs ? 1 : 0

  name = "${local.name}-flow-logs-role"

  assume_role_policy = jsonencode({
    Version = "2012-10-17"
    Statement = [
      {
        Action = "sts:AssumeRole"
        Effect = "Allow"
        Principal = {
          Service = "vpc-flow-logs.amazonaws.com"
        }
      }
    ]
  })

  tags = var.tags
}

resource "aws_iam_role_policy" "flow_logs" {
  count = var.enable_flow_logs ? 1 : 0

  name = "${local.name}-flow-logs-policy"
  role = aws_iam_role.flow_logs[0].id

  policy = jsonencode({
    Version = "2012-10-17"
    Statement = [
      {
        Action = [
          "logs:CreateLogGroup",
          "logs:CreateLogStream",
          "logs:PutLogEvents",
          "logs:DescribeLogGroups",
          "logs:DescribeLogStreams"
        ]
        Effect   = "Allow"
        Resource = "*"
      }
    ]
  })
}
