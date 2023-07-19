import { DataSourcePluginOptionsEditorProps, PluginType } from '@grafana/data';
import { RedshiftDataSourceOptions, RedshiftDataSourceSecureJsonData, RedshiftQuery } from '../types';
import { DataSource } from '../datasource';

export const mockDatasource = new DataSource({
  id: 1,
  uid: 'redshift-id',
  type: 'redshift-datasource',
  name: 'Redshift Data Source',
  jsonData: {},
  access: 'direct',
  readOnly: true,
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

export const mockDatasourceOptions: DataSourcePluginOptionsEditorProps<
  RedshiftDataSourceOptions,
  RedshiftDataSourceSecureJsonData
> = {
  options: {
    id: 1,
    uid: 'redshift-id',
    orgId: 1,
    name: 'Redshift',
    typeLogoUrl: '',
    type: '',
    typeName: '',
    access: '',
    url: '',
    user: '',
    basicAuth: false,
    basicAuthUser: '',
    database: '',
    isDefault: false,
    jsonData: {
      defaultRegion: 'us-east-2',
      useServerless: false,
      useManagedSecret: true,
    },
    secureJsonFields: {},
    readOnly: false,
    withCredentials: false,
  },
  onOptionsChange: jest.fn(),
};

export const mockQuery: RedshiftQuery = { rawSQL: 'select * from foo', refId: '', format: 0, fillMode: { mode: 0 } };
