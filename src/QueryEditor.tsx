import React, { useCallback, useEffect, useState } from 'react';
import { config } from '@grafana/runtime';
import { QueryEditorProps } from '@grafana/data';
import { QueryEditorHeader } from '@grafana/aws-sdk';
import { RedshiftDataSourceOptions, RedshiftQuery } from './types';
import { DataSource } from './datasource';
import { QueryEditorForm } from './QueryEditorForm';

export function QueryEditor(props: QueryEditorProps<DataSource, RedshiftQuery, RedshiftDataSourceOptions>) {
  const [dataIsStale, setDataIsStale] = useState(false);
  const { onChange } = props;

  useEffect(() => {
    setDataIsStale(false);
  }, [props.data]);

  const onChangeInternal = useCallback(
    (query: RedshiftQuery) => {
      setDataIsStale(true);
      onChange(query);
    },
    [onChange]
  );

  return (
    <>
      {props?.app !== 'explore' && (
        <QueryEditorHeader<DataSource, RedshiftQuery, RedshiftDataSourceOptions>
          {...props}
          enableRunButton={dataIsStale && !!props.query.rawSQL}
          showAsyncQueryButtons={config.featureToggles.redshiftAsyncQueryDataSupport}
          cancel={config.featureToggles.redshiftAsyncQueryDataSupport ? props.datasource.cancel : undefined}
        />
      )}
      <QueryEditorForm {...props} onChange={onChangeInternal} />
    </>
  );
}
