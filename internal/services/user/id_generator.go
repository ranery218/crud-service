package user

type IDGen interface {
	NewID() (string, error)
}
