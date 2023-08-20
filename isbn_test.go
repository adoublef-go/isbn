package isbn

import (
	"bytes"
	"context"
	"database/sql"
	"encoding/json"
	"os"
	"strings"
	"testing"

	"github.com/hyphengolang/prelude/testing/is"
	"github.com/jackc/pgx/v5"
	_ "github.com/mattn/go-sqlite3"
)

func TestIsbnValidation(t *testing.T) {
	t.Parallel()
	is := is.New(t)

	tt := []struct {
		desc string
		data string
		err  string
	}{
		{desc: "invalid: format", data: "978071670344a", err: "invalid ISBN format"},
		{desc: "invalid: length 12 != 13", data: "978071670344", err: "invalid ISBN length 12"},
		{desc: "valid: isbn 13", data: "9780716703440"},
		{desc: "invalid: isbn 13", data: "9780716703410", err: "invalid ISBN value"},
		{desc: "invalid: length 14 != 13", data: "97807167034403", err: "invalid ISBN length 14"},
		{desc: "valid: isbn 10", data: "0716703440"},
		{desc: "valid: isbn 10 w/ dashes", data: "0-7167-0344-0"},
		{desc: "valid: isbn 13 w/ dashes", data: "978-0-7167-0344-0"},
	}

	for _, tc := range tt {
		t.Run(tc.desc, func(t *testing.T) {
			isbn, err := Parse(tc.data)
			if err != nil {
				is.Equal(err.Error(), tc.err) // are errors equal
				return
			}

			is.Equal(len(isbn.String()), 13) // length of ISBN == 13
		})
	}
}

func TestIsbnJSON(t *testing.T) {
	t.Parallel()
	is := is.New(t)

	t.Run("working with json", func(t *testing.T) {
		data := `["9780716703440","0716703440"]`

		var isbns []ISBN
		err := json.NewDecoder(strings.NewReader(data)).Decode(&isbns)
		is.NoErr(err) // decoding isbn values

		var buf bytes.Buffer
		err = json.NewEncoder(&buf).Encode(&isbns[1])
		is.NoErr(err) // encoding isbn values

		// FIXME terrible check as doesn't really tell me anything
		is.Equal(len(buf.String()), len("9780716703440")+1) // buf is the same size as the input string
	})

}

func TestIsbnSql(t *testing.T) {
	t.Parallel()
	is := is.New(t)

	isbn, _ := Parse("9780716703440")
	db, _ := sql.Open("sqlite3", ":memory:")

	type testcase struct {
		ID   int
		Isbn *ISBN // nullable
	}

	tt := []struct {
		desc   string
		schema string
	}{
		{"isbn as BLOB", "CREATE TABLE \"__test__\" (id INTEGER PRIMARY KEY, isbn BLOB)"},
		{"isbn as TEXT", "CREATE TABLE \"__test__\" (id INTEGER PRIMARY KEY, isbn TEXT)"},
	}
	for _, tc := range tt {
		t.Run(tc.desc, func(t *testing.T) {
			_, err := db.Exec(tc.schema)
			is.NoErr(err) // migrate schema

			t.Cleanup(func() { db.Exec("DROP TABLE \"__test__\"") })

			_, err = db.Exec("INSERT INTO \"__test__\" (isbn) VALUES ($1)", isbn)
			is.NoErr(err) // insert value to database

			var tc testcase

			err = db.QueryRow("SELECT id, isbn FROM \"__test__\" WHERE id = 1").Scan(&tc.ID, &tc.Isbn)
			is.NoErr(err)                             // get entry from database
			is.Equal(tc.Isbn.String(), isbn.String()) // check that values are equal
		})
	}
}

func TestIsbnPsql(t *testing.T) {
	t.Parallel()
	is := is.New(t)

	db, err := pgx.Connect(context.Background(), os.ExpandEnv("postgres://${POSTGRES_USER}:${POSTGRES_PASS}@${POSTGRES_HOST}:${POSTGRES_PORT}"))
	is.NoErr(err) // connect to psql

	t.Cleanup(func() { db.Close(context.Background()) })

	isbn, _ := Parse("9780716703440")
	type testcase struct {
		ID    int
		Isbn  *ISBN     // nullable
		Title *[13]byte // nullable
	}

	t.Run("isbn as VARCHAR(13)", func(t *testing.T) {
		_, err := db.Exec(context.Background(), "CREATE TEMP TABLE \"__test__\" (id SERIAL PRIMARY KEY, isbn VARCHAR(13), title VARCHAR(255))")
		is.NoErr(err) // insert value to database

		_, err = db.Exec(context.Background(), "INSERT INTO \"__test__\" (isbn) VALUES ($1)", isbn)
		is.NoErr(err) // insert into table

		var tc testcase
		err = db.QueryRow(context.Background(), "SELECT id, isbn, title FROM \"__test__\" WHERE id = 1").Scan(&tc.ID, &tc.Isbn, &tc.Title)
		is.NoErr(err) // query from row

		is.Equal(tc.Isbn.String(), isbn.String())
	})
}
