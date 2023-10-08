import { ConfigSelect, ConnectionConfig, InlineInput, Divider } from '@grafana/aws-sdk';
import { DataSourcePluginOptionsEditorProps, SelectableValue } from '@grafana/data';
import { config, getBackendSrv } from '@grafana/runtime';
import React, { FormEvent, useEffect, useState } from 'react';
import { selectors } from 'selectors';

import {
  RedshiftDataSourceOptions,
  RedshiftDataSourceSecureJsonData,
  RedshiftDataSourceSettings,
  RedshiftManagedSecret,
} from '../types';
import { AuthTypeSwitch } from './AuthTypeSwitch';
import { Field, InlineField, Input, Switch } from '@grafana/ui';
import { ConfigSection } from '@grafana/experimental';

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

type Workgroup = {
  workgroupName: string;
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
  const newFormStylingEnabled = config.featureToggles.awsDatasourcesNewFormStyling;
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
        // get workgroupName from user input since it's not stored in Secrets Manager
        const workgroupName = props.options.jsonData.workgroupName;
        if (props.options.jsonData.useServerless && workgroupName) {
          getWorkgroupUrl(workgroupName).then((url) => {
            props.onOptionsChange({
              ...props.options,
              url,
              jsonData: {
                ...props.options.jsonData,
                workgroupName: workgroupName,
                dbUser: s.username,
              },
            });
          });
        } else {
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
        }
      });
    }
    // eslint-disable-next-line react-hooks/exhaustive-deps
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
    } catch (error) {
      console.error('error while fetching clusters', error);
      return [];
    }
  };

  // Workgroups
  const [workgroupEndpoint, setWorkgroupEndpoint] = useState('');
  const fetchWorkgroups = async () => {
    try {
      const res: Workgroup[] = await getBackendSrv().get(resourcesURL + '/workgroups');
      return res.map((w) => ({
        label: w.workgroupName,
        value: w.workgroupName,
        description: `${w.endpoint.address}:${w.endpoint.port}`,
      }));
    } catch (error) {
      console.error('error while fetching workgroups', error);
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

  const getWorkgroupUrl = async (workgroupName: string) => {
    const { jsonData } = props.options;
    if (workgroupName !== jsonData.workgroupName || workgroupEndpoint === '') {
      const workgroups = await fetchWorkgroups();
      const endpoint = workgroups.find((w) => w.value === workgroupName)?.description || workgroupName;
      setWorkgroupEndpoint(endpoint || workgroupName);
      return `${endpoint}/${jsonData.database || ''}`;
    }
    return `${workgroupEndpoint}/${jsonData.database || ''}`;
  };

  useEffect(() => {
    if (!props.options.jsonData.useServerless && props.options.jsonData.clusterIdentifier) {
      getClusterUrl(props.options.jsonData.clusterIdentifier);
    }
    if (props.options.jsonData.useServerless && props.options.jsonData.workgroupName) {
      getWorkgroupUrl(props.options.jsonData.workgroupName);
    }
    // eslint-disable-next-line react-hooks/exhaustive-deps
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

  const onChangeWorkgroupName = (e: SelectableValue<string> | null) => {
    const value = e?.value ?? '';
    const url = e?.description + '/' + props.options.jsonData.database ?? '';
    props.onOptionsChange({
      ...props.options,
      url,
      jsonData: {
        ...props.options.jsonData,
        workgroupName: value,
      },
    });
    setWorkgroupEndpoint(e?.description || e?.value || '');
  };

  const onChange = (resource: InputResourceType) => (e: FormEvent<HTMLInputElement>) => {
    const value = e.currentTarget.value;
    const endpoint = props.options.jsonData.useServerless ? workgroupEndpoint : clusterEndpoint;
    const url = resource === 'database' ? `${endpoint}/${value}` : props.options.url;
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
    <div className="width-30">
      <ConnectionConfig {...props} onOptionsChange={onOptionsChange} newFormStylingEnabled={newFormStylingEnabled} />
      {newFormStylingEnabled ? (
        <>
          <Divider />
          <ConfigSection title="Redshift Details">
            <AuthTypeSwitch
              key="managedSecret"
              useManagedSecret={useManagedSecret}
              onChangeAuthType={onChangeAuthType}
            />
            <Field
              label={selectors.components.ConfigEditor.UseServerless.input}
              data-testid={selectors.components.ConfigEditor.UseServerless.testID}
              htmlFor="useServerless"
            >
              <Switch
                {...props}
                id="useServerless"
                value={props.options.jsonData.useServerless}
                aria-label={selectors.components.ConfigEditor.UseServerless.input}
                onChange={(e) =>
                  props.onOptionsChange({
                    ...props.options,
                    jsonData: {
                      ...props.options.jsonData,
                      useServerless: e.currentTarget.checked,
                    },
                  })
                }
              />
            </Field>
            <Field
              label={selectors.components.ConfigEditor.ClusterID.input}
              data-testid={selectors.components.ConfigEditor.ClusterID.testID}
              hidden={props.options.jsonData.useServerless || useManagedSecret}
              htmlFor="clusterId"
            >
              <ConfigSelect
                {...props}
                id="clusterId"
                label={selectors.components.ConfigEditor.ClusterID.input}
                allowCustomValue={true}
                value={props.options.jsonData.clusterIdentifier ?? ''}
                onChange={onChangeClusterID}
                fetch={fetchClusters}
                saveOptions={saveOptions}
                newFormStylingEnabled={true}
              />
            </Field>

            <Field
              hidden={props.options.jsonData.useServerless || !useManagedSecret}
              label={selectors.components.ConfigEditor.ClusterIDText.input}
              disabled={true}
              htmlFor="clusterIdText"
            >
              <Input
                {...props}
                id="clusterIdText"
                aria-label={selectors.components.ConfigEditor.ClusterIDText.input}
                value={props.options.jsonData.clusterIdentifier ?? ''}
                onChange={() => {}}
                label={selectors.components.ConfigEditor.ClusterIDText.input}
                data-testid={selectors.components.ConfigEditor.ClusterIDText.testID}
              />
            </Field>

            <Field
              hidden={!props.options.jsonData.useServerless}
              label={selectors.components.ConfigEditor.WorkgroupText.input}
              data-testid={selectors.components.ConfigEditor.WorkgroupText.testID}
              htmlFor="workgroupName"
            >
              <ConfigSelect
                {...props}
                id="workgroupName"
                label={selectors.components.ConfigEditor.WorkgroupText.input}
                value={props.options.jsonData.workgroupName ?? ''}
                onChange={onChangeWorkgroupName}
                fetch={fetchWorkgroups}
                saveOptions={saveOptions}
                newFormStylingEnabled={true}
              />
            </Field>
            <Field
              label={selectors.components.ConfigEditor.ManagedSecret.input}
              hidden={!useManagedSecret}
              data-testid={selectors.components.ConfigEditor.ManagedSecret.testID}
              htmlFor="managedSecret"
            >
              <ConfigSelect
                {...props}
                id="managedSecret"
                label={selectors.components.ConfigEditor.ManagedSecret.input}
                value={props.options.jsonData.managedSecret?.arn ?? ''}
                onChange={onChangeManagedSecret}
                fetch={fetchSecrets}
                saveOptions={saveOptions}
                allowCustomValue
                newFormStylingEnabled={true}
              />
            </Field>

            <Field
              label={selectors.components.ConfigEditor.DatabaseUser.input}
              hidden={props.options.jsonData.useServerless && !useManagedSecret}
              disabled={useManagedSecret}
              htmlFor="dbUser"
            >
              <Input
                {...props}
                id="dbUser"
                aria-label={selectors.components.ConfigEditor.DatabaseUser.input}
                value={props.options.jsonData.dbUser ?? ''}
                onChange={onChange('dbUser')}
                data-testid={selectors.components.ConfigEditor.DatabaseUser.testID}
              />
            </Field>

            <Field label={selectors.components.ConfigEditor.Database.input}>
              <Input
                {...props}
                data-testid={selectors.components.ConfigEditor.Database.testID}
                value={props.options.jsonData.database ?? ''}
                aria-label={selectors.components.ConfigEditor.Database.input}
                onChange={onChange('database')}
              />
            </Field>

            <Field label={selectors.components.ConfigEditor.WithEvent.input} htmlFor="withEvent">
              <Switch
                {...props}
                id="withEvent"
                aria-label={selectors.components.ConfigEditor.WithEvent.input}
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
                data-testid={selectors.components.ConfigEditor.WithEvent.testID}
              />
            </Field>
          </ConfigSection>
        </>
      ) : (
        <>
          <h3>Redshift Details</h3>
          <AuthTypeSwitch key="managedSecret" useManagedSecret={useManagedSecret} onChangeAuthType={onChangeAuthType} />
          <InlineField
            {...props}
            label={selectors.components.ConfigEditor.UseServerless.input}
            labelWidth={28}
            style={{ alignItems: 'center' }}
          >
            <Switch
              value={props.options.jsonData.useServerless}
              onChange={(e) =>
                props.onOptionsChange({
                  ...props.options,
                  jsonData: {
                    ...props.options.jsonData,
                    useServerless: e.currentTarget.checked,
                  },
                })
              }
              data-testid={selectors.components.ConfigEditor.UseServerless.testID}
            />
          </InlineField>
          <ConfigSelect
            {...props}
            id="clusterId"
            allowCustomValue={true}
            value={props.options.jsonData.clusterIdentifier ?? ''}
            onChange={onChangeClusterID}
            fetch={fetchClusters}
            label={selectors.components.ConfigEditor.ClusterID.input}
            data-testid={selectors.components.ConfigEditor.ClusterID.testID}
            saveOptions={saveOptions}
            hidden={props.options.jsonData.useServerless || useManagedSecret}
          />
          <InlineInput
            {...props}
            value={props.options.jsonData.clusterIdentifier ?? ''}
            onChange={() => {}}
            label={selectors.components.ConfigEditor.ClusterIDText.input}
            data-testid={selectors.components.ConfigEditor.ClusterIDText.testID}
            disabled={true}
            hidden={props.options.jsonData.useServerless || !useManagedSecret}
          />
          <ConfigSelect
            {...props}
            id="workgroupName"
            value={props.options.jsonData.workgroupName ?? ''}
            onChange={onChangeWorkgroupName}
            fetch={fetchWorkgroups}
            label={selectors.components.ConfigEditor.WorkgroupText.input}
            data-testid={selectors.components.ConfigEditor.WorkgroupText.testID}
            saveOptions={saveOptions}
            hidden={!props.options.jsonData.useServerless}
          />
          <ConfigSelect
            {...props}
            id="managedSecret"
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
            value={props.options.jsonData.dbUser ?? ''}
            onChange={onChange('dbUser')}
            label={selectors.components.ConfigEditor.DatabaseUser.input}
            data-testid={selectors.components.ConfigEditor.DatabaseUser.testID}
            disabled={useManagedSecret}
            hidden={props.options.jsonData.useServerless && !useManagedSecret}
          />
          <InlineInput
            {...props}
            value={props.options.jsonData.database ?? ''}
            onChange={onChange('database')}
            label={selectors.components.ConfigEditor.Database.input}
            data-testid={selectors.components.ConfigEditor.Database.testID}
          />
          <InlineField
            {...props}
            label={selectors.components.ConfigEditor.WithEvent.input}
            labelWidth={28}
            style={{ alignItems: 'center' }}
          >
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
              data-testid={selectors.components.ConfigEditor.WithEvent.testID}
            />
          </InlineField>
        </>
      )}
    </div>
  );
}
