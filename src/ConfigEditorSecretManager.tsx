import React from 'react';
import { InlineField, Input } from '@grafana/ui';

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
      <InlineField label="Managed Secret ARN" labelWidth={28}>
        <Input
          data-test-id="managed-secret"
          css
          className="width-30"
          value={managedSecret ?? ''}
          onChange={onChangeSecret}
        />
      </InlineField>
      {/* TODO: Obtain this info from the secret and disable */}
      <InlineField label="Cluster Identifier" labelWidth={28}>
        <Input
          data-test-id="cluster-id"
          css
          className="width-30"
          value={clusterIdentifier ?? ''}
          onChange={onChangeCluster}
        />
      </InlineField>
      <InlineField label="Database" labelWidth={28}>
        <Input data-test-id="database" css className="width-30" value={database ?? ''} onChange={onChangeDB} />
      </InlineField>
      {/* TODO: Add db user info from the secret */}
    </>
  );
}
