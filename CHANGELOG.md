# Changelog

## 1.20.0

- Add PDC support in [#333](https://github.com/grafana/redshift-datasource/pull/333)
- Bump node dependencies with 29 updates, ignore react and react-dom major updates in dependabot in [#336](https://github.com/grafana/redshift-datasource/pull/336)
- Bump the all-go-dependencies group across 1 directory with 3 updates in [#335](https://github.com/grafana/redshift-datasource/pull/335)
- Add pre-commit hook in [#327](https://github.com/grafana/redshift-datasource/pull/327)

## 1.19.1

- Dependabot: Update dependencies in [#302](https://github.com/grafana/redshift-datasource/pull/302), [#303](https://github.com/grafana/redshift-datasource/pull/303), [#313](https://github.com/grafana/redshift-datasource/pull/313),
  [#308](https://github.com/grafana/redshift-datasource/pull/308), [#323](https://github.com/grafana/redshift-datasource/pull/323), [#322](https://github.com/grafana/redshift-datasource/pull/322):
  - path-to-regexp from 1.8.0 to 1.9.0 in the npm_and_yarn group
  - micromatch from 4.0.5 to 4.0.8 in the npm_and_yarn group
  - actions/checkout from 2 to 4
  - actions/setup-node from 3 to 4
  - tibdex/github-app-token from 1.8.0 to 2.1.0
  - github.com/aws/aws-sdk-go from 1.51.31 to 1.55.5
  - github.com/grafana/grafana-plugin-sdk-go from 0.254.0 to 0.259.4
  - github.com/grafana/grafana-aws-sdk from 0.31.3 to 0.31.4
  - github.com/grafana/sqlds/v4 from 4.1.2 to 4.1.4
  - braces from 3.0.2 to 3.0.3 in the npm_and_yarn group
  - github.com/stretchr/testify from 1.9.0 to 1.10.0(#321)
  - @emotion/css 11.13.4 11.13.5
  - @grafana/async-query-data from 0.2.0 to 0.3.0
  - @grafana/data from 11.2.2 to 11.3.1
  - @grafana/experimental from 2.1.2 to 2.1.4
  - @grafana/runtime from 11.2.2 to 11.3.1
  - tslib from 2.8.0 to 2.8.1
  - @babel/core from 7.25.8 to 7.26.0
  - @grafana/eslint-config from 7.0.0 to 8.0.0
  - @swc/core from 1.7.36 to 1.9.3
  - @swc/helpers from 0.5.13 to 0.5.15
  - @swc/jest from 0.2.36 to 0.2.37
  - @testing-library/jest-dom from 6.6.0 to 6.6.3
  - @types/jest from 29.5.13 to 29.5.14
  - @types/lodash from 4.17.10 to 4.17.13
  - @types/node from 22.7.5 to 22.10.1
  - cspell from 8.15.2 to 8.16.1
  - cypress from 7.7.0 to 13.16.0
  - prettier from 3.3.3 to 3.4.1
  - sass from 1.79.5 to 1.81.0
  - sass-loader from 16.0.2 to 16.0.3
  - typescript from 5.6.3 to 5.7.2
  - webpack from 5.95.0 to 5.96.1
  - cross-spawn from 7.0.3 to 7.0.6 in the npm_and_yarn group

## 1.19.0

- Chore: update dependencies [#300](https://github.com/grafana/redshift-datasource/pull/300)
- Chore: bump dependencies [#299](https://github.com/grafana/redshift-datasource/pull/299)
- Chore: Update plugin.json keywords [#298](https://github.com/grafana/redshift-datasource/pull/298)
- Add dependabot for grafana/plugin-sdk-go [#296](https://github.com/grafana/redshift-datasource/pull/296)
- Fix: don't check slice nilness before checking length [#294](https://github.com/grafana/redshift-datasource/pull/294)

## 1.18.0

- Add errorsource in [#292](https://github.com/grafana/redshift-datasource/pull/292)

## 1.17.0

- Update grafana/aws-sdk to get new regions

## 1.16.0

- Migrate to new form styling in config and query editors in [#287](https://github.com/grafana/redshift-datasource/pull/287)

## 1.15.2

- Fix: use ReadAuthSettings to get authSettings in [#288](https://github.com/grafana/redshift-datasource/pull/288)

## 1.15.1

- Upgrade grafana-aws-sdk to replace `GetSession` usages with `GetSessionWithAuthSettings` [#284](https://github.com/grafana/redshift-datasource/pull/284)

## 1.15.0

- Add keywords by @kevinwcyu in https://github.com/grafana/redshift-datasource/pull/273
- Add missing regions and use the region resource handler in the frontend by @iwysiu in https://github.com/grafana/redshift-datasource/pull/276
- Plugin.json: update schema reference URL by @leventebalogh in https://github.com/grafana/redshift-datasource/pull/277
- Fix E2E: Update region before sending the /secrets resource request by @idastambuk in https://github.com/grafana/redshift-datasource/pull/280
- Update for added context in grafana-aws-sdk by @njvrzm in https://github.com/grafana/redshift-datasource/pull/279

## New Contributors

- @leventebalogh made their first contribution in https://github.com/grafana/redshift-datasource/pull/277
- @njvrzm made their first contribution in https://github.com/grafana/redshift-datasource/pull/279

## 1.14.0

- Remove the redshiftAsyncQuerySupport feature toggle + styling improvements in https://github.com/grafana/redshift-datasource/pull/272

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
