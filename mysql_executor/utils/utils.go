package utils

import (
	"fmt"
	"github.com/pingcap/parser"
	"github.com/pingcap/parser/ast"
	_ "github.com/pingcap/tidb/types/parser_driver"
)

func ParseSql(sql string) ([]ast.StmtNode, error) {
	p := parser.New()
	stmts, _, err := p.Parse(sql, "", "")
	if err != nil {
		return nil, err
	}
	return stmts, nil
}

func ParseOneSql(sql string) (ast.StmtNode, error) {
	p := parser.New()
	stmt, err := p.ParseOneStmt(sql, "", "")
	if err != nil {
		fmt.Printf("parse error: %v\nsql: %v", err, sql)
		return nil, err
	}
	return stmt, nil
}
