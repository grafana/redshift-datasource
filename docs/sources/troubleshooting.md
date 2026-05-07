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

This document provides solutions to common issues you might encounter when installing, configuring, or using the Amazon Redshift data source. Sections are organized in the order you're likely to encounter issues, from installation through querying and alerting.

For configuration instructions, refer to [Configure the Amazon Redshift data source](https://grafana.com/docs/plugins/grafana-redshift-datasource/latest/configure/).

## Plugin installation errors

These errors occur when installing or updating the Redshift data source plugin.

### Plugin archive won't open or appears corrupt

**Symptoms:**

- The downloaded `.zip` file won't extract.
- The archive appears invalid or corrupt.

**Solutions:**

1. Verify that you downloaded the correct archive for your operating system and architecture. The plugin is distributed as platform-specific archives (for example, `grafana-redshift-datasource-<version>.linux_amd64.zip`). A Linux archive won't work on Windows or macOS.
1. Download the plugin directly from the [Grafana plugin catalog](https://grafana.com/grafana/plugins/grafana-redshift-datasource/) to ensure you have the correct file.
1. For self-managed Grafana, use the Grafana CLI to install the plugin instead of downloading manually:

   ```bash
   grafana cli plugins install grafana-redshift-datasource
   ```

### "Unsupported auth type" error

**Symptoms:**

- **Save & test** fails with `unsupported auth type default` or similar errors.
- Authentication options that should be available don't appear in the configuration UI.

**Solutions:**

1. Update the plugin to the latest version. Older plugin versions may not support all authentication methods (such as **AWS SDK Default**).
1. For self-managed Grafana, update using the CLI:

   ```bash
   grafana cli plugins update grafana-redshift-datasource
   ```

1. Restart Grafana after updating the plugin.

### "Save & test" fails with a 500 error

**Symptoms:**

- **Save & test** returns a 500 Internal Server Error.
- The plugin was recently installed or the Grafana instance was recently upgraded.

**Solutions:**

1. Update the plugin to the latest version. Older versions may be missing required features (such as the **External ID** field for Assume Role) and fail silently.
1. Verify plugin compatibility with your Grafana version. The Redshift data source requires **Grafana 10.4.0 or later**.
1. Check the Grafana server logs for detailed error messages:

   ```bash
   grep "redshift" /var/log/grafana/grafana.log
   ```

1. If the error persists after updating, remove and reinstall the plugin:

   ```bash
   grafana cli plugins remove grafana-redshift-datasource
   grafana cli plugins install grafana-redshift-datasource
   ```

1. Restart Grafana after reinstalling.

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

### Assume Role and STS errors

These errors occur when using the **Assume Role ARN** option to authenticate with a cross-account or delegated IAM role.

#### "User is not authorized to perform sts:AssumeRole"

**Symptoms:**

- **Save & test** fails with `AccessDenied: User is not authorized to perform sts:AssumeRole`.
- The error references the role ARN you configured.

**Solutions:**

1. Verify that the IAM role's trust policy allows the Grafana identity to assume the role. For Grafana Cloud, the trust policy must include the Grafana Cloud AWS account ID and your **External ID**:

   ```json
   {
     "Version": "2012-10-17",
     "Statement": [
       {
         "Effect": "Allow",
         "Principal": {
           "AWS": "arn:aws:iam::008923505280:root"
         },
         "Action": "sts:AssumeRole",
         "Condition": {
           "StringEquals": {
             "sts:ExternalId": "<YOUR_EXTERNAL_ID>"
           }
         }
       }
     ]
   }
   ```

1. Confirm the **External ID** in the data source configuration matches the value in the trust policy exactly.
1. Verify that the **Assume Role ARN** is correct and the role exists in the target account.
1. For self-managed Grafana, ensure the Grafana server's IAM identity has `sts:AssumeRole` permission for the target role.

#### "InvalidClientTokenId: The security token included in the request is invalid"

**Symptoms:**

- **Save & test** or queries fail with `InvalidClientTokenId`.
- The error occurs intermittently or suddenly without configuration changes.

**Solutions:**

1. If you haven't changed any credentials or configuration, this error can be a transient issue. Wait a few minutes and try again.
1. Verify that your access key and secret key are valid in the AWS console. Regenerate them if necessary.
1. If using Assume Role, confirm the trust policy and External ID are configured correctly.
1. For Grafana Cloud, if the issue persists after verifying credentials, contact [Grafana Support](https://grafana.com/support/) as it may require a service-side restart.

#### Intermittent IAM token expiry

**Symptoms:**

- Queries work initially but start failing with authentication errors after a period of time.
- Restarting the data source or re-saving the configuration temporarily resolves the issue.

**Solutions:**

1. This can occur when cached IAM tokens expire and aren't refreshed properly. Re-save the data source configuration by clicking **Save & test** to force a credential refresh.
1. If the issue recurs frequently, consider switching from Assume Role to direct **Access & secret key** authentication as a workaround.
1. Contact [Grafana Support](https://grafana.com/support/) if you have a support contract, as this may indicate a connection pool issue requiring engineering attention.

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

### Connect to Redshift in a VPC

If your Redshift cluster is in a private VPC without public internet access, the Grafana server or PDC agent must be able to reach the Redshift Data API through a private network path. The Redshift data source communicates with the Data API (not the cluster directly), so network access to `redshift-data.<region>.amazonaws.com` is required.

#### Required VPC endpoints

For self-managed Grafana in a private VPC, create [VPC interface endpoints](https://docs.aws.amazon.com/vpc/latest/privatelink/create-interface-endpoint.html) for the following services:

| Service                                      | Endpoint                                                   | Required for                          |
| -------------------------------------------- | ---------------------------------------------------------- | ------------------------------------- |
| **Redshift Data API**                        | `com.amazonaws.<region>.redshift-data`                     | All queries and resource selectors    |
| **Redshift**                                 | `com.amazonaws.<region>.redshift`                          | Listing clusters, temporary credentials |
| **Redshift Serverless**                      | `com.amazonaws.<region>.redshift-serverless`               | Serverless workgroup operations       |
| **Secrets Manager**                          | `com.amazonaws.<region>.secretsmanager`                    | Secrets Manager authentication        |
| **STS**                                      | `com.amazonaws.<region>.sts`                               | Assume Role authentication            |

{{< admonition type="note" >}}
You only need the endpoints for the services you use. For example, if you don't use Assume Role, you don't need the STS endpoint.
{{< /admonition >}}

#### VPC endpoint policies block API calls

**Symptoms:**

- Queries fail with access denied errors despite correct IAM permissions.
- Resource selectors (clusters, workgroups, secrets) don't load.
- Errors reference specific API actions being denied.

**Solutions:**

1. Check the VPC endpoint policy for your Redshift Data API endpoint. Restrictive policies may block required API actions like `redshift-data:ExecuteStatement` or `redshift-data:ListSchemas`.
1. Update the endpoint policy to allow all actions needed by the data source, or use a permissive policy during initial testing:

   ```json
   {
     "Statement": [
       {
         "Effect": "Allow",
         "Principal": "*",
         "Action": "*",
         "Resource": "*"
       }
     ]
   }
   ```

1. After confirming connectivity, restrict the policy to only the actions listed in the [IAM policies](https://grafana.com/docs/plugins/grafana-redshift-datasource/latest/configure/#iam-policies) section.

#### FedRAMP and restricted environments

In FedRAMP or other restricted AWS environments, additional constraints may apply:

1. Verify that the Grafana Cloud AWS account is allowlisted in your organization's service control policies (SCPs).
1. Ensure VPC endpoint policies permit traffic from the Grafana identity.
1. Confirm that the AWS region you're using is enabled in your account. Opt-in regions (such as `me-central-1`) require explicit enablement and may fail STS requests if not activated.
1. For Grafana Cloud in restricted environments, contact [Grafana Support](https://grafana.com/support/) for guidance on network configuration.

### Troubleshoot PDC connections with Redshift

These issues are specific to [Private data source connect (PDC)](https://grafana.com/docs/grafana-cloud/connect-externally-hosted/private-data-source-connect/) deployments where Grafana Cloud connects to a Redshift cluster in a private VPC. The PDC agent requires outbound access on both **port 443** (HTTPS to the Redshift Data API) and **port 22** (SSH to Grafana Cloud endpoints).

#### "Host unreachable" through PDC

**Symptoms:**

- **Save & test** fails with `socks connect tcp ... host unreachable`.
- The PDC agent appears connected in the UI, but queries fail.

**Solutions:**

1. Verify that the PDC agent host can reach the Redshift Data API endpoint. Run the following from the agent host:

   ```bash
   curl -v https://redshift-data.<region>.amazonaws.com
   ```

1. Check that firewall and security group rules allow **outbound** traffic on port 443 to `redshift-data.<region>.amazonaws.com`.
1. Verify that firewall rules also allow **outbound** traffic on port 22 to Grafana Cloud endpoints. PDC uses this port for the tunnel connection.
1. If the agent is in a private subnet, ensure a NAT gateway or VPC endpoint is configured for outbound internet access.

#### PDC agent connected but data source fails

**Symptoms:**

- The **Private data source connect network** drop-down shows the agent as connected (for example, "1 agent connected").
- **Save & test** still fails with connection errors.

**Solutions:**

1. Check that the firewall allows egress on **both** port 22 and port 443. A common misconfiguration is blocking port 22 while allowing 443, which lets the agent register but prevents data from flowing.
1. Verify that the agent is deployed with the correct manifest for your Grafana Cloud stack.
1. Review the agent logs for connection errors or certificate issues.

#### PDC certificate or credential errors

**Symptoms:**

- Agent logs show `failed to generate new certificate: key signing request failed: invalid credentials`.
- The agent fails to start or repeatedly restarts.

**Solutions:**

1. Verify that the agent was deployed using the correct installation method. If using a Docker-based deployment, try switching to the shell-based (binary) installation to rule out container-specific credential issues.
1. Re-download the agent manifest from **Grafana Cloud** > **Private data source connect** to get fresh credentials.
1. Ensure the agent host clock is synchronized (NTP). Certificate validation fails if the system clock is significantly out of sync.

#### Intermittent PDC connection drops

**Symptoms:**

- Queries fail randomly with `db query error — failed to connect to server`.
- Some panels load while others time out on the same dashboard.
- Errors resolve temporarily after refreshing.

**Solutions:**

1. If running the PDC agent in Docker, try switching to the binary version. Docker-based deployments can experience intermittent connectivity issues due to container networking.
1. Verify that the agent host has stable network connectivity and sufficient resources (CPU, memory).
1. Check for network-level connection limits or timeouts on firewalls or proxies between the agent and Grafana Cloud.

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

### Syntax errors

**Symptoms:**

- Query fails immediately with a syntax error.

**Solutions:**

1. Use the [Query Inspector](https://grafana.com/docs/plugins/grafana-redshift-datasource/latest/query-editor/#inspect-the-query) to view the fully interpolated SQL and check for macro expansion issues.
1. Copy the rendered query and run it directly in a Redshift SQL client to isolate the problem.
1. Verify that macro syntax is correct (for example, `$__timeFilter(column)` requires exactly one argument).

### Query performance and timeouts

Query timeouts are one of the most common issues with the Redshift data source. Timeouts can occur at multiple levels, and understanding the timeout hierarchy helps identify the cause.

#### Timeout hierarchy

Grafana enforces timeouts at several layers. A query fails when it exceeds any of these limits:

| Level                        | Default         | How to adjust                                                                                             |
| ---------------------------- | --------------- | --------------------------------------------------------------------------------------------------------- |
| **Data source query timeout**| 60 seconds      | Not directly configurable in the plugin. Optimize queries to complete within this window.                  |
| **Grafana instance timeout** | Varies          | Set `dataproxy.timeout` in the Grafana configuration file (self-managed only).                            |
| **Grafana Cloud global limit**| 5 minutes      | Not configurable. This is the maximum query duration for Grafana Cloud.                                    |

#### "Request did not complete within the allotted timeout"

**Symptoms:**

- Query fails with `Timeout: request did not complete within the allotted timeout`.
- Queries consistently fail at exactly 60 seconds or 5 minutes.

**Solutions:**

1. Optimize the query to reduce execution time:
   - Narrow the dashboard time range to scan less data.
   - Add `WHERE` clauses or filters to reduce the result set.
   - Use `LIMIT` to cap the number of returned rows.
   - Add appropriate sort keys and distribution keys in Redshift.
   - Break complex queries into smaller, focused queries.
1. Use [recording rules](https://grafana.com/docs/grafana/<GRAFANA_VERSION>/administration/recorded-queries/) to pre-aggregate expensive queries and serve results from a faster cache.
1. If queries consistently fail at exactly 60 seconds, check whether any intermediate network proxies (for example, Zscaler or corporate firewalls) impose their own timeout limits.

#### "Context deadline exceeded" on alert queries

**Symptoms:**

- Alert queries fail with `context deadline exceeded`.
- The same query works in the dashboard panel but fails during alert evaluation.

**Solutions:**

1. Switch the alert query format from **Time Series** to **Table**. The Time Series format requires additional data conversion that increases processing time and can push queries past the timeout.
1. Remove time-related columns that aren't needed for the alert condition. Alert queries only need to return a numeric value for threshold evaluation.
1. Simplify the query to return only the metric the alert condition evaluates.

#### Concurrent query failures

**Symptoms:**

- Multiple panels on a dashboard fail simultaneously.
- Errors occur when the dashboard auto-refreshes.
- Some panels succeed while others fail inconsistently.

**Solutions:**

1. Increase the dashboard auto-refresh interval to reduce the number of simultaneous queries sent to Redshift.
1. Reduce the number of panels that query Redshift on a single dashboard, or split the dashboard.
1. Check whether a network proxy (for example, Zscaler) is throttling or dropping concurrent connections.
1. Review your Redshift WLM (Workload Management) configuration for concurrent query limits that may cause queuing or rejections.

### Async query errors

These errors are specific to the asynchronous query execution model.

#### "ListStatements" or "CancelStatement" errors

**Symptoms:**

- Queries fail with IAM errors referencing `ListStatements` or `CancelStatement`.
- The stop button in the query editor doesn't work.

**Solutions:**

1. Add `redshift-data:ListStatements` and `redshift-data:CancelStatement` to your IAM policy. These are required for async query support.
1. Verify that the IAM policy resource scope includes your Redshift cluster or workgroup.

#### Query enters "Failed" or "Aborted" state

**Symptoms:**

- Query returns an error after running for some time.
- Error message references a failed or aborted statement.

**Solutions:**

1. Check the Redshift query history in the AWS console for detailed error messages. Navigate to **Amazon Redshift** > **Query monitoring** to view statement details.
1. Verify the SQL is valid by running it directly in a Redshift SQL client.
1. Check for Redshift resource limits such as WLM (Workload Management) queue capacity or concurrent query limits.
1. For aborted queries, check if another process or user canceled the statement.

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

## Unsupported SQL features

The Redshift data source uses the [Amazon Redshift Data API](https://docs.aws.amazon.com/redshift/latest/mgmt/data-api.html) to execute queries. The following SQL features are not supported:

- **Transactions:** The `BEGIN`, `COMMIT`, and `ROLLBACK` statements are not supported. Each query runs as an independent statement.
- **Prepared statements:** Parameterized queries using `PREPARE` and `EXECUTE` are not supported. Use standard SQL with Grafana macros and template variables instead.

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

## Transient platform errors

Some errors originate from the Grafana platform rather than your configuration. These typically affect all AWS data sources simultaneously and resolve without changes on your side.

### Identify a platform-level incident

If you experience sudden, widespread failures across multiple AWS data sources (Redshift, Athena, CloudWatch, and others), the issue is likely platform-related rather than configuration-related.

**Common symptoms:**

- All AWS data sources fail at the same time with `InvalidClientTokenId` or `connection refused` errors.
- Errors appear as `plugin.connectionUnavailable` in dashboards.
- No recent changes were made to credentials, IAM policies, or data source configuration.

**What to do:**

1. Check the [Grafana Cloud status page](https://status.grafana.com/) for active incidents.
1. If an incident is in progress, wait for resolution. No configuration changes are needed on your side.
1. Avoid regenerating credentials or modifying IAM policies during a platform incident, as this can complicate recovery.

### Prevent false-positive alerts during outages

Platform incidents can trigger alert rules that query AWS data sources, causing false-positive alert notifications.

To prevent this, set the alert error handling to **Keep Last State** so alerts maintain their current state when queries fail, rather than firing or resolving:

1. Open the alert rule.
1. In the **Configure no data and error handling** section, set **Error handling** to **Keep Last State**.
1. Click **Save rule and exit**.

This prevents transient platform errors from generating unnecessary alert notifications.

### Grafana Cloud release channels

If you encounter unexpected bugs (such as data source credential fields not saving correctly), your Grafana Cloud instance may be on the **fast** release channel, which receives updates earlier. Switching to the **steady** channel provides more stability:

1. Navigate to your Grafana Cloud portal.
1. Check the current release channel under your stack settings.
1. If you're on the **fast** channel and experiencing issues, contact [Grafana Support](https://grafana.com/support/) to switch to the **steady** channel.

## Get additional help

If you've tried the solutions on this page and still encounter issues:

1. Check the [Grafana Cloud status page](https://status.grafana.com/) to rule out active platform incidents.
1. Search the [Grafana community forums](https://community.grafana.com/) for similar issues and solutions.
1. Review the [Redshift data source GitHub issues](https://github.com/grafana/redshift-datasource/issues) for known bugs and feature requests.
1. Consult the [Amazon Redshift documentation](https://docs.aws.amazon.com/redshift/) for service-specific guidance.
1. Contact [Grafana Support](https://grafana.com/support/) if you have a support contract.
1. When reporting issues, include:
   - Grafana version and plugin version.
   - Error messages (redact sensitive information such as credentials and ARNs).
   - Steps to reproduce the issue.
   - Relevant configuration details (redact credentials).
