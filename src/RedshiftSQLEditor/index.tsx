import { SQLEditor } from '@grafana/experimental';
import React from 'react';
import { RedshiftQuery } from 'types';

interface RawEditorProps {
  query: RedshiftQuery;
  onChange: (q: RedshiftQuery) => void;
}

export default function RedshiftSQLEditor({ query, onChange }: RawEditorProps) {
  return <SQLEditor query={query.rawSQL} onChange={(rawSQL) => onChange({ ...query, rawSQL })}></SQLEditor>;
}
