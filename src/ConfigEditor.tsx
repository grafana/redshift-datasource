import React, { useState } from 'react';
import { Label, RadioButtonGroup } from '@grafana/ui';
import { DataSourcePluginOptionsEditorProps, onUpdateDatasourceJsonDataOption } from '@grafana/data';
import { RedshiftDataSourceOptions, RedshiftDataSourceSecureJsonData } from './types';
import { ConnectionConfig } from '@grafana/aws-sdk';
import { ConfigEditorTempCreds } from 'ConfigEditorTempCreds';
import { ConfigEditorSecretManager } from 'ConfigEditorSecretManager';

export type Props = DataSourcePluginOptionsEditorProps<RedshiftDataSourceOptions, RedshiftDataSourceSecureJsonData>;

export function ConfigEditor(props: Props) {
  const [useTempCreds, setUseTempCreds] = useState(!props.options.jsonData.managedSecret);
  const onChangeAuthType = (newUseTempCreds: boolean) => {
    setUseTempCreds(newUseTempCreds);
    // Clean up state for the non-used type
    if (newUseTempCreds) {
      props.onOptionsChange({
        ...props.options,
        jsonData: {
          ...props.options.jsonData,
          managedSecret: '',
        },
      });
    } else {
      props.onOptionsChange({
        ...props.options,
        jsonData: {
          ...props.options.jsonData,
          dbUser: '',
        },
      });
    }
  };
  return (
    <>
      <ConnectionConfig {...props} />
      <h6>Authentication</h6>
      <Label
        description={
          useTempCreds ? (
            <div style={{ marginTop: '10px', marginBottom: '10px', minWidth: '670px' }}>
              Use the <code>GetClusterCredentials</code> IAM permission and your database user to generate temporary
              access credentials.
            </div>
          ) : (
            <div style={{ marginTop: '10px', marginBottom: '10px' }}>Use a stored secret to authenticate access.</div>
          )
        }
      >
        <RadioButtonGroup
          options={[
            { label: 'Temporary credentials', value: true },
            { label: 'AWS Secrets Manager', value: false },
          ]}
          value={useTempCreds}
          onChange={onChangeAuthType}
        />
      </Label>
      {useTempCreds ? (
        <ConfigEditorTempCreds
          clusterIdentifier={props.options.jsonData.clusterIdentifier}
          database={props.options.jsonData.database}
          dbUser={props.options.jsonData.dbUser}
          onChangeDB={onUpdateDatasourceJsonDataOption(props, 'database')}
          onChangeDBUser={onUpdateDatasourceJsonDataOption(props, 'dbUser')}
          onChangeCluster={onUpdateDatasourceJsonDataOption(props, 'clusterIdentifier')}
        />
      ) : (
        <ConfigEditorSecretManager
          clusterIdentifier={props.options.jsonData.clusterIdentifier}
          database={props.options.jsonData.database}
          managedSecret={props.options.jsonData.managedSecret}
          onChangeDB={onUpdateDatasourceJsonDataOption(props, 'database')}
          onChangeSecret={onUpdateDatasourceJsonDataOption(props, 'managedSecret')}
          onChangeCluster={onUpdateDatasourceJsonDataOption(props, 'clusterIdentifier')}
        />
      )}
    </>
  );
}
