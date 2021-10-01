import React from 'react';
import { render, screen, waitFor } from '@testing-library/react';
import { SecretManager } from './SecretManager';
import { selectors } from '../selectors';
import { select } from 'react-select-event';

const props = {
  fetchSecrets: jest.fn(),
  fetchSecret: jest.fn(),
  onChangeClusterID: jest.fn(),
  onChangeDB: jest.fn(),
  onChangeSecret: jest.fn(),
  saveOptions: jest.fn(),
};

describe('SecretManager', () => {
  it('should display temporary credentials by default', async () => {
    render(<SecretManager {...props} />);
    expect(screen.getByText(selectors.components.ConfigEditor.ManagedSecret.input)).toBeInTheDocument();
    expect(screen.getByText(selectors.components.ConfigEditor.ClusterID.input)).toBeInTheDocument();
    expect(screen.getByText(selectors.components.ConfigEditor.Database.input)).toBeInTheDocument();
  });

  it('should fetch secrets', async () => {
    const secret = { label: 'foo', value: 'arn:foo' };
    const fetchSecrets = jest.fn().mockResolvedValue([secret]);
    render(<SecretManager {...props} fetchSecrets={fetchSecrets} />);
    screen.getByTestId(selectors.components.ConfigEditor.ManagedSecret.testID).click();

    const selectEl = screen.getByLabelText(selectors.components.ConfigEditor.ManagedSecret.input);
    expect(selectEl).toBeInTheDocument();
    await select(selectEl, secret.label, { container: document.body });

    expect(fetchSecrets).toHaveBeenCalled();
  });

  it('should fetch and display the secret info', async () => {
    const onChangeClusterID = jest.fn();
    const secret = { dbClusterIdentifier: 'clusterIdentifier', username: 'dbUser' };
    const fetchSecret = jest.fn().mockResolvedValue(secret);
    render(
      <SecretManager
        {...props}
        fetchSecret={fetchSecret}
        managedSecret={{ arn: 'foo', name: 'bar' }}
        onChangeClusterID={onChangeClusterID}
      />
    );

    expect(fetchSecret).toHaveBeenCalled();
    await waitFor(() => screen.getByDisplayValue(secret.username));
    expect(onChangeClusterID).toHaveBeenCalledWith(secret.dbClusterIdentifier);
  });
});
