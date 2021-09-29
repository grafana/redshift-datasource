import React from 'react';
import { InlineField, Input } from '@grafana/ui';
import { selectors } from 'selectors';

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
      <InlineField label={selectors.components.ConfigEditor.DatabaseUser.input} labelWidth={28}>
        <Input
          data-testid={selectors.components.ConfigEditor.DatabaseUser.testID}
          css
          className="width-30"
          value={dbUser ?? ''}
          onChange={onChangeDBUser}
        />
      </InlineField>
    </>
  );
}
