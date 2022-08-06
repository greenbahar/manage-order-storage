package mysql_db

import (
	"database/sql"
	"fmt"
	"github.com/greenbahar/manage-order-storage/order"

	//"github.com/go-sql-driver/mysql"
	"github.com/greenbahar/manage-order-storage/config/env"
)

const (
	insert_order_sql_statement = "INSERT INTO orders VALUES( ?, ?, ?)"
)

type MysqlRepo struct {
	db *sql.DB
}

func NewRepo() *MysqlRepo {
	sqlConf := env.GetMysqlConfig()
	db, err := sql.Open("mysql", sqlConf.MysqlUrl)
	if err != nil {
		panic(err)
	} else if err = db.Ping(); err != nil {
		panic(err)
	}
	defer db.Close()

	return &MysqlRepo{
		db: db,
	}
}

func (r *MysqlRepo) InsertOrder(orders []*order.Order) error {

	stmtIns, err := r.db.Prepare(insert_order_sql_statement)
	if err != nil {
		panic(err.Error())
	}
	defer stmtIns.Close() // Close the statement when we leave main() or the program terminates

	for _, ord := range orders {
		_, insertIrr := stmtIns.Exec(insert_order_sql_statement, ord.OrderId, ord.Title, ord.Price)
		if insertIrr != nil {
			fmt.Println("cannot insert the order into the mysql database. err: ", insertIrr)
			return insertIrr
		}
	}

	return nil
}
