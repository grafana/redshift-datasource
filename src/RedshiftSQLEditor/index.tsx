import { SQLEditor } from '@grafana/experimental';
import React, { useCallback } from 'react';
import { RedshiftQuery } from 'types';

interface RawEditorProps {
  query: RedshiftQuery;
  onChange: (q: RedshiftQuery) => void;
}

export default function RedshiftSQLEditor({ query, onChange }: RawEditorProps) {
  const onRawQueryChange = useCallback(
    (rawSQL: string) => {
      onChange({ ...query, rawSQL });
    },
    [onChange, query]
  );

  return <SQLEditor query={query.rawSQL} onChange={onRawQueryChange}></SQLEditor>;
}
