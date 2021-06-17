import { DataSourceInstanceSettings, ScopedVars } from '@grafana/data';
import { DataSourceWithBackend, getTemplateSrv } from '@grafana/runtime';
import { RedshiftDataSourceOptions, RedshiftQuery } from './types';

export class DataSource extends DataSourceWithBackend<RedshiftQuery, RedshiftDataSourceOptions> {
  constructor(instanceSettings: DataSourceInstanceSettings<RedshiftDataSourceOptions>) {
    super(instanceSettings);
  }

  applyTemplateVariables(query: RedshiftQuery, scopedVars: ScopedVars): RedshiftQuery {
    const templateSrv = getTemplateSrv();

    return {
      ...query,
      rawSQL: templateSrv.replace(query.rawSQL, scopedVars, 'singlequote'),
    };
  }
}
