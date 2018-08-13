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
func GetAccountById(id int , o orm.Ormer)  (*models.Account){
	account :=models.Account{Id: id}
	err:= o.Read(&account)
	if err == orm.ErrNoRows {
		fmt.Println("error: can`t find the player")
	} else if err == orm.ErrMissPK {
		fmt.Println("error: can`t find the primary key")
	} else {
		fmt.Println("query finished")
	}
	return &account
}


//get by account
func GetAccountByPlayer(player *models.Player , o orm.Ormer) (*models.Account){
	return GetAccountByPlayerId(player.Id , o)
}

//get by message id
func GetAccountByPlayerId(playerId int , o orm.Ormer) (*models.Account){
	var account *models.Account
	err := o.Raw("SELECT * FROM account WHERE player_id = ?" ,playerId).QueryRow(&account)
	if (err ==nil){
		fmt.Println("query successfully")
	}
	return account
}



/*----------------------------------------------------------------------------------------------*/
/*----------------------------------------------------------------------------------------------*/
/*------------------------------------------insert----------------------------------------------*/
/*----------------------------------------------------------------------------------------------*/
/*----------------------------------------------------------------------------------------------*/

//insert one account
func InsertAccount(account *models.Account , o orm.Ormer) (bool){
	res , err := o.Insert(account)
	if (err ==nil){
		fmt.Println("mysql row affected nums: " , res)
		return true
	}
	return false;
}

//insert a batch of accounts
func InsertAccounts(accounts []*models.Account , o orm.Ormer) (bool){
	for i:=0;i<len(accounts);i++{
		if (!InsertAccount(accounts[i] , o)){
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

//delete account by a account object
func DeleteAccount(account *models.Account , o orm.Ormer) (bool){
	return DeleteAccountById(account.Id , o)
}

//delete a account by id
func DeleteAccountById(accountId int , o orm.Ormer) (bool){
	res, err :=o.Raw("DELETE FROM account WHERE id = ?" , accountId).Exec()
	if (err ==nil){
		num, _ := res.RowsAffected()
		fmt.Println("mysql row affected nums: " , num)
		return true
	}
	return false
}

//delete a batch of accounts
func DeleteAccounts(accounts []*models.Account , o orm.Ormer) (bool){
	sql :="DELETE FROM account WHERE id IN ("
	for i:=0;i<len(accounts);i++{
		if (i>0){
			sql = sql + ","
		}
		sql = sql + strconv.Itoa(accounts[i].Id)
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
func UpdateAccountById(accountId int, newAccount *models.Account, o orm.Ormer)  (bool){
	oldAccount := models.Account{Id:accountId}
	if (o.Read(&oldAccount)==nil && oldAccount.Id==newAccount.Id){
		num , err := o.Update(newAccount)
		if (err ==nil){
			fmt.Println("mysql row affected nums: " , num)
			return true
		}
	}
	return false;
}


//update by account object
func UpdateAccount(oldAccount *models.Account , newAccount *models.Account , o orm.Ormer) (bool){
	return UpdateAccountById(oldAccount.Id , newAccount , o)
}



