import React, { ChangeEvent, PureComponent } from 'react';
import { DataSourcePluginOptionsEditorProps } from '@grafana/data';
import { RedshiftDataSourceOptions, RedshiftDataSourceSecureJsonData } from './types';

interface Props extends DataSourcePluginOptionsEditorProps<RedshiftDataSourceOptions> {}

interface State {}

export class ConfigEditor extends PureComponent<Props, State> {
  render() {
    return <div className="gf-form-group">Redshift config page</div>;
  }
}
