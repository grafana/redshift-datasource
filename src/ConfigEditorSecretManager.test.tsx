import React from 'react';
import { render, screen } from '@testing-library/react';
import { ConfigEditorSecretManager } from './ConfigEditorSecretManager';

const props = {
  onChangeCluster: jest.fn(),
  onChangeDB: jest.fn(),
  onChangeSecret: jest.fn(),
};

describe('ConfigEditorSecretManager', () => {
  it('should display temporary credentials by default', () => {
    render(<ConfigEditorSecretManager {...props} />);
    expect(screen.getByText('Managed Secret ARN')).toBeInTheDocument();
    expect(screen.getByText('Cluster Identifier')).toBeInTheDocument();
    expect(screen.getByText('Database')).toBeInTheDocument();
  });
});
