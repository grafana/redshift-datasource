import { SelectableValue } from '@grafana/data';
import { SegmentAsync } from '@grafana/ui';
import React from 'react';
import { SchemaInfo } from 'SchemaInfo';
import { RedshiftQuery } from 'types';

export type Resource = 'schema' | 'table' | 'column';

interface Props {
  resource: Resource;
  query: RedshiftQuery;
  schema: SchemaInfo;
  updateSchemaState: (query: RedshiftQuery) => void;
  value?: string;
}

function ResourceMacro({ resource, schema, updateSchemaState, query, value }: Props) {
  let placeholder = '';
  let current = '$__' + resource + ' = ';
  if (query[resource]) {
    current += query[resource];
  } else {
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
      updateSchemaState(newQuery);
    };
  };

  const loadOptions = {
    schema: schema.getSchemas,
    table: schema.getTables,
    column: schema.getColumns,
  };

  return (
    <SegmentAsync
      value={current}
      loadOptions={loadOptions[resource]}
      placeholder={placeholder}
      onChange={onChanged(resource)}
      allowCustomValue
    />
  );
}

export default ResourceMacro;
