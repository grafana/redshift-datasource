import { ColumnDefinition, SQLEditor, TableDefinition } from '@grafana/experimental';
import { getRedshiftCompletionProvider } from 'language/redshiftCompletionItemProvider';
// import { CodeEditor, Monaco } from '@grafana/ui';
import React, { useCallback, useEffect, useMemo, useRef } from 'react';
import { RedshiftQuery } from 'types';
import lang from '../language/definition';

type Props = {
  query: RedshiftQuery;
  getTables: (d?: string) => Promise<TableDefinition[]>;
  getColumns: (t: string) => Promise<ColumnDefinition[]>;
  getTableSchema: (l: string) => Promise<TableDefinition[]>;
  onChange: (value: RedshiftQuery, processQuery: boolean) => void;
  children?: (props: { formatQuery: () => void }) => React.ReactNode;
};

export function SQLEditorRaw({
  children,
  getColumns: apiGetColumns,
  getTables: apiGetTables,
  getTableSchema: apiGetTableSchema,
  onChange,
  query,
}: Props) {
  const getColumns = useRef<Props['getColumns']>(apiGetColumns);
  const getTables = useRef<Props['getTables']>(apiGetTables);
  const getTableSchema = useRef<Props['getTableSchema']>(apiGetTableSchema);
  const completionProvider = useMemo(
    () => getRedshiftCompletionProvider({ getTables, getColumns, getTableSchema }),
    []
  );

  useEffect(() => {
    getColumns.current = apiGetColumns;
    // getTables.current = apiGetTables;
  }, [apiGetColumns, apiGetTables]);

  const onRawQueryChange = useCallback(
    (rawSQL: string, processQuery: boolean) => {
      const newQuery = {
        ...query,
        rawQuery: true,
        rawSQL,
      };
      onChange(newQuery, processQuery);
    },
    [onChange, query]
  );

  return (
    <>
      <SQLEditor query={query.rawSQL} onChange={onRawQueryChange} language={{ ...lang, completionProvider }}>
        {/* <SQLEditor query={query.rawSQL} onChange={onRawQueryChange} language={{ completionProvider: () => {} }}> */}
        {children}
      </SQLEditor>
    </>
  );
}
