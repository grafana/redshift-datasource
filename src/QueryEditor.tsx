import { defaults } from 'lodash';

import React, { PureComponent } from 'react';
import { QueryEditorProps } from '@grafana/data';
import { DataSource } from './datasource';
import { defaultQuery, RedshiftDataSourceOptions, RedshiftQuery } from './types';
import { TextArea } from '@grafana/ui';

type Props = QueryEditorProps<DataSource, RedshiftQuery, RedshiftDataSourceOptions>;

export class QueryEditor extends PureComponent<Props> {
  render() {
    const { onRunQuery, onChange } = this.props;
    const { rawSQL } = defaults(this.props.query, defaultQuery);

    return (
      <div className="gf-form">
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
    );
  }
}
