package main

import (
	"./book_store"
	"fmt"
	"github.com/segmentio/encoding/json"
	"io/ioutil"
	"os/signal"
	"strconv"
	"syscall"
	//"github.com/go-pg/pg"
	//"github.com/go-pg/pg/orm"
	//"fmt"
	"github.com/gorilla/mux"
	"github.com/urfave/cli"
	//"log"
	"net/http"
	"os"
)

var (
	jsonFile = ""

	flags = []cli.Flag{
		&cli.StringFlag{
			Name:        "config",
			Aliases:     []string{"c"},
			Destination: &jsonFile,
		},
	}
)

func main() {

	app := cli.NewApp()
	app.Flags = flags
	app.Commands = cli.Commands{
		&cli.Command{
			Name:   "start",
			Usage:  "start the local server",
			Action: StartServer,
		},
	}
	app.Run(os.Args)

}
func PrintHello() {
	fmt.Println("Hello World")
}

func StartServer(d *cli.Context) error {

	bookStore, err := book_store.NewBookStore(jsonFile)
	if err != nil {
		panic(err)
	}

	endpoints := book_store.NewEndpointsFactory(bookStore)

	user := book_store.PostgreConfig{
		User:     "postgres",
		Password: "qwerty123",
		Port:     "8080",
		Host:     "0.0.0.0",
	}

	router := mux.NewRouter()
	router.Methods("GET").Path("/get/all").HandlerFunc(endpoints.ListBooks())
	router.Methods("GET").Path("/{id}").HandlerFunc(endpoints.GetBook("id"))
	router.Methods("PUT").Path("/{id}").HandlerFunc(endpoints.UpdateBook("id"))
	router.Methods("POST").Path("/").HandlerFunc(endpoints.CreateBook())
	router.Methods("DELETE").Path("/{id}").HandlerFunc(endpoints.DeleteBook("id"))

	db, err := book_store.NewPostgreBookStore(user)
	if err != nil {
		panic(err)
	}

	router.Methods("POST").Path("/db/").HandlerFunc(CreateBook(db))
	router.Methods("GET").Path("/db/get/all").HandlerFunc(ListBooks(db))
	router.Methods("GET").Path("/db/{id}").HandlerFunc(GetBook("id", db))
	router.Methods("DELETE").Path("/db/{id}").HandlerFunc(DeleteBook("id", db))
	router.Methods("PUT").Path("/db/{id}").HandlerFunc(UpdateBook("id", db))

	go func(rtr *mux.Router) {
		http.ListenAndServe("0.0.0.0:8000", rtr)
	}(router)

	c := make(chan os.Signal, 1)
	done := make(chan bool, 1)
	signal.Notify(c, os.Interrupt, syscall.SIGTERM)
	go func() {
		<-c
		done <- true
	}()
	<-done
	ExitWithSave(bookStore)
	os.Exit(1)
	return nil
}

func ListBooks(db book_store.BookStore) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		res, err := db.ListBooks()
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte("ErrorW: " + err.Error()))
			return
		}

		ress := make([][]byte, 0)
		for _, book := range res {
			resp, err := json.Marshal(book)
			if err != nil {
				w.Write([]byte(err.Error()))
				return
			}
			ress = append(ress, resp)
		}

		for _, r := range ress {
			w.Write(r)
			w.Write([]byte("\n"))
		}
		w.WriteHeader(http.StatusOK)

	}
}

func GetBook(idParam string, db book_store.BookStore) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		id, ok := vars[idParam]
		if !ok {
			w.WriteHeader(http.StatusBadRequest)
			w.Write([]byte("Book ID not found "))
			return
		}
		idd, _ := strconv.Atoi(id)
		res, err := db.GetBook(idd)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte("Sorry:( : " + err.Error()))
			return
		}
		data, err := json.Marshal(res)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte("Error: " + err.Error()))
			return
		}
		w.Write(data)
		w.WriteHeader(http.StatusOK)
	}
}

func UpdateBook(idParam string, db book_store.BookStore) func(w http.ResponseWriter, r *http.Request) {
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

		book := &book_store.Book{}
		if err := json.Unmarshal(data, book); err != nil {
			w.WriteHeader(http.StatusBadRequest)
			w.Write([]byte("Error: " + err.Error()))
			return
		}
		idd, _ := strconv.Atoi(id)
		_, err = db.UpdateBook(idd, book)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		w.Write([]byte("Account updated successfully"))
		w.WriteHeader(http.StatusOK)
	}
}

func DeleteBook(idParam string, db book_store.BookStore) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		id, ok := vars[idParam]
		if !ok {
			w.WriteHeader(http.StatusBadRequest)
			w.Write([]byte("Book ID not found "))
			return
		}
		idd, _ := strconv.Atoi(id)
		err := db.DeleteBook(idd)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte("Error: " + err.Error()))
			return
		}
		w.Write([]byte("Account deleted successfully"))
		w.WriteHeader(http.StatusOK)
	}
}

func CreateBook(db book_store.BookStore) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		data, err := ioutil.ReadAll(r.Body)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte("Error: " + err.Error()))
			return
		}
		book := &book_store.Book{}
		if err := json.Unmarshal(data, book); err != nil {
			w.WriteHeader(http.StatusBadRequest)
			w.Write([]byte("Error: " + err.Error()))
			return
		}
		_, err = db.Create(book)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
		}
		w.WriteHeader(http.StatusCreated)
		w.Write([]byte("Account created successfully"))

	}

}

func ExitWithSave(book book_store.BookStore) {
	err := book.SaveBooks(jsonFile)
	if err != nil {
		fmt.Println(err)
		return
	}

}
