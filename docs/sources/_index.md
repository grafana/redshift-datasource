---
description: Use the Amazon Redshift data source to query and visualize data from your Redshift clusters and serverless workgroups in Grafana.
keywords:
  - grafana
  - redshift
  - amazon redshift
  - aws
  - data source
  - data warehouse
  - sql
  - query
labels:
  products:
    - cloud
    - enterprise
    - oss
menuTitle: Amazon Redshift
title: Amazon Redshift data source
weight: 1350
review_date: 2026-05-06
---

# Amazon Redshift data source

The Amazon Redshift data source plugin lets you query and visualize data from Amazon Redshift in Grafana. Amazon Redshift is a fully managed, petabyte-scale cloud data warehouse that you can use to analyze structured and semi-structured data. With this plugin, you can write SQL queries against your Redshift provisioned clusters or Redshift Serverless workgroups and present results in Grafana dashboards.

## Supported features

The Redshift data source supports the following features:

| Feature     | Supported |
| ----------- | --------- |
| Metrics     | Yes       |
| Logs        | No        |
| Traces      | No        |
| Alerting    | Yes       |
| Annotations | Yes       |

Queries run asynchronously through the [Amazon Redshift Data API](https://docs.aws.amazon.com/redshift/latest/mgmt/data-api.html), which means long-running queries don't block or time out.

## Get started

The following pages help you set up and use the Amazon Redshift data source:

- [Configure the Amazon Redshift data source](https://grafana.com/docs/plugins/grafana-redshift-datasource/latest/configure/) to connect Grafana to your Redshift cluster or serverless workgroup.
- [Amazon Redshift query editor](https://grafana.com/docs/plugins/grafana-redshift-datasource/latest/query-editor/) to learn how to build SQL queries and use macros.
- [Template variables](https://grafana.com/docs/plugins/grafana-redshift-datasource/latest/template-variables/) to create dynamic, reusable dashboards.
- [Annotations](https://grafana.com/docs/plugins/grafana-redshift-datasource/latest/annotations/) to overlay event data on your graphs.
- [Alerting](https://grafana.com/docs/plugins/grafana-redshift-datasource/latest/alerting/) to create alert rules from Redshift queries.
- [Troubleshooting](https://grafana.com/docs/plugins/grafana-redshift-datasource/latest/troubleshooting/) for solutions to common issues.

## Additional features

After you configure the data source, you can:

- Use [Explore](https://grafana.com/docs/grafana/<GRAFANA_VERSION>/explore/) to run ad-hoc queries without building a dashboard.
- Add [transformations](https://grafana.com/docs/grafana/<GRAFANA_VERSION>/panels-visualizations/query-transform-data/transform-data/) to manipulate query results.
- Set up [alerting](https://grafana.com/docs/grafana/<GRAFANA_VERSION>/alerting/) rules based on Redshift queries.

## Async query data support

The Redshift data source uses an asynchronous query handling flow by default. Queries are handled over multiple requests -- starting the query, checking its status, and fetching the results -- instead of resolving over a single request. This is useful for queries that can potentially run for a long time and time out.

Your IAM policy must include the `redshift-data:ListStatements` and `redshift-data:CancelStatement` actions for async queries to work. Refer to the [IAM policies](https://grafana.com/docs/plugins/grafana-redshift-datasource/latest/configure/#iam-policies) section for details.

### Async query cache

To enable [query caching](https://grafana.com/docs/grafana/<GRAFANA_VERSION>/administration/data-source-management/#query-caching) for async queries, you need Grafana version 10.1 or later and the following feature toggles set to `true`:

- `useCachingService`
- `awsAsyncQueryCaching`

You also need to [configure query caching](https://grafana.com/docs/grafana/<GRAFANA_VERSION>/administration/data-source-management/#query-caching) for the specific Redshift data source instance.

## Pre-built dashboards

The Redshift data source includes three pre-built dashboards for monitoring your Redshift environment:

- **Redshift Monitoring** -- Tracks cluster performance metrics such as query execution times and throughput.
- **Redshift Privileges** -- Displays privilege and access configuration across your cluster.
- **Redshift Identities and Objects** -- Shows identity management and database object information.

These dashboards are based on similar dashboards from the [AWS Labs repository for Redshift](https://github.com/awslabs/amazon-redshift-monitoring).

To import a pre-built dashboard:

1. Navigate to the data source configuration page.
1. Select the **Dashboards** tab.
1. Click **Import** next to the dashboard you want to use.

For more information, refer to [Import a dashboard](https://grafana.com/docs/grafana/<GRAFANA_VERSION>/dashboards/build-dashboards/import-dashboards/).

## Plugin updates

Ensure that your plugin version is up to date so you have access to all current features and improvements. Navigate to **Administration** > **Plugins and data** > **Plugins** to check for updates.

{{< admonition type="note" >}}
Plugins are automatically updated in Grafana Cloud.
{{< /admonition >}}

## Related resources

- [Amazon Redshift documentation](https://docs.aws.amazon.com/redshift/)
- [Amazon Redshift Data API](https://docs.aws.amazon.com/redshift/latest/mgmt/data-api.html)
- [Redshift data source plugin on GitHub](https://github.com/grafana/redshift-datasource)
- [Grafana community forum](https://community.grafana.com/)
