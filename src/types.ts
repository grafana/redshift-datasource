import { DataQuery, SelectableValue } from '@grafana/data';
import { AwsAuthDataSourceJsonData, AwsAuthDataSourceSecureJsonData } from '@grafana/aws-sdk';

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

export interface RedshiftQuery extends DataQuery {
  rawSQL: string;
  format: FormatOptions;
  fillMode: { mode: FillValueOptions; value?: number };

  schema?: string;
  table?: string;
  column?: string;
}

export const defaultQuery: Partial<RedshiftQuery> = {
  rawSQL: '',
  format: FormatOptions.TimeSeries,
  fillMode: { mode: FillValueOptions.Previous },
};

/**
 * These are options configured for each DataSource instance
 */
export interface RedshiftDataSourceOptions extends AwsAuthDataSourceJsonData {
  clusterIdentifier?: string;
  database?: string;
  dbUser?: string;
  managedSecret?: string;
}

/**
 * Value that is used in the backend, but never sent over HTTP to the frontend
 */
export interface RedshiftDataSourceSecureJsonData extends AwsAuthDataSourceSecureJsonData {}
