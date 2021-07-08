import { defaults } from 'lodash';

import React, { PureComponent } from 'react';
import { QueryEditorProps, SelectableValue } from '@grafana/data';
import { DataSource } from './datasource';
import { defaultQuery, RedshiftDataSourceOptions, RedshiftQuery, SelectableFormatOptions } from './types';
import { CodeEditor, Alert, InlineField, Select, InlineFormLabel, SegmentAsync } from '@grafana/ui';
import { SchemaInfo } from 'SchemaInfo';
import { getTemplateSrv } from '@grafana/runtime';

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

  onSchemaChanged = (value: SelectableValue<string>) => {
    const query = {
      ...this.props.query,
      schema: value.value,
      table: undefined,
      column: undefined,
    };
    if (!query.schema) {
      delete query.schema;
    }
    this.updateSchemaState(query);
  };

  onTableChanged = (value: SelectableValue<string>) => {
    const query = {
      ...this.props.query,
      table: value.value,
      column: undefined,
    };
    if (!query.table) {
      delete query.table;
    }
    this.updateSchemaState(query);
  };

  onColumnChanged = (value: SelectableValue<string>) => {
    const query = {
      ...this.props.query,
      column: value.value,
    };
    if (!query.column) {
      delete query.column;
    }
    this.updateSchemaState(query);
  };

  renderResourceMacro = (
    resource: 'schema' | 'table' | 'column',
    schema: SchemaInfo,
    query?: string,
    value?: string
  ) => {
    let placehoder = '';
    let current = '$__' + resource + ' = ';
    if (query) {
      current += query;
    } else {
      placehoder = current + (value ?? '?');
      current = '';
    }

    let loadOptions;
    let onChange;
    switch (resource) {
      case 'schema':
        loadOptions = schema.getSchemas;
        onChange = this.onSchemaChanged;
        break;
      case 'table':
        loadOptions = schema.getTables;
        onChange = this.onTableChanged;
        break;
      case 'column':
        loadOptions = schema.getColumns;
        onChange = this.onColumnChanged;
        break;
    }
    return (
      <SegmentAsync
        value={current}
        loadOptions={loadOptions}
        placeholder={placehoder}
        onChange={onChange}
        allowCustomValue
      />
    );
  };

  render() {
    const { rawSQL, format } = defaults(this.props.query, defaultQuery);

    const { schema, schemaState } = this.state;
    console.log('foo');
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
              {this.renderResourceMacro('schema', schema, this.props.query.schema, schemaState.schema)}
              {this.renderResourceMacro('table', schema, this.props.query.table, schemaState.table)}
              {this.renderResourceMacro('column', schema, this.props.query.column, schemaState.column)}
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
