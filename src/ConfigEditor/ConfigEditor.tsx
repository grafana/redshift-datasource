import React, { useState, useEffect, FormEvent } from 'react';
import { DataSourcePluginOptionsEditorProps, SelectableValue } from '@grafana/data';
import {
  RedshiftDataSourceOptions,
  RedshiftDataSourceSecureJsonData,
  RedshiftDataSourceSettings,
  RedshiftManagedSecret,
} from '../types';
import { InlineInput, ConfigSelect, ConnectionConfig } from '@grafana/aws-sdk';
import { getBackendSrv } from '@grafana/runtime';
import { AuthTypeSwitch } from './AuthTypeSwitch';
import { selectors } from 'selectors';

export type Props = DataSourcePluginOptionsEditorProps<RedshiftDataSourceOptions, RedshiftDataSourceSecureJsonData>;

type Secret = { dbClusterIdentifier: string; username: string };

type InputResourceType = 'clusterIdentifier' | 'dbUser' | 'database';

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
  const fetchSecrets = async () => {
    const res: RedshiftManagedSecret[] = await getBackendSrv().get(resourcesURL + '/secrets');
    return res.map((r) => ({ label: r.name, value: r.arn, description: r.arn }));
  };
  const { arn } = props.options.jsonData.managedSecret || {};
  const fetchSecret = async (arn: string) => {
    const res: Secret = await getBackendSrv().post(resourcesURL + '/secret', { secretARN: arn });
    return res;
  };
  useEffect(() => {
    if (arn) {
      fetchSecret(arn).then((s) => {
        props.onOptionsChange({
          ...props.options,
          jsonData: {
            ...props.options.jsonData,
            clusterIdentifier: s.dbClusterIdentifier,
            dbUser: s.username,
          },
        });
      });
    }
  }, [arn]);

  // Description to show in datasource card
  const { dbUser, clusterIdentifier, database } = props.options.jsonData || {};
  useEffect(() => {
    console.log(`${dbUser}@${clusterIdentifier}/${database}`);
    if(dbUser || clusterIdentifier || database) {
      props.onOptionsChange({
        ...props.options,
        url: `${dbUser}@${clusterIdentifier}/${database}`,
      });
    }
  }, [dbUser, clusterIdentifier, database]);

  const onOptionsChange = (options: RedshiftDataSourceSettings) => {
    setSaved(false);
    props.onOptionsChange(options);
  };

  const onChangeManagedSecret = (e: SelectableValue<string> | null) => {
    const value = e?.value ?? '';
    const label = e?.label ?? '';
    props.onOptionsChange({
      ...props.options,
      jsonData: {
        ...props.options.jsonData,
        managedSecret: { arn: value, name: label },
      },
    });
  };
  const onChange = (resource: InputResourceType) => (e: FormEvent<HTMLInputElement>) => {
    const value = e.currentTarget.value;
    props.onOptionsChange({
      ...props.options,
      jsonData: {
        ...props.options.jsonData,
        [resource]: value,
      },
    });
  };

  return (
    <div className="gf-form-group">
      <ConnectionConfig {...props} onOptionsChange={onOptionsChange} />
      <h3>Redshift Details</h3>
      <AuthTypeSwitch key="managedSecret" useManagedSecret={useManagedSecret} onChangeAuthType={onChangeAuthType} />
      <ConfigSelect
        {...props}
        value={props.options.jsonData.managedSecret?.arn ?? ''}
        onChange={onChangeManagedSecret}
        fetch={fetchSecrets}
        label={selectors.components.ConfigEditor.ManagedSecret.input}
        data-testid={selectors.components.ConfigEditor.ManagedSecret.testID}
        saveOptions={saveOptions}
        hidden={!useManagedSecret}
      />
      <InlineInput
        {...props}
        value={props.options.jsonData.clusterIdentifier ?? ''}
        onChange={onChange('clusterIdentifier')}
        label={selectors.components.ConfigEditor.ClusterID.input}
        data-testid={selectors.components.ConfigEditor.ClusterID.testID}
        disabled={useManagedSecret}
      />
      <InlineInput
        {...props}
        value={props.options.jsonData.dbUser ?? ''}
        onChange={onChange('dbUser')}
        label={selectors.components.ConfigEditor.DatabaseUser.input}
        data-testid={selectors.components.ConfigEditor.DatabaseUser.testID}
        disabled={useManagedSecret}
      />
      <InlineInput
        {...props}
        value={props.options.jsonData.database ?? ''}
        onChange={onChange('database')}
        label={selectors.components.ConfigEditor.Database.input}
        data-testid={selectors.components.ConfigEditor.Database.testID}
      />
    </div>
  );
}
