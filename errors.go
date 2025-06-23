package gsrv

// Error implements constant errors
type Error string

func (e Error) Error() string {
	return string(e)
}
