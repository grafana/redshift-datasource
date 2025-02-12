import { SQLEditor as SQLCodeEditor } from '@grafana/plugin-ui';
import { getTemplateSrv } from '@grafana/runtime';
import { DataSource } from 'datasource';
import { getRedshiftCompletionProvider } from 'language/completionItemProvider';
import redshiftLanguageDefinition from 'language/definition';
import { SCHEMA_MACRO, TABLE_MACRO } from 'language/macros';
import React, { useRef, useMemo, useCallback, useEffect } from 'react';
import { RedshiftQuery } from 'types';

interface RawEditorProps {
  query: RedshiftQuery;
  onRunQuery?: () => void;
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
    const schemas: string[] = await datasource.postResource<string[]>('schemas').catch(() => []);
    return schemas.map((schema) => ({ name: schema, completion: schema }));
  }, [datasource]);

  const getTables = useCallback(
    async (schema?: string) => {
      const tables: string[] = await datasource
        .postResource<string[]>('tables', {
          // if schema is provided in the raw sql use that. if not, use schema defined in the query builder.
          schema: interpolate(schema) ?? queryRef.current.schema,
        })
        .catch(() => []);
      return tables.map((table) => ({ name: table, completion: table }));
    },
    [datasource]
  );

  const getColumns = useCallback(
    async (tableName?: string, schema?: string) => {
      const columns: string[] = await datasource
        .postResource<string[]>('columns', {
          // if schema and table have been provided in the raw sql use that. if not, use schema/table defined in the query builder.
          schema: interpolate(schema) ?? queryRef.current.schema,
          table: interpolate(tableName) ?? queryRef.current.table,
        })
        .catch(() => []);
      return columns.map((column) => ({ name: column, completion: column }));
    },
    [datasource]
  );

  const completionProvider = useMemo(
    () => getRedshiftCompletionProvider({ getTables, getColumns, getSchemas }),
    [getTables, getColumns, getSchemas]
  );

  return (
    <SQLCodeEditor
      query={query.rawSQL}
      onChange={(rawSQL) => onChange({ ...queryRef.current, rawSQL })}
      language={{
        ...redshiftLanguageDefinition,
        completionProvider,
      }}
    ></SQLCodeEditor>
  );
}
