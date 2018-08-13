package datautil

import (
	"github.com/astaxie/beego/orm"
	"fmt"
	"strconv"
	"go-with-friend-server/src/Server/DB/models"
)

/*----------------------------------------------------------------------------------------------*/
/*----------------------------------------------------------------------------------------------*/
/*------------------------------------------query-----------------------------------------------*/
/*----------------------------------------------------------------------------------------------*/
/*----------------------------------------------------------------------------------------------*/

//get by id
func GetGoodChoiceById(id int , o orm.Ormer)  (*models.GoodChoice){
	goodChoice :=models.GoodChoice{Id: id}
	err:= o.Read(&goodChoice)
	if err == orm.ErrNoRows {
		fmt.Println("error: can`t find the good choice")
	} else if err == orm.ErrMissPK {
		fmt.Println("error: can`t find the primary key")
	} else {
		fmt.Println("query finished")
	}
	return &goodChoice
}




//get by player
func GetGoodChoicesByPlayer(player *models.Player , o orm.Ormer) ([]*models.GoodChoice){
	return GetGoodChoicesByPlayerId(player.Id , o)
}

//get by player id
func GetGoodChoicesByPlayerId(playerId int , o orm.Ormer) ([]*models.GoodChoice){
	var goodChoices []*models.GoodChoice
	num , err := o.Raw("SELECT * FROM praised_goodchoice INNER JOIN good_choice ON praised_goodchoice.good_choice_id = good_choice.id WHERE player_id = ?" , playerId).QueryRows(&goodChoices)
	if (err ==nil){
		fmt.Println("good choice nums: " , num)
	}
	return goodChoices
}



//get by history
func GetGoodChoicesByHistory(history models.History , o orm.Ormer) ([]*models.GoodChoice){
	return GetGoodChoicesByHistoryId(history.Id , o)
}

func GetGoodChoicesByHistoryId(historyId int , o orm.Ormer) ([]*models.GoodChoice){
	var goodChoice []*models.GoodChoice
	num, err := o.Raw("SELECT * FROM good_choice WHERE histories_id = ?" ,historyId ).QueryRows(&goodChoice)
	if (err ==nil){
		fmt.Println("good choice nums: " , num)
	}
	return goodChoice
}

//get by praised number
func GetGoodChoicesByPraiseNumber(praiseNumber int , o orm.Ormer , comparator int) ([] *models.GoodChoice){
	var goodChoices []*models.GoodChoice
	comparatorStr := ""
	if (comparator==0){
		comparatorStr = "="
	} else if (comparator>0){
		comparatorStr = ">"
	}else {
		comparatorStr = "<"
	}
	num , err := o.Raw("SELECT * FROM good_choice WHERE praise_number " + comparatorStr + " ?" , praiseNumber).QueryRows(&goodChoices)
	if (err ==nil){
		fmt.Println("good choice nums: " , num)
	}
	return goodChoices
}

func GetGoodChoicesByPraiseNumberEqual(praiseNumber int , o orm.Ormer) ([] *models.GoodChoice){
	return GetGoodChoicesByPraiseNumber(praiseNumber , o , 0)
}

func GetGoodChoicesByPraiseNumberMoreThan(praiseNumber int , o orm.Ormer) ([] *models.GoodChoice){
	return GetGoodChoicesByPraiseNumber(praiseNumber , o , 1)
}

func GetGoodChoicesByPraiseNumberLessThan(praiseNumber int , o orm.Ormer) ([] *models.GoodChoice){
	return GetGoodChoicesByPraiseNumber(praiseNumber , o , -1)
}





/*----------------------------------------------------------------------------------------------*/
/*----------------------------------------------------------------------------------------------*/
/*------------------------------------------insert----------------------------------------------*/
/*----------------------------------------------------------------------------------------------*/
/*----------------------------------------------------------------------------------------------*/

//insert one good choice
func InsertGoodChoice(goodChoice *models.GoodChoice , o orm.Ormer) (bool){
	res , err := o.Insert(goodChoice)
	if (err ==nil){
		fmt.Println("mysql row affected nums: " , res)
		return true
	}
	return false;
}

//insert a batch of good choices
func InsertGoodChoices(goodChoices []*models.GoodChoice , o orm.Ormer) (bool){
	sql := "INSERT INTO good_choice VALUES"
	for i:=0;i<len(goodChoices);i++{
		goodChoice := goodChoices[i]
		if (i>0){
			sql = sql + ","
		}
		sql = sql + "("+strconv.Itoa(goodChoice.Id)
		sql = sql +"," +strconv.Itoa(goodChoice.Histories.Id)
		sql = sql +"," +strconv.Itoa(goodChoice.Pos)
		sql = sql +"," +strconv.Itoa(goodChoice.PraiseNumber)+")"
	}
	res, err :=o.Raw(sql).Exec()
	if (err ==nil){
		num, _ := res.RowsAffected()
		fmt.Println("mysql row affected nums: " , num)
		return true
	}
	return false;
}


//praise
func InsertPraiseChoice(goodChoiceId int , playerId int,  o orm.Ormer) (bool){
	sql := "INSERT INTO praised_goodchoice(`player_id`, `good_choice_id`) VALUES(? ,?)"
	res, err :=o.Raw(sql ,strconv.Itoa(goodChoiceId) , strconv.Itoa(playerId)).Exec()
	if (err ==nil){
		num, _ := res.RowsAffected()
		fmt.Println("mysql row affected nums: " , num)
		return true
	}
	return false;
}



/*----------------------------------------------------------------------------------------------*/
/*----------------------------------------------------------------------------------------------*/
/*------------------------------------------delete----------------------------------------------*/
/*----------------------------------------------------------------------------------------------*/
/*----------------------------------------------------------------------------------------------*/

//delete a good choice by an object
func DeleteGoodChoice(goodChoice *models.GoodChoice , o orm.Ormer) (bool){
	return DeleteGoodChoiceById(goodChoice.Id , o)
}

//delete a good choice by id
func DeleteGoodChoiceById(goodChoiceId int , o orm.Ormer) (bool){
	res, err :=o.Raw("DELETE FROM good_choice WHERE id = ?" , goodChoiceId).Exec()
	if (err ==nil){
		num, _ := res.RowsAffected()
		fmt.Println("mysql row affected nums: " , num)
		return true
	}
	return false
}




/*----------------------------------------------------------------------------------------------*/
/*----------------------------------------------------------------------------------------------*/
/*------------------------------------------update----------------------------------------------*/
/*----------------------------------------------------------------------------------------------*/
/*----------------------------------------------------------------------------------------------*/


//update by id
func UpdateGoodChoiceById(goodChoiceId int, newGoodChoice *models.GoodChoice, o orm.Ormer)  (bool){
	oldGoodChoice := models.GoodChoice{Id:goodChoiceId}
	if (o.Read(&oldGoodChoice)==nil && oldGoodChoice.Id==newGoodChoice.Id){
		num , err := o.Update(newGoodChoice)
		if (err ==nil){
			fmt.Println("mysql row affected nums: " , num)
			return true
		}
	}
	return false;
}


//update by good choice object
func UpdateGoodChoice(oldGoodChoice *models.GoodChoice , newGoodChoice *models.GoodChoice , o orm.Ormer) (bool){
	return UpdateGoodChoiceById(oldGoodChoice.Id , newGoodChoice , o)
}


