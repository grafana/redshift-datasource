import { defaults } from 'lodash';

import React, { PureComponent } from 'react';
import { QueryEditorProps } from '@grafana/data';
import { DataSource } from './datasource';
import { defaultQuery, RedshiftDataSourceOptions, RedshiftQuery, SelectableFormatOptions } from './types';
import { CodeEditor, InlineField, Select } from '@grafana/ui';

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
    const { onChange } = this.props;
    const { rawSQL, format } = defaults(this.props.query, defaultQuery);

    return (
      <>
        <CodeEditor
          height={'250px'}
          language="redshift"
          value={rawSQL || ''}
          onBlur={this.onRawSqlChange}
          onSave={this.onRawSqlChange}
          showMiniMap={false}
          showLineNumbers={true}
        />
        <InlineField label="Format as">
          <Select
            options={SelectableFormatOptions}
            value={format}
            onChange={({ value }) => onChange({ ...this.props.query, format: value! })}
          />
        </InlineField>
      </>
    );
  }
}
