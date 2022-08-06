/*
	we want to write an app to manage order process
	design a restful api that expose 1 endpoint
	POST api/order input : {“order_id”:10, ”price”:1000, ”title”:”burger”}
	in the receiving point send the received data to redis-db (queue)

->  create another app to read data from the above redis-db and process the order
	the process should goes like this :
	use mysql to save the retrieved data into orders table
	Acceptance criteria :
	1 – use multi layer config (file,environment,...)
	2 – use docker compose for external services (mysql,redis-db)
	2 – include build.sh file to build the apps
	3 – include run.sh to run the apps
	4 – writing unit tests is a plus
*/

package app

import (
	"github.com/greenbahar/manage-order-storage/Repository/mysql-db"
	"github.com/greenbahar/manage-order-storage/Repository/redis-db"
	"github.com/greenbahar/manage-order-storage/order-transfer-storage"
	"os"
)

func StartApplication() {
	redisRepo := redis_db.NewRepo()
	mysqlRepo := mysql_db.NewRepo()
	service := order_transfer_storage.NewService(redisRepo, mysqlRepo)

	mysql_db.CreateTables(mysqlRepo)

	if err := service.Run(); err != nil {
		os.Exit(1)
	}
}
