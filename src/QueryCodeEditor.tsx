import { defaults } from 'lodash';

import React, { useEffect } from 'react';
import { QueryEditorProps } from '@grafana/data';
import { DataSource } from './datasource';
import { defaultQuery, RedshiftDataSourceOptions, RedshiftQuery } from './types';
import { CodeEditor, CodeEditorSuggestionItem, InlineFormLabel } from '@grafana/ui';
import { getTemplateSrv } from '@grafana/runtime';
import ResourceMacro from 'ResourceMacro';
import { getSuggestions } from 'Suggestions';

type Props = QueryEditorProps<DataSource, RedshiftQuery, RedshiftDataSourceOptions>;

// getSuggestions result gets cached so we need to reference a var outside the component
// related issue: https://github.com/grafana/grafana/issues/39264
let suggestions: CodeEditorSuggestionItem[] = [];

export function QueryCodeEditor(props: Props) {
  const onChange = (value: RedshiftQuery) => {
    props.onChange(value);
    props.onRunQuery();
  };

  const onRawSqlChange = (rawSQL: string) => {
    props.onChange({
      ...props.query,
      rawSQL,
    });
    props.onRunQuery();
  };

  const { rawSQL } = defaults(props.query, defaultQuery);

  const loadSchemas = async () => {
    const schemas: string[] = await props.datasource.getResource('schemas');
    return schemas.map((schema) => ({ label: schema, value: schema })).concat({ label: '-- remove --', value: '' });
  };

  const loadTables = async () => {
    const tables: string[] = await props.datasource.postResource('tables', {
      schema: props.query.schema || '',
    });
    return tables.map((table) => ({ label: table, value: table })).concat({ label: '-- remove --', value: '' });
  };

  const loadColumns = async () => {
    const columns: string[] = await props.datasource.postResource('columns', {
      table: props.query.table,
    });
    return columns.map((column) => ({ label: column, value: column })).concat({ label: '-- remove --', value: '' });
  };

  const { table, column } = props.query;
  useEffect(() => {
    suggestions = getSuggestions({ table, column, templateSrv: getTemplateSrv() });
  }, [table, column]);

  return (
    <>
      <div className={'gf-form-inline'}>
        <InlineFormLabel width={8} className="query-keyword">
          Macros
        </InlineFormLabel>
        {ResourceMacro({
          resource: 'schema',
          query: props.query,
          loadOptions: loadSchemas,
          updateQuery: onChange,
        })}
        {ResourceMacro({
          resource: 'table',
          query: props.query,
          loadOptions: loadTables,
          updateQuery: onChange,
        })}
        {ResourceMacro({
          resource: 'column',
          query: props.query,
          loadOptions: loadColumns,
          updateQuery: onChange,
        })}
        <div className="gf-form gf-form--grow">
          <div className="gf-form-label gf-form-label--grow" />
        </div>
      </div>
      <CodeEditor
        height={'250px'}
        language={'redshift'}
        value={rawSQL}
        onBlur={onRawSqlChange}
        showMiniMap={false}
        showLineNumbers={true}
        getSuggestions={() => suggestions}
      />
    </>
  );
}
