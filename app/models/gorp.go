package models

import (
	"database/sql"

	"github.com/go-gorp/gorp"
	_ "github.com/go-sql-driver/mysql"
	db "github.com/revel/modules/db/app"
	r "github.com/revel/revel"
)

var (
	Dbm *gorp.DbMap
	Txn *gorp.Transaction
)

func InitDBM() {
	db.Init()
	Dbm = &gorp.DbMap{
		Db:      db.Db,
		Dialect: gorp.MySQLDialect{"InnoDB", "UTF8"},
	}
	// Dbm.TraceOn("[gorp]", r.RevelLog)
	Dbm.AddTableWithName(Book{}, "books").SetKeys(true, "ID")
}

type GorpController struct {
	*r.Controller
}

func (c *GorpController) Begin() r.Result {
	txn, err := Dbm.Begin()
	if err != nil {
		panic(err)
	}
	Txn = txn
	return nil
}

func (c *GorpController) Commit() r.Result {
	if Txn == nil {
		return nil
	}
	if err := Txn.Commit(); err != nil && err != sql.ErrTxDone {
		panic(err)
	}
	Txn = nil
	return nil
}

func (c *GorpController) Rollback() r.Result {
	if Txn == nil {
		return nil
	}
	if err := Txn.Rollback(); err != nil && err != sql.ErrTxDone {
		panic(err)
	}
	Txn = nil
	return nil
}
