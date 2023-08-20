package types

import (
	"testing"

	"github.com/hyphengolang/prelude/testing/is"
	"github.com/hyphengolang/prelude/types/email"
	"github.com/hyphengolang/prelude/types/password"
)

func TestEmail(t *testing.T) {
	t.Parallel()
	is := is.New(t)

	t.Run(`parse email from string`, func(t *testing.T) {
		_, err := email.Parse("foo@mail.com")
		is.NoErr(err) // parsing string to email

		_, err = email.Parse("foo.com")
		is.True(err != nil) // parsing string to email
	})
}

func TestPassword(t *testing.T) {
	t.Parallel()
	is := is.New(t)

	t.Run(`pares password from string`, func(t *testing.T) {
		_, err := password.Parse("p$4ssw_034")
		is.NoErr(err) // parsing string to email

		_, err = password.Parse("password")
		is.True(err != nil) // parsing string to email
	})
}
