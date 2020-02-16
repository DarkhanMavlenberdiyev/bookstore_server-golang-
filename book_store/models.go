package book_store

type BookStore interface {
	SaveBooks(filename string) error

	Create(book *Book) (*Book, error)

	GetBook(id int) (*Book, error)

	ListBooks() ([]*Book, error)

	UpdateBook(id int, book *Book) (*Book, error)

	DeleteBook(id int) error
}

type Book struct {
	ID            int    `json:"id"`
	Title         string `json:"title,omitempty"`
	Author        string `json:"author,omitempty"`
	Description   string `json:"description,omitempty"`
	NumberOfPages int    `json:"number_of_pages"`
}
