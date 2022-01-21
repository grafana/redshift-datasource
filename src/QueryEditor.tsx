import { FillValueSelect, FormatSelect, QueryCodeEditor, ResourceSelector } from '@grafana/aws-sdk';
import { QueryEditorProps, SelectableValue } from '@grafana/data';
import { InlineSegmentGroup } from '@grafana/ui';
import React from 'react';
import { selectors } from 'selectors';
import { getSuggestions } from 'Suggestions';

import { DataSource } from './datasource';
import { FormatOptions, RedshiftDataSourceOptions, RedshiftQuery, SelectableFormatOptions } from './types';

type Props = QueryEditorProps<DataSource, RedshiftQuery, RedshiftDataSourceOptions>;

type QueryProperties = 'schema' | 'table' | 'column';

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

  const onChange = (prop: QueryProperties) => (e: SelectableValue<string> | null) => {
    const newQuery = { ...props.query };
    const value = e?.value;
    newQuery[prop] = value;
    props.onChange(newQuery);
    if (props.onRunQuery) {
      props.onRunQuery();
    }
  };

  return (
    <>
      <InlineSegmentGroup>
        <div className="gf-form-group">
          <h6>Macros</h6>
          <ResourceSelector
            onChange={onChange('schema')}
            fetch={fetchSchemas}
            value={props.query.schema || null}
            tooltip="Use the selected schema with the $__schema macro"
            label={selectors.components.ConfigEditor.schema.input}
            data-testid={selectors.components.ConfigEditor.schema.testID}
            labelWidth={11}
            className="width-12"
          />
          <ResourceSelector
            onChange={onChange('table')}
            fetch={fetchTables}
            value={props.query.table || null}
            dependencies={[props.query.schema]}
            tooltip="Use the selected table with the $__table macro"
            label={selectors.components.ConfigEditor.table.input}
            data-testid={selectors.components.ConfigEditor.table.testID}
            labelWidth={11}
            className="width-12"
          />
          <ResourceSelector
            onChange={onChange('column')}
            fetch={fetchColumns}
            value={props.query.column || null}
            dependencies={[props.query.table]}
            tooltip="Use the selected column with the $__column macro"
            label={selectors.components.ConfigEditor.column.input}
            data-testid={selectors.components.ConfigEditor.column.testID}
            labelWidth={11}
            className="width-12"
          />
          <h6>Frames</h6>
          <FormatSelect
            query={props.query}
            options={SelectableFormatOptions}
            onChange={props.onChange}
            onRunQuery={props.onRunQuery}
          />
          {props.query.format === FormatOptions.TimeSeries && (
            <FillValueSelect query={props.query} onChange={props.onChange} onRunQuery={props.onRunQuery} />
          )}
        </div>
        <div style={{ minWidth: '400px', marginLeft: '10px', flex: 1 }}>
          <QueryCodeEditor
            language="redshift"
            query={props.query}
            onChange={props.onChange}
            onRunQuery={props.onRunQuery}
            getSuggestions={getSuggestions}
          />
        </div>
      </InlineSegmentGroup>
    </>
  );
}
