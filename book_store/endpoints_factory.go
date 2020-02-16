package book_store

import (
	"encoding/json"
	"github.com/gorilla/mux"
	"io/ioutil"
	"net/http"
	"strconv"
)

type Endpoints interface {
	GetBook(idParam string) func(w http.ResponseWriter, r *http.Request)
	CreateBook() func(w http.ResponseWriter, r *http.Request)
	ListBooks() func(w http.ResponseWriter, r *http.Request)
	UpdateBook(idParam string) func(w http.ResponseWriter, r *http.Request)
	DeleteBook(idParam string) func(w http.ResponseWriter, r *http.Request)
}

func NewEndpointsFactory(bookstore BookStore) *endpointsFactory {
	return &endpointsFactory{bookStore: bookstore}
}

type endpointsFactory struct {
	bookStore BookStore
}

func (ef *endpointsFactory) GetBook(idParam string) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		id, ok := vars[idParam]
		if !ok {
			w.WriteHeader(http.StatusBadRequest)
			w.Write([]byte("Book ID not found "))
			return
		}
		idd, _ := strconv.Atoi(id)
		book, err := ef.bookStore.GetBook(idd)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte("I'm sorry: " + err.Error()))
			return
		}
		data, err := json.Marshal(book)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte("Error: " + err.Error()))
			return
		}
		w.Write(data)
		w.WriteHeader(http.StatusOK)
	}
}

func (ef *endpointsFactory) CreateBook() func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		data, err := ioutil.ReadAll(r.Body)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte("Error: " + err.Error()))
			return
		}
		book := &Book{}
		if err := json.Unmarshal(data, book); err != nil {
			w.WriteHeader(http.StatusBadRequest)
			w.Write([]byte("Error: " + err.Error()))
			return
		}
		result, err := ef.bookStore.Create(book)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte("Error: " + err.Error()))
			return
		}
		response, err := json.Marshal(result)
		if err != nil {
			w.Write([]byte("Error: " + err.Error()))
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		w.Write(response)
		w.WriteHeader(http.StatusCreated)
	}
}

func (ef *endpointsFactory) ListBooks() func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		books, err := ef.bookStore.ListBooks()
		if err != nil {
			w.Write([]byte(err.Error()))
			return
		}
		response := make([][]byte, 0)
		for _, book := range books {
			resp, err := json.Marshal(book)
			if err != nil {
				w.Write([]byte(err.Error()))
				return
			}
			response = append(response, resp)
		}

		for _, r := range response {
			w.Write(r)
			w.Write([]byte("\n"))
		}
	}
}

func (ef *endpointsFactory) UpdateBook(idParam string) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		id, ok := vars[idParam]
		if !ok {
			w.WriteHeader(http.StatusBadRequest)
			w.Write([]byte("Book ID not found "))
			return
		}
		data, err := ioutil.ReadAll(r.Body)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte("Error: " + err.Error()))
			return
		}

		book := &Book{}
		if err := json.Unmarshal(data, book); err != nil {
			w.WriteHeader(http.StatusBadRequest)
			w.Write([]byte("Error: " + err.Error()))
			return
		}
		idd, _ := strconv.Atoi(id)
		res, err := ef.bookStore.UpdateBook(idd, book)
		if err != nil {
			w.Write([]byte("Error: " + err.Error()))
			w.WriteHeader(http.StatusBadRequest)
			return
		}
		response, err := json.Marshal(res)
		if err != nil {
			w.Write([]byte("Error: " + err.Error()))
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		w.Write(response)
		w.WriteHeader(http.StatusCreated)
	}
}

func (ef *endpointsFactory) DeleteBook(idParam string) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		id, ok := vars[idParam]
		if !ok {
			w.WriteHeader(http.StatusBadRequest)
			w.Write([]byte("Error: Not Found"))
			return
		}
		idd, _ := strconv.Atoi(id)
		err := ef.bookStore.DeleteBook(idd)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte("Error" + err.Error()))
			return
		}
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("Book is deleted successfully"))
	}
}
