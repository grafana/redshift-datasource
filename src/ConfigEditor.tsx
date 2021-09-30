import React, { useState } from 'react';
import { Label, RadioButtonGroup } from '@grafana/ui';
import { DataSourcePluginOptionsEditorProps, onUpdateDatasourceJsonDataOption, SelectableValue } from '@grafana/data';
import { RedshiftDataSourceOptions, RedshiftDataSourceSecureJsonData, RedshiftDataSourceSettings } from './types';
import { ConnectionConfig } from '@grafana/aws-sdk';
import { ConfigEditorTempCreds } from 'ConfigEditorTempCreds';
import { ConfigEditorSecretManager } from 'ConfigEditorSecretManager';
import { getBackendSrv } from '@grafana/runtime';

export type Props = DataSourcePluginOptionsEditorProps<RedshiftDataSourceOptions, RedshiftDataSourceSecureJsonData>;

export function ConfigEditor(props: Props) {
  const [useTempCreds, setUseTempCreds] = useState(!props.options.jsonData.managedSecret);
  const baseURL = `/api/datasources/${props.options.id}`;
  const resourcesURL = `${baseURL}/resources`;
  const { jsonData } = props.options;
  const [saved, setSaved] = useState(!!jsonData.defaultRegion);
  const saveOptions = async () => {
    if (saved) {
      return;
    }
    await getBackendSrv()
      .put(baseURL, props.options)
      .then((result: { datasource: RedshiftDataSourceSettings }) => {
        props.onOptionsChange({
          ...props.options,
          version: result.datasource.version,
        });
      });
    setSaved(true);
  };
  // Secrets
  const fetchSecrets = async () => {
    const res: Array<{ name: string; arn: string }> = await getBackendSrv().get(resourcesURL + '/secrets');
    return res.map((r) => ({ label: r.name, value: r.arn, description: r.arn }));
  };
  const fetchSecret = async (arn: string) => {
    const res: { dbClusterIdentifier: string; username: string } = await getBackendSrv().post(
      resourcesURL + '/secret',
      { secretARN: arn }
    );
    return res;
  };

  const setClusterID = (id: string) => {
    props.onOptionsChange({
      ...props.options,
      jsonData: {
        ...props.options.jsonData,
        clusterIdentifier: id,
      },
    });
  };

  const onOptionsChange = (options: RedshiftDataSourceSettings) => {
    // clean up related state
    setSaved(false);
    props.onOptionsChange(options);
  };

  const onSecretChange = (e: SelectableValue<string> | null) => {
    const name = e === null ? e : e.label || '';
    const arn = e === null ? e : e.description || '';
    const managedSecret = !name || !arn ? undefined : { name, arn };
    props.onOptionsChange({
      ...props.options,
      jsonData: {
        ...props.options.jsonData,
        managedSecret,
      },
    });
  };

  const onChangeAuthType = (newUseTempCreds: boolean) => {
    setUseTempCreds(newUseTempCreds);
    // Clean up state
    props.onOptionsChange({
      ...props.options,
      jsonData: {
        ...props.options.jsonData,
        managedSecret: undefined,
        dbUser: undefined,
        clusterIdentifier: undefined,
        database: undefined,
      },
    });
  };
  return (
    <>
      <ConnectionConfig {...props} onOptionsChange={onOptionsChange} />
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
          managedSecret={jsonData.managedSecret}
          secretsDisabled={!jsonData.defaultRegion}
          fetchSecrets={fetchSecrets}
          fetchSecret={fetchSecret}
          onChangeDB={onUpdateDatasourceJsonDataOption(props, 'database')}
          onChangeSecret={onSecretChange}
          setClusterID={setClusterID}
          saveOptions={saveOptions}
        />
      )}
    </>
  );
}
