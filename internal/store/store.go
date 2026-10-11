package store

import "errors"

var (
	ErrNotFound      = errors.New("chop not found")
	ErrAlreadyExists = errors.New("chop already exists")
)

type SortField string

const (
	SortByName    SortField = "name"
	SortByColor   SortField = "color"
	SortByOrdered SortField = "ordered"
)

type Chop struct {
	ID        int64
	Name      string
	MainColor string
	Ordered   bool
	Items     []ListItem
}

type ListItem struct {
	ID               int64
	Text             string
	TextColor        string
	Description      string
	DescriptionColor string
	Position         int
}

type Store interface {
	CreateChop(name, color string, ordered bool) error
	AddItem(chop string, item ListItem) error
	Get(name string) (Chop, error)
	List(sort SortField) ([]Chop, error)
	DeleteChop(name string) error
	RmItem(chop string, position int) error
	Close() error
}
