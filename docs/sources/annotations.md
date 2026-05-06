---
description: Add annotations to Grafana dashboards using the Amazon Redshift data source.
keywords:
  - grafana
  - redshift
  - amazon redshift
  - annotations
  - events
labels:
  products:
    - cloud
    - enterprise
    - oss
menuTitle: Annotations
title: Amazon Redshift annotations
weight: 400
review_date: 2026-05-06
---

# Amazon Redshift annotations

Annotations let you overlay event information on top of graphs. You can use SQL queries against your Redshift data to mark important events, deployments, incidents, or other time-based occurrences directly on your dashboard panels.

For general information about annotations, refer to [Annotate visualizations](https://grafana.com/docs/grafana/<GRAFANA_VERSION>/dashboards/build-dashboards/annotate-visualizations/).

## Before you begin

- [Configure the Amazon Redshift data source](https://grafana.com/docs/plugins/grafana-redshift-datasource/latest/configure/).
- Verify that your Redshift database contains the event data you want to annotate.

## Add an annotation query

To add an annotation query to a dashboard:

1. Open the dashboard you want to annotate.
1. Click **Dashboard settings** (the gear icon).
1. Select **Annotations** in the left menu.
1. Click **Add annotation query**.
1. Select the Amazon Redshift data source.
1. Enter a SQL query that returns the required columns.
1. Click **Apply**.

## Column conventions

Your annotation query must return specific columns for Grafana to render annotations correctly. The annotation query editor provides the same SQL editor, resource selectors, and macro support as the standard [query editor](https://grafana.com/docs/plugins/grafana-redshift-datasource/latest/query-editor/).

| Column    | Required | Description                                                                                              |
| --------- | -------- | -------------------------------------------------------------------------------------------------------- |
| `time`    | Yes      | The date/time field for the annotation. Can be a column with a native SQL date/time data type or an epoch value. |
| `timeend` | No       | The end date/time for region annotations that span a time range. Uses the same format as `time`.          |
| `text`    | Yes      | The description text displayed when you hover over the annotation.                                        |
| `tags`    | No       | A comma-separated string used to tag and filter annotations.                                              |

## Example annotation query

The following query creates annotations for high-humidity events:

```sql
SELECT
  time AS time,
  environment AS tags,
  humidity AS text
FROM
  $__table
WHERE
  $__timeFilter(time) AND humidity > 95
```

This query:

- Uses the `$__table` macro to reference the table selected in the resource selector.
- Filters results to the dashboard time range with `$__timeFilter(time)`.
- Creates an annotation for each record where humidity exceeds 95.
- Tags each annotation with the `environment` value.

## Region annotations

To create annotations that span a time range (region annotations), include both a `time` and `timeend` column:

```sql
SELECT
  start_time AS time,
  end_time AS timeend,
  description AS text,
  severity AS tags
FROM
  incidents
WHERE
  $__timeFilter(start_time)
```

Region annotations display as shaded areas on graphs, useful for marking maintenance windows, incidents, or other events with a duration.

{{< admonition type="note" >}}
Annotation queries run asynchronously through the Redshift Data API, just like regular queries. Complex annotation queries or queries that return large result sets can affect dashboard load times.
{{< /admonition >}}

## Use template variables in annotations

You can reference [template variables](https://grafana.com/docs/plugins/grafana-redshift-datasource/latest/template-variables/) in your annotation queries to make them dynamic:

```sql
SELECT
  event_time AS time,
  event_description AS text,
  event_type AS tags
FROM
  events
WHERE
  $__timeFilter(event_time) AND region = '${region}'
```
