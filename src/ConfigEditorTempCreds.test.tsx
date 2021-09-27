import React from 'react';
import { render, screen } from '@testing-library/react';
import { ConfigEditorTempCreds } from './ConfigEditorTempCreds';

const props = {
  onChangeCluster: jest.fn(),
  onChangeDB: jest.fn(),
  onChangeDBUser: jest.fn(),
};

describe('ConfigEditorTempCreds', () => {
  it('should display temporary credentials by default', () => {
    render(<ConfigEditorTempCreds {...props} />);
    expect(screen.getByText('Cluster Identifier')).toBeInTheDocument();
    expect(screen.getByText('Database')).toBeInTheDocument();
    expect(screen.getByText('DB User')).toBeInTheDocument();
  });
});
