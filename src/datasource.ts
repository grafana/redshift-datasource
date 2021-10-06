import { DataSourceInstanceSettings, ScopedVars } from '@grafana/data';
import { DataSourceWithBackend, getTemplateSrv } from '@grafana/runtime';
import { RedshiftVariableSupport } from 'variables';
import { RedshiftDataSourceOptions, RedshiftQuery } from './types';

export class DataSource extends DataSourceWithBackend<RedshiftQuery, RedshiftDataSourceOptions> {
  constructor(instanceSettings: DataSourceInstanceSettings<RedshiftDataSourceOptions>) {
    super(instanceSettings);
    this.variables = new RedshiftVariableSupport(this);
  }

  // This will support annotation queries for 7.2+
  annotations = {};

  /**
   * Do not execute queries that do not exist yet
   */
  filterQuery(query: RedshiftQuery): boolean {
    return !!query.rawSQL;
  }

  applyTemplateVariables(query: RedshiftQuery, scopedVars: ScopedVars): RedshiftQuery {
    const templateSrv = getTemplateSrv();
    return {
      ...query,
      rawSQL: templateSrv.replace(query.rawSQL, scopedVars, 'singlequote'),
    };
  }
}
