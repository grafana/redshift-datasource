import React from 'react';
import { QueryCodeEditor } from '@grafana/aws-sdk';
import { getSuggestions } from 'Suggestions';
import { getTemplateSrv } from '@grafana/runtime';
import { RedshiftQuery, RedshiftDataSourceOptions } from './types';
import { QueryEditorProps } from '@grafana/data';
import { DataSource } from 'datasource';

export function VariableQueryCodeEditor(props: QueryEditorProps<DataSource, RedshiftQuery, RedshiftDataSourceOptions>) {
  return <QueryCodeEditor {...props} getSuggestions={getSuggestions} getTemplateSrv={getTemplateSrv} />;
}
