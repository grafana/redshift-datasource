import { DataSourceSettings, SelectableValue } from '@grafana/data';
import { AwsAuthDataSourceJsonData, AwsAuthDataSourceSecureJsonData, SQLQuery } from '@grafana/aws-sdk';

export enum FormatOptions {
  TimeSeries,
  Table,
}

export enum FillValueOptions {
  Previous,
  Null,
  Value,
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

export const SelectableFillValueOptions: Array<SelectableValue<FillValueOptions>> = [
  {
    label: 'Previous Value',
    value: FillValueOptions.Previous,
  },
  {
    label: 'NULL',
    value: FillValueOptions.Null,
  },
  {
    label: 'Value',
    value: FillValueOptions.Value,
  },
];

export interface RedshiftQuery extends SQLQuery {
  format: FormatOptions;

  schema?: string;
  table?: string;
  column?: string;
}

export interface RedshiftManagedSecret {
  name: string;
  arn: string;
}

export const defaultKey = '__default';

export const defaultQuery: Partial<RedshiftQuery> = {
  rawSQL: '',
  format: FormatOptions.TimeSeries,
  fillMode: { mode: FillValueOptions.Previous },
};

/**
 * These are options configured for each DataSource instance
 */
export interface RedshiftDataSourceOptions extends AwsAuthDataSourceJsonData {
  useManagedSecret?: boolean;
  clusterIdentifier?: string;
  database?: string;
  dbUser?: string;
  managedSecret?: {
    name: string;
    arn: string;
  };
}

/**
 * Value that is used in the backend, but never sent over HTTP to the frontend
 */
export interface RedshiftDataSourceSecureJsonData extends AwsAuthDataSourceSecureJsonData {}

export type RedshiftDataSourceSettings = DataSourceSettings<
  RedshiftDataSourceOptions,
  RedshiftDataSourceSecureJsonData
>;
