import { render, screen } from '@testing-library/react';
import React from 'react';
import ResourceMacro, { Resource } from 'ResourceMacro';
import { SchemaInfo } from 'SchemaInfo';
import { mockDatasource, mockQuery, mockSchemaInfo } from '__mocks__/datasource';
import userEvent from '@testing-library/user-event';

const defaultProps = {
  resource: 'table' as Resource,
  query: mockQuery,
  schema: mockSchemaInfo,
  updateSchemaState: jest.fn(),
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

  it('should load the resource options', () => {
    const schema = new SchemaInfo(mockDatasource, mockQuery);
    schema.getTables = jest.fn().mockReturnValue({ then: jest.fn() });
    render(<ResourceMacro {...defaultProps} schema={schema} />);
    const node = screen.getByText('$__table = ?');
    userEvent.click(node);
    expect(schema.getTables).toHaveBeenCalled();
  });

  it('should change the selected option', async () => {
    const updateSchemaState = jest.fn();
    const schema = new SchemaInfo(mockDatasource, mockQuery);
    schema.getTables = jest.fn().mockResolvedValue([
      { label: 'foo', value: 'foo' },
      { label: 'bar', value: 'bar' },
    ]);
    render(
      <ResourceMacro
        {...defaultProps}
        query={{ ...defaultProps.query, table: 'foo' }}
        schema={schema}
        updateSchemaState={updateSchemaState}
      />
    );
    // TODO: investigate why this throws a console.log error in our test suite
    userEvent.click(screen.getByText('$__table = foo'));
    expect(schema.getTables).toHaveBeenCalled();
    await screen.findByText('bar');
    userEvent.click(screen.getByText('bar'));
    expect(updateSchemaState).toHaveBeenCalledWith({ ...mockQuery, table: 'bar' });
  });
});
