# Changelog

## 1.13.3

- Upgrade @grafana/async-query-data from 0.1.10 to 0.1.11 https://github.com/grafana/redshift-datasource/pull/269

## 1.13.2

- Update grafana/aws-sdk-go to 0.20.0 https://github.com/grafana/redshift-datasource/pull/268

## 1.13.1

- Bump go.opentelemetry.io/contrib/instrumentation/net/http/httptrace/otelhttptrace from 0.37.0 to 0.44.0 by @dependabot in https://github.com/grafana/redshift-datasource/pull/257
- Upgrade grafana-plugin-sdk-go; add underscore, debug to package resolutions by @fridgepoet in https://github.com/grafana/redshift-datasource/pull/265

**Full Changelog**: https://github.com/grafana/redshift-datasource/compare/v1.13.0...v1.13.1

## 1.13.0

- Migrate Query and config editors to new form styling under feature toggle [#255](https://github.com/grafana/redshift-datasource/pull/255)
- Support Node 18 [#249](https://github.com/grafana/redshift-datasource/pull/249)
- Fix datasource type in provisioning docs in [#246](https://github.com/grafana/redshift-datasource/pull/246)

## 1.12.2

- Fix async queries by not calling ListStatements in GetQueryID [#252](https://github.com/grafana/redshift-datasource/pull/252)

## 1.12.1

- upgrade @grafana/aws-sdk to fix a bug in temporary credentials

## 1.12.0

- Update grafana-aws-sdk to v0.19.1 to add `il-central-1` to opt-in region list

## 1.11.1

- Upgrade @grafana/async-query-data to reduce minimum query time https://github.com/grafana/redshift-datasource/pull/237

## 1.11.0

- Upgrade grafana/aws-sdk-react dependency [#239](https://github.com/grafana/redshift-datasource/pull/236)
- Fix connection error when changing access and secret key [#235](https://github.com/grafana/redshift-datasource/pull/235)
- Support async query caching [#233](https://github.com/grafana/redshift-datasource/pull/233)

## 1.10.0

- Add support for Redshift Serverless https://github.com/grafana/redshift-datasource/pull/228 by @yota-p

## 1.9.0

- Upgrade @grafana/aws-sdk to v0.0.47 to support numeric values when applying template variables to SQL queries
- Fix async queries and expressions https://github.com/grafana/redshift-datasource/pull/225

## 1.8.4

- Upgrade Readme.md re: Grafana 10 https://github.com/grafana/redshift-datasource/pull/224

## 1.8.3

- Upgrade grafana/aws-sdk-react to 0.0.46 https://github.com/grafana/redshift-datasource/pull/223

## 1.8.2

- Update grafana-aws-sdk version to include new region in opt-in region list https://github.com/grafana/grafana-aws-sdk/pull/80
- Security: Upgrade Go in build process to 1.20.4
- Update grafana-plugin-sdk-go version to 0.161.0 to avoid a potential http header problem. https://github.com/grafana/athena-datasource/issues/233

## 1.8.1

- Update async-query-data with a fix for errors in [#220](https://github.com/grafana/redshift-datasource/pull/220)

## 1.8.0

- Update backend dependencies

## 1.7.0

- Fix converting rows with FLOAT, FLOAT4, and BOOLEAN by @iwysiu in [#213](https://github.com/grafana/redshift-datasource/pull/213)
- Add header component to Query Editor by @idastambuk in [#214](https://github.com/grafana/redshift-datasource/pull/214)
- Use organization ISSUE_COMMANDS_TOKEN with reduced scope by @iwysiu in [#210](https://github.com/grafana/redshift-datasource/pull/210)

## 1.6.0

- Remove run and cancel buttons in annotations editor in https://github.com/grafana/redshift-datasource/pull/206

## 1.5.0

- Migrate to create-plugin by @iwysiu in https://github.com/grafana/redshift-datasource/pull/195
- Update code coverage in workflow to latest by @idastambuk in https://github.com/grafana/redshift-datasource/pull/198
- Update @grafana/aws-sdk by @kevinwcyu in https://github.com/grafana/redshift-datasource/pull/199
- Update @grafana/ packages by @idastambuk in https://github.com/grafana/redshift-datasource/pull/201
- Upgrade grafana-aws-sdk to v0.12.0 by @fridgepoet in https://github.com/grafana/redshift-datasource/pull/202

## 1.4.1

- Hide the stop button when async query data support is not enabled https://github.com/grafana/redshift-datasource/pull/196

## 1.4.0

- Add Async Query Data Support https://github.com/grafana/redshift-datasource/pull/177

## 1.3.3

- Update @grafana dependencies to v8.5.10 https://github.com/grafana/redshift-datasource/pull/194

## 1.3.2

- Security: Upgrade Go in build process to 1.19.3

## 1.3.1

- Security: Upgrade Go in build process to 1.19.2

## 1.3.0

- Upgrade to grafana-aws-sdk v0.11.0 by @fridgepoet in https://github.com/grafana/redshift-datasource/pull/183

## 1.2.0

- Add database security monitoring dashboards by @yota-p in https://github.com/grafana/redshift-datasource/pull/175

## 1.1.0

- Add support for context aware autocompletion by @sunker in https://github.com/grafana/redshift-datasource/pull/174

## 1.0.7

- Bug fix for auth issues with when using keys and dependency upgrades (https://github.com/grafana/redshift-datasource/pull/165)
- Updates to code coverage

## 1.0.6

### What's Changed

- Update grafana-aws-sdk by @andresmgot in https://github.com/grafana/redshift-datasource/pull/146
- Autocomplete: Render SQL editor in case feature toggle is enabled by @sunker in https://github.com/grafana/redshift-datasource/pull/151
- fix: WLM panels query fix by @vgkowski in https://github.com/grafana/redshift-datasource/pull/152
- Custom redshift language by @sunker in https://github.com/grafana/redshift-datasource/pull/154
- Align Monaco language with official language ref by @sunker in https://github.com/grafana/redshift-datasource/pull/156

**Full Changelog**: https://github.com/grafana/redshift-datasource/compare/v1.0.5...v1.0.6

## 1.0.5

- Reduces backoff time factor to retrieve results.
- Upgrades internal dependencies.

## 1.0.4

- Add details in the datasource card #130
- Enable WithEvent to send an event to the AWS EventBridge #132

## 1.0.3

Fixes bugs for Endpoint and Assume Role settings.

## 1.0.2

Fixes a bug preventing from getting null values in a query.

## 1.0.1

Fixes a bug preventing from creating several data sources of the plugin in the same instance.

## 1.0.0

Initial release.

## 0.4.1

Improved curated dashboard.

## 0.4.0

Allow to authenticate using AWS Secret Manager. More bug fixes.

## 0.3.0

Third preview release. Includes curated dashboard.

## 0.2.0

Second release.

## 0.1.0

Initial release.
