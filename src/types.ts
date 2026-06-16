import {
  AwsAuthDataSourceJsonData,
  AwsAuthDataSourceSecureJsonData,
  FillValueOptions,
  SQLQuery,
} from '@grafana/aws-sdk';
import { DataSourceSettings, SelectableValue } from '@grafana/data';

export enum FormatOptions {
  TimeSeries,
  Table,
}

export const SelectableFormatOptions: Array<SelectableValue<FormatOptions>> = [
  {
    label: 'Time Series',
    value: FormatOptions.TimeSeries,
  },
  {
    label: 'Table',
    value: FormatOptions.Table,
  },
];

export interface RedshiftQuery extends SQLQuery {
  format: FormatOptions;

  schema?: string;
  table?: string;
  column?: string;

  queryID?: string;
}

export interface RedshiftManagedSecret {
  name: string;
  arn: string;
}

export const defaultKey = '__default';

export const defaultQuery: Partial<RedshiftQuery> = {
  rawSQL: '',
  format: FormatOptions.Table,
  fillMode: { mode: FillValueOptions.Previous },
};

/**
 * These are options configured for each DataSource instance
 */
export interface RedshiftDataSourceOptions extends AwsAuthDataSourceJsonData {
  withEvent?: boolean;
  useManagedSecret?: boolean;
  useServerless?: boolean;
  workgroupName?: string;
  clusterIdentifier?: string;
  database?: string;
  dbUser?: string;
  managedSecret?: {
    name: string;
    arn: string;
  };
  enableSecureSocksProxy?: boolean;
}

/**
 * Value that is used in the backend, but never sent over HTTP to the frontend
 */
export interface RedshiftDataSourceSecureJsonData extends AwsAuthDataSourceSecureJsonData {}

export type RedshiftDataSourceSettings = DataSourceSettings<
  RedshiftDataSourceOptions,
  RedshiftDataSourceSecureJsonData
>;

export interface RedshiftRunningQueryInfo {
  queryID?: string;
  shouldCancel?: boolean;
}
