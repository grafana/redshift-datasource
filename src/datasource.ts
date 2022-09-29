import { applySQLTemplateVariables, filterSQLQuery } from '@grafana/aws-sdk';
import { DataSourceInstanceSettings, ScopedVars } from '@grafana/data';
import { getTemplateSrv, config } from '@grafana/runtime';
import { RedshiftVariableSupport } from 'variables';
import { AsyncDatasourceWithBackend } from './AsyncDatasourceWithBackend';

import { RedshiftDataSourceOptions, RedshiftQuery } from './types';

export class DataSource extends AsyncDatasourceWithBackend<RedshiftQuery, RedshiftDataSourceOptions> {
  constructor(instanceSettings: DataSourceInstanceSettings<RedshiftDataSourceOptions>) {
    super(instanceSettings, config.featureToggles.redshiftAsyncQueryDataSupport);
    this.variables = new RedshiftVariableSupport(this);
  }

  // This will support annotation queries for 7.2+
  annotations = {};

  filterQuery = filterSQLQuery;

  applyTemplateVariables = (query: RedshiftQuery, scopedVars: ScopedVars) =>
    applySQLTemplateVariables(query, scopedVars, getTemplateSrv);
}
