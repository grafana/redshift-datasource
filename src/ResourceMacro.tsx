import { SelectableValue } from '@grafana/data';
import { SegmentAsync } from '@grafana/ui';
import React from 'react';
import { RedshiftQuery } from 'types';

export type Resource = 'schema' | 'table' | 'column';

interface Props {
  resource: Resource;
  query: RedshiftQuery;
  updateQuery: (q: RedshiftQuery) => void;
  loadOptions: (query?: string) => Promise<Array<SelectableValue<string>>>;
}

function ResourceMacro({ resource, query, updateQuery, loadOptions }: Props) {
  let placeholder = '';
  let current = '$__' + resource + ' = ';
  if (query[resource]) {
    current += query[resource];
  } else {
    let value = query[resource];
    if (!value && resource === 'schema') {
      // Use the public schema by default
      value = 'public';
    }
    placeholder = current + (value ?? '?');
    current = '';
  }

  const onChanged = (resource: Resource) => {
    return (value: SelectableValue<string>) => {
      const newQuery = {
        ...query,
        [resource]: value.value,
      };
      if (!newQuery[resource]) {
        delete newQuery[resource];
      }
      if (resource === 'schema') {
        // Clean up table and column since a new schema is set
        newQuery.table = undefined;
        newQuery.column = undefined;
      }
      if (resource === 'table') {
        // Clean up column since a new table is set
        newQuery.column = undefined;
      }
      updateQuery(newQuery);
    };
  };

  return (
    <SegmentAsync
      value={current}
      loadOptions={loadOptions}
      placeholder={placeholder}
      onChange={onChanged(resource)}
      allowCustomValue
    />
  );
}

export default ResourceMacro;
