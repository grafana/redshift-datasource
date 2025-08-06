# Building and releasing

## How to build the Redshift data source plugin locally

## Dependencies

Make sure you have the following dependencies installed first:

- [Git](https://git-scm.com/)
- [Go](https://golang.org/dl/) (see [go.mod](../go.mod#L3) for minimum required version)
- [Mage](https://magefile.org/)
- [Node.js (Long Term Support)](https://nodejs.org)
- [Yarn](https://yarnpkg.com)

## Frontend

1. Install dependencies

   ```bash
   yarn install --pure-lockfile
   ```

2. Build plugin in development mode or run in watch mode

   ```bash
   yarn dev
   ```

3. Build plugin in production mode

   ```bash
   yarn build
   ```

## Backend

1. Build the backend binaries

   ```bash
   mage -v
   ```

## Development with local Grafana

Checkout the guide available at https://grafana.com/docs/grafana/latest/developers/plugins/development-with-local-grafana/.

## Setting up a go workspace

Setting up go workspace can be helpful when making changes across modules like `grafana-aws-sdk` and `sqlds` and wanting to see those changes reflected in the Redshift data source.

From https://go.dev/blog/get-familiar-with-workspaces:

> Workspaces in Go 1.18 let you work on multiple modules simultaneously without having to edit go.mod files for each module. Each module within a workspace is treated as a main module when resolving dependencies.
>
> Previously, to add a feature to one module and use it in another module, you needed to either publish the changes to the first module, or edit the go.mod file of the dependent module with a replace directive for your local, unpublished module changes. In order to publish without errors, you had to remove the replace directive from the dependent module’s go.mod file after you published the local changes to the first module.

1. Make a new directory somewhere, for example `redshift_workspace`
2. `cd redshift_workspace`
3. `git clone https://github.com/grafana/redshift-datasource.git`
4. `git clone https://github.com/grafana/grafana-aws-sdk`
5. `git clone https://github.com/grafana/sqlds`
6. `go work init ./redshift-datasource ./grafana-aws-sdk ./sqlds`
7. Make modifications in any of these directories and build the backend in `redshift-datasource` with `mage` as usual. The changes in these directories will be taken into account.

If you build Grafana locally, you can for example symlink `redshift-datasource` to your clone of `github.com/grafana/grafana`'s `data/plugins` directory, e.g. `cd <path to your Grafana repo>/github.com/grafana/grafana/data/plugins && ln -s <path to your workspaces>/redshift_workspace/redshift-datasource redshift-datasource`

## E2E Tests

1. `yarn playwright install --with-deps`
1. `yarn server`
1. `yarn e2e`

## Build a release for the Redshift data source plugin

You need to have commit rights to the GitHub repository to publish a release.

1. Update the version number in the `package.json` file.
2. Update the `CHANGELOG.md` with the changes contained in the release.
3. Commit the changes to master and push to GitHub.
4. Follow the release process that you can find [here](https://enghub.grafana-ops.net/docs/default/component/grafana-plugins-platform/plugins-ci-github-actions/010-plugins-ci-github-actions/#cd_1)

# Plugin Technical Documentation

## What is AWS Redshift?

Amazon Redshift is a fully managed, petabyte-scale data warehouse service in the cloud. You can start with just a few hundred gigabytes of data and scale to a petabyte or more. This enables you to use your data to acquire new insights for your business and customers.

Redshift is based on an older version of PostgreSQL (fork of 8.0.2). It still has a lot in common and it is possible to use the Postgres data source in core Grafana to query Redshift. However, in the last couple of years it has started to diverge more and more. [While the SQL language in Redshift has started to diverge from Postgres](https://docs.aws.amazon.com/redshift/latest/dg/c_redshift-and-postgres-sql.html), the biggest differences have been in the underlying fundamentals of how the data is stored - PostgreSQL is a row-store database while RedShift is column-stored database.

While it's possible to use the PostgreSQL datasource to query Redshift, this plugins provides better integration with AWS authentication and provides Redshift exclusive features like its query editor, which uses Redshift built-in functions and specific syntax highlight.

## Authentication

The [awsds](https://github.com/grafana/grafana-aws-sdk/tree/main/pkg/awsds) package is currently used by all the Grafana AWS data sources. It contains logic for creating, caching and providing aws-sdk-go sessions based upon the selected type of authentication and its parameters. At the moment, four different types of authentication are supported - Workspace IAM role (currently only enabled in Amazon Managed Grafana), AWS SDK Default, Credential file auth, and Keys & secrets.

The Redshift configuration page renders the @grafana/aws-sdk ConnectionConfig, which shows all the common fields such as type of authentication, assume role details etc. However, in order to test the connection to AWS using temporary credentials, we need three additional fields - cluster identifier, database and db user. These fields are not needed when creating the aws-sdk-go session, but are params to each new SQL query. Consequently, we could add these fields in the query editor so that [they could be changed for each query](https://github.com/grafana/redshift-datasource/issues/42) (similar to how the region can be changed for each query in the CloudWatch data source).

### AWS Policy

There are four predefined (AWS managed) policies for Amazon Redshift.

- AmazonRedshiftReadOnlyAccess – Grants read-only access to all Amazon Redshift resources for the AWS account.
- AmazonRedshiftFullAccess – Grants full access to all Amazon Redshift resources for the AWS account.
- AmazonRedshiftQueryEditor – Grants full access to the Query Editor on the Amazon Redshift console.
- AmazonRedshiftDataFullAccess – Grants full access to the Amazon Redshift Data API operations and resources for the AWS account.

The intent of the Redshift query editor in Grafana is **not** to be used to modify values in the data source, so in the [documentation](./README.md#iam-policies) we make it clear that the role that is running the data source in Grafana should only have the AmazonRedshiftReadOnlyAccess policy attached to it.

## Architecture

The idiomatic way to use a SQL, or SQL-like, database in Go is through the [database/sql package](https://golang.org/pkg/database/sql/). The sql package provides a generic interface around SQL databases. One main benefit of using this pattern for data fetching is that we are reusing building blocks from other SQL-like data source plugins in Grafana.

### grafana/sqlds and sqlutil

From the [sqlds](https://github.com/grafana/sqlds) readme:

_sqlds stands for SQL Datasource._

_Most SQL-driven datasources, like Postgres, MySQL, and MSSQL share extremely similar codebases._

_The sqlds package is intended to remove the repetition of these datasources and centralize the datasource logic. The only thing that the datasources themselves should have to define is connecting to the database, and what driver to use, and the plugin frontend._

Furthermore, sqlds allows each datasource to implement its own fillmode, macros and string converters.

Internally, sqlds is using [sqlutil](https://github.com/grafana/grafana-plugin-sdk-go/tree/master/data/sqlutil) which is a package in `grafana-plugin-sdk-go`. `sqlutil` exposes utility functions for converting database/sql rows into data frames.

### Redshift driver

The database/sql package can only be used in conjunction with a database driver. The AWS Redshift team offers support for two drivers - [JDBC](https://docs.aws.amazon.com/redshift/latest/mgmt/configure-jdbc-connection.html) and [ODBC](https://docs.aws.amazon.com/redshift/latest/mgmt/configure-odbc-connection.html). JDBC, which is the recommended driver, doesn’t have a golang version and it has a dependency to a Java runtime. There are golang drivers available for ODBC, but they still require an ODBC manager to be installed on the machine that is running grafana. Also we’d have to fork the golang driver in order to integrate it with our AWS sdk for handling authentication. For these reasons, we decided to use the [Amazon Redshift Data API](https://docs.aws.amazon.com/redshift/latest/mgmt/data-api.html) instead.

This plugin implements our own sql driver on top of the Amazon Redshift data api. To meet the criteria of a sql driver, a few interfaces have been implemented. The database/sql api is big, but we haven't added support for everything - only the parts that we need.

#### Amazon Redshift Data API

In contrast with many database/sql drivers, the redshift data api doesn't require a persistent connection to the cluster. Instead, it provides a secure HTTP endpoint and integration with AWS SDKs. Therefore, there’ll be no need for connection pooling.

The procedure of fetching data based on a SQL query has the following order:

1. Run a SQL statement

Request example:

```console
aws redshift-data execute-statement
    --region us-west-2
    --db-user myuser
    --cluster-identifier mycluster-test
    --database dev
    --sql "select * from stl_query limit 1"
```

Response example:

```json
{
  "ClusterIdentifier": "mycluster-test",
  "CreatedAt": 1598306924.632,
  "Database": "dev",
  "DbUser": "myuser",
  "Id": "aaaaaaaa-bbbb-cccc-dddd-eeeeeeeee"
}
```

2. Describe metadata about an SQL statement

Request example:

```console
aws redshift-data describe-statement
    --id aaaaaaaa-bbbb-cccc-dddd-eeeeeeeee
    --region us-west-2
```

Response example;

```json
{
  "ClusterIdentifier": "mycluster-test",
  "CreatedAt": 1598306924.632,
  "Duration": 1095981511,
  "Id": "aaaaaaaa-bbbb-cccc-dddd-eeeeeeeee",
  "QueryString": "select * from stl_query limit 1",
  "RedshiftPid": 20859,
  "RedshiftQueryId": 48879,
  "ResultRows": 1,
  "ResultSize": 4489,
  "Status": "FINISHED",
  "UpdatedAt": 1598306926.667
}
```

3. Fetch the results of an SQL statement

Request example:

```console
aws redshift-data get-statement-result
    --id aaaaaaaa-bbbb-cccc-dddd-eeeeeeeee
    --region us-west-2
```

Response example:

```json
{
    "ColumnMetadata": [
        {
            "isCaseSensitive": false,
            "isCurrency": false,
            "isSigned": true,
            "label": "userid",
            "length": 0,
            "name": "userid",
            "nullable": 0,
            "precision": 10,
            "scale": 0,
            "schemaName": "",
            "tableName": "stll_query",
            "typeName": "int4"
        },
   …
 "Records": [
        [
            {
                "longValue": 1
            },
]
}
```

For step 2, initially it keeps on calling `DescribeStatement` until the status is finished or errored. In the longer term, this procedure should be non-blocking by using polling or something like that.

The redshift data api also has commands for list-databases, list-schemas, list-tables and describe-tables. But for the sake of simplicity simplicity we are using SQL queries to retrieve this metadata for the data source UI.
