import React, { useEffect, useState } from 'react';
import { InlineField, Input } from '@grafana/ui';
import { selectors } from 'selectors';
import { ResourceSelector } from 'ResourceSelector';
import { SelectableValue } from '@grafana/data';
import { RedshiftDataSourceOptions } from 'types';

export type Props = {
  clusterIdentifier?: string;
  database?: string;
  managedSecret?: RedshiftDataSourceOptions['managedSecret'];
  secretsDisabled?: boolean;
  fetchSecrets: () => Promise<Array<string | SelectableValue<string>>>;
  fetchSecret: (arn: string) => Promise<{ dbClusterIdentifier: string; username: string }>;
  setClusterID: (id: string) => void;
  onChangeDB: React.FormEventHandler<HTMLInputElement>;
  onChangeSecret: (e: SelectableValue<string> | null) => void;
  saveOptions: () => Promise<void>;
};

export function ConfigEditorSecretManager(props: Props) {
  const {
    setClusterID,
    onChangeDB,
    onChangeSecret,
    clusterIdentifier,
    database,
    managedSecret,
    secretsDisabled,
    fetchSecrets,
    fetchSecret,
    saveOptions,
  } = props;
  // The DB user is not stored in the JSON data since is not used
  const [dbUser, setDBUser] = useState('');
  useEffect(() => {
    if (managedSecret) {
      fetchSecret(managedSecret.arn).then((s) => {
        setClusterID(s.dbClusterIdentifier);
        setDBUser(s.username);
      });
    }
  }, [managedSecret, setClusterID, fetchSecret]);
  return (
    <>
      <ResourceSelector
        resource="ManagedSecret"
        onChange={onChangeSecret}
        fetch={fetchSecrets}
        value={managedSecret?.name || null}
        saveOptions={saveOptions}
        disabled={secretsDisabled}
        labelWidth={28}
        className="width-30"
      />
      <InlineField label={selectors.components.ConfigEditor.ClusterID.input} labelWidth={28} disabled={true}>
        <Input
          data-testid={selectors.components.ConfigEditor.ClusterID.testID}
          css
          className="width-30"
          value={clusterIdentifier ?? ''}
        />
      </InlineField>
      <InlineField label={selectors.components.ConfigEditor.DatabaseUser.input} labelWidth={28} disabled={true}>
        <Input
          data-testid={selectors.components.ConfigEditor.DatabaseUser.testID}
          css
          className="width-30"
          value={dbUser}
        />
      </InlineField>
      <InlineField label={selectors.components.ConfigEditor.Database.input} labelWidth={28}>
        <Input
          data-testid={selectors.components.ConfigEditor.Database.testID}
          css
          className="width-30"
          value={database ?? ''}
          onChange={onChangeDB}
        />
      </InlineField>
    </>
  );
}
