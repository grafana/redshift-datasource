import { SchemaInfo } from 'SchemaInfo';
import { mockDatasource, mockQuery } from '__mocks__/datasource';

const ds = mockDatasource;
const q = mockQuery;

describe('SchemaInfo', () => {
  describe('constructor', () => {
    it("should select the 'public' schema by default", () => {
      const schema = new SchemaInfo(ds, q);
      expect(schema.state.schema).toEqual('public');
    });
  });

  describe('getSuggestions', () => {
    const macros = [
      '$__timeEpoch',
      '$__timeFilter',
      '$__timeFrom',
      '$__timeTo',
      '$__timeGroup',
      '$__unixEpochFilter',
      '$__unixEpochGroup',
      '$__schema',
      '$__table',
      '$__column',
    ];
    it('should return the list of macros', () => {
      const schema = new SchemaInfo(ds, q);
      const sugs = schema.getSuggestions();
      expect(sugs.map((s) => s.label)).toEqual(macros);
    });

    it('should return the list of template variables', () => {
      const templateSrv = {
        getVariables: jest.fn().mockReturnValue([{ name: 'foo' }, { name: 'bar' }]),
        replace: jest.fn(),
      };
      const schema = new SchemaInfo(ds, q, templateSrv);
      const sugs = schema.getSuggestions();
      expect(sugs.map((s) => s.label)).toEqual(macros.concat('$foo', '$bar'));
    });
  });

  describe('updateState', () => {
    it('updates the schema in the state', () => {
      const schema = new SchemaInfo(ds, q);
      schema.updateState({ schema: 'foo' });
      expect(schema.state).toEqual({ ...q, schema: 'foo' });
    });

    it('cleans up tables and columns if the schema changes', () => {
      const schema = new SchemaInfo(ds, q);
      schema.tables = [{ label: 'foo', value: 'foo' }];
      schema.columns = [{ label: 'bar', value: 'bar' }];
      schema.updateState({ schema: 'foobar' });
      expect(schema.tables).toBeUndefined();
      expect(schema.columns).toBeUndefined();
    });

    it('sets a table in the state', () => {
      const schema = new SchemaInfo(ds, q);
      schema.updateState({ table: 'foo' });
      expect(schema.state.table).toEqual('foo');
    });

    it('cleans up columns if the table changes', () => {
      const schema = new SchemaInfo(ds, q);
      schema.columns = [{ label: 'bar', value: 'bar' }];
      schema.updateState({ table: 'foobar' });
      expect(schema.columns).toBeUndefined();
    });

    it('sets a column in the state', () => {
      const schema = new SchemaInfo(ds, q);
      schema.updateState({ column: 'foo' });
      expect(schema.state.column).toEqual('foo');
    });

    it('uses the templateSrv to replace values', () => {
      const templateSrv = {
        getVariables: jest.fn(),
        replace: (s: string) => s.replace('$', ''),
      };
      const schema = new SchemaInfo(ds, q, templateSrv);
      schema.updateState({ schema: '$foobar', table: '$foo', column: '$bar' });
      expect(schema.state).toEqual({ ...q, schema: 'foobar', table: 'foo', column: 'bar' });
    });
  });

  describe('getSchemas', () => {
    it('should return cached schemas', async () => {
      const schema = new SchemaInfo(ds, q);
      const schemas = [{ label: 'foo', value: 'foo' }];
      schema.schemas = schemas;
      const res = await schema.getSchemas();
      expect(res).toEqual(schemas);
    });

    it('should get schemas as a resource', async () => {
      ds.getResource = jest.fn().mockResolvedValue(['foo', 'bar']);
      const schema = new SchemaInfo(ds, q);
      const res = await schema.getSchemas();
      expect(res).toEqual([
        { label: 'foo', value: 'foo' },
        { label: 'bar', value: 'bar' },
        { label: '-- remove --', value: '' },
      ]);
    });
  });

  describe('getTables', () => {
    it('should return cached tables', async () => {
      const schema = new SchemaInfo(ds, q);
      const tables = [{ label: 'foo', value: 'foo' }];
      schema.tables = tables;
      const res = await schema.getTables();
      expect(res).toEqual(tables);
    });

    it('should get tables as a resource', async () => {
      ds.postResource = jest.fn().mockResolvedValue(['foo', 'bar']);
      const schema = new SchemaInfo(ds, q);
      const res = await schema.getTables();
      expect(res).toEqual([
        { label: 'foo', value: 'foo' },
        { label: 'bar', value: 'bar' },
        { label: '-- remove --', value: '' },
      ]);
    });

    it('should get tables as a resource', async () => {
      ds.postResource = jest.fn().mockResolvedValue(['foo', 'bar']);
      const schema = new SchemaInfo(ds, q);
      await schema.getTables();
      expect(ds.postResource).toHaveBeenCalledWith('tables', { schema: 'public' });
    });
  });

  describe('getColumns', () => {
    it('should return cached columns', async () => {
      const schema = new SchemaInfo(ds, q);
      const columns = [{ label: 'foo', value: 'foo' }];
      schema.columns = columns;
      const res = await schema.getColumns();
      expect(res).toEqual(columns);
    });

    it('should get columns as a resource', async () => {
      ds.postResource = jest.fn().mockResolvedValue(['foo', 'bar']);
      const schema = new SchemaInfo(ds, q);
      schema.state.table = 'foobar';
      const res = await schema.getColumns();
      expect(res).toEqual([
        { label: 'foo', value: 'foo' },
        { label: 'bar', value: 'bar' },
        { label: '-- remove --', value: '' },
      ]);
      expect(ds.postResource).toHaveBeenCalledWith('columns', { table: 'foobar' });
    });

    it('should return empty if the table is not set', async () => {
      const schema = new SchemaInfo(ds, q);
      const res = await schema.getColumns();
      expect(res).toEqual([{ label: 'table not configured', value: '' }]);
    });
  });
});
