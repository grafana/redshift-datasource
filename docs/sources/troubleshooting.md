---
description: Troubleshoot common issues with the Amazon Redshift data source in Grafana.
keywords:
  - grafana
  - redshift
  - amazon redshift
  - troubleshooting
  - errors
  - authentication
  - query
  - connection
labels:
  products:
    - cloud
    - enterprise
    - oss
menuTitle: Troubleshooting
title: Troubleshoot Amazon Redshift data source issues
weight: 500
review_date: 2026-05-06
---

# Troubleshoot Amazon Redshift data source issues

This page provides solutions to common issues you might encounter when configuring or using the Amazon Redshift data source. For configuration instructions, refer to [Configure the Amazon Redshift data source](https://grafana.com/docs/plugins/grafana-redshift-datasource/latest/configure/).

## Authentication errors

These errors occur when AWS credentials are invalid, missing, or don't have the required permissions.

### "Access denied" or "Authorization failed"

**Symptoms:**

- **Save & test** fails with authorization errors.
- Queries return access denied messages.
- Resource drop-downs (clusters, workgroups, secrets) don't load.

**Possible causes and solutions:**

| Cause                    | Solution                                                                                                                                                              |
| ------------------------ | --------------------------------------------------------------------------------------------------------------------------------------------------------------------- |
| Missing IAM permissions  | Verify that your IAM identity has all the permissions listed in the [IAM policies](https://grafana.com/docs/plugins/grafana-redshift-datasource/latest/configure/#iam-policies) section. |
| Invalid credentials      | Verify credentials in the AWS console. Regenerate the access key and secret key if necessary.                                                                          |
| Expired credentials      | Create new credentials and update the data source configuration.                                                                                                       |
| Wrong region             | Verify that the **Default Region** setting matches the region where your Redshift cluster or workgroup is deployed.                                                    |
| Wrong authentication provider | Ensure the selected **Authentication Provider** matches how your credentials are configured (SDK default, credentials file, or access and secret key).             |

### "GetClusterCredentials" or "GetCredentials" errors

**Symptoms:**

- **Save & test** fails when using temporary credentials.
- Error messages reference `GetClusterCredentials` or `GetCredentials`.

**Solutions:**

1. Verify that your IAM policy includes `redshift:GetClusterCredentials` (for provisioned clusters) or `redshift-serverless:GetCredentials` (for Serverless workgroups).
1. Confirm that the **Database User** field is set correctly for provisioned clusters.
1. For Serverless, ensure the **Workgroup** field is configured and the workgroup exists in the selected region.

## Connection errors

These errors occur when Grafana can't reach your Redshift environment.

### "Connection refused" or timeout errors

**Symptoms:**

- Data source test times out.
- Queries fail with network errors.
- Intermittent connection issues.

**Solutions:**

1. Verify network connectivity from the Grafana server to the Redshift Data API endpoints.
1. Check that firewall rules allow outbound HTTPS (port 443) to `redshift-data.<region>.amazonaws.com`.
1. If your Redshift cluster is in a private VPC, ensure the Grafana server has network access through VPC peering, a VPN, or AWS PrivateLink.
1. For Grafana Cloud accessing private resources, configure [Private data source connect](https://grafana.com/docs/grafana-cloud/connect-externally-hosted/private-data-source-connect/).

### Custom endpoint issues

**Symptoms:**

- Connection fails after setting a custom **Endpoint** value.

**Solutions:**

1. Verify the endpoint URL format is correct for the AWS service you're targeting.
1. Ensure the endpoint is reachable from the Grafana server.
1. Remove the custom endpoint to use the default AWS endpoint and verify the connection works without it.

## Secrets Manager errors

These errors relate to the AWS Secrets Manager authentication method.

### Secrets don't appear in the drop-down

**Symptoms:**

- The **Managed Secret** drop-down is empty.
- Expected secrets aren't listed.

**Possible causes and solutions:**

| Cause                             | Solution                                                                                                                                                                                      |
| --------------------------------- | --------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------- |
| Missing `RedshiftQueryOwner` tag  | Only secrets tagged with `RedshiftQueryOwner` appear in the list. Add this tag to your secret in AWS Secrets Manager. Refer to [AWS documentation](https://docs.aws.amazon.com/redshift/latest/mgmt/data-api-access.html#data-api-secrets). |
| Missing IAM permissions           | Verify that your IAM policy includes `secretsmanager:ListSecrets` and `secretsmanager:GetSecretValue` with the appropriate resource conditions.                                                |
| Wrong region                      | Secrets are region-specific. Verify that the **Default Region** matches the region where your secrets are stored.                                                                               |

### Secret resolves but connection fails

**Symptoms:**

- A secret is selected and the **Cluster Identifier** and **Database User** fields populate, but **Save & test** fails.

**Solutions:**

1. Verify that the credentials stored in the secret are still valid.
1. Ensure the secret format matches the expected Redshift format. Refer to [Storing database credentials in AWS Secrets Manager](https://docs.aws.amazon.com/redshift/latest/mgmt/data-api-access.html#data-api-secrets).
1. For Serverless with Secrets Manager, verify that the **Workgroup** field is configured separately, as the workgroup name isn't stored in the secret.

## Query errors

These errors occur when running queries against Redshift.

### "No data" or empty results

**Symptoms:**

- Query executes without error but returns no data.
- Charts show a "No data" message.

**Possible causes and solutions:**

| Cause                          | Solution                                                                                             |
| ------------------------------ | ---------------------------------------------------------------------------------------------------- |
| Time range doesn't contain data | Expand the dashboard time range or verify that data exists for the selected period in Redshift.       |
| Wrong schema, table, or column | Verify the resource selectors or SQL references match actual objects in your Redshift database.       |
| Permissions issue              | Verify that the database user has `SELECT` permissions on the referenced tables.                     |
| Macro expansion issue          | Use the [Query Inspector](https://grafana.com/docs/plugins/grafana-redshift-datasource/latest/query-editor/#inspect-the-query) to check the fully rendered SQL query. |

### Query timeout

**Symptoms:**

- Query runs for a long time and then fails.
- Error message references a timeout.

**Solutions:**

1. Narrow the dashboard time range to reduce the amount of data scanned.
1. Add `WHERE` clauses or filters to reduce the result set.
1. Use `LIMIT` to cap the number of returned rows.
1. Optimize your query by adding appropriate sort keys or distribution keys in Redshift.
1. Break complex queries into smaller, focused queries.

### Syntax errors

**Symptoms:**

- Query fails immediately with a syntax error.

**Solutions:**

1. Use the [Query Inspector](https://grafana.com/docs/plugins/grafana-redshift-datasource/latest/query-editor/#inspect-the-query) to view the fully interpolated SQL and check for macro expansion issues.
1. Copy the rendered query and run it directly in a Redshift SQL client to isolate the problem.
1. Verify that macro syntax is correct (for example, `$__timeFilter(column)` requires exactly one argument).

## Async query errors

These errors are specific to the asynchronous query execution model.

### "ListStatements" or "CancelStatement" errors

**Symptoms:**

- Queries fail with IAM errors referencing `ListStatements` or `CancelStatement`.
- The stop button in the query editor doesn't work.

**Solutions:**

1. Add `redshift-data:ListStatements` and `redshift-data:CancelStatement` to your IAM policy. These are required for async query support.
1. Verify that the IAM policy resource scope includes your Redshift cluster or workgroup.

## Template variable errors

These errors occur when using template variables with the data source.

### Variables return no values

**Solutions:**

1. Verify the data source connection is working by running **Save & test** in the data source settings.
1. Check that parent variables (for cascading variables) have valid selections.
1. Open the variable's query in the variable editor and click **Run query** to test it independently.
1. Verify the IAM identity has permissions to query the tables referenced in the variable query.

### Variables are slow to load

**Solutions:**

1. Set the variable refresh to **On dashboard load** instead of **On time range change** to reduce how often the query runs.
1. Add `LIMIT` clauses to variable queries to reduce result set sizes.
1. Narrow the scope of variable queries to specific schemas or tables.

## Debug logs

To capture detailed error information for troubleshooting:

1. Set the Grafana log level to `debug` in your Grafana configuration file:

   ```ini
   [log]
   level = debug
   ```

1. Review logs in `/var/log/grafana/grafana.log` (or your configured log location).
1. Look for entries containing `redshift` for request and response details.
1. Reset the log level to `info` after troubleshooting to avoid excessive log volume.

## Get additional help

If you've tried the solutions on this page and still encounter issues:

1. Check the [Grafana community forums](https://community.grafana.com/) for similar issues and solutions.
1. Review the [Redshift data source GitHub issues](https://github.com/grafana/redshift-datasource/issues) for known bugs and feature requests.
1. Consult the [Amazon Redshift documentation](https://docs.aws.amazon.com/redshift/) for service-specific guidance.
1. Contact [Grafana Support](https://grafana.com/support/) if you have a support contract.
1. When reporting issues, include:
   - Grafana version and plugin version.
   - Error messages (redact sensitive information such as credentials and ARNs).
   - Steps to reproduce the issue.
   - Relevant configuration details (redact credentials).
