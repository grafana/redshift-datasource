import React from 'react';
import { InlineField, Input } from '@grafana/ui';

export type Props = {
  clusterIdentifier?: string;
  database?: string;
  dbUser?: string;
  onChangeCluster: React.FormEventHandler<HTMLInputElement>;
  onChangeDB: React.FormEventHandler<HTMLInputElement>;
  onChangeDBUser: React.FormEventHandler<HTMLInputElement>;
};

export function ConfigEditorTempCreds(props: Props) {
  const { onChangeCluster, onChangeDB, onChangeDBUser, clusterIdentifier, database, dbUser } = props;
  return (
    <>
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
      <InlineField label="DB User" labelWidth={28}>
        <Input aria-label="DB User" css className="width-30" value={dbUser ?? ''} onChange={onChangeDBUser} />
      </InlineField>
    </>
  );
}
