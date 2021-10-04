import React, { useState } from 'react';
import { DataSourcePluginOptionsEditorProps } from '@grafana/data';
import {
  RedshiftDataSourceOptions,
  RedshiftDataSourceSecureJsonData,
  RedshiftDataSourceSettings,
  RedshiftManagedSecret,
} from '../types';
import { ConnectionConfig } from '@grafana/aws-sdk';
import { TempCreds } from './TempCreds';
import { SecretManager } from './SecretManager';
import { getBackendSrv } from '@grafana/runtime';
import { AuthTypeSwitch } from './AuthTypeSwitch';

export type Props = DataSourcePluginOptionsEditorProps<RedshiftDataSourceOptions, RedshiftDataSourceSecureJsonData>;

export function ConfigEditor(props: Props) {
  const baseURL = `/api/datasources/${props.options.id}`;
  const resourcesURL = `${baseURL}/resources`;
  const { jsonData } = props.options;
  const [saved, setSaved] = useState(!!jsonData.defaultRegion);
  const saveOptions = async () => {
    if (saved) {
      return;
    }
    const result: { datasource: RedshiftDataSourceSettings } = await getBackendSrv().put(baseURL, props.options);
    props.onOptionsChange({
      ...props.options,
      version: result.datasource.version,
    });
    setSaved(true);
  };

  // Auth type
  const [useManagedSecret, setUseManagedSecret] = useState(!!props.options.jsonData.useManagedSecret);
  const onChangeAuthType = (newAuthType: boolean) => {
    setUseManagedSecret(newAuthType);
    props.onOptionsChange({
      ...props.options,
      jsonData: {
        ...props.options.jsonData,
        useManagedSecret: newAuthType,
      },
    });
  };

  // Secrets
  const [managedSecret, setManagedSecret] = useState(jsonData.managedSecret);
  const fetchSecrets = async () => {
    const res: RedshiftManagedSecret[] = await getBackendSrv().get(resourcesURL + '/secrets');
    return res.map((r) => ({ label: r.name, value: r.arn, description: r.arn }));
  };
  const fetchSecret = async (arn: string) => {
    const res: { dbClusterIdentifier: string; username: string } = await getBackendSrv().post(
      resourcesURL + '/secret',
      { secretARN: arn }
    );
    return res;
  };
  const onSecretChange = (managedSecret?: RedshiftManagedSecret) => {
    setManagedSecret(managedSecret);
    props.onOptionsChange({
      ...props.options,
      jsonData: {
        ...props.options.jsonData,
        managedSecret,
      },
    });
  };

  // ClusterID
  const [clusterIdentifier, setClusterIdentifier] = useState(jsonData.clusterIdentifier);
  const onClusterIdentifierChange = (id?: string) => {
    setClusterIdentifier(id);
    props.onOptionsChange({
      ...props.options,
      jsonData: {
        ...props.options.jsonData,
        clusterIdentifier: id,
      },
    });
  };

  // Database
  const [database, setDatabase] = useState(jsonData.database);
  const onDatabaseChange = (db?: string) => {
    setDatabase(db);
    props.onOptionsChange({
      ...props.options,
      jsonData: {
        ...props.options.jsonData,
        database: db,
      },
    });
  };

  // DB user
  const [dbUser, setDBUser] = useState(jsonData.dbUser);
  const onDBUserChange = (user?: string) => {
    setDBUser(user);
    props.onOptionsChange({
      ...props.options,
      jsonData: {
        ...props.options.jsonData,
        dbUser: user,
      },
    });
  };

  const onOptionsChange = (options: RedshiftDataSourceSettings) => {
    setSaved(false);
    props.onOptionsChange(options);
  };

  return (
    <>
      <ConnectionConfig {...props} onOptionsChange={onOptionsChange} />
      <h6>Authentication</h6>
      <AuthTypeSwitch useManagedSecret={useManagedSecret} onChangeAuthType={onChangeAuthType} />
      {useManagedSecret ? (
        <SecretManager
          clusterIdentifier={clusterIdentifier}
          database={database}
          managedSecret={managedSecret}
          secretsDisabled={!jsonData.defaultRegion}
          fetchSecrets={fetchSecrets}
          fetchSecret={fetchSecret}
          onChangeDB={onDatabaseChange}
          onChangeSecret={onSecretChange}
          onChangeClusterID={onClusterIdentifierChange}
          saveOptions={saveOptions}
        />
      ) : (
        <TempCreds
          clusterIdentifier={clusterIdentifier}
          database={database}
          dbUser={dbUser}
          onChangeDB={onDatabaseChange}
          onChangeDBUser={onDBUserChange}
          onChangeCluster={onClusterIdentifierChange}
        />
      )}
    </>
  );
}
