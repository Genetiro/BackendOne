package transport

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"path/filepath"
	"text/template"
)

type Result struct {
	Link   string
	Code   string
	Status string
}

func home(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/" {
		http.NotFound(w, r)
		return
	}

	path := filepath.Join("internal", "html", "home.html")
	ts, err := template.ParseFiles(path)
	if err != nil {
		log.Println(err.Error())
		http.Error(w, "Internal Server Error", 500)
		return
	}
	err = ts.Execute(w, nil)
	if err != nil {
		log.Println(err.Error())
		http.Error(w, "Internal Server Error", 500)
	}

	result := Result{}
	if r.Method == "POST" {
		if !isValidUrl(r.FormValue("s")) {
			fmt.Println("Что-то не так")
			result.Status = "Bad format"
			result.Link = ""
		} else {
			result.Link = r.FormValue("s")
			result.Code = shorting()
			db, err := sql.Open("sqlite3", "links.db")
			if err != nil {
				panic(err)
			}
			defer db.Close()
			db.Exec("insert into linkshort (link, short) values ($1, $2)", result.Link, result.Code)
			result.Status = "Successfully shorting link"
		}
	}
	ts.Execute(w, result)
}
