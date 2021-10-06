import { render, screen, waitFor } from '@testing-library/react';
import React from 'react';
import ResourceMacro, { Resource } from 'ResourceMacro';
import { mockQuery } from '__mocks__/datasource';
import { select } from 'react-select-event';

const defaultProps = {
  resource: 'table' as Resource,
  query: mockQuery,
  updateQuery: jest.fn(),
  loadOptions: jest.fn(),
};

describe('ResourceMacro', () => {
  it('should render a question mark if the value is not set', () => {
    render(<ResourceMacro {...defaultProps} />);
    expect(screen.getByText('$__table = ?')).toBeInTheDocument();
  });

  it('should render the resource value if set', () => {
    render(<ResourceMacro {...defaultProps} query={{ ...defaultProps.query, table: 'foo' }} />);
    expect(screen.getByText('$__table = foo')).toBeInTheDocument();
  });

  it('should load the resource options', async () => {
    const loadOptions = jest.fn().mockResolvedValue([]);
    render(<ResourceMacro {...defaultProps} loadOptions={loadOptions} />);
    const node = screen.getByText('$__table = ?');
    node.click();
    expect(loadOptions).toHaveBeenCalled();
    await waitFor(() => screen.getByText('No options found'));
  });

  it('should change the selected option', async () => {
    const loadOptions = jest.fn().mockResolvedValue([
      { label: 'foo', value: 'foo' },
      { label: 'bar', value: 'bar' },
    ]);
    const updateQuery = jest.fn();
    render(
      <ResourceMacro
        {...defaultProps}
        query={{ ...defaultProps.query, table: 'foo' }}
        loadOptions={loadOptions}
        updateQuery={updateQuery}
      />
    );
    const selectEl = screen.getByText('$__table = foo');
    expect(selectEl).toBeInTheDocument();
    selectEl.click();
    await select(selectEl, 'bar', { container: document.body });

    expect(updateQuery).toHaveBeenCalledWith({ ...mockQuery, table: 'bar' });
  });
});
