import React, { useState } from 'react';
import { QueryEditorProps } from '@grafana/data';
import { DataSource } from './datasource';
import {
  defaultQuery,
  FormatOptions,
  RedshiftDataSourceOptions,
  RedshiftQuery,
  SelectableFormatOptions,
  SelectableFillValueOptions,
  FillValueOptions,
} from './types';
import { InlineField, Select, Input, InlineSegmentGroup } from '@grafana/ui';
import { QueryCodeEditor } from 'QueryCodeEditor';
import { ResourceSelector } from 'ConfigEditor/ResourceSelector';

type Props = QueryEditorProps<DataSource, RedshiftQuery, RedshiftDataSourceOptions>;

export function QueryEditor(props: Props) {
  const queryWithDefaults = {
    ...defaultQuery,
    ...props.query,
  };
  const { format, fillMode } = { ...defaultQuery, ...props.query };
  const [fillValue, setFillValue] = useState(fillMode.value || 0);

  const onChange = (value: RedshiftQuery) => {
    props.onChange(value);
    props.onRunQuery();
  };

  const onFillValueChange = ({ currentTarget }: React.FormEvent<HTMLInputElement>) => {
    setFillValue(currentTarget.valueAsNumber);
  };

  // Schema selector
  const [schema, setSchema] = useState<string | undefined>(queryWithDefaults.schema);
  const fetchSchemas = async () => {
    const schemas: string[] = await props.datasource.getResource('schemas');
    return schemas.map((schema) => ({ label: schema, value: schema })).concat({ label: '-- remove --', value: '' });
  };
  const onSchemaChange = (newSchema?: string) => {
    setSchema(newSchema);
    props.onChange({
      ...queryWithDefaults,
      schema: newSchema,
      table: undefined,
      column: undefined,
    });
    props.onRunQuery();
  };

  // Tables selector
  const [table, setTable] = useState<string | undefined>(queryWithDefaults.table);
  const fetchTables = async () => {
    const tables: string[] = await props.datasource.postResource('tables', {
      schema: props.query.schema || '',
    });
    return tables.map((table) => ({ label: table, value: table })).concat({ label: '-- remove --', value: '' });
  };
  const onTableChange = (newTable?: string) => {
    setTable(newTable);
    props.onChange({
      ...queryWithDefaults,
      table: newTable || undefined,
      column: undefined,
    });
    props.onRunQuery();
  };

  // Columns selector
  const [column, setColumn] = useState<string | undefined>(queryWithDefaults.column);
  const fetchColumns = async () => {
    const columns: string[] = await props.datasource.postResource('columns', {
      schema: props.query.schema,
      table: props.query.table,
    });
    return columns.map((column) => ({ label: column, value: column })).concat({ label: '-- remove --', value: '' });
  };
  const onColumnChange = (newColumn?: string) => {
    setColumn(newColumn);
    props.onChange({
      ...queryWithDefaults,
      column: newColumn,
    });
    props.onRunQuery();
  };

  return (
    <>
      <InlineSegmentGroup>
        <div className="gf-form-group">
          <h6>Macros</h6>
          <ResourceSelector
            resource="schema"
            value={schema || null}
            fetch={fetchSchemas}
            onChange={(e) => onSchemaChange(e?.value)}
            tooltip="Use the selected schema with the $__schema macro"
            labelWidth={11}
            className="width-12"
          />
          <ResourceSelector
            resource="table"
            value={table || null}
            fetch={fetchTables}
            onChange={(e) => onTableChange(e?.value)}
            dependencies={[schema]}
            tooltip="Use the selected table with the $__table macro"
            labelWidth={11}
            className="width-12"
          />
          <ResourceSelector
            resource="column"
            value={column || null}
            fetch={fetchColumns}
            onChange={(e) => onColumnChange(e?.value)}
            // TODO: Add schema as dependency
            dependencies={[table]}
            tooltip="Use the selected column with the $__column macro"
            labelWidth={11}
            className="width-12"
          />
          <h6>Frames</h6>
          <InlineField label="Format as" labelWidth={11}>
            <Select
              options={SelectableFormatOptions}
              value={format}
              onChange={({ value }) => onChange({ ...props.query, format: value || FormatOptions.TimeSeries })}
            />
          </InlineField>
          <InlineField label="Fill value" tooltip="value to fill missing points">
            <Select
              aria-label="Fill value"
              options={SelectableFillValueOptions}
              value={fillMode.mode}
              onChange={({ value }) =>
                onChange({
                  ...props.query,
                  fillMode: { mode: value || FillValueOptions.Previous, value: fillValue },
                })
              }
            />
          </InlineField>
          {fillMode.mode === FillValueOptions.Value && (
            <InlineField label="Value" labelWidth={11}>
              <Input
                type="number"
                css
                value={fillValue}
                onChange={onFillValueChange}
                onBlur={() =>
                  onChange({
                    ...props.query,
                    fillMode: { mode: FillValueOptions.Value, value: fillValue },
                  })
                }
              />
            </InlineField>
          )}
        </div>
        <div style={{ minWidth: '400px', marginLeft: '10px', flex: 1 }}>
          <QueryCodeEditor {...props} />
        </div>
      </InlineSegmentGroup>
    </>
  );
}
