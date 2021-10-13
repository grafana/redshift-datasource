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
import { InlineField, Select, Input, InlineFieldRow } from '@grafana/ui';
import { QueryCodeEditor } from 'QueryCodeEditor';

type Props = QueryEditorProps<DataSource, RedshiftQuery, RedshiftDataSourceOptions>;

export function QueryEditor(props: Props) {
  const [fillValue, setFillValue] = useState(0);

  const onChange = (value: RedshiftQuery) => {
    props.onChange(value);
    props.onRunQuery();
  };

  const onFillValueChange = ({ currentTarget }: React.FormEvent<HTMLInputElement>) => {
    setFillValue(currentTarget.valueAsNumber);
  };

  const { format, fillMode } = { ...props.query, ...defaultQuery };

  return (
    <>
      <QueryCodeEditor {...props} />
      <InlineField label="Format as">
        <Select
          options={SelectableFormatOptions}
          value={format}
          onChange={({ value }) => onChange({ ...props.query, format: value || FormatOptions.TimeSeries })}
        />
      </InlineField>
      <InlineFieldRow>
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
          <InlineField label="Value">
            <Input
              type="number"
              css=""
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
      </InlineFieldRow>
    </>
  );
}
