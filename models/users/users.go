package users

import (
	"fmt"
	"math/rand"
	"time"

	"github.com/astaxie/beego"
	"github.com/astaxie/beego/orm"
	"github.com/ninefive/hims/models"
	"github.com/ninefive/hims/utils"
)

type Users struct {
	Id       int64         `orm:"pk;column(userid);"`
	Profile  *UsersProfile `orm:"rel(one);"`
	Username string
	Password string
	Avatar   string
	Status   int
}

type UsersProfile struct {
	Id          int64 `orm:"pk;column(userid);"`
	Realname    string
	Sex         int
	Birth       string
	Email       string
	Webchat     string
	Qq          string
	Phone       string
	Tel         string
	Address     string
	Emercontact string
	Emerphone   string
	Departid    int64
	Positionid  int64
	Lognum      int
	Ip          string
	Lasted      int64
}

func (this *Users) TableName() string {
	return models.TableName("users")
}

func init() {
	orm.RegisterModel(new(Users))
	orm.RegisterModelWithPrefix("hims_", new(UsersProfile))
}

//登録
func LoginUser(username, password string) (err error, user Users) {
	o := orm.NewOrm()
	qs := o.QueryTable(models.TableName("users"))
	cond := orm.NewCondition()

	cond = cond.And("username", username)
	pwdmd5 := utils.Md5(password)
	cond = cond.And("password", pwdmd5)
	cond = cond.And("status", 1)

	qs = qs.SetCond(cond)
	var users Users
	err = qs.Limit(1).One(&users, "userid", "username", "avatar")
	fmt.Println(err)
	if err == nil {
		o.Raw("UPDATE hims_users_profile SET lasted = ?,lognum=lognum+? WHERE userid = ?", time.Now().Unix(), 1, users.Id).Exec()
	}
	return err, users
}

//得到用戶信息
func GetUser(id int64) (Users, error) {
	var user Users
	var err error
	o := orm.NewOrm()

	user = Users{Id: id}
	err = o.Read(&user)
	if err == orm.ErrNoRows {
		return user, nil
	}
	return user, err
}

func GetRealName(id int64) string {
	var err error
	var realname string

	err = utils.GetCache("GetRealname.id."+fmt.Sprintf("%d", id), &realname)
	if err != nil {
		cache_expire, _ := beego.AppConfig.Int("cache_expire")
		var user UsersProfile
		o := orm.NewOrm()
		o.QueryTable(models.TableName("users_profile")).Filter("userid", id).One(&user, "realname")
		realname = user.Realname
		utils.SetCache("GetRealname.id."+fmt.Sprintf("%d", id), realname, cache_expire)
	}
	return realname
}

func GetUserEmail(id int64) string {
	var err error
	var email string

	err = utils.GetCache("GetUserEmail.id."+fmt.Sprintf("%d", id), &email)
	if err != nil {
		cache_expire, _ := beego.AppConfig.Int("cache_expire")
		var user UsersProfile
		o := orm.NewOrm()
		o.QueryTable(models.TableName("users_profile")).Filter("userid", id).One(&user, "email")
		email = user.Email
		utils.SetCache("GetUserEmail.id."+fmt.Sprintf("%d", id), email, cache_expire)
	}
	return email
}

func GetAvatarUserid(id int64) string {
	var err error
	var avatar string

	err = utils.GetCache("GetAvatarUserid.id."+fmt.Sprintf("%d", id), &avatar)
	if err != nil {
		cache_expire, _ := beego.AppConfig.Int("cache_expire")
		var user Users
		o := orm.NewOrm()
		o.QueryTable(models.TableName("users")).Filter("userid", id).One(&user, "avatar")
		avatar = user.Avatar
		utils.SetCache("GetAvatarUserid.id."+fmt.Sprintf("%d", id), avatar, cache_expire)
	}
	if avatar == "" {
		return fmt.Sprintf("/static/img/avatar/%d.jpg", rand.Intn(5))
	}
	return avatar
}

func GetPositionsNameForUserid(id int64) string {
	var err error
	var position string

	err = utils.GetCache("GetPositionsNameForUserid.id."+fmt.Sprintf("%d", id), &position)
	if err != nil {
		cache_expire, _ := beego.AppConfig.Int("cache_expire")
		var user UsersProfile
		o := orm.NewOrm()
		o.QueryTable(models.TableName("users_profile")).Filter("userid", id).One(&user, "positionid")
		position = GetPositionsName(user.Positionid)
		utils.SetCache("GetPositionsNameForUserid.id."+fmt.Sprintf("%d", id), position, cache_expire)
	}
	return position
}

func GetDepartmentsNameForUserid(id int64) string {
	var err error
	var depart string
	err = utils.GetCache("GetDepartmentsNameForUserid.id."+fmt.Sprintf("%d", id), &depart)
	if err != nil {
		cache_expire, _ := beego.AppConfig.Int("cache_expire")
		var user UsersProfile
		o := orm.NewOrm()
		o.QueryTable(models.TableName("users_profile")).Filter("userid", id).One(&user, "departid")
		depart = GetDepartsName(user.Departid)
		utils.SetCache("GetDepartmentsNameForUserid.id."+fmt.Sprintf("%d", id), depart, cache_expire)
	}
	return depart
}

//得到用戶詳情信息
func GetProfile(id int64) (UsersProfile, error) {
	var pro UsersProfile
	var err error
	o := orm.NewOrm()

	pro = UsersProfile{Id: id}
	err = o.Read(&pro)

	if err == orm.ErrNoRows {
		return pro, nil
	}
	return pro, err
}

//修改個人信息
func UpdateProfile(id int64, updPro UsersProfile) error {
	var pro UsersProfile
	o := orm.NewOrm()
	pro = UsersProfile{Id: id}

	pro.Realname = updPro.Realname
	pro.Sex = updPro.Sex
	pro.Birth = updPro.Birth
	pro.Email = updPro.Email
	pro.Webchat = updPro.Webchat
	pro.Qq = updPro.Qq
	pro.Phone = updPro.Phone
	pro.Tel = updPro.Tel
	pro.Address = updPro.Address
	pro.Emercontact = updPro.Emercontact
	pro.Emerphone = updPro.Emerphone

	if updPro.Departid > 0 {
		pro.Departid = updPro.Departid
		pro.Positionid = updPro.Positionid
		_, err := o.Update(&pro)
		return err
	} else {
		_, err := o.Update(&pro, "realname", "sex", "birth", "email", "webchat", "qq", "phone", "tel", "address", "emercontact", "emerphone")
		return err
	}
}

//修改用戶
func UpdateUser(id int64, updUser Users) error {
	var user Users
	o := orm.NewOrm()
	user = Users{Id: id}

	user.Username = updUser.Username
	if updUser.Password != "" {
		user.Password = utils.Md5(updUser.Password)
		_, err := o.Update(&user, "username", "password")
		return err
	} else {
		_, err := o.Update(&user, "username")
		return err
	}
}

//修改密碼
func UpdatePassword(id int64, oldPwd, newPwd string) error {
	o := orm.NewOrm()

	user := Users{Id: id}
	err := o.Read(&user)
	if nil != err {
		return err
	} else {
		if user.Password == utils.Md5(oldPwd) {
			user.Password = utils.Md5(newPwd)
			_, err := o.Update(&user)
			return err
		} else {
			return fmt.Errorf("驗證出錯")
		}
	}
}
