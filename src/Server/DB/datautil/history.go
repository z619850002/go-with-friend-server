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
func GetHistoryById(historyId int , o orm.Ormer) (*models.History){
	history :=models.History{Id: historyId}
	err:= o.Read(&history)
	if err == orm.ErrNoRows {
		fmt.Println("error: can`t find the battle")
	} else if err == orm.ErrMissPK {
		fmt.Println("error: can`t find the primary key")
	} else {
		fmt.Println("query finished")
	}
	return &history
}


//get by player
func GetPlayerHistoriesByPlayer(player *models.Player , o orm.Ormer) ([]*models.History){
	return GetPlayerHistoriesByPlayerId(player.Id , o)
}

func GetCollectedHistoriesByPlayer(player models.Player ,o orm.Ormer) ([]*models.History){
	return GetCollectedHistoriesByPlayerId(player.Id , o)
}

//get by player id
func GetPlayerHistoriesByPlayerId(playerId int , o orm.Ormer) ([]*models.History){
	var histories []*models.History
	num , err := o.Raw("SELECT * FROM history WHERE player1_id = ? OR player2_id = ?" , playerId , playerId).QueryRows(&histories)
	if (err ==nil){
		fmt.Println("histories nums: " , num)
	}
	return histories
}

func GetCollectedHistoriesByPlayerId(playerId int , o orm.Ormer) ([]*models.History){
	var histories []*models.History
	num , err := o.Raw("SELECT * FROM collected_history INNER JOIN history ON collected_history.history_id = history.id WHERE player_id = ?" , playerId).QueryRows(&histories)
	if (err ==nil){
		fmt.Println("histories nums: " , num)
	}
	return histories
}

//get by bystanders
//match = 0      equal
//match = 1		 more than
//match = -1	 less than
func GetHistoriesByBystanders(bystanders int , o orm.Ormer , comparator int) ([] *models.History){
	var histories []*models.History
	comparatorStr := ""
	if (comparator==0){
		comparatorStr = "="
	} else if (comparator>0){
		comparatorStr = ">"
	}else {
		comparatorStr = "<"
	}
	num , err := o.Raw("SELECT * FROM history WHERE bystanders " + comparatorStr + " ?" , bystanders).QueryRows(&histories)
	if (err ==nil){
		fmt.Println("histories nums: " , num)
	}
	return histories
}

func GetHistoriesByBystandersEqual(bystanders int , o orm.Ormer) ([] *models.History){
	return GetHistoriesByBystanders(bystanders , o , 0)
}

func GetHistoriesByBystandersMoreThan(bystanders int , o orm.Ormer) ([] *models.History){
	return GetHistoriesByBystanders(bystanders , o , 1)
}

func GetHistoriesByBystandersLessThan(bystanders int , o orm.Ormer) ([] *models.History){
	return GetHistoriesByBystanders(bystanders , o , -1)
}

//get by like
func GetHistoriesByLike(like int , o orm.Ormer , comparator int) ([] *models.History){
	var histories []*models.History
	comparatorStr := ""
	if (comparator==0){
		comparatorStr = "="
	} else if (comparator>0){
		comparatorStr = ">"
	}else {
		comparatorStr = "<"
	}
	num , err := o.Raw("SELECT * FROM history WHERE like " + comparatorStr + " ?" , like).QueryRows(&histories)
	if (err ==nil){
		fmt.Println("histories nums: " , num)
	}
	return histories
}

func GetHistoriesByLikeEqual(like int , o orm.Ormer) ([] *models.History){
	return GetHistoriesByBystanders(like , o , 0)
}

func GetHistoriesByLikeMoreThan(like int , o orm.Ormer) ([] *models.History){
	return GetHistoriesByBystanders(like , o , 1)
}

func GetHistoriesByLikeLessThan(like int , o orm.Ormer) ([] *models.History){
	return GetHistoriesByBystanders(like , o , -1)
}

//get by winner

