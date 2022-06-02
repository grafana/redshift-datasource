import { SQLEditor } from '@grafana/experimental';
import { DataSource } from 'datasource';
import React, { useCallback } from 'react';
import { RedshiftQuery } from 'types';

interface RawEditorProps {
  query: RedshiftQuery;
  onChange: (q: RedshiftQuery, processQuery: boolean) => void;
  datasource: DataSource;
}

export default function RedshiftSQLEditor({ query, onChange, datasource }: RawEditorProps) {
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

  return <SQLEditor query={query.rawSQL} onChange={onRawQueryChange}></SQLEditor>;
}
