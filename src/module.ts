import { DataSourcePlugin } from '@grafana/data';

import { ConfigEditor } from './ConfigEditor';
import { DataSource } from './datasource';
import { QueryEditor } from './QueryEditor';
import { RedshiftDataSourceOptions, RedshiftQuery } from './types';

export const plugin = new DataSourcePlugin<DataSource, RedshiftQuery, RedshiftDataSourceOptions>(DataSource)
  .setConfigEditor(ConfigEditor)
  .setQueryEditor(QueryEditor);
