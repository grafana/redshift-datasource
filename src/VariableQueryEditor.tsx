import React from 'react';
import { RedshiftQuery, RedshiftDataSourceOptions } from './types';
import { QueryEditorProps } from '@grafana/data';
import { DataSource } from 'datasource';
import { QueryEditorForm } from 'QueryEditorForm';

export function VariableQueryCodeEditor(props: QueryEditorProps<DataSource, RedshiftQuery, RedshiftDataSourceOptions>) {
  return <QueryEditorForm {...props}></QueryEditorForm>;
}
