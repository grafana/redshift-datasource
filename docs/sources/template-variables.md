---
description: Use template variables with the Amazon Redshift data source to create dynamic, reusable dashboards in Grafana.
keywords:
  - grafana
  - redshift
  - amazon redshift
  - template variables
  - variables
  - dynamic dashboards
labels:
  products:
    - cloud
    - enterprise
    - oss
menuTitle: Template variables
title: Amazon Redshift template variables
weight: 300
review_date: 2026-05-06
---

# Amazon Redshift template variables

Template variables let you create dynamic, reusable dashboards. Instead of hard-coding values like schema names, table names, or filter criteria, you can use variables that change based on user selections or query results.

For an introduction to variables in Grafana, refer to [Templates and variables](https://grafana.com/docs/grafana/<GRAFANA_VERSION>/dashboards/variables/).

## Before you begin

- [Configure the Amazon Redshift data source](https://grafana.com/docs/plugins/grafana-redshift-datasource/latest/configure/).
- Familiarize yourself with [Grafana template variables](https://grafana.com/docs/grafana/<GRAFANA_VERSION>/dashboards/variables/).

## Supported variable types

| Variable type | Supported |
| ------------- | --------- |
| Query         | Yes       |
| Custom        | Yes       |
| Data source   | Yes       |

## Create a query variable

Query variables populate their options by running an SQL query against your Redshift data source. Any value returned from a Redshift table can be used as a variable. The variable query editor provides the same SQL editor, resource selectors, and macro support as the regular query editor.

To create a query variable:

1. Navigate to **Dashboard settings** > **Variables**.
1. Click **Add variable**.
1. Select **Query** as the variable type.
1. Select the Amazon Redshift data source.
1. Enter an SQL query in the query editor.
1. Click **Run query** to preview the variable options.
1. Click **Apply**.

### Query examples

Return a list of schemas:

```sql
SELECT schema_name FROM information_schema.schemata;
```

Return a list of tables from a specific schema:

```sql
SELECT table_name FROM information_schema.tables WHERE table_schema = 'public';
```

Return distinct values from a column:

```sql
SELECT DISTINCT region FROM sales_data;
```

Return event categories from the Redshift sample database:

```sql
SELECT catname FROM public.category ORDER BY catname;
```

### Custom display names

To display a user-friendly label while using a different value in queries, alias the columns as `text` and `value`:

```sql
SELECT hostname AS text, id AS value FROM servers;
```

In this example, the variable drop-down shows hostnames, but the underlying value used in queries is the `id`.

{{< admonition type="caution" >}}
Avoid selecting large numbers of values in a single variable query, as this can cause performance issues.
{{< /admonition >}}

## Use variables in queries

After you create a variable, reference it in your Redshift queries using [variable syntax](https://grafana.com/docs/grafana/<GRAFANA_VERSION>/dashboards/variables/variable-syntax/). Grafana replaces the variable with the selected value before sending the query to Redshift.

For example, if you create a variable named `schema`:

```sql
SELECT * FROM ${schema}.my_table WHERE $__timeFilter(created_at);
```

### Multi-value variables

When a variable allows multiple selections, use the variable inside an `IN` clause:

```sql
SELECT * FROM sales WHERE region IN (${region});
```

Grafana automatically formats the selected values as a comma-separated list.

### Combine variables with macros

You can use template variables alongside Grafana macros in the same query:

```sql
SELECT
  $__timeGroup(saletime, '1h'),
  sum(pricepaid) AS total_sales
FROM
  ${schema}.sales
WHERE
  $__timeFilter(saletime) AND catname IN (${category})
GROUP BY
  1
ORDER BY
  1 ASC;
```

## Create cascading variables

You can create cascading (dependent) variables where one variable's options depend on another variable's selection. For example, create a `schema` variable first, then create a `table` variable that queries tables in the selected schema:

```sql
SELECT table_name FROM information_schema.tables WHERE table_schema = '${schema}';
```

When the user changes the `schema` selection, the `table` options automatically refresh.
