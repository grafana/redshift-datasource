import { defaults } from 'lodash';

import React, { PureComponent } from 'react';
import { QueryEditorProps } from '@grafana/data';
import { DataSource } from './datasource';
import { defaultQuery, RedshiftDataSourceOptions, RedshiftQuery, SelectableFormatOptions } from './types';
import { Select, TextArea } from '@grafana/ui';

type Props = QueryEditorProps<DataSource, RedshiftQuery, RedshiftDataSourceOptions>;

export class QueryEditor extends PureComponent<Props> {
  render() {
    const { onRunQuery, onChange } = this.props;
    const { rawSQL, format } = defaults(this.props.query, defaultQuery);

    return (
      <>
        <div>
          <TextArea
            css
            style={{ width: '100%' }}
            name="Query"
            className="slate-query-field"
            value={rawSQL}
            rows={10}
            placeholder="Enter a Redshift SQL query"
            onBlur={onRunQuery}
            onChange={e => onChange({ ...this.props.query, rawSQL: e.currentTarget.value })}
          />
        </div>
        <div>
          <Select
            options={SelectableFormatOptions}
            value={format}
            onChange={({ value }) => onChange({ ...this.props.query, format: value! })}
          />
        </div>
      </>
    );
  }
}
