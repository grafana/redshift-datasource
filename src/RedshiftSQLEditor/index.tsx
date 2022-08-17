import { SQLEditor } from '@grafana/experimental';
import { getTemplateSrv } from '@grafana/runtime';
import { DataSource } from 'datasource';
import { getRedshiftCompletionProvider } from 'language/completionItemProvider';
import redshiftLanguageDefinition from 'language/definition';
import { SCHEMA_MACRO, TABLE_MACRO } from 'language/macros';
import React, { useRef, useMemo, useCallback } from 'react';
import { RedshiftQuery } from 'types';

interface RawEditorProps {
  query: RedshiftQuery;
  onChange: (q: RedshiftQuery) => void;
  datasource: DataSource;
}

export default function RedshiftSQLEditor({ query, datasource, onChange }: RawEditorProps) {
  const interpolate = (value: string | undefined) => {
    if (!value) {
      return value;
    }

    value = value.replace(SCHEMA_MACRO, query.schema ?? '');
    value = value.replace(TABLE_MACRO, query.table ?? '');
    value = getTemplateSrv().replace(value);

    return value;
  };

  const getSchemas = useCallback(async () => {
    const schemas: string[] = await datasource.postResource('schemas');
    return schemas.map((schema) => ({ name: schema, completion: schema }));
  }, [query.schema]);

  const getTables = useCallback(
    async (schema?: string) => {
      const tables: string[] = await datasource.postResource('tables', {
        // if schema is provided in the raw sql use that. if not, use schema defined in the query builder.
        schema: interpolate(schema) ?? query.schema,
      });
      return tables.map((table) => ({ name: table, completion: table }));
    },
    [query.schema]
  );

  const getColumns = useCallback(
    async (tableName?: string, schema?: string) => {
      const columns: string[] = await datasource.postResource('columns', {
        // if schema and table have been provided in the raw sql use that. if not, use schema/table defined in the query builder.
        schema: interpolate(schema) ?? query.schema,
        table: interpolate(tableName) ?? query.table,
      });
      return columns.map((column) => ({ name: column, completion: column }));
    },
    [query.schema]
  );

  const getSchemasRef = useRef(getSchemas);
  const getTablesRef = useRef(getTables);
  const getColumnsRef = useRef(getColumns);
  const completionProvider = useMemo(
    () =>
      getRedshiftCompletionProvider({ getTables: getTablesRef, getColumns: getColumnsRef, getSchemas: getSchemasRef }),
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
