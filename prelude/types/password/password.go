package password

import (
	gpv "github.com/wagslane/go-password-validator"
	"golang.org/x/crypto/bcrypt"
)

const minEntropy float64 = 40.0 // during production, this value needs to be > 40

type Password string

func Parse(s string) (p Password, err error) {
	p = Password(s)
	return p, p.Validate()
}

func MustParse(s string) Password {
	e, err := Parse(s)
	if err != nil {
		panic(err)
	}
	return e
}

func (p Password) String() string { return string(p) }

func (p Password) Validate() error {
	return gpv.Validate(p.String(), minEntropy)
}

func (p Password) IsValid() bool { return p.Validate() == nil }

func (p *Password) UnmarshalJSON(b []byte) error {
	*p = Password(b[1 : len(b)-1])
	return p.Validate()
}

func (p Password) MarshalJSON() (b []byte, err error) {
	return []byte(`"` + p.String() + `"`), nil
}

func (p Password) Hash() (PasswordHash, error) {
	if err := p.Validate(); err != nil {
		return nil, err
	}

	return bcrypt.GenerateFromPassword([]byte(p), bcrypt.DefaultCost)
}

func (p Password) MustHash() PasswordHash {
	h, err := p.Hash()
	if err != nil {
		panic(err)
	}
	return h
}

type PasswordHash []byte

func (h PasswordHash) String() string { return string(h) }

func (h PasswordHash) Compare(cmp string) error {
	return bcrypt.CompareHashAndPassword(h, []byte(cmp))
}
