package main

import (
	"fmt"
	"github.com/casbin/casbin/v2"
	"github.com/casbin/casbin/v2/model"
	gormadapter "github.com/casbin/gorm-adapter/v3"
	"github.com/noovertime7/dailytest/dao_demo/database"
	"log"
)

var config = `[request_definition]
r = sub, obj, act

[policy_definition]
p = sub, obj, act

[role_definition]
g = _, _

[policy_effect]
e = some(where (p.eft == allow))

[matchers]
m = g(r.sub, p.sub) && r.obj == p.obj && r.act == p.act || r.sub == "admin"`

var (
	sub = "admin" // 想要访问资源的用户。
	obj = "book"  // 将被访问的资源。
	act = "read"
)

func main() {
	db, err := gormadapter.NewAdapterByDB(database.Gorm)
	if err != nil {
		log.Fatalf("error: adapter: %s", err)
	}
	m, err := model.NewModelFromString(config)
	if err != nil {
		log.Fatalf("error: model: %s", err)
	}
	e, err := casbin.NewEnforcer(m, db)
	if err != nil {
		log.Fatalf("error: enforcer: %s", err)
	}
	//data1 := [][]string{
	//	{"bob", "管理员", "book", "edit"},
	//	{"alice", "普通用户", "book", "read"},
	//}
	//addProxy(e, data1)
	ok, err := e.AddGroupingPolicies([][]string{
		{
			"test", "管理员",
		},
		{
			"bob", "普通用户",
		},
	})
	if err != nil {
		log.Fatalln("AddGroupingPolicies err ", err)
	}
	addProxy(e)
	fmt.Println(ok)
	isOK, err := e.Enforce(sub, obj, act)
	if err != nil {
		log.Fatalln("Enforce err ", err)
	}

	if isOK {
		fmt.Println("通过")
	} else {
		fmt.Println("不通过")
	}

	fmt.Println(e.GetFilteredPolicy(0, "alice "))
	fmt.Println(e.GetGroupingPolicy())

}

func addProxy(e *casbin.Enforcer) {
	policy, err := e.AddPolicy("管理员", "book", "read")
	if err != nil {
		log.Fatalln(err)
	}
	fmt.Println("给管理员添加权限成功", policy)
}
