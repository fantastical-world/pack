package pack

type Error string

func (pe Error) Error() string { return string(pe) }

const ErrThingDoesNotExist = Error("thing does not exist")

//Packable is something that can be stored in a Pack.
type Packable interface {
	Pack() (string, []byte)
	Unpack([]byte)
}

//Pack is a place to store things.
type Pack interface {
	Save(location string, thing Packable) error
	Get(location, name string) ([]byte, error)
	Delete(location, name string) error
	List(location string) ([]string, error)
	ListMeta(location string) ([]interface{}, error)
}
