import React from 'react';
import { Label, RadioButtonGroup } from '@grafana/ui';

export type Props = {
  useManagedSecret: boolean;
  onChangeAuthType: (v: boolean) => void;
};

export function AuthTypeSwitch({ useManagedSecret, onChangeAuthType }: Props) {
  return (
    <Label
      description={
        useManagedSecret ? (
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
        ) : (
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
        )
      }
    >
      <RadioButtonGroup
        options={[
          { label: 'Temporary credentials', value: false },
          { label: 'AWS Secrets Manager', value: true },
        ]}
        value={useManagedSecret}
        onChange={onChangeAuthType}
      />
    </Label>
  );
}
