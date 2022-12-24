package repository

type Item interface {
	Add(id string, url string) error
	Get(id string) (string, error)
}

type Repository struct {
	Item
}

func NewRepository() *Repository {
	return &Repository{
		Item: NewItemLocalStorageRepository(),
	}
}
