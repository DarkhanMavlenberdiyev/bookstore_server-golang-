package book_store

import (
	"github.com/go-pg/pg"
)

type PostgreConfig struct {
	User     string
	Password string
	Port     string
	Host     string
}

type postgreStore struct {
	db *pg.DB
}

func NewPostgreBookStore(config PostgreConfig) (BookStore, error) {
	db := pg.Connect(&pg.Options{
		Addr:     config.Host + ":" + config.Port,
		User:     "postgres",
		Password: config.Password,
	})

	return &postgreStore{db: db}, nil
}

func (ps *postgreStore) Create(book *Book) (*Book, error) {
	err := ps.db.Insert(book)
	return book, err
}

func (ps *postgreStore) SaveBooks(filename string) error {
	return nil
}

func (ps *postgreStore) GetBook(id int) (*Book, error) {
	book := &Book{ID: id}
	err := ps.db.Select(book)
	if err != nil {
		return nil, err
	}
	return book, nil

}

func (ps *postgreStore) ListBooks() ([]*Book, error) {
	var books []*Book
	err := ps.db.Model(&books).Select()

	if err != nil {
		return nil, err
	}
	return books, nil

}

func (ps *postgreStore) UpdateBook(id int, book *Book) (*Book, error) {
	book.ID = id
	err := ps.db.Update(book)
	return book, err
}

func (ps *postgreStore) DeleteBook(id int) error {
	book := &Book{ID: id}
	err := ps.db.Delete(book)
	if err != nil {
		return err
	}
	return nil
}
