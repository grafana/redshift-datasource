import { defaults } from 'lodash';

import React, { PureComponent } from 'react';
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
import { CodeEditor, Alert, InlineField, Select, InlineFormLabel, Input, InlineFieldRow } from '@grafana/ui';
import { SchemaInfo } from 'SchemaInfo';
import { getTemplateSrv } from '@grafana/runtime';
import ResourceMacro from 'ResourceMacro';

type Props = QueryEditorProps<DataSource, RedshiftQuery, RedshiftDataSourceOptions>;

interface State {
  schema: SchemaInfo;
  schemaState?: Partial<RedshiftQuery>;
  fillValue: number;
}

export class QueryEditor extends PureComponent<Props, State> {
  state: State = {
    schema: new SchemaInfo(this.props.datasource, this.props.query, getTemplateSrv()),
    fillValue: 0,
  };

  componentDidMount = () => {
    const { schema } = this.state;
    this.setState({ schemaState: schema.state });
    schema.preload();
  };

  onChange = (value: RedshiftQuery) => {
    this.props.onChange(value);
    this.props.onRunQuery();
  };

  onRawSqlChange = (rawSQL: string) => {
    this.props.onChange({
      ...this.props.query,
      rawSQL,
    });
    this.props.onRunQuery();
  };

  updateSchemaState = (query: RedshiftQuery) => {
    const schemaState = this.state.schema.updateState(query);
    this.setState({ schemaState });

    this.props.onChange(query);
    this.props.onRunQuery();
  };

  isPanelEditor = () => {
    // If there can be more than one query, it's a panel editor
    return !!this.props.queries;
  };

  onFillValueChange = ({ currentTarget }: React.FormEvent<HTMLInputElement>) => {
    this.setState({ fillValue: currentTarget.valueAsNumber });
  };

  render() {
    const { rawSQL, format, fillMode } = defaults(this.props.query, defaultQuery);

    const { schema, schemaState } = this.state;
    return (
      <>
        <Alert title="" severity="info">
          To save and re-run the query, press ctrl/cmd+S.
        </Alert>
        <div className={'gf-form-inline'}>
          <InlineFormLabel width={8} className="query-keyword">
            Macros
          </InlineFormLabel>
          {schema && schemaState && (
            <>
              {ResourceMacro({
                resource: 'schema',
                schema,
                query: this.props.query,
                value: schemaState.schema,
                updateSchemaState: this.updateSchemaState,
              })}
              {ResourceMacro({
                resource: 'table',
                schema,
                query: this.props.query,
                value: schemaState.table,
                updateSchemaState: this.updateSchemaState,
              })}
              {ResourceMacro({
                resource: 'column',
                schema,
                query: this.props.query,
                value: schemaState.column,
                updateSchemaState: this.updateSchemaState,
              })}
            </>
          )}
          <div className="gf-form gf-form--grow">
            <div className="gf-form-label gf-form-label--grow" />
          </div>
        </div>
        {schema && (
          <CodeEditor
            height={'250px'}
            language={'redshift'}
            value={rawSQL}
            onBlur={this.onRawSqlChange}
            onSave={this.onRawSqlChange}
            showMiniMap={false}
            showLineNumbers={true}
            getSuggestions={schema.getSuggestions}
          />
        )}
        {this.isPanelEditor() && (
          <>
            <InlineField label="Format as">
              <Select
                options={SelectableFormatOptions}
                value={format}
                onChange={({ value }) =>
                  this.onChange({ ...this.props.query, format: value || FormatOptions.TimeSeries })
                }
              />
            </InlineField>
            <InlineFieldRow>
              <InlineField label="Fill value" tooltip="value to fill missing points">
                <Select
                  aria-label="Fill value"
                  options={SelectableFillValueOptions}
                  value={fillMode.mode}
                  onChange={({ value }) =>
                    this.onChange({
                      ...this.props.query,
                      fillMode: { mode: value || FillValueOptions.Previous, value: this.state.fillValue },
                    })
                  }
                />
              </InlineField>
              {fillMode.mode === FillValueOptions.Value && (
                <InlineField label="Value">
                  <Input
                    type="number"
                    css=""
                    value={this.state.fillValue}
                    onChange={this.onFillValueChange}
                    onBlur={() =>
                      this.onChange({
                        ...this.props.query,
                        fillMode: { mode: FillValueOptions.Value, value: this.state.fillValue },
                      })
                    }
                  />
                </InlineField>
              )}
            </InlineFieldRow>
          </>
        )}
      </>
    );
  }
}
