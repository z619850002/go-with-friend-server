package RegistAndLogin

import (
	"log"

	"github.com/golang/protobuf/proto"
	"fmt"
	"strconv"
	"go-with-friend-server/src/Server/DB"
	"go-with-friend-server/src/Server/Service"
	"encoding/binary"
)

var imageUrl = "http://172.26.163.124:10086/"

func Addprotocolnum(num uint16, data []byte)[]byte{
	// add the protocol to the header
	m := make([]byte, 2+len(data))
	// 默认使用大端序
	binary.BigEndian.PutUint16(m, uint16(num))
	copy(m[2:], data)
	fmt.Println(m)
	//return the response
	return m
}

func Login(replay []byte) []byte{
	//change the json to struct
	message := &Service.Player_LoginReq{}
	err := proto.Unmarshal(replay, message)
	if err != nil {
		log.Fatal("unmarshaling error: ", err)
	}
	flag,id := DB.CheckAccount("username", message.UserName, message.PassWord)
	response := Service.Player_LoginResp{}
	if flag == 0 {
		response.Result = Service.Player_LoginResp_Login_Success
	} else if flag == 1 {
		response.Result = Service.Player_LoginResp_Login_Failure
	}else if flag == 2 {
		response.Result = Service.Player_LoginResp_Login_WrongPasswordOrAccount
	}else if flag == 3 {
		response.Result = Service.Player_LoginResp_Login_NonexistentAccount
	}
	response.Id = strconv.Itoa(id)
	data, err := proto.Marshal(&response)
	if err != nil {
		panic(err)
	}
	//return addprotocolnum(1,data)
	return Addprotocolnum(uint16(Service.ModuleId_LoginResp),data)
}


//get player infomation
func GetPlayerInfo(replay []byte)[]byte{
	id := &Service.Player_PlayerInfoReq{}
	err := proto.Unmarshal(replay, id)
	if err != nil {
		log.Fatal("unmarshaling error: ", err)
	}
	nickname, avatar, avataredge := DB.Getplayerinfo(id.Id)
	player := Service.Player_PlayerInfoResp{Nickname:nickname, Avatar:avatar, Avataredge:avataredge}
	data, err := proto.Marshal(&player)
	if err != nil {
		panic(err)
	}

	return Addprotocolnum(uint16(Service.ModuleId_PlayerInfoResp),data)
}


//get theirturn
func GetTheirturn(replay []byte)[]byte{
	id := &Service.Player_PlayerInfoReq{}
	err := proto.Unmarshal(replay, id)
	if err != nil {
		log.Fatal("unmarshaling error: ", err)
	}
	lists := DB.GetTheirturn(id.Id)
	theirturn := Service.Player_TheirturnResp{}
	for i:=0; i< len(lists); i++{
		game := Service.PlayerGameUnderway{}
		//设置id
		s, _ := strconv.ParseInt(lists[i][0].(string), 10, 64)
		game.Id = s
		//设置对手id
		s, _ = strconv.ParseInt(lists[i][1].(string), 10, 64)
		game.OpponentId = s
		game.OpponentNickname = lists[i][2].(string)
		game.OpponentProfile = imageUrl +lists[i][3].(string)
		//设置观看人数
		s, _ = strconv.ParseInt(lists[i][4].(string), 10, 64)
		game.Bystanders = s
		//设置步数
		s, _ = strconv.ParseInt(lists[i][5].(string), 10, 64)
		game.Step = s
		//设置收藏数量
		s, _ = strconv.ParseInt(lists[i][6].(string), 10, 64)
		game.Likenum = s
		game.Laststeptime = lists[i][7].(string)
		//设置类型
		s, _ = strconv.ParseInt(lists[i][8].(string), 10, 64)
		game.Type = s
		//设置评论数
		s, _ = strconv.ParseInt(lists[i][9].(string), 10, 64)
		game.Commentnum = s
		//设置是否在线
		if(lists[i][10].(string) == "1"){
			game.OpponentStatus = true
		} else {
			game.OpponentStatus = false
		}
		temp := append(theirturn.Theirturn, &game)
		theirturn.Theirturn = temp
	}
	data, err := proto.Marshal(&theirturn)
	if err != nil {
		panic(err)
	}
	//return addprotocolnum(7,data)
	return Addprotocolnum(uint16(Service.ModuleId_TheirturnResp),data)
}

