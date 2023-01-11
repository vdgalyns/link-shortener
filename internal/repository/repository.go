package repository

type Link interface {
	Get(hash string) (string, error)
	Add(hash string, url string) error
}

type Repository struct {
	Link
}

func NewRepository() *Repository {
	return &Repository{
		Link: NewLinkLocalRepository(),
	}
}
