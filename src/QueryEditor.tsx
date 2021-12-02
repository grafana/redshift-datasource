import React from 'react';
import { QueryEditorProps } from '@grafana/data';
import { DataSource } from './datasource';
import { RedshiftDataSourceOptions, RedshiftQuery, SelectableFormatOptions } from './types';
import { InlineSegmentGroup } from '@grafana/ui';
import { QueryCodeEditor, FormatSelect, QuerySelect, FillValueSelect } from '@grafana/aws-sdk';
import { selectors } from 'selectors';
import { getTemplateSrv } from '@grafana/runtime';
import { getSuggestions } from 'Suggestions';

type Props = QueryEditorProps<DataSource, RedshiftQuery, RedshiftDataSourceOptions>;

export function QueryEditor(props: Props) {
  const fetchSchemas = async () => {
    const schemas: string[] = await props.datasource.getResource('schemas');
    return schemas.map((schema) => ({ label: schema, value: schema })).concat({ label: '-- remove --', value: '' });
  };
  const fetchTables = async () => {
    const tables: string[] = await props.datasource.postResource('tables', {
      schema: props.query.schema || '',
    });
    return tables.map((table) => ({ label: table, value: table })).concat({ label: '-- remove --', value: '' });
  };
  const fetchColumns = async () => {
    const columns: string[] = await props.datasource.postResource('columns', {
      schema: props.query.schema,
      table: props.query.table,
    });
    return columns.map((column) => ({ label: column, value: column })).concat({ label: '-- remove --', value: '' });
  };
  return (
    <>
      <InlineSegmentGroup>
        <div className="gf-form-group">
          <h6>Macros</h6>
          <QuerySelect
            query={props.query}
            queryPropertyPath="schema"
            fetch={fetchSchemas}
            onChange={props.onChange}
            label={selectors.components.ConfigEditor.schema.input}
            data-testid={selectors.components.ConfigEditor.schema.testID}
            tooltip="Use the selected schema with the $__schema macro"
            onRunQuery={props.onRunQuery}
          />
          <QuerySelect
            query={props.query}
            queryPropertyPath="table"
            fetch={fetchTables}
            onChange={props.onChange}
            label={selectors.components.ConfigEditor.table.input}
            data-testid={selectors.components.ConfigEditor.table.testID}
            tooltip="Use the selected table with the $__table macro"
            onRunQuery={props.onRunQuery}
            dependencies={[props.query.schema || '']}
          />
          <QuerySelect
            query={props.query}
            queryPropertyPath="column"
            fetch={fetchColumns}
            onChange={props.onChange}
            label={selectors.components.ConfigEditor.column.input}
            data-testid={selectors.components.ConfigEditor.column.testID}
            tooltip="Use the selected column with the $__column macro"
            onRunQuery={props.onRunQuery}
            dependencies={[props.query.table || '']}
          />
          <h6>Frames</h6>
          <FormatSelect
            query={props.query}
            options={SelectableFormatOptions}
            onChange={props.onChange}
            onRunQuery={props.onRunQuery}
          />
          <FillValueSelect query={props.query} onChange={props.onChange} onRunQuery={props.onRunQuery} />
        </div>
        <div style={{ minWidth: '400px', marginLeft: '10px', flex: 1 }}>
          <QueryCodeEditor
            query={props.query}
            onChange={props.onChange}
            onRunQuery={props.onRunQuery}
            getSuggestions={getSuggestions}
            getTemplateSrv={getTemplateSrv}
          />
        </div>
      </InlineSegmentGroup>
    </>
  );
}
