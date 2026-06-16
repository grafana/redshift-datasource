import React from 'react';
import { Label, RadioButtonGroup, useStyles2 } from '@grafana/ui';
import { GrafanaTheme2 } from '@grafana/data';
import { css } from '@emotion/css';

export type Props = {
  useManagedSecret: boolean;
  onChangeAuthType: (v: boolean) => void;
};

export function AuthTypeSwitch({ useManagedSecret, onChangeAuthType }: Props) {
  const styles = useStyles2(getStyles);
  return (
    <Label
      className={styles.label}
      description={
        useManagedSecret ? (
          <div className={styles.gap}>
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
          <div className={styles.gap}>
            Use
            <a
              href="https://docs.aws.amazon.com/redshift/latest/APIReference/API_GetClusterCredentials.html"
              target="_blank"
              rel="noreferrer"
            >
              <code>GetClusterCredentials</code>
            </a>
            or
            <a
              href="https://docs.aws.amazon.com/redshift-serverless/latest/APIReference/API_GetCredentials.html"
              target="_blank"
              rel="noreferrer"
            >
              <code>GetCredentials</code>
            </a>
            to generate temporary database username and password.{' '}
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
        className={styles.buttonGroup}
        options={[
          { label: 'Temporary credentials', value: false },
          { label: 'AWS Secrets Manager', value: true },
        ]}
        value={useManagedSecret}
        onChange={onChangeAuthType}
        fullWidth
      />
    </Label>
  );
}

const getStyles = (theme: GrafanaTheme2) => ({
  label: css({
    label: {
      display: 'block',
    },
    width: '100%',
    lineHeight: theme.typography.body.lineHeight,
  }),
  buttonGroup: css({
    flexGrow: 1,
  }),
  gap: css({
    marginTop: theme.spacing(2),
    marginBottom: theme.spacing(2),
  }),
});
