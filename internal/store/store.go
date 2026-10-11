package store

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
}

type Store interface {
	CreateChop(name, color string, ordered bool) error
	AddItem(chop string, item ListItem) error
	Get(name string) (Chop, error)
	List() ([]Chop, error)
	DeleteChop(name string) error
	Close() error
}
