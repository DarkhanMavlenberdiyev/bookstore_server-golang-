package book_store

import (
	"bufio"
	"encoding/json"
	"github.com/pkg/errors"
	"io/ioutil"
	"os"
)

func NewBookStore(filename string) (BookStore, error) {
	file, err := os.Open(filename)
	if err != nil {
		os.Create(filename)

	}
	file.Write([]byte("[]"))
	buffer := bufio.NewReader(file)
	data, err := ioutil.ReadAll(buffer)
	if err != nil {
		return nil, err
	}
	var books []*Book
	if err := json.Unmarshal(data, &books); err != nil {
		return nil, err
	}
	defer file.Close()
	return &bookStoreClass{books}, nil
}

type bookStoreClass struct {
	books []*Book
}

func (bsc *bookStoreClass) GetBook(id int) (*Book, error) {
	for _, v := range bsc.books {
		if v.ID == id {
			return v, nil
		}
	}
	return nil, errors.New("Not found ")
}

func (bsc *bookStoreClass) Create(book *Book) (*Book, error) {
	bsc.books = append(bsc.books, book)
	return book, nil
}

func (bsc *bookStoreClass) SaveBooks(filename string) error {
	file, err := os.Create(filename)
	if err != nil {
		return err
	}
	data, err := json.Marshal(bsc.books)
	if err != nil {
		return err
	}
	_, err = file.Write(data)
	if err != nil {
		return err
	}
	return nil
}
func (bsc *bookStoreClass) ListBooks() ([]*Book, error) {
	return bsc.books, nil
}

func (bsc *bookStoreClass) UpdateBook(id int, book *Book) (*Book, error) {

	for i, v := range bsc.books {
		if v.ID == id {
			bsc.books[i].Author = book.Author
			bsc.books[i].Description = book.Description
			bsc.books[i].NumberOfPages = book.NumberOfPages
			return bsc.books[i], nil
		}
	}
	return nil, errors.New("Can not update")
}

func (bsc *bookStoreClass) DeleteBook(id int) error {
	for i, v := range bsc.books {
		if v.ID == id {
			bsc.books[i] = bsc.books[len(bsc.books)-1]

			bsc.books[len(bsc.books)-1] = nil
			bsc.books = bsc.books[:len(bsc.books)-1]
			return nil
		}
	}
	return errors.New("Not Found")
}
