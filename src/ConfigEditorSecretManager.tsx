import React from 'react';
import { InlineField, Input } from '@grafana/ui';
import { selectors } from 'selectors';

export type Props = {
  clusterIdentifier?: string;
  database?: string;
  managedSecret?: string;
  onChangeCluster: React.FormEventHandler<HTMLInputElement>;
  onChangeDB: React.FormEventHandler<HTMLInputElement>;
  onChangeSecret: React.FormEventHandler<HTMLInputElement>;
};

export function ConfigEditorSecretManager(props: Props) {
  const { onChangeCluster, onChangeDB, onChangeSecret, clusterIdentifier, database, managedSecret } = props;
  return (
    <>
      <InlineField label={selectors.components.ConfigEditor.ManagedSecret.input} labelWidth={28}>
        <Input
          data-testid={selectors.components.ConfigEditor.ManagedSecret.testID}
          css
          className="width-30"
          value={managedSecret ?? ''}
          onChange={onChangeSecret}
        />
      </InlineField>
      {/* TODO: Obtain this info from the secret and disable */}
      <InlineField label={selectors.components.ConfigEditor.ClusterID.input} labelWidth={28}>
        <Input
          data-testid={selectors.components.ConfigEditor.ClusterID.testID}
          css
          className="width-30"
          value={clusterIdentifier ?? ''}
          onChange={onChangeCluster}
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
      {/* TODO: Add db user info from the secret */}
    </>
  );
}
