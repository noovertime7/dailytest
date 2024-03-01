package main

import (
	"fmt"
	"github.com/sirupsen/logrus"
	"mysql_executor/executor"
	"mysql_executor/utils"
)

var sql = `
update  t_workflow set  description = '666' where id = 12;
update  t_workflow set  description = '666' where id = 14;
update  t_workflow set  description = '666' where id = 13;
`

func main() {
	logger := logrus.NewEntry(logrus.New())
	e, err := executor.NewExecutor(logger,
		&executor.DSN{
			Host:         "127.0.0.1",
			Port:         "3306",
			User:         "root",
			Password:     "123456",
			DatabaseName: "kubemanage",
		},
		"test",
	)
	if err != nil {
		panic(err)
	}
	fmt.Println(e.ShowSchemaTables("kubemanage"))

	nodes, err := utils.ParseSql(sql)
	if err != nil {
		panic(err)
	}
	qs := make([]string, 0, len(nodes))
	for _, executeSQL := range nodes {
		qs = append(qs, executeSQL.Text())
	}

	res, err := e.Db.Transact(qs...)
	if err != nil {
		panic(err)
	}
	for sql, r := range res {
		LastInsertId, _ := r.LastInsertId()
		RowsAffected, _ := r.RowsAffected()
		fmt.Printf("执行sql:%s, LastInsertId:%d, RowsAffected:%d\n", sql, LastInsertId, RowsAffected)
	}
}
