import { applySQLTemplateVariables, filterSQLQuery } from '@grafana/aws-sdk';
import { DatasourceWithAsyncBackend } from '@grafana/async-query-data';
import { DataSourceInstanceSettings, ScopedVars } from '@grafana/data';
import { getTemplateSrv } from '@grafana/runtime';
import { RedshiftVariableSupport } from 'variables';

import { RedshiftDataSourceOptions, RedshiftQuery, defaultQuery } from './types';
import { RedshiftAnnotationsSupport } from './annotations';

export class DataSource extends DatasourceWithAsyncBackend<RedshiftQuery, RedshiftDataSourceOptions> {
  constructor(instanceSettings: DataSourceInstanceSettings<RedshiftDataSourceOptions>) {
    super(instanceSettings);
    this.variables = new RedshiftVariableSupport(this);
  }

  getDefaultQuery(): Partial<RedshiftQuery> {
    return defaultQuery;
  }

  // This will support annotation queries for 7.2+
  annotations = RedshiftAnnotationsSupport;

  filterQuery = filterSQLQuery;

  applyTemplateVariables = (query: RedshiftQuery, scopedVars: ScopedVars) =>
    applySQLTemplateVariables(query, scopedVars, getTemplateSrv);
}
