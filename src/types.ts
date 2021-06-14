import { DataQuery, SelectableValue } from '@grafana/data';
import { AwsAuthDataSourceJsonData, AwsAuthDataSourceSecureJsonData } from '@grafana/aws-sdk';

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

export interface RedshiftQuery extends DataQuery {
  rawSQL: string;
  format: FormatOptions;
}

export const defaultQuery: Partial<RedshiftQuery> = {
  rawSQL: '',
  format: FormatOptions.TimeSeries,
};

/**
 * These are options configured for each DataSource instance
 */
export interface RedshiftDataSourceOptions extends AwsAuthDataSourceJsonData {
  clusterIdentifier?: string;
  database?: string;
  dbUser?: string;
}

/**
 * Value that is used in the backend, but never sent over HTTP to the frontend
 */
export interface RedshiftDataSourceSecureJsonData extends AwsAuthDataSourceSecureJsonData {}
