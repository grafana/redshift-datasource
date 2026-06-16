import '@testing-library/jest-dom';

import { render, screen, waitFor } from '@testing-library/react';
import userEvent from '@testing-library/user-event';
import React from 'react';
import { select } from 'react-select-event';
import { FillValueOptions } from '@grafana/aws-sdk';
import { FormatOptions } from 'types';
import * as experimental from '@grafana/plugin-ui';

import { mockDatasource, mockQuery } from './__mocks__/datasource';
import { QueryEditorForm } from './QueryEditorForm';

jest.mock('@grafana/plugin-ui', () => ({
  ...jest.requireActual<typeof experimental>('@grafana/plugin-ui'),
  SQLEditor: function SQLEditor() {
    return <></>;
  },
}));

const ds = mockDatasource;
const q = mockQuery;

beforeEach(() => {
  ds.getResource = jest.fn().mockResolvedValue([]);
  ds.postResource = jest.fn().mockResolvedValue([]);
});

const props = {
  datasource: ds,
  query: q,
  onChange: jest.fn(),
  onRunQuery: jest.fn(),
};

describe('QueryEditorForm', () => {
  it('should request select schemas but not execute the query', async () => {
    const onChange = jest.fn();
    const onRunQuery = jest.fn();
    ds.getResource = jest.fn().mockResolvedValue(['foo']);
    render(
      <QueryEditorForm
        {...props}
        onChange={onChange}
        onRunQuery={onRunQuery}
        query={{ ...props.query, schema: 'bar' }}
      />
    );

    const selectEl = screen.getByLabelText('Schema');
    expect(selectEl).toBeInTheDocument();

    await select(selectEl, 'foo', { container: document.body });

    expect(ds.getResource).toHaveBeenCalledWith('schemas');
    expect(onChange).toHaveBeenCalledWith({
      ...q,
      schema: 'foo',
    });
    expect(onRunQuery).not.toHaveBeenCalled();
  });

  it('should request select tables but not execute the query', async () => {
    const onChange = jest.fn();
    const onRunQuery = jest.fn();
    ds.postResource = jest.fn().mockResolvedValue(['foo']);
    render(
      <QueryEditorForm
        {...props}
        onChange={onChange}
        onRunQuery={onRunQuery}
        query={{ ...props.query, schema: 'bar' }}
      />
    );

    const selectEl = screen.getByLabelText('Table');
    expect(selectEl).toBeInTheDocument();

    await select(selectEl, 'foo', { container: document.body });

    expect(ds.postResource).toHaveBeenCalledWith('tables', { schema: 'bar' });
    expect(onChange).toHaveBeenCalledWith({
      ...q,
      schema: 'bar',
      table: 'foo',
    });
    expect(onRunQuery).not.toHaveBeenCalled();
  });

  it('should request select column but not execute the query', async () => {
    const onChange = jest.fn();
    const onRunQuery = jest.fn();
    ds.postResource = jest.fn().mockResolvedValue(['foo']);
    render(
      <QueryEditorForm
        {...props}
        onChange={onChange}
        onRunQuery={onRunQuery}
        query={{ ...props.query, table: 'bar' }}
      />
    );

    const selectEl = screen.getByLabelText('Column');
    expect(selectEl).toBeInTheDocument();

    await select(selectEl, 'foo', { container: document.body });

    expect(ds.postResource).toHaveBeenCalledWith('columns', { table: 'bar' });
    expect(onChange).toHaveBeenCalledWith({
      ...q,
      table: 'bar',
      column: 'foo',
    });
    expect(onRunQuery).not.toHaveBeenCalled();
  });

  it('should include the Format As input', async () => {
    render(<QueryEditorForm {...props} />);
    await openFormatCollapse();
    await waitFor(() => screen.getByText('Format data frames as'));
  });

  it('should skip the fill mode input if the format is not TimeSeries', async () => {
    const onChange = jest.fn();
    render(
      <QueryEditorForm
        {...props}
        query={{ ...props.query, format: FormatOptions.Table }}
        queries={[]}
        onChange={onChange}
      />
    );
    await openFormatCollapse();
    await waitFor(() => screen.getByText('Format data frames as'));
    const selectEl = screen.queryByLabelText('Fill with');
    expect(selectEl).not.toBeInTheDocument();
  });

  it('should allow to change the fill mode', async () => {
    const onChange = jest.fn();
    render(<QueryEditorForm {...props} queries={[]} onChange={onChange} />);
    await openFormatCollapse();
    const selectEl = screen.getByLabelText('Fill with');
    expect(selectEl).toBeInTheDocument();

    await select(selectEl, 'NULL', { container: document.body });

    expect(onChange).toHaveBeenCalledWith({
      ...q,
      fillMode: { mode: FillValueOptions.Null },
    });
  });
});

async function openFormatCollapse() {
  const collapseTitle = screen.getByTestId('collapse-title');
  return userEvent.click(collapseTitle);
}
