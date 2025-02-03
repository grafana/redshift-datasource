import {
  ColumnDefinition,
  getStandardSQLCompletionProvider,
  LanguageCompletionProvider,
  TableDefinition,
  TableIdentifier,
  SchemaDefinition,
} from '@grafana/plugin-ui';
import { MACROS } from './macros';

interface CompletionProviderGetterArgs {
  getSchemas: () => Promise<TableDefinition[]>;
  getTables: (d?: string) => Promise<SchemaDefinition[]>;
  getColumns: (table: string, schema?: string) => Promise<ColumnDefinition[]>;
}

export const getRedshiftCompletionProvider: (args: CompletionProviderGetterArgs) => LanguageCompletionProvider =
  ({ getSchemas, getTables, getColumns }) =>
  (monaco, language) => {
    return {
      // get standard SQL completion provider which will resolve functions and macros
      ...(language && getStandardSQLCompletionProvider(monaco, language)),
      triggerCharacters: ['.', ' ', '$', ',', '(', "'"],
      schemas: {
        resolve: getSchemas,
      },
      tables: {
        resolve: (t?: TableIdentifier | null) => {
          return getTables(t?.schema);
        },
      },
      columns: {
        resolve: (t?: TableIdentifier | null) => getColumns(t?.table ?? '', t?.schema ?? ''),
      },
      supportedMacros: () => MACROS,
    };
  };
