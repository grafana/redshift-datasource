import { DataQuery, DataSourceJsonData } from '@grafana/data';

export interface RedshiftQuery extends DataQuery {}

export const defaultQuery: Partial<RedshiftQuery> = {};

/**
 * These are options configured for each DataSource instance
 */
export interface RedshiftDataSourceOptions extends DataSourceJsonData {
  path?: string;
}

/**
 * Value that is used in the backend, but never sent over HTTP to the frontend
 */
export interface RedshiftDataSourceSecureJsonData {}
