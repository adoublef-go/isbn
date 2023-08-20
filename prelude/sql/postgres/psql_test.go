package psql

import (
	"context"
	"os"
	"testing"

	"github.com/hyphengolang/prelude/testing/is"
	"github.com/jackc/pgx/v5"
)

var psql = os.ExpandEnv("postgres://${POSTGRES_USER}:${POSTGRES_PASS}@${POSTGRES_HOST}:${POSTGRES_PORT}")

var conn *pgx.Conn
var err error

func init() {
	migrate := `
create temp table person (
	id serial primary key,
	name text not null,
	age int not null check (age > 0)
);

insert into person (name, age) values ('John', 23);
insert into person (name, age) values ('Mary', 25);
`

	conn, err = pgx.Connect(context.Background(), psql)
	if err != nil {
		panic(err)
	}

	if _, err := conn.Exec(context.Background(), migrate); err != nil {
		panic(err)
	}
}

func TestGenericQueries(t *testing.T) {
	t.Parallel()

	type testcase struct {
		ID   int
		Name string
		Age  int
	}

	is := is.New(t)

	t.Cleanup(func() { conn.Close(context.Background()) })

	t.Run(`select entry from database`, func(t *testing.T) {
		var b testcase
		err := QueryRow(conn, "SELECT * FROM person WHERE id = $1", func(row pgx.Row) error { return row.Scan(&b.ID, &b.Name, &b.Age) }, 1)
		is.NoErr(err) // can make query

		is.Equal(b.Name, "John") // user.Name == "John"
	})

	t.Run(`select all from database`, func(t *testing.T) {
		items, err := Query(conn, "select * from person where age > $1", func(rows pgx.Rows, i *testcase) error {
			return rows.Scan(&i.ID, &i.Name, &i.Age)
		}, 24)

		is.NoErr(err) // can make query

		is.Equal(len(items), 1) // two entries in the test database
	})

	t.Run(`insert into database`, func(t *testing.T) {
		args := pgx.NamedArgs{
			"name": "Maxi",
			"age":  31,
		}
		err := Exec(conn, "insert into person (name,age) values (@name,@age)", args)
		is.NoErr(err) // no error making exec
	})

}
