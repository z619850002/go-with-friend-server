package MainPageInfo

import (
	"fmt"
	"encoding/binary"
	"github.com/golang/protobuf/proto"
	"log"
	"strconv"
	"go-with-friend-server/src/Server/Service"
	"github.com/astaxie/beego/orm"
	"go-with-friend-server/src/Server/DB/datautil"
	"go-with-friend-server/src/Server/DB/models"
)

//增加协议号码
func addprotocolnum(num uint16, data []byte)[]byte{
	// add the protocol to the header
	m := make([]byte, 2+len(data))
	// 默认使用大端序
	binary.BigEndian.PutUint16(m, uint16(1))
	copy(m[2:], data)
	fmt.Println(m)
	//return the response
	return m
}


//获取到对战历史
func GetMyHistory(replay []byte) []byte{
	message := Service.Player_MyHistoryReq{}
	//change the json to struct
	o := orm.NewOrm()
	err := proto.Unmarshal(replay, &message)
	if err != nil {
		log.Fatal("unmarshaling error: ", err)
	}

	playerId ,_ := strconv.Atoi(message.Id)

	histories := datautil.GetPlayerHistoriesByPlayerId(playerId , o)

	response := Service.Player_MyHistoryResp{}

	response.Myhistories = getResponseByHistories(playerId , histories , o)

	data, err := proto.Marshal(&response)
	if err != nil {
		panic(err)
	}
	return addprotocolnum(1,data)
}

//获取收藏对战历史
func GetCollectedHistory(replay []byte) []byte{
	message := Service.Player_MyHistoryReq{}
	//change the json to struct
	o := orm.NewOrm()
	err := proto.Unmarshal(replay, &message)
	if err != nil {
		log.Fatal("unmarshaling error: ", err)
	}

	playerId ,_ := strconv.Atoi(message.Id)

	histories := datautil.GetCollectedHistoriesByPlayerId(playerId , o)

	response := Service.Player_MyHistoryResp{}

	response.Myhistories = getResponseByHistories(playerId , histories , o)

	data, err := proto.Marshal(&response)
	if err != nil {
		panic(err)
	}
	return addprotocolnum(1,data)
}


func getResponseByHistories(playerId int , histories []*models.History , o orm.Ormer) []*Service.PlayerHistoryInfo{
	myHistories := []*Service.PlayerHistoryInfo{}
	myPlayer := datautil.GetPlayerById(playerId , o)

	for i:=0;i<len(histories);i++{
		history :=histories[i]
		if (history.Player1.Id != myPlayer.Id){
			history.Player1 = myPlayer
			history.Player2 = datautil.GetPlayerById(history.Player2.Id , o)
		}else {
			history.Player2 = myPlayer
			history.Player1 = datautil.GetPlayerById(history.Player1.Id , o)
		}
		myHistories = append(myHistories, bindHistoryInfoByHistory(history , o))
	}
	return myHistories
}




//use the history in getting historyInfo
func bindHistoryInfoByHistory(history *models.History , o orm.Ormer) *Service.PlayerHistoryInfo{
	info :=Service.PlayerHistoryInfo{}
	//bind id
	info.Id = int64(history.Id)

	//bind players
	info.Player1Id = int64(history.Player1.Id)
	info.Player2Id = int64(history.Player2.Id)
	info.Player1Nickname = history.Player1.NickName
	info.Player2Nickname = history.Player2.NickName
	info.Player1Profile = history.Player1.Avatar
	info.Player2Profile = history.Player2.Avatar

	//bind int values
	info.Winner = int64(history.Winner)
	info.Type = int64(history.Type)
	info.Bystanders = int64(history.Bystanders)
	info.Likenum = int64(history.Like)
	info.Step = int64(history.Step)
	info.Commentnum = int64(history.Commentsnumber)
	info.Player1Score = int64(history.Player1Score)
	info.Player2Score = int64(history.Player2Score)

	//bind string values
	info.Finishtime = history.LastTurnTime.String()
	info.Board = history.ChessBoard

	return &info
}
