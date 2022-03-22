import { ConfigSelect, ConnectionConfig, InlineInput } from '@grafana/aws-sdk';
import { DataSourcePluginOptionsEditorProps, SelectableValue } from '@grafana/data';
import { getBackendSrv } from '@grafana/runtime';
import React, { FormEvent, useEffect, useState } from 'react';
import { selectors } from 'selectors';

import {
  RedshiftDataSourceOptions,
  RedshiftDataSourceSecureJsonData,
  RedshiftDataSourceSettings,
  RedshiftManagedSecret,
} from '../types';
import { AuthTypeSwitch } from './AuthTypeSwitch';
import { InlineField, Switch } from '@grafana/ui';

export type Props = DataSourcePluginOptionsEditorProps<RedshiftDataSourceOptions, RedshiftDataSourceSecureJsonData>;

type Secret = { dbClusterIdentifier: string; username: string };

type Cluster = {
  clusterIdentifier: string;
  endpoint: {
    address: string;
    port: number;
  };
  database: string;
};

type InputResourceType = 'dbUser' | 'database';

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
        getClusterUrl(s.dbClusterIdentifier).then((url) => {
          props.onOptionsChange({
            ...props.options,
            url,
            jsonData: {
              ...props.options.jsonData,
              clusterIdentifier: s.dbClusterIdentifier,
              dbUser: s.username,
            },
          });
        });
      });
    }
  }, [arn]);

  // Clusters
  const [clusterEndpoint, setClusterEndpoint] = useState('');
  const fetchClusters = async () => {
    try {
      const res: Cluster[] = await getBackendSrv().get(resourcesURL + '/clusters');
      return res.map((c) => ({
        label: c.clusterIdentifier,
        value: c.clusterIdentifier,
        description: `${c.endpoint.address}:${c.endpoint.port}`,
      }));
    } catch {
      return [];
    }
  };

  const getClusterUrl = async (clusterID: string) => {
    const { jsonData } = props.options;
    if (clusterID !== jsonData.clusterIdentifier || clusterEndpoint === '') {
      const clusters = await fetchClusters();
      const endpoint = clusters.find((c) => c.value === clusterID)?.description || clusterID;
      setClusterEndpoint(endpoint || clusterID);
      return `${endpoint}/${jsonData.database || ''}`;
    }
    return `${clusterEndpoint}/${jsonData.database || ''}`;
  };

  useEffect(() => {
    if (props.options.jsonData.clusterIdentifier) {
      getClusterUrl(props.options.jsonData.clusterIdentifier);
    }
  }, []);

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
  const onChangeClusterID = (e: SelectableValue<string> | null) => {
    const value = e?.value ?? '';
    const url = e?.description + '/' + props.options.jsonData.database ?? '';
    props.onOptionsChange({
      ...props.options,
      url,
      jsonData: {
        ...props.options.jsonData,
        clusterIdentifier: value,
      },
    });
    setClusterEndpoint(e?.description || e?.value || '');
  };
  const onChange = (resource: InputResourceType) => (e: FormEvent<HTMLInputElement>) => {
    const value = e.currentTarget.value;
    const url = resource === 'database' ? `${clusterEndpoint}/${value}` : props.options.url;
    props.onOptionsChange({
      ...props.options,
      url,
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
      <ConfigSelect
        {...props}
        allowCustomValue={true}
        value={props.options.jsonData.clusterIdentifier ?? ''}
        onChange={onChangeClusterID}
        fetch={fetchClusters}
        label={selectors.components.ConfigEditor.ClusterID.input}
        data-testid={selectors.components.ConfigEditor.ClusterID.testID}
        saveOptions={saveOptions}
        hidden={useManagedSecret}
      />
      <InlineInput
        {...props}
        value={props.options.jsonData.clusterIdentifier ?? ''}
        onChange={() => {}}
        label={selectors.components.ConfigEditor.ClusterIDText.input}
        data-testid={selectors.components.ConfigEditor.ClusterIDText.testID}
        disabled={true}
        hidden={!useManagedSecret}
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
      <InlineField label={selectors.components.ConfigEditor.WithEvent.input} labelWidth={28} style={{ alignItems: 'center' }}>
        <Switch
          value={props.options.jsonData.withEvent ?? false}
          onChange={(e) =>
            props.onOptionsChange({
              ...props.options,
              jsonData: {
                ...props.options.jsonData,
                withEvent: e.currentTarget.checked,
              },
            })
          }
          css={undefined}
          data-testid={selectors.components.ConfigEditor.WithEvent.testID}
        />
      </InlineField>
    </div>
  );
}
