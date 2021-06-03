import { DataSourcePlugin } from '@grafana/data';
import { DataSource } from './datasource';
import { ConfigEditor } from './ConfigEditor';
import { QueryEditor } from './QueryEditor';
import { RedshiftQuery, RedshiftDataSourceOptions } from './types';

export const plugin = new DataSourcePlugin<DataSource, RedshiftQuery, RedshiftDataSourceOptions>(DataSource)
  .setConfigEditor(ConfigEditor)
  .setQueryEditor(QueryEditor);
