package users

import (
	"fmt"

	"github.com/astaxie/beego"
	"github.com/astaxie/beego/orm"
	"github.com/ninefive/hims/models"
	"github.com/ninefive/hims/utils"
)

type Departs struct {
	Id     int64 `orm:"pk;column(departid);"`
	Name   string
	Desc   string
	Status int
}

func (this *Departs) TableName() string {
	return models.TableName("departs")
}

func init() {
	orm.RegisterModel(new(Departs))
}

func GetDeparts(id int64) (Departs, error) {
	var depart Departs
	var err error
	o := orm.NewOrm()

	depart = Departs{Id: id}
	err := o.Read(&depart)

	if err == orm.ErrNoRows {
		return depart, nil
	}
	return depart, err
}

func GetDepartsName(id int64) string {
	var err error
	var name string
	err = utils.GetCache("GetDepartsName.id."+fmt.Sprintf("%d", id), &name)
	if err != nil {
		cache_expire, _ := beego.AppConfig.Int("cache_expire")
		var depart Departs
		o := orm.NewOrm()
		o.QueryTable(models.TableName("departs")).Filter("departid", id).One(&depart, "name")
		name = depart.Name
		utils.SetCache("GetDepartsName.id."+fmt.Sprintf("%d", id), name, cache_expire)
	}
	return name
}
