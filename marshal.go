package isbn

func (isbn ISBN) MarshalJSON() ([]byte, error) {
	return isbn[:], nil
}

func (isbn *ISBN) UnmarshalJSON(b []byte) (err error) {
	*isbn, err = ParseBytes(b[1 : len(b)-1])
	return err
}
