import React from 'react';
import { render, screen } from '@testing-library/react';
import { ConfigEditorSecretManager } from './ConfigEditorSecretManager';
import { selectors } from 'selectors';

const props = {
  onChangeCluster: jest.fn(),
  onChangeDB: jest.fn(),
  onChangeSecret: jest.fn(),
};

describe('ConfigEditorSecretManager', () => {
  it('should display temporary credentials by default', () => {
    render(<ConfigEditorSecretManager {...props} />);
    expect(screen.getByText(selectors.components.ConfigEditor.ManagedSecret.input)).toBeInTheDocument();
    expect(screen.getByText(selectors.components.ConfigEditor.ClusterID.input)).toBeInTheDocument();
    expect(screen.getByText(selectors.components.ConfigEditor.Database.input)).toBeInTheDocument();
  });
});
