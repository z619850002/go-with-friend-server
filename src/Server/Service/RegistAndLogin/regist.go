package RegistAndLogin

import (
	"log"
	"github.com/golang/protobuf/proto"
	"encoding/binary"
	"fmt"
	"github.com/z619850002/go-with-friend-server/src/Server/Service"
	"github.com/z619850002/go-with-friend-server/src/Server/DB"
)

func Regist(replay []byte) []byte{
	//change the json to struct
	message := &Service.Player_RegisterReq{}
	err := proto.Unmarshal(replay, message)
	if err != nil {
		fmt.Println("regist fail")
		registReponse(1)
		log.Fatal("unmarshaling error: ", err)
	}
	//return registReponse( conn, DB.Newplayer("username", message.UserName, message.PassWord))
	return registReponse( DB.CreatePlayer("username", message.UserName, message.PassWord))
}

//respond to client
func registReponse(flag int) []byte{
	//create the response
	response :=  Service.Player_RegisterResp{}
	if flag == 0 {
		response.Result = Service.Player_RegisterResp_Register_Success
	} else if flag ==1 {
		response.Result = Service.Player_RegisterResp_Register_Failure
	} else if flag == 2{
		response.Result = Service.Player_RegisterResp_Register_Duplicate
	}
	data, err := proto.Marshal(&response)
	if err != nil {
		panic(err)
	}
	// add the protocol to the header
	m := make([]byte, 2+len(data))
	// 默认使用大端序
	binary.BigEndian.PutUint16(m, uint16(3))
	copy(m[2:], data)
	fmt.Println(m)

	//return the response
	return m
}


