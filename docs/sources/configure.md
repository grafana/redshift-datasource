---
description: Configure the Amazon Redshift data source in Grafana, including authentication, provisioning, and Terraform setup.
keywords:
  - grafana
  - redshift
  - amazon redshift
  - aws
  - data source
  - configure
  - authentication
  - provisioning
  - terraform
  - iam
  - secrets manager
  - serverless
labels:
  products:
    - cloud
    - enterprise
    - oss
menuTitle: Configure
title: Configure the Amazon Redshift data source
weight: 100
review_date: 2026-05-06
---

# Configure the Amazon Redshift data source

This document explains how to configure the Amazon Redshift data source in Grafana, including authentication methods, connection settings, provisioning, and Terraform configuration.

## Before you begin

Before you configure the data source, ensure you have the following:

- **Grafana permissions:** An organization administrator role in Grafana.
- **AWS account:** An AWS account with either a Redshift provisioned cluster or a Redshift Serverless workgroup.
- **IAM permissions:** An IAM identity (user or role) with the required Redshift Data API permissions. Refer to [IAM policies](https://grafana.com/docs/plugins/grafana-redshift-datasource/latest/configure/#iam-policies) for the minimum required permissions.
- **Database credentials:** Either temporary credentials through IAM or a managed secret stored in AWS Secrets Manager.

## Key concepts

If you're new to Amazon Redshift or AWS, the following terms are used throughout this documentation:

| Term                     | Description                                                                                                                                          |
| ------------------------ | ---------------------------------------------------------------------------------------------------------------------------------------------------- |
| **IAM policy**           | A JSON document attached to an AWS identity that grants permissions to perform specific API actions.                                                  |
| **Assume role**          | An AWS mechanism that lets one identity take on temporary credentials for another role, often used for cross-account access.                          |
| **Redshift Serverless**  | A deployment option for Amazon Redshift that automatically provisions and scales capacity without requiring you to manage clusters.                   |
| **Workgroup**            | A named group of compute resources in Redshift Serverless that processes queries.                                                                    |
| **Cluster identifier**   | The unique name assigned to a Redshift provisioned cluster.                                                                                          |
| **AWS Secrets Manager**  | An AWS service that stores and manages database credentials, API keys, and other secrets.                                                            |
| **Amazon EventBridge**   | An AWS event bus service. The Redshift data source can optionally send Data API events to EventBridge for monitoring.                                 |

## Add the data source

To add the Amazon Redshift data source to Grafana:

1. Click **Connections** in the left-side menu.
1. Click **Add new connection**.
1. Type **Amazon Redshift** in the search bar.
1. Select **Amazon Redshift**.
1. Click **Add new data source**.

## Configure connection settings

The connection section configures how Grafana authenticates with AWS. For detailed information about AWS authentication options, refer to [AWS authentication](https://grafana.com/docs/grafana/<GRAFANA_VERSION>/datasources/aws-cloudwatch/aws-authentication/).

| Setting                       | Description                                                                                                                |
| ----------------------------- | -------------------------------------------------------------------------------------------------------------------------- |
| **Authentication Provider**   | The method Grafana uses to authenticate with AWS. Options: **AWS SDK Default**, **Credentials file**, **Access & secret key**. |
| **Access Key ID**             | The AWS access key ID. Required when using **Access & secret key** authentication.                                          |
| **Secret Access Key**         | The AWS secret access key. Required when using **Access & secret key** authentication.                                      |
| **Credentials Profile Name**  | The name of the credentials profile to use from the `~/.aws/credentials` file. Leave blank for the default profile.         |
| **Assume Role ARN**           | The ARN of an IAM role to assume. Use this for cross-account access.                                                        |
| **External ID**               | An external ID to use when assuming a role in another account, if the role requires one.                                    |
| **Session Token**             | A temporary session token for AWS authentication. Required when using temporary AWS credentials from STS.                   |
| **Endpoint**                  | A custom endpoint URL for the AWS service. Leave blank to use the default endpoint.                                         |
| **Default Region**            | The AWS region where your Redshift cluster or workgroup is deployed.                                                        |

## Configure Redshift details

The Redshift details section configures how Grafana connects to your specific Redshift environment. These settings appear under the **Redshift Details** heading in the configuration page.

### Database authentication

Choose how Grafana authenticates with your Redshift database. There are two options:

- **Temporary credentials** -- Generates temporary database credentials using the AWS `GetClusterCredentials` API (for provisioned clusters) or `GetCredentials` API (for Redshift Serverless). This is the default.
- **AWS Secrets Manager** -- Uses database credentials stored in AWS Secrets Manager. Secrets must be tagged with the `RedshiftQueryOwner` tag to appear in the drop-down. Refer to [Storing database credentials in AWS Secrets Manager](https://docs.aws.amazon.com/redshift/latest/mgmt/data-api-access.html#data-api-secrets) for instructions on creating compatible secrets.

{{< admonition type="note" >}}
When using AWS Secrets Manager with Redshift Serverless, the workgroup name isn't stored in the secret. You must configure the **Workgroup** field separately.
{{< /admonition >}}

### Redshift details settings

| Setting                                 | Description                                                                                                                                    |
| --------------------------------------- | ---------------------------------------------------------------------------------------------------------------------------------------------- |
| **Serverless**                          | Toggle to use a Redshift Serverless workgroup instead of a provisioned cluster.                                                                |
| **Cluster Identifier**                  | The Redshift provisioned cluster to connect to. Hidden when **Serverless** is enabled. Automatically set when using AWS Secrets Manager.       |
| **Workgroup**                           | The Redshift Serverless workgroup to connect to. Only visible when **Serverless** is enabled.                                                  |
| **Managed Secret**                      | The AWS Secrets Manager secret containing database credentials. Only visible when using **AWS Secrets Manager** authentication.                 |
| **Database User**                       | The database user for authentication. Automatically set when using AWS Secrets Manager. Hidden for Serverless with temporary credentials.       |
| **Database**                            | The name of the database to connect to within the cluster or workgroup.                                                                        |
| **Send events to Amazon EventBridge**   | Toggle to send Redshift Data API events to Amazon EventBridge for monitoring purposes.                                                         |

## Private data source connect (Grafana Cloud only)

If you use Grafana Cloud and your Redshift cluster is in a private VPC that isn't directly accessible, you can use [Private data source connect (PDC)](https://grafana.com/docs/grafana-cloud/connect-externally-hosted/private-data-source-connect/) to securely connect Grafana Cloud to your data source without opening your network to inbound traffic.

When PDC is enabled for your Grafana Cloud instance, a **Private data source connect** section appears below the Redshift Details configuration. Select the **Private data source connect network** where your Redshift cluster is available. Click **Manage private data source connect networks** to configure networks.

{{< admonition type="note" >}}
Private data source connect is available exclusively in Grafana Cloud. For self-managed Grafana instances, use VPC peering, a VPN, or AWS PrivateLink to connect to private Redshift clusters.
{{< /admonition >}}

For setup instructions, refer to [Private data source connect](https://grafana.com/docs/grafana-cloud/connect-externally-hosted/private-data-source-connect/).

## IAM policies

Grafana needs IAM permissions to query Redshift through the Data API. You can attach the required permissions to an IAM user or role. Refer to [Creating IAM policies](https://docs.aws.amazon.com/IAM/latest/UserGuide/access_policies_create.html) for instructions. AWS also provides [predefined managed policies for Redshift](https://docs.aws.amazon.com/redshift/latest/mgmt/redshift-iam-access-control-identity-based.html#redshift-policy-resources.managed-policies).

The following policy provides the minimum permissions required:

```json
{
  "Version": "2012-10-17",
  "Statement": [
    {
      "Sid": "AllowRedshiftDataAPI",
      "Effect": "Allow",
      "Action": [
        "redshift-data:ListTables",
        "redshift-data:DescribeTable",
        "redshift-data:GetStatementResult",
        "redshift-data:DescribeStatement",
        "redshift-data:ListStatements",
        "redshift-data:ExecuteStatement",
        "redshift-data:CancelStatement",
        "redshift-data:ListSchemas",
        "redshift:GetClusterCredentials",
        "redshift:DescribeClusters",
        "redshift-serverless:ListWorkgroups",
        "redshift-serverless:GetCredentials",
        "secretsmanager:ListSecrets"
      ],
      "Resource": "*"
    },
    {
      "Sid": "AllowReadingRedshiftQuerySecrets",
      "Effect": "Allow",
      "Action": ["secretsmanager:GetSecretValue"],
      "Resource": "*",
      "Condition": {
        "Null": {
          "secretsmanager:ResourceTag/RedshiftQueryOwner": "false"
        }
      }
    }
  ]
}
```

The second statement restricts `GetSecretValue` access to only secrets that have the `RedshiftQueryOwner` tag, which is required for Secrets Manager integration.

{{< admonition type="note" >}}
For async query support, your IAM policy must include the `redshift-data:ListStatements` and `redshift-data:CancelStatement` actions.
{{< /admonition >}}

## Verify the connection

Click **Save & test** to verify that the data source is configured correctly. Grafana runs a `SELECT 1` query against your Redshift environment to confirm the connection is working. A success message indicates the authentication, network connectivity, and database access are all configured correctly.

## Provision the data source

You can define and configure the Redshift data source using YAML files as part of Grafana's provisioning system. For more information about provisioning, refer to [Provision Grafana](https://grafana.com/docs/grafana/<GRAFANA_VERSION>/administration/provisioning/#data-sources).

The following examples show common provisioning configurations.

### AWS SDK default

```yaml
apiVersion: 1

datasources:
  - name: Redshift
    type: grafana-redshift-datasource
    jsonData:
      authType: default
      defaultRegion: us-east-2
      clusterIdentifier: my-redshift-cluster
      database: mydb
      dbUser: admin
```

### Credentials profile

```yaml
apiVersion: 1

datasources:
  - name: Redshift
    type: grafana-redshift-datasource
    jsonData:
      authType: credentials
      defaultRegion: us-east-2
      profile: secondary
      clusterIdentifier: my-redshift-cluster
      database: mydb
      dbUser: admin
```

### Access & secret key

```yaml
apiVersion: 1

datasources:
  - name: Redshift
    type: grafana-redshift-datasource
    jsonData:
      authType: keys
      defaultRegion: us-east-2
      clusterIdentifier: my-redshift-cluster
      database: mydb
      dbUser: admin
    secureJsonData:
      accessKey: <YOUR_ACCESS_KEY>
      secretKey: <YOUR_SECRET_KEY>
```

To use temporary AWS credentials from STS, add a `sessionToken` to `secureJsonData`:

```yaml
    secureJsonData:
      accessKey: <YOUR_ACCESS_KEY>
      secretKey: <YOUR_SECRET_KEY>
      sessionToken: <YOUR_SESSION_TOKEN>
```

### Assume role

```yaml
apiVersion: 1

datasources:
  - name: Redshift
    type: grafana-redshift-datasource
    jsonData:
      authType: default
      defaultRegion: us-east-2
      assumeRoleArn: arn:aws:iam::123456789012:role/redshift-access
      clusterIdentifier: my-redshift-cluster
      database: mydb
      dbUser: admin
```

### Redshift Serverless

```yaml
apiVersion: 1

datasources:
  - name: Redshift Serverless
    type: grafana-redshift-datasource
    jsonData:
      authType: default
      defaultRegion: us-east-2
      useServerless: true
      workgroupName: my-workgroup
      database: mydb
```

### Private data source connect

```yaml
apiVersion: 1

datasources:
  - name: Redshift
    type: grafana-redshift-datasource
    jsonData:
      authType: default
      defaultRegion: us-east-2
      clusterIdentifier: my-redshift-cluster
      database: mydb
      dbUser: admin
      enableSecureSocksProxy: true
```

## Provision with Terraform

You can use the [Grafana Terraform provider](https://registry.terraform.io/providers/grafana/grafana/latest/docs) to provision the Redshift data source as code. The following examples demonstrate common configurations using the `grafana_data_source` resource.

### Access & secret key with Terraform

```hcl
resource "grafana_data_source" "redshift" {
  type = "grafana-redshift-datasource"
  name = "Amazon Redshift"

  json_data_encoded = jsonencode({
    authType          = "keys"
    defaultRegion     = "us-east-2"
    clusterIdentifier = "my-redshift-cluster"
    database          = "mydb"
    dbUser            = "admin"
  })

  secure_json_data_encoded = jsonencode({
    accessKey = var.aws_access_key
    secretKey = var.aws_secret_key
  })
}
```

To use temporary AWS credentials from STS, add `sessionToken` to `secure_json_data_encoded`:

```hcl
  secure_json_data_encoded = jsonencode({
    accessKey    = var.aws_access_key
    secretKey    = var.aws_secret_key
    sessionToken = var.aws_session_token
  })
```

### Assume role with Terraform

```hcl
resource "grafana_data_source" "redshift" {
  type = "grafana-redshift-datasource"
  name = "Amazon Redshift"

  json_data_encoded = jsonencode({
    authType          = "default"
    defaultRegion     = "us-east-2"
    assumeRoleArn     = "arn:aws:iam::123456789012:role/redshift-access"
    clusterIdentifier = "my-redshift-cluster"
    database          = "mydb"
    dbUser            = "admin"
  })
}
```

### Redshift Serverless with Terraform

```hcl
resource "grafana_data_source" "redshift_serverless" {
  type = "grafana-redshift-datasource"
  name = "Redshift Serverless"

  json_data_encoded = jsonencode({
    authType      = "default"
    defaultRegion = "us-east-2"
    useServerless = true
    workgroupName = "my-workgroup"
    database      = "mydb"
  })
}
```

For more information, refer to the [`grafana_data_source` resource](https://registry.terraform.io/providers/grafana/grafana/latest/docs/resources/data_source) in the Grafana Terraform provider documentation.
