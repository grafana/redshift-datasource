import React from 'react';
import { render, screen } from '@testing-library/react';
import { QueryCodeEditor } from './QueryCodeEditor';
import { mockDatasource, mockQuery } from './__mocks__/datasource';
import '@testing-library/jest-dom';
import { select } from 'react-select-event';

const ds = mockDatasource;
const q = mockQuery;

const props = {
  datasource: ds,
  query: q,
  onChange: jest.fn(),
  onRunQuery: jest.fn(),
};

beforeEach(() => {
  ds.getResource = jest.fn().mockResolvedValue([]);
  ds.postResource = jest.fn().mockResolvedValue([]);
});

describe('QueryCodeEditor', () => {
  it('should list and select an schema', async () => {
    ds.getResource = jest.fn().mockResolvedValue(['public', 'foo']);
    const onChange = jest.fn();
    render(<QueryCodeEditor {...props} queries={[]} onChange={onChange} />);
    const selectEl = screen.getByText('$__schema = public');
    expect(selectEl).toBeInTheDocument();
    selectEl.click();

    await select(selectEl, 'foo', { container: document.body });

    expect(onChange).toHaveBeenCalledWith({
      ...q,
      schema: 'foo',
    });
  });

  it('should list and select a table', async () => {
    ds.postResource = jest.fn().mockResolvedValue(['foo', 'bar']);
    const onChange = jest.fn();
    render(<QueryCodeEditor {...props} queries={[]} onChange={onChange} />);
    const selectEl = screen.getByText('$__table = ?');
    expect(selectEl).toBeInTheDocument();
    selectEl.click();

    await select(selectEl, 'foo', { container: document.body });

    expect(onChange).toHaveBeenCalledWith({
      ...q,
      table: 'foo',
    });
  });

  it('should list and select a column', async () => {
    ds.postResource = jest.fn().mockResolvedValue(['foo', 'bar']);
    const onChange = jest.fn();
    render(<QueryCodeEditor {...props} queries={[]} onChange={onChange} />);
    const selectEl = screen.getByText('$__column = ?');
    expect(selectEl).toBeInTheDocument();
    selectEl.click();

    await select(selectEl, 'foo', { container: document.body });

    expect(onChange).toHaveBeenCalledWith({
      ...q,
      column: 'foo',
    });
  });
});
