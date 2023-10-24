import { FillValueSelect, FormatSelect, ResourceSelector } from '@grafana/aws-sdk';
import { QueryEditorProps, SelectableValue } from '@grafana/data';
import { CollapsableSection, InlineSegmentGroup } from '@grafana/ui';
import React from 'react';
import { selectors } from 'selectors';
import SQLEditor from './SQLEditor';

import { DataSource } from './datasource';
import { FormatOptions, RedshiftDataSourceOptions, RedshiftQuery, SelectableFormatOptions } from './types';
import { config } from '@grafana/runtime';
import { EditorField, EditorFieldGroup, EditorRow, EditorRows } from '@grafana/experimental';
import { css } from '@emotion/css';

type Props = QueryEditorProps<DataSource, RedshiftQuery, RedshiftDataSourceOptions>;

type QueryProperties = 'schema' | 'table' | 'column';

export function QueryEditorForm(props: Props) {
  const newFormStylingEnabled = config.featureToggles.awsDatasourcesNewFormStyling;
  const styles = getStyles;

  const fetchSchemas = async () => {
    const schemas: string[] = await props.datasource.getResource('schemas');
    return schemas.map((schema) => ({ label: schema, value: schema })).concat({ label: '-- remove --', value: '' });
  };

  const fetchTables = async () => {
    const tables: string[] = await props.datasource.postResource('tables', {
      schema: props.query.schema || '',
    });
    return tables.map((table) => ({ label: table, value: table })).concat({ label: '-- remove --', value: '' });
  };

  const fetchColumns = async () => {
    const columns: string[] = await props.datasource.postResource('columns', {
      schema: props.query.schema,
      table: props.query.table,
    });
    return columns.map((column) => ({ label: column, value: column })).concat({ label: '-- remove --', value: '' });
  };

  const onChange = (prop: QueryProperties) => (e: SelectableValue<string> | null) => {
    const newQuery = { ...props.query };
    const value = e?.value;
    newQuery[prop] = value;
    props.onChange(newQuery);
  };

  return (
    <>
      {newFormStylingEnabled ? (
        <EditorRows>
          <EditorRow>
            <EditorFieldGroup>
              <EditorField
                width={20}
                label={selectors.components.ConfigEditor.schema.input}
                tooltip="Use the selected schema with the $__schema macro"
                data-testid={selectors.components.ConfigEditor.schema.testID}
                htmlFor="schema"
              >
                <ResourceSelector
                  newFormStylingEnabled={true}
                  id="schema"
                  label={selectors.components.ConfigEditor.schema.input}
                  onChange={onChange('schema')}
                  fetch={fetchSchemas}
                  value={props.query.schema || null}
                />
              </EditorField>
              <EditorField
                width={20}
                label={selectors.components.ConfigEditor.table.input}
                tooltip="Use the selected table with the $__table macro"
                data-testid={selectors.components.ConfigEditor.table.testID}
                htmlFor="table"
              >
                <ResourceSelector
                  newFormStylingEnabled={true}
                  id="table"
                  label={selectors.components.ConfigEditor.table.input}
                  onChange={onChange('table')}
                  fetch={fetchTables}
                  value={props.query.table || null}
                  dependencies={[props.query.schema]}
                />
              </EditorField>
              <EditorField
                width={20}
                label={selectors.components.ConfigEditor.column.input}
                tooltip="Use the selected column with the $__column macro"
                data-testid={selectors.components.ConfigEditor.column.testID}
                htmlFor="column"
              >
                <ResourceSelector
                  newFormStylingEnabled={true}
                  id="column"
                  label={selectors.components.ConfigEditor.column.input}
                  onChange={onChange('column')}
                  fetch={fetchColumns}
                  value={props.query.column || null}
                  dependencies={[props.query.table]}
                />
              </EditorField>
            </EditorFieldGroup>
          </EditorRow>
          <EditorRow>
            <div className={styles.collapseRow}>
              <CollapsableSection
                className={styles.collapse}
                label={
                  <p className={styles.collapseTitle} data-testid="collapse-title">
                    Format
                  </p>
                }
                isOpen={false}
              >
                <EditorFieldGroup>
                  <EditorField label="Format data frames as" htmlFor="formatAs" width={20}>
                    <FormatSelect
                      newFormStylingEnabled={true}
                      id="formatAs"
                      query={props.query}
                      options={SelectableFormatOptions}
                      onChange={props.onChange}
                    />
                  </EditorField>

                  {props.query.format === FormatOptions.TimeSeries && (
                    <FillValueSelect newFormStylingEnabled={true} query={props.query} onChange={props.onChange} />
                  )}
                </EditorFieldGroup>
              </CollapsableSection>
            </div>
          </EditorRow>
          <EditorRow>
            <div style={{ width: '100%' }}>
              <SQLEditor query={props.query} onChange={props.onChange} datasource={props.datasource} />
            </div>
          </EditorRow>
        </EditorRows>
      ) : (
        <InlineSegmentGroup>
          <div className="gf-form-group">
            <h6>Macros</h6>
            <ResourceSelector
              id="schema"
              onChange={onChange('schema')}
              fetch={fetchSchemas}
              value={props.query.schema || null}
              tooltip="Use the selected schema with the $__schema macro"
              label={selectors.components.ConfigEditor.schema.input}
              data-testid={selectors.components.ConfigEditor.schema.testID}
              labelWidth={11}
              className="width-12"
            />
            <ResourceSelector
              id="table"
              onChange={onChange('table')}
              fetch={fetchTables}
              value={props.query.table || null}
              dependencies={[props.query.schema]}
              tooltip="Use the selected table with the $__table macro"
              label={selectors.components.ConfigEditor.table.input}
              data-testid={selectors.components.ConfigEditor.table.testID}
              labelWidth={11}
              className="width-12"
            />
            <ResourceSelector
              id="column"
              onChange={onChange('column')}
              fetch={fetchColumns}
              value={props.query.column || null}
              dependencies={[props.query.table]}
              tooltip="Use the selected column with the $__column macro"
              label={selectors.components.ConfigEditor.column.input}
              data-testid={selectors.components.ConfigEditor.column.testID}
              labelWidth={11}
              className="width-12"
            />
            <h6>Frames</h6>
            <FormatSelect query={props.query} options={SelectableFormatOptions} onChange={props.onChange} />
            {props.query.format === FormatOptions.TimeSeries && (
              <FillValueSelect query={props.query} onChange={props.onChange} />
            )}
          </div>
          <div style={{ minWidth: '400px', marginLeft: '10px', flex: 1 }}>
            <SQLEditor query={props.query} onChange={props.onChange} datasource={props.datasource} />
          </div>
        </InlineSegmentGroup>
      )}
    </>
  );
}

const getStyles = {
  collapse: css({
    alignItems: 'flex-start',
    paddingTop: 0,
  }),
  collapseTitle: css({
    fontSize: 14,
    fontWeight: 500,
    marginBottom: 0,
  }),
  collapseRow: css({
    display: 'flex',
    flexDirection: 'column',
    '>div': {
      alignItems: 'baseline',
      justifyContent: 'flex-end',
    },
    '*[id^="collapse-content-"]': {
      padding: 'unset',
    },
  }),
};
