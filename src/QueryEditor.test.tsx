import React from 'react';
import { render, screen } from '@testing-library/react';
import { QueryEditor } from './QueryEditor';
import { mockDatasource, mockQuery } from './__mocks__/datasource';
import '@testing-library/jest-dom';

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
    expect(screen.getByText('$__schema = public')).toBeInTheDocument();
    expect(screen.getByText('$__table = ?')).toBeInTheDocument();
    expect(screen.getByText('$__column = ?')).toBeInTheDocument();
  });

  it('should not include the Format As input if the query editor does not support multiple queries', async () => {
    render(<QueryEditor {...props} queries={undefined} />);
    expect(screen.queryByText('Format as')).not.toBeInTheDocument();
  });

  it('should include the Format As input', async () => {
    render(<QueryEditor {...props} queries={[]} />);
    expect(screen.queryByText('Format as')).toBeInTheDocument();
  });
});
