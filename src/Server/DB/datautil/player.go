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
func GetPlayerById(id int , o orm.Ormer)  (*models.Player){
	player :=models.Player{Id: id}
	err:= o.Read(&player)
	if err == orm.ErrNoRows {
		fmt.Println("error: can`t find the player")
	} else if err == orm.ErrMissPK {
		fmt.Println("error: can`t find the primary key")
	} else {
		fmt.Println("query finished")
	}
	return &player
}




//get by account
func GetPlayerByAccount(account *models.Account , o orm.Ormer) (*models.Player){
	return GetPlayerByAccountId(account.Id , o)
}

//get by message id
func GetPlayerByAccountId(accountId int , o orm.Ormer) (*models.Player){
	var player *models.Player
	err := o.Raw("SELECT * FROM account INNER JOIN player ON account.player_id = player.id WHERE account.id = ?" ,accountId).QueryRow(&player)
	if (err ==nil){
		fmt.Println("query successfully")
	}
	return player
}

//get by history
func GetPlayer1ByHistory(history models.History , o orm.Ormer) (*models.Player){
	return GetPlayer1ByHistoryId(history.Id , o)
}

func GetPlayer1ByHistoryId(historyId int , o orm.Ormer) (*models.Player){
	var player *models.Player
	err := o.Raw("SELECT * FROM history INNER JOIN player ON history.player1_id = player.id WHERE history.id = ?" ,historyId ).QueryRow(&player)
	if (err ==nil){
		fmt.Println("query successfully")
	}
	return player
}

func GetPlayer2ByHistory(history models.History , o orm.Ormer) (*models.Player){
	return GetPlayer2ByHistoryId(history.Id , o)
}

func GetPlayer2ByHistoryId(historyId int , o orm.Ormer) (*models.Player){
	var player *models.Player
	err := o.Raw("SELECT * FROM history INNER JOIN player ON history.player2_id = player.id WHERE history.id = ?" ,historyId ).QueryRow(&player)
	if (err ==nil){
		fmt.Println("query successfully")
	}
	return player
}

//get by nickname
func GetPlayersByNickName(nickname string , o orm.Ormer) ([]*models.Player){
	var players	[]*models.Player
	num , err := o.Raw("SELECT * FROM player WHERE nickname = ?" , nickname).QueryRows(players)
	if (err ==nil){
		fmt.Println("mysql row affected nums: " , num)
	}
	return players
}





/*----------------------------------------------------------------------------------------------*/
/*----------------------------------------------------------------------------------------------*/
/*------------------------------------------insert----------------------------------------------*/
/*----------------------------------------------------------------------------------------------*/
/*----------------------------------------------------------------------------------------------*/

//insert one player
func InsertPlayer(player *models.Player , o orm.Ormer) (bool){
	res , err := o.Insert(player)
	if (err ==nil){
		fmt.Println("mysql row affected nums: " , res)
		return true
	}
	return false;
}

//insert a batch of battles
func InsertPlayers(players []*models.Player , o orm.Ormer) (bool){
	for i:=0;i<len(players);i++{
		if (!InsertPlayer(players[i] , o)){
			return false;
		}
	}
	return true
}




/*----------------------------------------------------------------------------------------------*/
/*----------------------------------------------------------------------------------------------*/
/*------------------------------------------delete----------------------------------------------*/
/*----------------------------------------------------------------------------------------------*/
/*----------------------------------------------------------------------------------------------*/

//delete player by a player object
func DeletePlayer(player *models.Player , o orm.Ormer) (bool){
	return DeletePlayerById(player.Id , o)
}

//delete a player by id
func DeletePlayerById(playerId int , o orm.Ormer) (bool){
	res, err :=o.Raw("DELETE FROM player WHERE id = ?" , playerId).Exec()
	if (err ==nil){
		num, _ := res.RowsAffected()
		fmt.Println("mysql row affected nums: " , num)
		return true
	}
	return false
}

//delete a batch of players
func DeletePlayers(players []*models.Player , o orm.Ormer) (bool){
	sql :="DELETE FROM player WHERE id IN ("
	for i:=0;i<len(players);i++{
		if (i>0){
			sql = sql + ","
		}
		sql = sql + strconv.Itoa(players[i].Id)
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
func UpdatePlayerById(playerId int, newPlayer *models.Player, o orm.Ormer)  (bool){
	oldPlayer := models.Player{Id:playerId}
	if (o.Read(&oldPlayer)==nil && oldPlayer.Id==newPlayer.Id){
		num , err := o.Update(newPlayer)
		if (err ==nil){
			fmt.Println("mysql row affected nums: " , num)
			return true
		}
	}
	return false;
}


//update by battle
func UpdatePlayer(oldPlayer *models.Player , newPlayer *models.Player , o orm.Ormer) (bool){
	return UpdatePlayerById(oldPlayer.Id , newPlayer , o)
}




