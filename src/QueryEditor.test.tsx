import React from 'react';
import { render, screen, waitFor } from '@testing-library/react';
import { QueryEditor } from './QueryEditor';
import { mockDatasource, mockQuery } from './__mocks__/datasource';
import '@testing-library/jest-dom';
import { select } from 'react-select-event';
import { FillValueOptions } from 'types';

const ds = mockDatasource;
const q = mockQuery;
ds.getResource = jest.fn().mockResolvedValue([]);
ds.postResource = jest.fn().mockResolvedValue([]);

const props = {
  datasource: ds,
  query: q,
  onChange: jest.fn(),
  onRunQuery: jest.fn(),
};

describe('QueryEditor', () => {
  it('should render Macros input', async () => {
    render(<QueryEditor {...props} />);
    await waitFor(() => screen.getByText('$__schema = public'));
    expect(screen.getByText('$__table = ?')).toBeInTheDocument();
    expect(screen.getByText('$__column = ?')).toBeInTheDocument();
  });

  it('should include the Format As input', async () => {
    render(<QueryEditor {...props} />);
    await waitFor(() => screen.getByText('Format as'));
  });

  it('should allow to change the fill mode', async () => {
    const onChange = jest.fn();
    render(<QueryEditor {...props} queries={[]} onChange={onChange} />);
    const selectEl = screen.getByLabelText('Fill value');
    expect(selectEl).toBeInTheDocument();

    await select(selectEl, 'NULL', { container: document.body });

    expect(onChange).toHaveBeenCalledWith({
      ...q,
      fillMode: { mode: FillValueOptions.Null, value: 0 },
    });
  });
});
