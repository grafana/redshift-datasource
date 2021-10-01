import React from 'react';
import { render, screen } from '@testing-library/react';
import { TempCreds } from './TempCreds';
import { selectors } from '../selectors';

const props = {
  onChangeCluster: jest.fn(),
  onChangeDB: jest.fn(),
  onChangeDBUser: jest.fn(),
};

describe('TempCreds', () => {
  it('should display temporary credentials by default', () => {
    render(<TempCreds {...props} />);
    expect(screen.getByText(selectors.components.ConfigEditor.ClusterID.input)).toBeInTheDocument();
    expect(screen.getByText(selectors.components.ConfigEditor.ClusterID.input)).toBeInTheDocument();
    expect(screen.getByText(selectors.components.ConfigEditor.Database.input)).toBeInTheDocument();
  });
});