func GetHistoriesByWinner(winner models.Player , o orm.Ormer) ([]* models.History){
	return GetHistoriesByWinnerId(winner.Id , o)
}

func GetHistoriesByWinnerId(winnerPlayerId int , o orm.Ormer) ([] *models.History){
	var histories []*models.History
	num , err := o.Raw("SELECT *  FROM history  WHERE (winner = 1 AND player1_id = ?) OR (winner = 2 AND player2_id = ?)" , winnerPlayerId , winnerPlayerId).QueryRows(&histories)
	if (err ==nil){
		fmt.Println("histories nums: " , num)
	}
	return histories
}

//get by loser

func GetHistoriesByLoser(loser models.Player , o orm.Ormer) ([]* models.History){
	return GetHistoriesByLoserId(loser.Id , o)
}

func GetHistoriesByLoserId(loserPlayerId int , o orm.Ormer) ([] *models.History){
	var histories []*models.History
	num , err := o.Raw("SELECT *  FROM history  WHERE (winner = 1 AND player2_id = ?) OR (winner = 2 AND player1_id = ?)" , loserPlayerId , loserPlayerId).QueryRows(&histories)
	if (err ==nil){
		fmt.Println("histories nums: " , num)
	}
	return histories
}

//get by type
func GetHistoriesByType(t int , o orm.Ormer) ([] *models.History){
	var histories []*models.History
	num , err := o.Raw("SELECT * FROM history WHERE type = ?" , t).QueryRows(&histories)
	if (err ==nil){
		fmt.Println("histories nums: " , num)
	}
	return histories
}

/*----------------------------------------------------------------------------------------------*/
/*----------------------------------------------------------------------------------------------*/
/*------------------------------------------insert----------------------------------------------*/
/*----------------------------------------------------------------------------------------------*/
/*----------------------------------------------------------------------------------------------*/

//insert one battle
func InsertHistory(history *models.History , o orm.Ormer) (bool){
	num , err:= o.Insert(history)
	if (err ==nil){
		fmt.Println("mysql row affected nums: " , num)
		return true
	}
	return false;
}

//collect history
func CollectHistory(historyId int , playerId int , o orm.Ormer) (bool){
	sql := "INSERT INTO collected_history(`player_id`, `history_id`) VALUES(? ,?)"
	res, err :=o.Raw(sql ,strconv.Itoa(playerId) , strconv.Itoa(historyId)).Exec()
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

//delete battle by a battle object
func DeleteHistory(history *models.History , o orm.Ormer) (bool){
	return DeleteHistoryById(history.Id , o)
}

//delete battle by a battle id
func DeleteHistoryById(historyId int , o orm.Ormer) (bool){
	res, err :=o.Raw("DELETE FROM history WHERE id = ?" , historyId).Exec()
	if (err ==nil){
		num, _ := res.RowsAffected()
		fmt.Println("mysql row affected nums: " , num)
		return true
	}
	return false
}

//delete a batch of battle
func DeleteHistories(histories []*models.History , o orm.Ormer) (bool){
	sql :="DELETE FROM history WHERE ID IN ("
	for i:=0;i<len(histories);i++{
		if (i>0){
			sql = sql + ","
		}
		sql = sql + strconv.Itoa(histories[i].Id)
	}
	sql = sql + ")"
	res, err :=o.Raw(sql).Exec()
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
func UpdateHistoryById(historyId int, newHistory *models.History, o orm.Ormer)  (bool){
	oldHistory := models.History{Id:historyId}
	if (o.Read(&oldHistory)==nil && oldHistory.Id==newHistory.Id){
		num , err := o.Update(newHistory)
		if (err ==nil){
			fmt.Println("mysql row affected nums: " , num)
			return true
		}
	}
	return false;
}


//update by battle
func UpdateHistory(oldHistory *models.History , newHistory *models.History , o orm.Ormer) (bool){
	return UpdateHistoryById(oldHistory.Id , newHistory , o)
}


