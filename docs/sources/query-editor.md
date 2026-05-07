---
description: Use the Amazon Redshift query editor to write SQL queries and visualize data in Grafana.
keywords:
  - grafana
  - redshift
  - amazon redshift
  - query editor
  - sql
  - macros
  - time series
  - table
labels:
  products:
    - cloud
    - enterprise
    - oss
menuTitle: Query editor
title: Amazon Redshift query editor
weight: 200
review_date: 2026-05-06
---

# Amazon Redshift query editor

The Amazon Redshift query editor lets you write SQL queries to retrieve and visualize data from your Redshift clusters or serverless workgroups. The editor provides a code editor with context-aware autocompletion, resource selectors for common database objects, and support for Grafana macros.

## Before you begin

- [Configure the Amazon Redshift data source](https://grafana.com/docs/plugins/grafana-redshift-datasource/latest/configure/).
- Verify that your IAM identity has the required permissions to query the target database.
- Familiarize yourself with [Amazon Redshift SQL](https://docs.aws.amazon.com/redshift/latest/dg/c_redshift-sql.html) syntax and supported functions.

## Resource selectors

The query editor includes three drop-down selectors at the top of the editor that let you choose database objects. Each selector populates a corresponding Grafana macro you can use in your SQL queries:

| Selector     | Macro        | Description                                                      |
| ------------ | ------------ | ---------------------------------------------------------------- |
| **Schema**   | `$__schema`  | Selects a database schema. If not set, Redshift uses the default `search_path` (typically `public`). |
| **Table**    | `$__table`   | Selects a table from the chosen schema.                          |
| **Column**   | `$__column`  | Selects a column from the chosen table.                          |

These selectors are optional. You can write queries without using them.

## SQL code editor

The query editor uses a Monaco-based code editor with the following features:

- **Syntax highlighting** for Redshift SQL.
- **Context-aware autocompletion** that suggests schemas, tables, columns, and SQL keywords as you type.
- **Macro expansion** so you can use Grafana macros like `$__timeFilter` directly in your queries.

To write a query, click inside the editor and start typing SQL. Press `Ctrl+Space` (or `Cmd+Space` on macOS) to trigger autocompletion.

## Format options

Expand the **Format** section to configure how query results are formatted. The following options are available:

| Format          | Description                                                                                              |
| --------------- | -------------------------------------------------------------------------------------------------------- |
| **Table**       | Returns results as a table. This is the default format and works with any query.                         |
| **Time Series** | Returns results as time-series data frames for use with time-series visualizations like the Graph panel.  |

### Time-series requirements

To use the **Time Series** format, your query must meet the following requirements:

- Include a column with a `date` or `datetime` data type.
- Order the date column in ascending order using `ORDER BY column ASC`.
- Include at least one numeric column.

### Fill value

When the format is set to **Time Series**, the **Fill value** option controls how missing data points are rendered. This affects whether lines in graphs appear connected or disconnected. Choose a fill mode based on how you want gaps in data to display:

- **Previous** (default) -- Fills missing values with the last known value.
- **Null** -- Leaves gaps in the data.
- **Value** -- Fills missing values with a specific number you define.

## Macros

The query editor supports Grafana macros that simplify time-series queries. Macros are expanded to Redshift-compatible SQL before the query is executed.

| Macro                              | Description                                                                       | Output example                                                     |
| ---------------------------------- | --------------------------------------------------------------------------------- | ------------------------------------------------------------------ |
| `$__timeEpoch(column)`             | Converts a timestamp column to a Unix epoch and renames it to `time`.             | `extract(epoch from dateColumn) as "time"`                         |
| `$__timeFilter(column)`            | Filters results to the dashboard's selected time range.                           | `time BETWEEN '2017-07-18T11:15:52Z' AND '2017-07-18T12:15:52Z'`  |
| `$__timeFrom()`                    | Returns the start of the current time range.                                      | `'2017-07-18T11:15:52Z'`                                          |
| `$__timeTo()`                      | Returns the end of the current time range.                                        | `'2017-07-18T11:15:52Z'`                                          |
| `$__timeGroup(column, 'interval')` | Groups timestamps into fixed intervals so there's one data point per interval.    | `floor(extract(epoch from time)/60)*60 AS "time"`                  |
| `$__schema`                        | Replaced with the schema selected in the resource selector.                       | `public`                                                           |
| `$__table`                         | Replaced with the table selected in the resource selector.                        | `sales`                                                            |
| `$__column`                        | Replaced with the column selected in the resource selector.                       | `date`                                                             |
| `$__unixEpochFilter(column)`       | Filters results by the time range using Unix timestamps.                          | `column >= 1624406400 AND column <= 1624410000`                    |
| `$__unixEpochGroup(column, 'interval')` | Groups Unix epoch timestamps into fixed intervals.                           | `floor(time/60)*60 AS "time"`                                      |

## Query examples

The following examples demonstrate common query patterns.

### Table query

This query retrieves columns from a table and is best suited for the **Table** visualization:

```sql
SELECT column_1, column_2 FROM my_schema.my_table LIMIT 100;
```

You can also use resource selectors with macros:

```sql
SELECT * FROM $__schema.$__table LIMIT 100;
```

### Time-series query

This query calculates total sales commission grouped by time and is suited for time-series visualizations:

```sql
SELECT
  $__timeGroup(saletime, '1h'),
  sum(commission) AS total_commission,
  eventname
FROM
  public.sales
  JOIN public.event USING (eventid)
WHERE
  $__timeFilter(saletime)
GROUP BY
  1, eventname
ORDER BY
  1 ASC;
```

### Use case: monitor query performance

To track Redshift query execution time over time:

```sql
SELECT
  $__timeGroup(starttime, '5m'),
  avg(elapsed) / 1000000 AS avg_duration_seconds,
  query_group
FROM
  stl_query
WHERE
  $__timeFilter(starttime)
GROUP BY
  1, query_group
ORDER BY
  1 ASC;
```

### Use case: track storage usage

To monitor table sizes within a schema:

```sql
SELECT
  "table" AS table_name,
  size AS size_mb,
  tbl_rows AS row_count
FROM
  svv_table_info
WHERE
  schema = 'public'
ORDER BY
  size DESC;
```

## Inspect the query

Because Grafana macros aren't valid Redshift SQL, the actual query sent to Redshift differs from what you write in the editor. To view the fully interpolated query:

1. Open the panel editor.
1. Click **Query Inspector**.
1. Select the **Query** tab to see the rendered SQL.

You can copy the rendered query and run it directly in any Redshift SQL client for debugging.
