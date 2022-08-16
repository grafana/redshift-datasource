import { SQLEditor } from '@grafana/experimental';
import { DataSource } from 'datasource';
import { getRedshiftCompletionProvider } from 'language/completionItemProvider';
import redshiftLanguageDefinition from 'language/definition';
import React, { useRef, useMemo, useCallback } from 'react';
import { RedshiftQuery } from 'types';

interface RawEditorProps {
  query: RedshiftQuery;
  onChange: (q: RedshiftQuery) => void;
  datasource: DataSource;
}

export default function RedshiftSQLEditor({ query, datasource, onChange }: RawEditorProps) {
  const getTables = useCallback(
    async (schema?: string) => {
      const tables: string[] = await datasource.postResource('tables', {
        // if schema is provided in the raw sql use that. if not, use schema defined in the query builder.
        schema: schema ?? query.schema,
      });
      return tables.map((table) => ({ name: table, completion: table }));
    },
    [query.schema]
  );

  const getColumns = useCallback(
    async (tableName?: string, schema?: string) => {
      const columns: string[] = await datasource.postResource('columns', {
        // if schema and table have been provided in the raw sql use that. if not, use schema/table defined in the query builder.
        schema: schema ?? query.schema,
        table: tableName ?? query.table,
      });
      return columns.map((column) => ({ name: column, completion: column }));
    },
    [query.schema]
  );

  const getColumnsRef = useRef(getColumns);
  const getTablesRef = useRef(getTables);
  const completionProvider = useMemo(
    () => getRedshiftCompletionProvider({ getTables: getTablesRef, getColumns: getColumnsRef }),
    []
  );

  return (
    <SQLEditor
      query={query.rawSQL}
      onChange={(rawSQL) => onChange({ ...query, rawSQL })}
      language={{
        ...redshiftLanguageDefinition,
        completionProvider,
      }}
    ></SQLEditor>
  );
}
