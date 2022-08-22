import { SQLEditor as SQLCodeEditor } from '@grafana/experimental';
import { getTemplateSrv } from '@grafana/runtime';
import { DataSource } from 'datasource';
import { getRedshiftCompletionProvider } from 'language/completionItemProvider';
import redshiftLanguageDefinition from 'language/definition';
import { SCHEMA_MACRO, TABLE_MACRO } from 'language/macros';
import React, { useRef, useMemo, useCallback, useEffect } from 'react';
import { RedshiftQuery } from 'types';

interface RawEditorProps {
  query: RedshiftQuery;
  onRunQuery: () => void;
  onChange: (q: RedshiftQuery) => void;
  datasource: DataSource;
}

export default function SQLEditor({ query, datasource, onRunQuery, onChange }: RawEditorProps) {
  const queryRef = useRef<RedshiftQuery>(query);
  useEffect(() => {
    queryRef.current = query;
  }, [query]);

  const interpolate = (value: string | undefined) => {
    if (!value) {
      return value;
    }

    value = value.replace(SCHEMA_MACRO, queryRef.current.schema ?? '');
    value = value.replace(TABLE_MACRO, queryRef.current.table ?? '');
    value = getTemplateSrv().replace(value);

    return value;
  };

  const getSchemas = useCallback(async () => {
    const schemas: string[] = await datasource.postResource('schemas');
    return schemas.map((schema) => ({ name: schema, completion: schema }));
  }, [datasource]);

  const getTables = useCallback(
    async (schema?: string) => {
      const tables: string[] = await datasource.postResource('tables', {
        // if schema is provided in the raw sql use that. if not, use schema defined in the query builder.
        schema: interpolate(schema) ?? queryRef.current.schema,
      });
      return tables.map((table) => ({ name: table, completion: table }));
    },
    [datasource]
  );

  const getColumns = useCallback(
    async (tableName?: string, schema?: string) => {
      const columns: string[] = await datasource.postResource('columns', {
        // if schema and table have been provided in the raw sql use that. if not, use schema/table defined in the query builder.
        schema: interpolate(schema) ?? queryRef.current.schema,
        table: interpolate(tableName) ?? queryRef.current.table,
      });
      return columns.map((column) => ({ name: column, completion: column }));
    },
    [datasource]
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
    <SQLCodeEditor
      query={query.rawSQL}
      onBlur={() => onRunQuery()}
      onChange={(rawSQL) => onChange({ ...queryRef.current, rawSQL })}
      language={{
        ...redshiftLanguageDefinition,
        completionProvider,
      }}
    ></SQLCodeEditor>
  );
}
