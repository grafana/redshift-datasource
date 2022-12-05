import '@testing-library/jest-dom';

import * as runtime from '@grafana/runtime';
import { fireEvent, render, screen, waitFor } from '@testing-library/react';
import React from 'react';
import { select } from 'react-select-event';
import { FillValueOptions } from '@grafana/aws-sdk';
import { FormatOptions } from 'types';
import * as experimental from '@grafana/experimental';

import { mockDatasource, mockQuery } from './__mocks__/datasource';
import { QueryEditor } from './QueryEditor';

jest.mock('@grafana/experimental', () => ({
  ...jest.requireActual<typeof experimental>('@grafana/experimental'),
  SQLEditor: function SQLEditor() {
    return <></>;
  },
}));

jest.mock('@grafana/runtime', () => ({
  ...jest.requireActual<typeof runtime>('@grafana/runtime'),
  config: {
    featureToggles: {
      redshiftAsyncQueryDataSupport: true,
    },
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

describe('QueryEditor', () => {
  it('should request select schemas but not execute the query', async () => {
    const onChange = jest.fn();
    const onRunQuery = jest.fn();
    ds.getResource = jest.fn().mockResolvedValue(['foo']);
    render(
      <QueryEditor {...props} onChange={onChange} onRunQuery={onRunQuery} query={{ ...props.query, schema: 'bar' }} />
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
      <QueryEditor {...props} onChange={onChange} onRunQuery={onRunQuery} query={{ ...props.query, schema: 'bar' }} />
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
      <QueryEditor {...props} onChange={onChange} onRunQuery={onRunQuery} query={{ ...props.query, table: 'bar' }} />
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
    render(<QueryEditor {...props} />);
    await waitFor(() => screen.getByText('Format as'));
  });

  it('should skip the fill mode input if the format is not TimeSeries', async () => {
    const onChange = jest.fn();
    render(
      <QueryEditor
        {...props}
        query={{ ...props.query, format: FormatOptions.Table }}
        queries={[]}
        onChange={onChange}
      />
    );
    await waitFor(() => screen.getByText('Format as'));
    const selectEl = screen.queryByLabelText('Fill value');
    expect(selectEl).not.toBeInTheDocument();
  });

  it('should allow to change the fill mode', async () => {
    const onChange = jest.fn();
    render(<QueryEditor {...props} queries={[]} onChange={onChange} />);
    const selectEl = screen.getByLabelText('Fill value');
    expect(selectEl).toBeInTheDocument();

    await select(selectEl, 'NULL', { container: document.body });

    expect(onChange).toHaveBeenCalledWith({
      ...q,
      fillMode: { mode: FillValueOptions.Null },
    });
  });

  it('run button should be disabled if the query is not valid', () => {
    render(<QueryEditor {...props} query={{ ...props.query, rawSQL: '' }} />);
    const runButton = screen.getByRole('button', { name: 'Run' });
    expect(runButton).toBeDisabled();
  });

  it('should run queries when the run button is clicked', () => {
    const onChange = jest.fn();
    const onRunQuery = jest.fn();
    render(<QueryEditor {...props} onRunQuery={onRunQuery} onChange={onChange} />);
    const runButton = screen.getByRole('button', { name: 'Run' });
    expect(runButton).toBeInTheDocument();

    expect(onRunQuery).not.toBeCalled();
    fireEvent.click(runButton);
    expect(onRunQuery).toBeCalledTimes(1);
  });

  it('stop button should be disabled until run button is clicked', () => {
    render(<QueryEditor {...props} />);
    const runButton = screen.getByRole('button', { name: 'Run' });
    const stopButton = screen.getByRole('button', { name: 'Stop' });
    expect(stopButton).toBeInTheDocument();
    expect(stopButton).toBeDisabled();
    fireEvent.click(runButton);
    expect(stopButton).not.toBeDisabled();
  });
});
