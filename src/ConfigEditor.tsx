import React, { PureComponent } from 'react';
import { InlineField, Input } from '@grafana/ui';
import { DataSourcePluginOptionsEditorProps, onUpdateDatasourceJsonDataOption } from '@grafana/data';
import { RedshiftDataSourceOptions, RedshiftDataSourceSecureJsonData } from './types';
import { ConnectionConfig } from '@grafana/aws-sdk';

export type Props = DataSourcePluginOptionsEditorProps<RedshiftDataSourceOptions, RedshiftDataSourceSecureJsonData>;

export class ConfigEditor extends PureComponent<Props> {
  render() {
    const { clusterIdentifier, database, dbUser } = this.props.options.jsonData;
    return (
      <>
        <ConnectionConfig {...this.props} />
        <InlineField label="Cluster Identifier" labelWidth={28}>
          <Input
            data-test-id="cluster-id"
            css
            className="width-30"
            value={clusterIdentifier ?? ''}
            onChange={onUpdateDatasourceJsonDataOption(this.props, 'clusterIdentifier')}
          />
        </InlineField>
        <InlineField label="Database" labelWidth={28}>
          <Input
            data-test-id="database"
            css
            className="width-30"
            value={database ?? ''}
            onChange={onUpdateDatasourceJsonDataOption(this.props, 'database')}
          />
        </InlineField>
        <InlineField label="DB User" labelWidth={28}>
          <Input
            data-test-id="dbuser"
            css
            className="width-30"
            value={dbUser ?? ''}
            onChange={onUpdateDatasourceJsonDataOption(this.props, 'dbUser')}
          />
        </InlineField>
      </>
    );
  }
}
