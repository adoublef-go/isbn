package email

import "net/mail"

type Email string

func Parse(s string) (e Email, err error) {
	e = Email(s)
	return e, e.Validate()
}

func MustParse(s string) Email {
	e, err := Parse(s)
	if err != nil {
		panic(err)
	}
	return e
}

func (e Email) String() string { return string(e) }

func (e Email) Validate() error {
	_, err := mail.ParseAddress(string(e))
	return err
}

func (e Email) IsValid() bool { return e.Validate() == nil }

func (e *Email) UnmarshalJSON(b []byte) error {
	*e = Email(b[1 : len(b)-1])
	return e.Validate()
}
