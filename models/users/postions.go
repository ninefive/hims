package users

import (
	"fmt"

	"github.com/astaxie/beego"
	"github.com/astaxie/beego/orm"
	"github.com/ninefive/hims/models"
	"github.com/ninefive/hims/utils"
)

type Positions struct {
	Id     int64 `orm:"pk;column(positionid);"`
	Name   string
	Desc   string
	Status int
}

func (this *Positions) TableName() string {
	return models.TableName("positions")
}

func init() {
	orm.RegisterModel(new(Positions))
}

func GetPositions(id int64) (Positions, error) {
	var pos Positions
	var err error
	o := orm.NewOrm()

	pos = Positions{Id: id}
	err = o.Read(&pos)

	if err == orm.ErrNoRows {
		return pos, nil
	}
	return pos, err
}

func GetPositionsName(id int64) string {
	var err error
	var name string
	err = utils.GetCache("GetPositionsName.id."+fmt.Sprintf("%d", id), &name)
	if err != nil {
		cache_expire, _ := beego.AppConfig.Int("cache_expire")
		var pos Positions
		o := orm.NewOrm()
		o.QueryTable(models.TableName("positions")).Filter("positionid", id).One(&pos, "name")
		name = pos.Name
		utils.SetCache("GetPositionsName.id."+fmt.Sprintf("%d", id), name, cache_expire)
	}
	return name
}
