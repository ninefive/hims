package users

import (
	"github.com/astaxie/beego/orm"
	"github.com/ninefive/hims/models"
)

type Notices struct {
	Id      int64 `orm:"pk;column(noticeid);"`
	Title   string
	Content string
	Created int64
	Status  int
}

func (this *Notices) TableName() string {
	return models.TableName("notices")
}

func init() {
	orm.RegisterModel(new(Notices))
}

func GetNotices(id int64) (Notices, error) {
	var not Notices
	var err error
	o := orm.NewOrm()

	not = Notices{Id: id}
	err = o.Read(&not)

	if err == orm.ErrNoRows {
		return not, nil
	}
	return not, err
}

func UpdateNotices(id int64, upd Notices) error {
	var not Notices
	o := orm.NewOrm()
	not = Notices{Id: id}

	not.Title = upd.Title
	not.Content = upd.Content
	_, err := o.Update(&not, "title", "content")
	return err
}
