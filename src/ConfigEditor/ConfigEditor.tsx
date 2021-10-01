import React, { useState } from 'react';
import { Label, RadioButtonGroup } from '@grafana/ui';
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

  // Auth type
  const [useTempCreds, setUseTempCreds] = useState(props.options.jsonData.useTemporaryCredentials ?? true);
  const onChangeAuthType = (newUseTempCreds: boolean) => {
    setUseTempCreds(newUseTempCreds);
    props.onOptionsChange({
      ...props.options,
      jsonData: {
        ...props.options.jsonData,
        useTemporaryCredentials: newUseTempCreds,
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
  const [clusterIdentifier, setClusterclusterIdentifier] = useState(jsonData.clusterIdentifier);
  const onClusterIdentifierChange = (id?: string) => {
    setClusterclusterIdentifier(id);
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
      <Label
        description={
          useTempCreds ? (
            <div style={{ marginTop: '10px', marginBottom: '10px', minWidth: '670px' }}>
              Use the <code>GetClusterCredentials</code> IAM permission and your database user to generate temporary
              access credentials.{' '}
              <a
                href="https://docs.aws.amazon.com/redshift/latest/mgmt/generating-user-credentials.html"
                target="_blank"
                rel="noreferrer"
              >
                Learn more
              </a>
            </div>
          ) : (
            <div style={{ marginTop: '10px', marginBottom: '10px' }}>
              Use a stored secret to authenticate access.{' '}
              <a
                href="https://docs.aws.amazon.com/secretsmanager/latest/userguide/intro.html"
                target="_blank"
                rel="noreferrer"
              >
                Learn more
              </a>
            </div>
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
        <TempCreds
          clusterIdentifier={clusterIdentifier}
          database={database}
          dbUser={dbUser}
          onChangeDB={onDatabaseChange}
          onChangeDBUser={onDBUserChange}
          onChangeCluster={onClusterIdentifierChange}
        />
      ) : (
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
      )}
    </>
  );
}
