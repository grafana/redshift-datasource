import {
  ColumnDefinition,
  getStandardSQLCompletionProvider,
  LanguageCompletionProvider,
  TableDefinition,
  TableIdentifier,
  SchemaDefinition,
  // getStandardSQLCompletionProvider,
} from '@grafana/experimental';
import { MACROS } from './macros';

interface CompletionProviderGetterArgs {
  getSchemas: React.MutableRefObject<() => Promise<TableDefinition[]>>;
  getTables: React.MutableRefObject<(d?: string) => Promise<SchemaDefinition[]>>;
  getColumns: React.MutableRefObject<(table: string, schema?: string) => Promise<ColumnDefinition[]>>;
}

export const getRedshiftCompletionProvider: (args: CompletionProviderGetterArgs) => LanguageCompletionProvider =
  ({ getSchemas, getTables, getColumns }) =>
  (monaco, language) => {
    return {
      // get standard SQL completion provider which will resolve functions and macros
      ...(language && getStandardSQLCompletionProvider(monaco, language)),
      triggerCharacters: ['.', ' ', '$', ',', '(', "'"],
      schemas: {
        resolve: async () => getSchemas.current(),
      },
      tables: {
        resolve: async (t: TableIdentifier) => {
          return await getTables.current(t?.schema);
        },
      },
      columns: {
        resolve: async (t: TableIdentifier) => getColumns.current(t.table!, t.schema),
      },
      supportedMacros: () => MACROS,
    };
  };
