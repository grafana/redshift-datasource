# Redshift data source for Grafana

The Redshift data source plugin allows you to query and visualize Redshift data metrics from within Grafana.

This topic explains options, variables, querying, and other options specific to this data source. Refer to [Add a data source](https://grafana.com/docs/grafana/latest/datasources/add-a-data-source/) for instructions on how to add a data source to Grafana.

## Configure the data source in Grafana

To access data source settings, hover your mouse over the **Configuration** (gear) icon, then click **Data Sources**, and then click the Amazon Redshift data source.

| Name                                | Description                                                                                                                                                                                                                                                                                                                                                  |
| ----------------------------------- | ------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------ |
| `Name`                              | The data source name. This is how you refer to the data source in panels and queries.                                                                                                                                                                                                                                                                        |
| `Default`                           | Default data source means that it will be pre-selected for new panels.                                                                                                                                                                                                                                                                                       |
| `Authentication Provider`           | Specify the provider to get credentials.                                                                                                                                                                                                                                                                                                                     |
| `Access Key ID`                     | If `Access & secret key` is selected, specify the Access Key of the security credentials to use.                                                                                                                                                                                                                                                             |
| `Secret Access Key`                 | If `Access & secret key` is selected, specify the Secret Key of the security credentials to use.                                                                                                                                                                                                                                                             |
| `Credentials Profile Name`          | Specify the name of the profile to use (if you use `~/.aws/credentials` file), leave blank for default.                                                                                                                                                                                                                                                      |
| `Assume Role Arn` (optional)        | Specify the ARN of the role to assume.                                                                                                                                                                                                                                                                                                                       |
| `External ID` (optional)            | If you are assuming a role in another account, that has been created with an external ID, specify the external ID here.                                                                                                                                                                                                                                      |
| `Endpoint` (optional)               | Optionally, specify a custom endpoint for the service.                                                                                                                                                                                                                                                                                                       |
| `Default Region`                    | Region in which the cluster is deployed.                                                                                                                                                                                                                                                                                                                     |
| `AWS Secrets Manager`               | To authenticate with Amazon Redshift using AWS Secrets Manager.                                                                                                                                                                                                                                                                                              |
| `Temporary credentials`             | To authenticate with Amazon Redshift using temporary database credentials.                                                                                                                                                                                                                                                                                   |
| `Serverless`                        | To use a Redshift Serverless workgroup.                                                                                                                                                                                                                                                                                                                      |
| `Cluster Identifier`                | Redshift Provisioned Cluster to use (automatically set if using AWS Secrets Manager).                                                                                                                                                                                                                                                                        |
| `Workgroup`                         | Redshift Serverless Workgroup to use.                                                                                                                                                                                                                                                                                                                        |
| `Managed Secret`                    | When using AWS Secrets Manager, select the secret containing the credentials to access the database. Note that Provisioned and Serverless stores credentials in a different format. Refer to [Storing database credentials in AWS Secrets Manager](https://docs.aws.amazon.com/redshift/latest/mgmt/data-api-access.html#data-api-secrets) for instructions. |
| `Database User`                     | User of the database. Automatically set if using AWS Secrets Manager.                                                                                                                                                                                                                                                                                        |
| `Database`                          | Name of the database within the cluster or workgroup.                                                                                                                                                                                                                                                                                                        |
| `Send events to Amazon EventBridge` | To send Data API events to Amazon EventBridge for monitoring purpose.                                                                                                                                                                                                                                                                                        |

## Authentication

For authentication options and configuration details, see [AWS authentication](https://grafana.com/docs/grafana/next/datasources/aws-cloudwatch/aws-authentication/) topic.

### IAM policies

Grafana needs permissions granted via IAM to be able to read Redshift metrics. You can attach these permissions to IAM roles and utilize Grafana's built-in support for assuming roles. Note that you will need to [configure the required policy](https://docs.aws.amazon.com/IAM/latest/UserGuide/access_policies_create.html) before adding the data source to Grafana. [You can check some predefined policies by AWS here](https://docs.aws.amazon.com/redshift/latest/mgmt/redshift-iam-access-control-identity-based.html#redshift-policy-resources.managed-policies).

Here is a minimal policy example:

```json
{
  "Version": "2012-10-17",
  "Statement": [
    {
      "Sid": "AllowReadingMetricsFromRedshift",
      "Effect": "Allow",
      "Action": [
        "redshift-data:ListTables",
        "redshift-data:DescribeTable",
        "redshift-data:GetStatementResult",
        "redshift-data:DescribeStatement",
        "redshift-data:ListStatements",
        "redshift-data:ListSchemas",
        "redshift-data:ExecuteStatement",
        "redshift-data:CancelStatement",
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

## Query Redshift data

The provided query editor is a standard SQL query editor. Grafana includes some macros to help with writing more complex timeseries queries.

#### Macros

| Macro                        | Description                                                                                                                      | Output example                                                   |
| ---------------------------- | -------------------------------------------------------------------------------------------------------------------------------- | ---------------------------------------------------------------- |
| `$__timeEpoch(column)`       | `$__timeEpoch` will be replaced by an expression to convert to a UNIX timestamp and rename the column to time                    | `UNIX_TIMESTAMP(dateColumn) as "time"`                           |
| `$__timeFilter(column)`      | `$__timeFilter` creates a conditional that filters the data (using `column`) based on the time range of the panel                | `time BETWEEN '2017-07-18T11:15:52Z' AND '2017-07-18T11:15:52Z'` |
| `$__timeFrom()`              | `$__timeFrom` outputs the current starting time of the range of the panel with quotes                                            | `'2017-07-18T11:15:52Z'`                                         |
| `$__timeTo()`                | `$__timeTo` outputs the current ending time of the range of the panel with quotes                                                | `'2017-07-18T11:15:52Z'`                                         |
| `$__timeGroup(column, '1m')` | `$__timeGroup` groups timestamps so that there is only 1 point for every period on the graph                                     | `floor(extract(epoch from time)/60)*60 AS "time"`                |
| `$__schema`                  | `$__schema` uses the selected schema                                                                                             | `public`                                                         |
| `$__table`                   | `$__table` outputs a table from the given `$__schema` (it uses the `public` schema by default)                                   | `sales`                                                          |
| `$__column`                  | `$__column` outputs a column from the current `$__table`                                                                         | `date`                                                           |
| `$__unixEpochFilter(column)` | `$__unixEpochFilter` be replaced by a time range filter using the specified column name with times represented as Unix timestamp | `column >= 1624406400 AND column <= 1624410000`                  |
| `$__unixEpochGroup(column)`  | `$__unixEpochGroup` is the same as $\_\_timeGroup but for times stored as Unix timestamp                                         | `floor(time/60)*60 AS "time"`                                    |

#### Table Visualization

Most queries in Redshift will be best represented by a table visualization. Any query will display data in a table. If it can be queried, then it can be put in a table.

This example returns results for a table visualization:

```sql
SELECT {column_1}, {column_2} FROM {table};
```

#### Timeseries / Graph visualizations

For timeseries / graph visualizations, there are a few requirements:

- A column with a `date` or `datetime` type must be selected
- The `date` column must be in ascending order (using `ORDER BY column ASC`)
- A numeric column must also be selected

To make a more reasonable graph, be sure to use the `$__timeFilter` and `$__timeGroup` macros.

Example timeseries query:

```sql
SELECT
  avg(execution_time) AS average_execution_time,
  $__timeGroup(start_time, 'hour'),
  query_type
FROM
  account_usage.query_history
WHERE
  $__timeFilter(start_time)
group by
  query_type,start_time
order by
  start_time,query_type ASC;
```

##### Fill value

When data frames are formatted as time series, you can choose how missing values should be filled. This in turn affects how they are rendered: with connected or disconnected values. To configure this value, change the "Fill Value" in the query editor.

#### Inspecting the query

Because Grafana supports macros that Redshift does not, the fully rendered query, which can be copy/pasted directly into Redshift, is visible in the Query Inspector. To view the full interpolated query, click the Query Inspector button, and the full query will be visible under the "Query" tab.

### Templates and variables

To add a new Redshift query variable, refer to [Add a query variable](https://grafana.com/docs/grafana/latest/variables/variable-types/add-query-variable/). Use your Redshift data source as your data source for the following available queries:

Any value queried from a Redshift table can be used as a variable. Be sure to avoid selecting too many values, as this can cause performance issues.

To display a custom display name for a variable, you can use a query such as `SELECT hostname AS text, id AS value FROM MyTable`. In this case the variable value field must be a string type or cast to a string type.

After creating a variable, you can use it in your Redshift queries by using [Variable syntax](https://grafana.com/docs/grafana/latest/variables/syntax/). For more information about variables, refer to [Templates and variables](https://grafana.com/docs/grafana/latest/variables/).

### Annotations

[Annotations](https://grafana.com/docs/grafana/latest/dashboards/annotations/) allow you to overlay rich event information on top of graphs. You can add annotations by clicking on panels or by adding annotation queries via the Dashboard menu / Annotations view.

**Example query to automatically add annotations:**

```sql
SELECT
  time as time,
  environment as tags,
  humidity as text
FROM
  $__table
WHERE
  $__timeFilter(time) and humidity > 95
```

The following table represents the values of the columns taken into account to render annotations:

| Name      | Description                                                                                                                       |
| --------- | --------------------------------------------------------------------------------------------------------------------------------- |
| `time`    | The name of the date/time field. Could be a column with a native SQL date/time data type or epoch value.                          |
| `timeend` | Optional name of the end date/time field. Could be a column with a native SQL date/time data type or epoch value. (Grafana v6.6+) |
| `text`    | Event description field.                                                                                                          |
| `tags`    | Optional field name to use for event tags as a comma separated string.                                                            |

## Provision Redshift data source

You can configure the Redshift data source using configuration files with Grafana's provisioning system. For more information, refer to the [provisioning docs page](https://grafana.com/docs/grafana/latest/administration/provisioning/).

Here are some provisioning examples.

### Using AWS SDK (default)

```yaml
apiVersion: 1
datasources:
  - name: Redshift
    type: redshift
    jsonData:
      authType: default
      defaultRegion: eu-west-2
```

### Using credentials' profile name (non-default)

```yaml
apiVersion: 1

datasources:
  - name: Redshift
    type: redshift
    jsonData:
      authType: credentials
      defaultRegion: eu-west-2
      profile: secondary
```

### Using `accessKey` and `secretKey`

```yaml
apiVersion: 1

datasources:
  - name: Redshift
    type: grafana-redshift-datasource
    jsonData:
      authType: keys
      defaultRegion: eu-west-2
    secureJsonData:
      accessKey: '<your access key>'
      secretKey: '<your secret key>'
```

### Using AWS SDK Default and ARN of IAM Role to Assume

```yaml
apiVersion: 1
datasources:
  - name: Redshift
    type: grafana-redshift-datasource
    jsonData:
      authType: default
      assumeRoleArn: arn:aws:iam::123456789012:root
      defaultRegion: eu-west-2
```

## Pre-configured Redshift dashboards

Redshift data source ships with a pre-configured dashboard for some advanced monitoring parameters. This curated dashboard is based on similar dashboards in the [AWS Labs repository for Redshift](https://github.com/awslabs/amazon-redshift-monitoring). Check it out for more details.

Follow these [instructions](https://grafana.com/docs/grafana/latest/dashboards/export-import/#importing-a-dashboard) for importing a dashboard in Grafana.

Imported dashboards can be found in Configuration > Data Sources > select your Redshift data source > select the Dashboards tab to see available pre-made dashboards.

## Get the most out of the plugin

- Add [Annotations](https://grafana.com/docs/grafana/latest/dashboards/annotations/).
- Configure and use [Templates and variables](https://grafana.com/docs/grafana/latest/variables/).
- Add [Transformations](https://grafana.com/docs/grafana/latest/panels/transformations/).
- Set up alerting; refer to [Alerts overview](https://grafana.com/docs/grafana/latest/alerting/).

## Async Query Data Support

Async Query Data support enables an asynchronous query handling flow. With Async Query Data support enabled, queries will be handled over multiple requests (starting, checking its status, and fetching the results) instead of having a query be started and resolved over a single request. This is useful for queries that can potentially run for a long time and timeout. You'll need to ensure the IAM policy used by Grafana allows the following actions `redshift-data:ListStatements` and `redshift-data:CancelStatement`.

Async query data support is enabled by default in all Redshift datasources.

### Async Query Caching

To enable [query caching](https://grafana.com/docs/grafana/latest/administration/data-source-management/#query-caching) for async queries, you need to be on Grafana version 10.1 or above, and to set the feature toggles `useCachingService` and `awsAsyncQueryCaching` to `true`. You'll also need to [configure query caching](https://grafana.com/docs/grafana/latest/administration/data-source-management/#query-caching) for the specific Redshift datasource.

### Plugin repository

You can request new features, report issues, or contribute code directly through the [Redshift Data Source Github repository](https://github.com/grafana/redshift-datasource)