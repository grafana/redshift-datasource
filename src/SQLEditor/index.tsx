import { DataSource } from 'datasource';
import React, { useCallback } from 'react';
import { RedshiftQuery } from 'types';

import { SQLEditorRaw } from './SQLEditor';

interface RawEditorProps {
  query: RedshiftQuery;
  onChange: (q: RedshiftQuery, processQuery: boolean) => void;
  datasource: DataSource;
}

export default function SQLEditor({ query, onChange, datasource }: RawEditorProps) {
  const getTables = useCallback(
    async (d?: string) => {
      const tables: string[] = await datasource.postResource('tables', {
        schema: query.schema || '',
      });
      return tables.map((table) => ({ name: table, completion: table }));
    },
    [query.schema]
  );

  const getColumnsColumns = useCallback(
    async (d?: string) => {
      const columns: string[] = await datasource.postResource('columns', {
        schema: query.schema,
        table: query.table,
      });
      return columns.map((column) => ({ name: column, completion: column }));
    },
    [query.schema]
  );

  const getSchema = useCallback(
    async (d?: string) => {
      const schemas: string[] = await datasource.getResource('schemas');
      return schemas.map((schema) => ({ name: schema, completion: schema }));
    },
    [query.schema]
  );

  return (
    <SQLEditorRaw
      query={query}
      onChange={onChange}
      getTables={getTables}
      getColumns={getColumnsColumns}
      getTableSchema={getSchema}
    />
  );
}
