package datautil

import (
	"github.com/astaxie/beego/orm"
	"fmt"
	"go-with-friend-server/src/Server/DB/models"
)

/*----------------------------------------------------------------------------------------------*/
/*----------------------------------------------------------------------------------------------*/
/*------------------------------------------query-----------------------------------------------*/
/*----------------------------------------------------------------------------------------------*/
/*----------------------------------------------------------------------------------------------*/

//get by id
func GetLevelById(id int , o orm.Ormer) (*models.Level){
	level :=models.Level{Id: id}
	err:= o.Read(&level)
	if err == orm.ErrNoRows {
		fmt.Println("error: can`t find the player")
	} else if err == orm.ErrMissPK {
		fmt.Println("error: can`t find the primary key")
	} else {
		fmt.Println("query finished")
	}
	return &level
}



func GetLevelByPlayerId(playerId int , o orm.Ormer) (*models.Level){
	var level *models.Level
	err := o.Raw("SELECT * FROM level WHERE player_id = ?" , playerId).QueryRow(&level)
	if (err==nil){
		fmt.Println("query successfully")
	} else {
		fmt.Println(err)
	}
	return level
}



//
///*----------------------------------------------------------------------------------------------*/
///*----------------------------------------------------------------------------------------------*/
///*------------------------------------------insert----------------------------------------------*/
///*----------------------------------------------------------------------------------------------*/
///*----------------------------------------------------------------------------------------------*/
//

func InsertLevelByPlayerId(playerId int , o orm.Ormer) (bool){
	level := models.Level{Player:&models.Player{Id:playerId}}
	practiceStr := []byte("")
	for i:=0;i<500;i++{
		practiceStr = append(practiceStr, '0' )
	}
	level.PracticeInfo = string(practiceStr)
	res , err := o.Insert(&level)
	if (err ==nil){
		fmt.Println("mysql row affected nums: " , res)
		return true
	}else {
		fmt.Println(err)
	}
	return false;
}



//
///*----------------------------------------------------------------------------------------------*/
///*----------------------------------------------------------------------------------------------*/
///*------------------------------------------delete----------------------------------------------*/
///*----------------------------------------------------------------------------------------------*/
///*----------------------------------------------------------------------------------------------*/
//

//
//
///*----------------------------------------------------------------------------------------------*/
///*----------------------------------------------------------------------------------------------*/
///*------------------------------------------update----------------------------------------------*/
///*----------------------------------------------------------------------------------------------*/
///*----------------------------------------------------------------------------------------------*/
//


//add the fall number
func UpdateLevelAddFallNumber(playerId int , num int , o orm.Ormer) (bool){
	level := GetLevelByPlayerId(playerId , o)
	level.FallNumber = level.FallNumber + int64(num)
	n , err := o.Update(level)
	if (err == nil){
		fmt.Println(n,"rows are affected")
		return true
	}else{
		fmt.Println(err)
	}
	return false
}

//add the take number
func UpdateLevelAddTakeNumber(playerId int , num int , o orm.Ormer) (bool){
	level := GetLevelByPlayerId(playerId , o)
	level.TakeNumber = level.TakeNumber + int64(num)
	n , err := o.Update(level)
	if (err == nil){
		fmt.Println(n,"rows are affected")
		return true
	}else{
		fmt.Println(err)
	}
	return false
}

//add the practice information

func UpdateLevelAddPracticeInfo(playerId int , finishId []int , o orm.Ormer) (bool){
	level := GetLevelByPlayerId(playerId , o)
	practiceInfoStr := []byte(level.PracticeInfo)
	for i:=0;i<len(finishId);i++{
		practiceInfoStr[finishId[i]] = '1'
	}
	level.PracticeInfo = string(practiceInfoStr)
	n , err := o.Update(level)
	if (err ==nil){
		println(n , "rows are affected")
		return true
	}
	return false
}


func UpdateLevelAddFriendsBattle(playerId int , num int , o orm.Ormer) (bool){
	level := GetLevelByPlayerId(playerId , o)
	level.FriendsBattle = level.FriendsBattle + num
	n , err := o.Update(level)
	if (err ==nil){
		println(n , "rows are affected")
		return true
	}
	return false
}

func UpdateLevelAddSmartBattle(playerId int , num int , o orm.Ormer) (bool){
	level := GetLevelByPlayerId(playerId , o)
	level.SmartBattle = level.SmartBattle + num
	n , err := o.Update(level)
	if (err ==nil){
		println(n , "rows are affected")
		return true
	}
	return false
}

func UpdateLevelAddCircumstance(playerId int , num int , o orm.Ormer) (bool){
	level := GetLevelByPlayerId(playerId , o)
	level.Circumstance = level.Circumstance + num
	n , err := o.Update(level)
	if (err ==nil){
		println(n , "rows are affected")
		return true
	}
	return false
}

