import React from 'react';
import { render, screen } from '@testing-library/react';
import { ConfigEditorTempCreds } from './ConfigEditorTempCreds';
import { selectors } from 'selectors';

const props = {
  onChangeCluster: jest.fn(),
  onChangeDB: jest.fn(),
  onChangeDBUser: jest.fn(),
};

describe('ConfigEditorTempCreds', () => {
  it('should display temporary credentials by default', () => {
    render(<ConfigEditorTempCreds {...props} />);
    expect(screen.getByText(selectors.components.ConfigEditor.ClusterID.input)).toBeInTheDocument();
    expect(screen.getByText(selectors.components.ConfigEditor.ClusterID.input)).toBeInTheDocument();
    expect(screen.getByText(selectors.components.ConfigEditor.Database.input)).toBeInTheDocument();
  });
});
