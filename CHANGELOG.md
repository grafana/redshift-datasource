# Changelog
 ## 1.8.1

- Update async-query-data with a fix for errors (#220) [#220](https://github.com/grafana/redshift-datasource/pull/220)

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
