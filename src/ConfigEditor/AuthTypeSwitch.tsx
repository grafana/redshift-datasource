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
          <div style={{ marginTop: '10px', marginBottom: '10px', minWidth: '670px' }}>
            Use database username and password stored in Secrets Manager.{' '}
            <a
              href="https://docs.aws.amazon.com/redshift/latest/mgmt/data-api-access.html#data-api-secrets"
              target="_blank"
              rel="noreferrer"
            >
              Learn more
            </a>
          </div>
        ) : (
          <div style={{ marginTop: '10px', marginBottom: '10px', minWidth: '670px' }}>
            Use the IAM permission to generate temporary database username and password.{' '}
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
          { label: 'AWS Secrets Manager', value: true },
          { label: 'Temporary credentials', value: false },
        ]}
        value={useManagedSecret}
        onChange={onChangeAuthType}
      />
    </Label>
  );
}
