import { defaults } from 'lodash';

import React from 'react';
import { QueryEditorProps } from '@grafana/data';
import { DataSource } from './datasource';
import { defaultQuery, RedshiftDataSourceOptions, RedshiftQuery } from './types';
import { CodeEditor } from '@grafana/ui';
import { getTemplateSrv } from '@grafana/runtime';
import { getSuggestions } from 'Suggestions';

type Props = QueryEditorProps<DataSource, RedshiftQuery, RedshiftDataSourceOptions>;

export function QueryCodeEditor(props: Props) {
  const onRawSqlChange = (rawSQL: string) => {
    props.onChange({
      ...props.query,
      rawSQL,
    });
    props.onRunQuery();
  };

  const { rawSQL } = defaults(props.query, defaultQuery);

  return (
    <CodeEditor
      height={'231px'}
      language={'redshift'}
      value={rawSQL}
      onBlur={onRawSqlChange}
      showMiniMap={false}
      showLineNumbers={true}
      getSuggestions={() => getSuggestions({ query: props.query, templateSrv: getTemplateSrv() })}
    />
  );
}
