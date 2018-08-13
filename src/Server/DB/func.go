package DB

import (
	"fmt"
	"github.com/astaxie/beego/orm"
	_ "github.com/go-sql-driver/mysql"
	"crypto/md5"
	"github.com/z619850002/go-with-friend-server/src/Server/DB/models"
	"github.com/z619850002/go-with-friend-server/src/Server/DB/datautil"
	)





func init() {
	//datasource := "root:ztjztj120@tcp(127.0.0.1:3306)/chess"
	datasource := "root:12345678@/gowithfriend?charset=utf8"



	// set default database
	orm.RegisterDriver("mysql", orm.DRMySQL)
	orm.RegisterDataBase("default", "mysql", datasource, 30)
	// register model
	orm.RegisterModel( new(models.Invitation), new(models.Account), new(models.Player), new(models.History), new(models.GoodChoice), new(models.Message) , new(models.Level))
	// create table
	orm.RunSyncdb("default", false, true)



	o := orm.NewOrm()
	o.Using("default")
	
}
func getmd5(str string)string{
	data := []byte(str)
	has := md5.Sum(data)
	return  fmt.Sprintf("%x", has)
}

//regist
func CreatePlayer(kind, username, pwd string) int{
	o := orm.NewOrm()
	//check the nickname
	students := make([]models.Account,1)
	//sql := "select * from account where identity_type = \""+kind+"\" and identifier = \""+username+"\""
	sql := "select * from account where identity_type = ? and identifier = ?"
	fmt.Println(sql)
	num, _ := o.Raw(sql,kind,username).QueryRows(&students)
	if num == 0 {
		//玩家用户名合法，可以存入数据库中
		player := models.Player{NickName:username, Avatar:"default", AvatarEdge:"default"}
		o.Insert(&player)
		//插入auths表
		account := models.Account{Player:&player, IdentityType:kind, Identifier:username, Credential:getmd5(pwd)}
		o.Insert(&account)
		fmt.Println("regist success")
		return 0
	} else {
		fmt.Println("regist fail : duplicate account")
		return 2
	}
}



//login
func CheckAccount(kind,username,pwd string) (int,int){
	o := orm.NewOrm()
	//check the nickname
	students := make([]models.Account,1)
	//sql := "select * from account where identity_type = \""+kind+"\" and identifier = \"" + username + "\""
	sql := "select * from account where identity_type = ? and identifier = ?"
	fmt.Println(sql)
	num, _ := o.Raw(sql,kind,username).QueryRows(&students)
	if num == 0 {
		fmt.Println("login fail : nonexistent account！")
		return 3,0
	} else {
		fmt.Println("correct account！")
		//sql = "select * from account where identity_type = \""+kind+"\" and identifier = \""+username+"\" and credential = \""+getmd5(pwd)+"\""
		sql = "select * from account where identity_type = ? and identifier = ? and credential = ?"
		fmt.Println(sql)
		num, _ = o.Raw(sql,kind,username,getmd5(pwd)).QueryRows(&students)
		if num == 0{
			fmt.Println("login fail : unmatched account and password")
			return 2,0
		} else {
			fmt.Println("login success")
			return 0,students[0].Id
		}
	}
}


//get player infomation
func Getplayerinfo(id string)(string ,string, string) {
	o := orm.NewOrm();
	var player []orm.ParamsList
	//sql := "select * from player where id = \"" + id + "\""
	sql := "select p.nick_name, p.avatar, p.avatar_edge from player p where p.id = ?"
	fmt.Println(sql)
	//可能需要首先验证id
	fmt.Println(id)
	_ , err := o.Raw(sql,id).ValuesList(&player)
	if err != nil {
		fmt.Println(err)
		return  "","",""
	} else {
		fmt.Println(player[0][0].(string))
		fmt.Println(player[0][1].(string))
		fmt.Println(player[0][2].(string))
		return 	player[0][0].(string), player[0][1].(string), player[0][2].(string)
	}
}


func GetTheirturn(id string)[]orm.ParamsList{
	fmt.Println("theirturn")
	o := orm.NewOrm();
	var lists1 []orm.ParamsList
	var lists2 []orm.ParamsList
	//自己所持的棋子颜色是黑色，即自己为player1，turn就为白棋下，即：BATTLE_TURN_PLAYER2 （false）
	sql := "select h.id,p.id ,p.nick_name,p.avatar,h.bystanders,h.step,h.like,h.last_turn_time,h.type,h.commentsnumber,p.isonline " +
		"from history h left join player p on p.id = h.player2_id " +
		"where h.winner = ? and h.turn = ? and h.player1_id = ?"
	o.Raw(sql,datautil.HISTORY_WINNER_PROCESSING,datautil.BATTLE_TURN_PLAYER2,id).ValuesList(&lists1)
	//自己所持的棋子颜色是白色，即自己为player1，turn就为黑棋下，即：BATTLE_TURN_PLAYER1 （true）
	sql = "select h.id,p.id ,p.nick_name,p.avatar,h.bystanders,h.step,h.like,h.last_turn_time,h.type,h.commentsnumber,p.isonline " +
		"from history h left join player p on p.id = h.player1_id " +
		"where h.winner = ? and h.turn = ? and h.player2_id = ?"
	o.Raw(sql,datautil.HISTORY_WINNER_PROCESSING,datautil.BATTLE_TURN_PLAYER1,id).ValuesList(&lists2)
	//第三种情况就是在invitation表中
	lists := append(lists1,lists2...)

	return lists
}


func Getmyturn(id string)[]orm.ParamsList{
	fmt.Println("getmyturn")
	o := orm.NewOrm();
	var lists []orm.ParamsList
	//自己所持的棋子颜色是黑色，即自己为player2，turn就为白棋下，即：BATTLE_TURN_PLAYER2 （false）
	sql := "select h.id,p.id ,p.nick_name,p.avatar,h.bystanders,h.step,h.like,h.last_turn_time,h.type,h.commentsnumber,p.isonline " +
		"from history h inner join player p on p.id = h.player2_id " +
		"where h.winner = ? and h.turn = ? and h.player1_id = ? " +
		" union "+
		"select h.id,p.id ,p.nick_name,p.avatar,h.bystanders,h.step,h.like,h.last_turn_time,h.type,h.commentsnumber,p.isonline " +
		"from history h left join player p on p.id = h.player1_id " +
		"where h.winner = ? and h.turn = ? and h.player2_id = ?"
	o.Raw(sql,datautil.HISTORY_WINNER_PROCESSING,datautil.BATTLE_TURN_PLAYER1,id,datautil.HISTORY_WINNER_PROCESSING,datautil.BATTLE_TURN_PLAYER2,id).ValuesList(&lists)
	return lists
}


func Getpending(id string)[]orm.ParamsList{
	fmt.Println("get pending")
	o := orm.NewOrm();
	var lists []orm.ParamsList
	//查询自己被邀请的所有的棋局
	sql :=  "select p2.id,p2.nick_name,p2.avatar,i.size,i.turn,p2.isonline,p2.dan,i.time,i.type " +
		"from invitation i, player p1, player p2 " +
		"where i.inviter_id = p2.id and i.invitee_id = p1.id and p1.id = ?"
	o.Raw(sql,id).ValuesList(&lists)
	return lists
}

