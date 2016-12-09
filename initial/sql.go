package initial

import (
	"fmt"

	"github.com/astaxie/beego"
	"github.com/astaxie/beego/orm"
	_ "github.com/go-sql-driver/mysql"
)

func InitSql() {
	user := beego.AppConfig.String("mysqluser")
	password := beego.AppConfig.String("mysqlpassword")
	host := beego.AppConfig.String("mysqlhost")
	port, err := beego.AppConfig.Int("mysqlport")
	dbname := beego.AppConfig.String("mysqldb")
	if err != nil {
		port = 3306
	}
	if beego.AppConfig.String("runmode") == "dev" {
		orm.Debug = true
	}
	orm.RegisterDriver("mysql", orm.DRMySQL)
	orm.RegisterDataBase("default", "mysql", fmt.Sprintf("%s:%s@%tcp(%s:%d)/%s?charset=utf8", user, password, host, port, dbname))
}
