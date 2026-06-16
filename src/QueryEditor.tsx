import React from 'react';
import { QueryEditorProps } from '@grafana/data';
import { QueryEditorHeader } from '@grafana/aws-sdk';
import { RedshiftDataSourceOptions, RedshiftQuery } from './types';
import { DataSource } from './datasource';
import { QueryEditorForm } from './QueryEditorForm';

export function QueryEditor(props: QueryEditorProps<DataSource, RedshiftQuery, RedshiftDataSourceOptions>) {
  return (
    <>
      {props?.app !== 'explore' && (
        <QueryEditorHeader<DataSource, RedshiftQuery, RedshiftDataSourceOptions>
          {...props}
          enableRunButton={!!props.query.rawSQL}
          showAsyncQueryButtons
          cancel={props.datasource.cancel}
        />
      )}
      <QueryEditorForm {...props} />
    </>
  );
}
