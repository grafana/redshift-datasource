import React from 'react';
import { render, screen, waitFor } from '@testing-library/react';
import { ConfigEditorSecretManager } from './ConfigEditorSecretManager';
import { selectors } from 'selectors';
import { select } from 'react-select-event';

const props = {
  fetchSecrets: jest.fn(),
  fetchSecret: jest.fn(),
  setClusterID: jest.fn(),
  onChangeDB: jest.fn(),
  onChangeSecret: jest.fn(),
  saveOptions: jest.fn(),
};

describe('ConfigEditorSecretManager', () => {
  it('should display temporary credentials by default', async () => {
    render(<ConfigEditorSecretManager {...props} />);
    expect(screen.getByText(selectors.components.ConfigEditor.ManagedSecret.input)).toBeInTheDocument();
    expect(screen.getByText(selectors.components.ConfigEditor.ClusterID.input)).toBeInTheDocument();
    expect(screen.getByText(selectors.components.ConfigEditor.Database.input)).toBeInTheDocument();
  });

  it('should fecth secrets', async () => {
    const secret = { label: 'foo', value: 'arn:foo' };
    const fetchSecrets = jest.fn().mockResolvedValue([secret]);
    render(<ConfigEditorSecretManager {...props} fetchSecrets={fetchSecrets} />);
    screen.getByTestId(selectors.components.ConfigEditor.ManagedSecret.testID).click();

    const selectEl = screen.getByLabelText(selectors.components.ConfigEditor.ManagedSecret.input);
    expect(selectEl).toBeInTheDocument();
    await select(selectEl, secret.label, { container: document.body });

    expect(fetchSecrets).toHaveBeenCalled();
  });

  it('should fetch and display the secret info', async () => {
    const setClusterID = jest.fn();
    const secret = { dbClusterIdentifier: 'clusterIdentifier', username: 'dbUser' };
    const fetchSecret = jest.fn().mockResolvedValue(secret);
    render(
      <ConfigEditorSecretManager
        {...props}
        fetchSecret={fetchSecret}
        managedSecret={{ arn: 'foo', name: 'bar' }}
        setClusterID={setClusterID}
      />
    );

    expect(fetchSecret).toHaveBeenCalled();
    await waitFor(() => screen.getByDisplayValue(secret.username));
    expect(setClusterID).toHaveBeenCalledWith(secret.dbClusterIdentifier);
  });
});
