import { DataSourceVariableSupport, VariableSupportType } from '@grafana/data';

import { DataSource } from './datasource';
import { RedshiftQuery } from './types';

export class RedshiftVariableSupport extends DataSourceVariableSupport<DataSource, RedshiftQuery> {
  constructor() {
    super();
  }

  getType() {
    return VariableSupportType.Datasource;
  }
}
