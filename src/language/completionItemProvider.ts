import {
  ColumnDefinition,
  getStandardSQLCompletionProvider,
  LanguageCompletionProvider,
  TableDefinition,
  TableIdentifier,
  SchemaDefinition,
} from '@grafana/experimental';
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
        resolve: () => getSchemas(),
      },
      tables: {
        resolve: (t: TableIdentifier) => {
          return getTables(t?.schema);
        },
      },
      columns: {
        resolve: (t: TableIdentifier) => getColumns(t.table!, t.schema),
      },
      supportedMacros: () => MACROS,
    };
  };
