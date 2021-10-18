import { defaults } from 'lodash';

import React, { useRef, useEffect } from 'react';
import { QueryEditorProps } from '@grafana/data';
import { DataSource } from './datasource';
import { defaultQuery, RedshiftDataSourceOptions, RedshiftQuery } from './types';
import { CodeEditor, CodeEditorSuggestionItem } from '@grafana/ui';
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
  // Use a reference for suggestions because a bug in CodeEditor getSuggestions
  // https://github.com/grafana/grafana/issues/40121
  const suggestionsRef = useRef<CodeEditorSuggestionItem[]>([]);
  useEffect(() => {
    suggestionsRef.current = getSuggestions(
      getTemplateSrv(),
      props.query.schema,
      props.query.table,
      props.query.column
    );
  }, [props.query.schema, props.query.table, props.query.column]);

  return (
    <CodeEditor
      height={'231px'}
      language={'redshift'}
      value={rawSQL}
      onBlur={onRawSqlChange}
      showMiniMap={false}
      showLineNumbers={true}
      getSuggestions={() => suggestionsRef.current}
    />
  );
}
