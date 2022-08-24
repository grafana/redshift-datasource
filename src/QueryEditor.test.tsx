import '@testing-library/jest-dom';

import * as runtime from '@grafana/runtime';
import { render, screen, waitFor } from '@testing-library/react';
import React from 'react';
import { select } from 'react-select-event';
import { FillValueOptions } from '@grafana/aws-sdk';
import { FormatOptions } from 'types';
import * as experimental from '@grafana/experimental';

import { mockDatasource, mockQuery } from './__mocks__/datasource';
import { QueryEditor } from './QueryEditor';

jest.spyOn(runtime, 'getTemplateSrv').mockImplementation(() => ({
  getVariables: jest.fn().mockReturnValue([]),
  replace: jest.fn(),
  containsTemplate: jest.fn(),
  updateTimeRange: jest.fn(),
}));

jest.mock('@grafana/experimental', () => ({
  ...jest.requireActual<typeof experimental>('@grafana/experimental'),
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

describe('QueryEditor', () => {
  it('should request select schemas and execute the query', async () => {
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

  it('should request select tables and execute the query', async () => {
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

  it('should request select column and execute the query', async () => {
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
});
