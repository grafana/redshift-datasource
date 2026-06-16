import { DataQueryRequest, DataQueryResponse, CustomVariableSupport } from '@grafana/data';
import { assign } from 'lodash';
import { Observable } from 'rxjs';
import { VariableQueryCodeEditor } from 'VariableQueryEditor';
import { DataSource } from './datasource';
import { RedshiftQuery, defaultQuery } from './types';

export class RedshiftVariableSupport extends CustomVariableSupport<DataSource, RedshiftQuery, RedshiftQuery> {
  constructor(private readonly datasource: DataSource) {
    super();
    this.datasource = datasource;
    this.query = this.query.bind(this);
  }

  editor = VariableQueryCodeEditor;

  query(request: DataQueryRequest<RedshiftQuery>): Observable<DataQueryResponse> {
    // fill query params with default data
    assign(request.targets, [{ ...defaultQuery, ...request.targets[0], refId: 'A' }]);
    return this.datasource.query(request);
  }
}
