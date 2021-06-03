import { defaults } from 'lodash';

import React, { PureComponent } from 'react';
import { LegacyForms } from '@grafana/ui';
import { QueryEditorProps } from '@grafana/data';
import { DataSource } from './datasource';
import { defaultQuery, RedshiftDataSourceOptions, RedshiftQuery } from './types';

type Props = QueryEditorProps<DataSource, RedshiftQuery, RedshiftDataSourceOptions>;

export class QueryEditor extends PureComponent<Props> {
  render() {
    const query = defaults(this.props.query, defaultQuery);
    console.log(query);
    return <div className="gf-form"></div>;
  }
}
