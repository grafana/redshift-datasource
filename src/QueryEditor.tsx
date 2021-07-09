import { defaults } from 'lodash';

import React, { PureComponent } from 'react';
import { QueryEditorProps } from '@grafana/data';
import { DataSource } from './datasource';
import { defaultQuery, RedshiftDataSourceOptions, RedshiftQuery, SelectableFormatOptions } from './types';
import { CodeEditor, Alert, InlineField, Select, InlineFormLabel } from '@grafana/ui';
import { SchemaInfo } from 'SchemaInfo';
import { getTemplateSrv } from '@grafana/runtime';
import ResourceMacro from 'ResourceMacro';

type Props = QueryEditorProps<DataSource, RedshiftQuery, RedshiftDataSourceOptions>;

interface State {
  schema?: SchemaInfo;
  schemaState?: Partial<RedshiftQuery>;
}

export class QueryEditor extends PureComponent<Props, State> {
  state: State = {};

  componentDidMount = () => {
    const { datasource, query } = this.props;

    const schema = new SchemaInfo(datasource, query, getTemplateSrv());
    this.setState({ schema: schema, schemaState: schema.state });

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
    const schemaState = this.state.schema!.updateState(query);
    this.setState({ schemaState });

    this.props.onChange(query);
    this.props.onRunQuery();
  };

  render() {
    const { rawSQL, format } = defaults(this.props.query, defaultQuery);

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
            language="redshift"
            value={rawSQL}
            onBlur={this.onRawSqlChange}
            onSave={this.onRawSqlChange}
            showMiniMap={false}
            showLineNumbers={true}
            getSuggestions={schema.getSuggestions}
          />
        )}
        <InlineField label="Format as">
          <Select
            options={SelectableFormatOptions}
            value={format}
            onChange={({ value }) => this.onChange({ ...this.props.query, format: value! })}
          />
        </InlineField>
      </>
    );
  }
}
