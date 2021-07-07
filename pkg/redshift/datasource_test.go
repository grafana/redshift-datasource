package redshift

import (
	"context"
	"fmt"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/google/go-cmp/cmp"
)

func TestSchemas(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer db.Close()

	ds := RedshiftDatasource{db: db}
	schemaName := "foo"
	mock.ExpectQuery("SELECT nspname FROM pg_namespace").
		WillReturnRows(sqlmock.NewRows([]string{"schema"}).AddRow(schemaName))

	schemas, err := ds.Schemas(context.TODO())
	if err != nil {
		t.Errorf("unexpected error %v", err)
	}
	expectedSchemas := []string{schemaName}
	if !cmp.Equal(schemas, expectedSchemas) {
		t.Errorf("unexpected result: %v", cmp.Diff(schemas, expectedSchemas))
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expectations: %s", err)
	}
}

func TestTables(t *testing.T) {
	tests := []struct {
		description    string
		schemaName     string
		expectedTables []string
		expectedErr    error
	}{
		{
			description:    "should return tables",
			schemaName:     "foobar",
			expectedTables: []string{"foo", "bar"},
			expectedErr:    nil,
		},
		{
			description:    "should use public schema by default",
			schemaName:     "",
			expectedTables: []string{"foo", "bar"},
			expectedErr:    nil,
		},
		{
			description:    "should fail if the schema name is not supported",
			schemaName:     "'*'",
			expectedTables: []string{},
			expectedErr:    fmt.Errorf("unsupported schema name '*'"),
		},
	}
	for _, test := range tests {
		t.Run(test.description, func(t *testing.T) {
			db, mock, err := sqlmock.New()
			if err != nil {
				t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
			}
			defer db.Close()

			ds := RedshiftDatasource{db: db}
			schema := test.schemaName
			if schema == "" {
				schema = "public"
			}
			rows := sqlmock.NewRows([]string{"tables"})
			for _, table := range test.expectedTables {
				rows.AddRow(table)
			}
			mock.ExpectQuery(fmt.Sprintf("SELECT table_name FROM information_schema.tables WHERE table_schema='%s'", schema)).
				WillReturnRows(rows)

			tables, err := ds.Tables(context.TODO(), test.schemaName)
			if err != nil {
				if test.expectedErr == nil || (err.Error() != test.expectedErr.Error()) {
					t.Errorf("unexpected error %v", cmp.Diff(err.Error(), test.expectedErr.Error()))
				}
			} else {
				if !cmp.Equal(tables, test.expectedTables) {
					t.Errorf("unexpected result: %v", cmp.Diff(tables, test.expectedTables))
				}

				if err := mock.ExpectationsWereMet(); err != nil {
					t.Errorf("there were unfulfilled expectations: %s", err)
				}
			}
		})
	}
}

func TestColumns(t *testing.T) {
	tests := []struct {
		description  string
		tableName    string
		expectedCols []string
		expectedErr  error
	}{
		{
			description:  "should return columns",
			tableName:    "foobar",
			expectedCols: []string{"foo", "bar"},
			expectedErr:  nil,
		},
		{
			description:  "should fail if the table name is not supported",
			tableName:    "'*'",
			expectedCols: []string{},
			expectedErr:  fmt.Errorf("unsupported table name '*'"),
		},
	}
	for _, test := range tests {
		t.Run(test.description, func(t *testing.T) {
			db, mock, err := sqlmock.New()
			if err != nil {
				t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
			}
			defer db.Close()

			ds := RedshiftDatasource{db: db}
			rows := sqlmock.NewRows([]string{"table"})
			for _, col := range test.expectedCols {
				rows.AddRow(col)
			}
			mock.ExpectQuery(fmt.Sprintf("SELECT column_name FROM INFORMATION_SCHEMA.COLUMNS WHERE table_name = '%s'", test.tableName)).
				WillReturnRows(rows)

			columns, err := ds.Columns(context.TODO(), test.tableName)
			if err != nil {
				if test.expectedErr == nil || (err.Error() != test.expectedErr.Error()) {
					t.Errorf("unexpected error %v", cmp.Diff(err.Error(), test.expectedErr.Error()))
				}
			} else {
				if !cmp.Equal(columns, test.expectedCols) {
					t.Errorf("unexpected result: %v", cmp.Diff(columns, test.expectedCols))
				}

				if err := mock.ExpectationsWereMet(); err != nil {
					t.Errorf("there were unfulfilled expectations: %s", err)
				}
			}
		})
	}
}
