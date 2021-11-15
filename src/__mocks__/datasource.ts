import { DataSourcePluginOptionsEditorProps, PluginType } from '@grafana/data';
import { RedshiftDataSourceOptions, RedshiftDataSourceSecureJsonData, RedshiftQuery } from '../types';
import { DataSource } from '../datasource';

export const mockDatasource = new DataSource({
  id: 1,
  uid: 'redshift-id',
  type: 'redshift-datasource',
  name: 'Redshift Data Source',
  access: 'proxy',
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
    password: '',
    user: '',
    basicAuth: false,
    basicAuthPassword: '',
    basicAuthUser: '',
    database: '',
    isDefault: false,
    jsonData: {
      defaultRegion: 'us-east-2',
    },
    secureJsonFields: {},
    readOnly: false,
    withCredentials: false,
  },
  onOptionsChange: jest.fn(),
};

export const mockQuery: RedshiftQuery = { rawSQL: 'select * from foo', refId: '', format: 0, fillMode: { mode: 0 } };
