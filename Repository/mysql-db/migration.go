package mysql_db

func CreateTables(mysqlRepo *MysqlRepo) {
	creatStmt, stmtErr := mysqlRepo.db.Prepare(
		`CREATE TABLE IF NOT EXISTS orders (
    				id BIGINT NOT NULL AUTO_INCREMENT PRIMARY KEY,
    				order_id BIGINT NOT NULL,
    				price	int32	NOT NULL,
    				title	char	NOT NULL )`,
	)
	if stmtErr != nil {
		panic(stmtErr)
	}

	_, err := creatStmt.Exec()
	if err != nil {
		panic(err)
	}
}
