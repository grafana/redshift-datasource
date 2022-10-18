import React from 'react';
import { RedshiftQuery, RedshiftDataSourceOptions } from './types';
import { QueryEditorProps } from '@grafana/data';
import { DataSource } from 'datasource';
import { QueryEditor } from 'QueryEditor';

export function VariableQueryCodeEditor(props: QueryEditorProps<DataSource, RedshiftQuery, RedshiftDataSourceOptions>) {
  return <QueryEditor {...props} hideRunQueryButtons></QueryEditor>;
}