//get pending
func GetPending(replay []byte)[]byte{
	id := &Service.Player_PlayerInfoReq{}
	err := proto.Unmarshal(replay, id)
	if err != nil {
		log.Fatal("unmarshaling error: ", err)
	}

	lists := DB.Getpending(id.Id)
	pending := Service.Player_PendingResp{}
	for i:=0; i< len(lists); i++{
		game := Service.PlayerPending{}
		s, _ := strconv.ParseInt(lists[i][0].(string), 10, 64)
		game.OpponentId = s
		game.OpponentNickname = lists[i][1].(string)
		game.OpponentProfile = imageUrl + lists[i][2].(string)
		s, _ = strconv.ParseInt(lists[i][3].(string), 10, 64)
		game.Boardsize = s
		//设置颜色
		if(lists[i][4].(string) == "1"){
			game.Color = true
		} else {
			game.Color = false
		}
		//设置对手状态
		if(lists[i][5].(string) == "1"){
			game.OpponentStatus = true
		} else {
			game.OpponentStatus = false
		}
		//设置段位
		s, _ = strconv.ParseInt(lists[i][6].(string), 10, 64)
		game.OpponentDan = s
		//设置时间
		game.Time = lists[i][7].(string)
		game.Type = lists[i][8].(string)

		temp := append(pending.Pending, &game)
		pending.Pending = temp
	}
	data, err := proto.Marshal(&pending)
	if err != nil {
		panic(err)
	}
	//return addprotocolnum(9,data)
	return Addprotocolnum(uint16(Service.ModuleId_PendingResp),data)
}

//get myturn
func GetMyturn(replay []byte)[]byte{
	id := &Service.Player_PlayerInfoReq{}
	err := proto.Unmarshal(replay, id)
	if err != nil {
		log.Fatal("unmarshaling error: ", err)
	}
	//获得查询结果
	lists := DB.Getmyturn(id.Id)
	//将查询结果封装导Service.Player_MyturnResp对象当中去
	myturn := Service.Player_MyturnResp{}
	//排序一下

	for i:=0; i< len(lists); i++{
		game := Service.PlayerGameUnderway{}
		//设置id
		s, _ := strconv.ParseInt(lists[i][0].(string), 10, 64)
		game.Id = s
		//设置对手id
		s, _ = strconv.ParseInt(lists[i][1].(string), 10, 64)
		game.OpponentId = s
		game.OpponentNickname = lists[i][2].(string)
		game.OpponentProfile = imageUrl + lists[i][3].(string)
		//设置观看人数
		s, _ = strconv.ParseInt(lists[i][4].(string), 10, 64)
		game.Bystanders = s
		//设置步数
		s, _ = strconv.ParseInt(lists[i][5].(string), 10, 64)
		game.Step = s
		//设置收藏数量
		s, _ = strconv.ParseInt(lists[i][6].(string), 10, 64)
		game.Likenum = s
		game.Laststeptime = lists[i][7].(string)
		//设置类型
		s, _ = strconv.ParseInt(lists[i][8].(string), 10, 64)
		game.Type = s
		//设置评论数
		s, _ = strconv.ParseInt(lists[i][9].(string), 10, 64)
		game.Commentnum = s
		//设置是否在线
		if(lists[i][10].(string) == "1"){
			game.OpponentStatus = true
		} else {
			game.OpponentStatus = false
		}


		temp := append(myturn.Myturn, &game)
		myturn.Myturn = temp
	}
	fmt.Println(myturn)
	data, err := proto.Marshal(&myturn)
	if err != nil {
		panic(err)
	}
	//return addprotocolnum(11,data)
	return Addprotocolnum(uint16(Service.ModuleId_MyturnResp),data)
}

