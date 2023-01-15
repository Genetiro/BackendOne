package database

import (
	"database/sql"
	"errors"
	"fmt"
	"log"

	_ "github.com/mattn/go-sqlite3"
)

type Db struct {
	db *sql.DB
}
type ListDb struct {
	Id    int64
	Link  string
	Short string
}

var ErrDuplicate = errors.New("record already exists")

func (dbase *Db) CloseDb() error {
	return dbase.db.Close()

}

func NewDB(dbFile string) (*Db, error) {
	sqlDB, err := sql.Open("sqlite3", dbFile)
	if err != nil {
		return nil, err
	}
	db := Db{
		db: sqlDB,
	}
	return &db, nil
}

func (dbase *Db) GetShortLinks() ([]ListDb, error) {
	rows, err := dbase.db.Query("SELECT * FROM linkshort")
	if err != nil {
		log.Println(err)
	}
	defer rows.Close()
	ListLinks := []ListDb{}
	for rows.Next() {
		l := ListDb{}
		err := rows.Scan(&l.Id, &l.Link, &l.Short)
		if err != nil {
			fmt.Println(err)
			continue
		}
		ListLinks = append(ListLinks, l)
	}
	return ListLinks, nil
}

func (dbase *Db) DeleteShort(shrt string) error {
	_, err := dbase.db.Exec("DELETe FROM linkshort WHERE short = $1", shrt)
	if err != nil {
		log.Println(err)
	}
	return nil
}

func (dbase *Db) GetByShort(shrt string) (*ListDb, error) {
	row := dbase.db.QueryRow("SELECT Id, Link, Short FROM linkshort WHERE short = $1", shrt)

	l := ListDb{}
	err := row.Scan(&l.Id, &l.Link, &l.Short)
	if err != nil {
		log.Println("can't found ", err)
	}
	return &l, nil
}

func (dbase *Db) CreateShort(shrt ListDb) (*ListDb, error) {

	res, err := dbase.db.Exec("INSERT INTO linkshort (Link, Short) values ($1, $2)", shrt.Link, shrt.Short)
	if err != nil {
		// var sqliteErr sqlite3.Error
		// if errors.As(err, &sqliteErr) {
		// 	if errors.Is(sqliteErr.ExtendedCode, sqlite3.ErrConstraintUnique) {
		// 		return nil, ErrDuplicate
		// 	}
		// }
		return nil, err
	}
	id, err := res.LastInsertId()
	if err != nil {
		return nil, err
	}
	shrt.Id = id
	return &shrt, nil

}
