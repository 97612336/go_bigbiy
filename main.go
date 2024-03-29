package main

import (
	"go_bigbiy/handlers"
	"go_bigbiy/util"
	"net/http"
	_ "github.com/go-sql-driver/mysql"
	"flag"
	"fmt"
)

func init() {
	util.DB = util.Get_sql_db()
}

func main() {
	//设置路由
	handlers.MyUrls()
	//设置端口号
	var port string
	flag.StringVar(&port,"port","8081","listen port")
	flag.Parse()
	fmt.Println(port)
	//设置监听端口
	err := http.ListenAndServe(":"+port, nil)
	//启动程序
	util.CheckErr(err)

}
