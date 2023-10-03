import { FillValueSelect, FormatSelect, ResourceSelector } from '@grafana/aws-sdk';
import { QueryEditorProps, SelectableValue } from '@grafana/data';
import React from 'react';
import { selectors } from 'selectors';
import SQLEditor from './SQLEditor';

import { DataSource } from './datasource';
import { FormatOptions, RedshiftDataSourceOptions, RedshiftQuery, SelectableFormatOptions } from './types';
import { EditorField, EditorFieldGroup, EditorRow, EditorRows } from '@grafana/experimental';

type Props = QueryEditorProps<DataSource, RedshiftQuery, RedshiftDataSourceOptions>;

type QueryProperties = 'schema' | 'table' | 'column';

export function QueryEditorForm(props: Props) {
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

  const onChange = (prop: QueryProperties) => (e: SelectableValue<string> | null) => {
    const newQuery = { ...props.query };
    const value = e?.value;
    newQuery[prop] = value;
    props.onChange(newQuery);
  };

  return (
    <EditorRows>
      <EditorRow>
        <EditorFieldGroup>
          <EditorField
            width={15}
            label={selectors.components.ConfigEditor.schema.input}
            tooltip="Use the selected schema with the $__schema macro"
            data-testid={selectors.components.ConfigEditor.schema.testID}
            htmlFor="schema"
          >
            <ResourceSelector
              id="schema"
              aria-label={selectors.components.ConfigEditor.schema.input}
              onChange={onChange('schema')}
              fetch={fetchSchemas}
              value={props.query.schema || null}
            />
          </EditorField>
          <EditorField
            width={15}
            label={selectors.components.ConfigEditor.table.input}
            tooltip="Use the selected table with the $__table macro"
            data-testid={selectors.components.ConfigEditor.table.testID}
            htmlFor="table"
          >
            <ResourceSelector
              id="table"
              aria-label={selectors.components.ConfigEditor.table.input}
              onChange={onChange('table')}
              fetch={fetchTables}
              value={props.query.table || null}
              dependencies={[props.query.schema]}
            />
          </EditorField>
          <EditorField
            width={15}
            label={selectors.components.ConfigEditor.column.input}
            tooltip="Use the selected column with the $__column macro"
            data-testid={selectors.components.ConfigEditor.column.testID}
            htmlFor="column"
          >
            <ResourceSelector
              id="column"
              aria-label={selectors.components.ConfigEditor.column.input}
              onChange={onChange('column')}
              fetch={fetchColumns}
              value={props.query.column || null}
              dependencies={[props.query.table]}
            />
          </EditorField>
        </EditorFieldGroup>
        <EditorFieldGroup>
          <EditorField label="Format frames as" htmlFor="formatAs" width={20}>
            <FormatSelect
              id="formatAs"
              query={props.query}
              options={SelectableFormatOptions}
              onChange={props.onChange}
            />
          </EditorField>
          {props.query.format === FormatOptions.TimeSeries && (
            <FillValueSelect query={props.query} onChange={props.onChange} />
          )}
        </EditorFieldGroup>
      </EditorRow>
      <EditorRow>
        <div style={{ width: '100%' }}>
          <SQLEditor query={props.query} onChange={props.onChange} datasource={props.datasource} />
        </div>
      </EditorRow>
    </EditorRows>
  );
}
