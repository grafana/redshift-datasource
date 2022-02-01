import { DataSourceSettings, SelectableValue } from '@grafana/data';
import { AwsAuthDataSourceJsonData, AwsAuthDataSourceSecureJsonData, SQLQuery } from '@grafana/aws-sdk';
import { FillValueOptions } from '@grafana/aws-sdk/dist/sql/QueryEditor/FillValueSelect';

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
