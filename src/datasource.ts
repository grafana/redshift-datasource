import { DataSourceInstanceSettings } from '@grafana/data';
import { DataSourceWithBackend } from '@grafana/runtime';
import { RedshiftDataSourceOptions, RedshiftQuery } from './types';

export class DataSource extends DataSourceWithBackend<RedshiftQuery, RedshiftDataSourceOptions> {
  constructor(instanceSettings: DataSourceInstanceSettings<RedshiftDataSourceOptions>) {
    super(instanceSettings);
  }
}
