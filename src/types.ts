import { DataQuery } from '@grafana/data';
import { AwsAuthDataSourceJsonData, AwsAuthDataSourceSecureJsonData } from '@grafana/aws-sdk';

export interface RedshiftQuery extends DataQuery {}

export const defaultQuery: Partial<RedshiftQuery> = {};

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
