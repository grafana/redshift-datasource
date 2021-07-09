import { defaults } from 'lodash';

import React from 'react';
import { QueryEditorProps } from '@grafana/data';
import { DataSource } from './datasource';
import { defaultQuery, RedshiftDataSourceOptions, RedshiftQuery, SelectableFormatOptions } from './types';
import { CodeEditor, Alert, InlineField, Select } from '@grafana/ui';
import { SchemaInfo } from 'SchemaInfo';
import { getTemplateSrv } from '@grafana/runtime';

type Props = QueryEditorProps<DataSource, RedshiftQuery, RedshiftDataSourceOptions>;

export function QueryEditor(props: Props) {
  const { rawSQL, format } = defaults(props.query, defaultQuery);
  const schema = new SchemaInfo(getTemplateSrv());

  const onChange = (value: RedshiftQuery) => {
    props.onChange(value);
    props.onRunQuery();
  };

  const onRawSqlChange = (rawSQL: string) => {
    const { onChange, query, onRunQuery } = props;
    onChange({ ...query, rawSQL });
    onRunQuery();
  };

  return (
    <>
      <Alert title="" severity="info">
        To save and re-run the query, press ctrl/cmd+S.
      </Alert>
      <CodeEditor
        height={'250px'}
        // TODO: Use language="redshift" once Grafana v7.x is deprecated.
        language="sql"
        value={rawSQL}
        onBlur={onRawSqlChange}
        onSave={onRawSqlChange}
        showMiniMap={false}
        showLineNumbers={true}
        getSuggestions={schema.getSuggestions}
      />
      <InlineField label="Format as">
        <Select
          options={SelectableFormatOptions}
          value={format}
          onChange={({ value }) => onChange({ ...props.query, format: value! })}
        />
      </InlineField>
    </>
  );
}
