import { DataQueryRequest, DataQueryResponse, DataSourceInstanceSettings, ScopedVars } from '@grafana/data';
import { DataSourceWithBackend, getTemplateSrv } from '@grafana/runtime';
import { Observable } from 'rxjs';
import { RedshiftVariableSupport } from 'variables';
import { RedshiftDataSourceOptions, RedshiftQuery } from './types';

export class DataSource extends DataSourceWithBackend<RedshiftQuery, RedshiftDataSourceOptions> {
  constructor(instanceSettings: DataSourceInstanceSettings<RedshiftDataSourceOptions>) {
    super(instanceSettings);
    this.variables = new RedshiftVariableSupport();
  }

  // This will support annotation queries for 7.2+
  annotations = {};

  /**
   * Do not execute queries that do not exist yet
   */
  filterQuery(query: RedshiftQuery): boolean {
    return !!query.rawSQL;
  }

  query(request: DataQueryRequest<RedshiftQuery>): Observable<DataQueryResponse> {
    // What is this about? Due to a bug in the templating query system, data source variables doesn't get assigned ref id.
    // This leads to bad things to therefore we need to assign a dummy value in case it's undefined.
    // The implementation of this method can be removed completely once we upgrade to a version of grafana/data that has this https://github.com/grafana/grafana/pull/35923
    request.targets = request.targets.map((q) => ({ ...q, refId: q.refId ?? 'variable-query' }));
    return super.query(request);
  }

  applyTemplateVariables(query: RedshiftQuery, scopedVars: ScopedVars): RedshiftQuery {
    const templateSrv = getTemplateSrv();
    return {
      ...query,
      rawSQL: templateSrv.replace(query.rawSQL, scopedVars, 'singlequote'),
    };
  }
}
