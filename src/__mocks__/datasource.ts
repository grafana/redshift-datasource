import { PluginType } from '@grafana/data';
import { RedshiftQuery } from '../types';
import { DataSource } from '../datasource';

export const mockDatasource = new DataSource({
  id: 1,
  uid: 'redshift-id',
  type: 'redshift-datasource',
  name: 'Redshift Data Source',
  jsonData: {},
  meta: {
    id: 'redshift-datasource',
    name: 'Redshift Data Source',
    type: PluginType.datasource,
    module: '',
    baseUrl: '',
    info: {
      description: '',
      screenshots: [],
      updated: '',
      version: '',
      logos: {
        small: '',
        large: '',
      },
      author: {
        name: '',
      },
      links: [],
    },
  },
});

export const mockQuery: RedshiftQuery = { rawSQL: 'select * from foo', refId: '', format: 0 };
