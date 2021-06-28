import { defaults } from 'lodash';

import React, { PureComponent } from 'react';
import { QueryEditorProps } from '@grafana/data';
import { DataSource } from './datasource';
import { defaultQuery, RedshiftDataSourceOptions, RedshiftQuery, SelectableFormatOptions } from './types';
import { CodeEditor, Alert, InlineField, Select } from '@grafana/ui';

type Props = QueryEditorProps<DataSource, RedshiftQuery, RedshiftDataSourceOptions>;

export class QueryEditor extends PureComponent<Props> {
  onRawSqlChange = (rawSQL: string) => {
    const { onChange, query, onRunQuery } = this.props;
    onChange({ ...query, rawSQL });
    onRunQuery();
  };

  onChange = (value: RedshiftQuery) => {
    this.props.onChange(value);
    this.props.onRunQuery();
  };

  render() {
    const { rawSQL, format } = defaults(this.props.query, defaultQuery);

    return (
      <>
        <Alert title="" severity="info">To save and re-run the query, press ctrl/cmd+S.</InfoBox>
        <CodeEditor
          height={'250px'}
          language="redshift"
          value={rawSQL}
          onBlur={this.onRawSqlChange}
          onSave={this.onRawSqlChange}
          showMiniMap={false}
          showLineNumbers={true}
        />
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
