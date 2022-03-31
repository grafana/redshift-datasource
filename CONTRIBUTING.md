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

   or

   ```bash
   yarn watch
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

This guide allows you to setup a development environment where you run Grafana and your plugin locally. With this, you will be able to see your changes as you add them.

### Run Grafana in your host

If you have git, Go and the required version of NodeJS in your system, you can clone and run Grafana locally:

1. Clone [Grafana](https://github.com/grafana/grafana)

   a. [optional] Checkout to the specific version you want to target (e.g. `git checkout v8.3.4`)

2. Install its dependencies

   ```bash
   yarn install
   ```

3. Grafana will look for plugins, by default, on its `data/plugins` directory. You can create a symbolic link to your plugin repository to detect new changes:

   ```bash
   ln -s <plugin-path>/dist data/plugins/<plugin-name>
   ```

4. (Optional) If the step above doesn't work for you (e.g. you are running on Windows), you can also modify the default path in the Grafana configuration (that can be found at `conf/custom.ini`) and point to the directory with your plugin:

   ```ini
   [paths]
   plugins = <path-to-your-plugin-parent-directory>
   ```

### Run Grafana with docker-compose

Another possibility is to run Grafana with docker-compose so it runs in a container. For doing so, create the docker-compose file in your plugin directory:

```yaml
version: '3.7'

services:
  grafana:
    # Change latest with your target version, if needed
    image: grafana/grafana:latest
    ports:
      - '3000:3000'
    volumes:
      # Use your plugin folder (e.g. redshift-datasource)
      - ./:/var/lib/grafana/plugins/<plugin-folder>
    environment:
      - TERM=linux
      - GF_LOG_LEVEL=debug
      - GF_DATAPROXY_LOGGING=true
      # Use your plugin name (e.g. grafana-redshift-datasource)
      - GF_PLUGINS_ALLOW_LOADING_UNSIGNED_PLUGINS=<plugin-name>
```

### Run your plugin

Finally start your plugin in development mode:

1. Build your plugin backend and start the frontend in watch mode (if you haven't done that already):

   ```bash
   mage -v
   yarn watch
   ```

2. Start Grafana backend and frontend:

   2.1 For a local copy of Grafana:

   ```bash
   make run
   ```

   ```bash
   yarn start
   ```

   2.2 For docker-compose:

   ```bash
   docker-compose up
   ```

After this, you should be able to see your plugin listed in Grafana and test your changes. Note that any change in the fronted will require you to refresh your browser while changes in the backend may require to rebuild your plugin binaries and reload the plugin (`mage && mage reloadPlugin` for local development or `docker-compose up` again if you are using docker-compose).

## Build a release for the Redshift data source plugin

You need to have commit rights to the GitHub repository to publish a release.

1. Update the version number in the `package.json` file.
2. Update the `CHANGELOG.md` with the changes contained in the release.
3. Commit the changes to master and push to GitHub.
4. Follow the Drone release process that you can find [here](https://github.com/grafana/integrations-team/wiki/Plugin-Release-Process#drone-release-process)

# Plugin Technical Documentation

## What is AWS Redshift?

Amazon Redshift is a fully managed, petabyte-scale data warehouse service in the cloud. You can start with just a few hundred gigabytes of data and scale to a petabyte or more. This enables you to use your data to acquire new insights for your business and customers.

Redshift is based on an older version of PostgreSQL (fork of 8.0.2). It still has a lot in common and it is possible to use the Postgres data source in core Grafana to query Redshift. However, in the last couple of years it has started to diverge more and more. [While the SQL language in Redshift has started to diverge from Postgres](https://docs.aws.amazon.com/redshift/latest/dg/c_redshift-and-postgres-sql.html), the biggest differences have been in the underlying fundamentals of how the data is stored - PostgreSQL is a row-store database while RedShift is column-stored database.

While it's possible to use the PostgreSQL datasource to query Redshift, this plugins provides better integration with AWS authentication and provides Redshift exclusive features like its query editor, which uses Redshift built-in functions and specific syntax hightlight.

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
