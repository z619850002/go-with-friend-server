package Image

import (
	"encoding/binary"
	"fmt"
	"github.com/golang/protobuf/proto"
	"log"
	"strconv"
	"github.com/astaxie/beego/orm"
	"github.com/z619850002/go-with-friend-server/src/Server/Service"
	"github.com/z619850002/go-with-friend-server/src/Server/DB/datautil"
)

var imageUrl = "http://172.26.163.124:10086/"


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


//default image id
//8ECCBBB175CB28A2
//DA1C9761BEE03F03

func DownloadImg(replay []byte) []byte{
	//change the json to struct
	message := &Service.Player_AvatarReq{}
	err := proto.Unmarshal(replay, message)
	if err != nil {
		log.Fatal("unmarshaling error: ", err)
	}
	playerId , _ := strconv.Atoi(message.Id)
	avatarId := datautil.GetPlayerById(playerId , orm.NewOrm()).Avatar


	response := Service.Player_AvatarResp{Url:string(imageUrl + avatarId)}

	//if flag == 0 {
	//	response.Result = Service.Player_LoginResp_Login_Success
	//} else if flag == 1 {
	//	response.Result = Service.Player_LoginResp_Login_Failure
	//}else if flag == 2 {
	//	response.Result = Service.Player_LoginResp_Login_WrongPasswordOrAccount
	//}else if flag == 3 {
	//	response.Result = Service.Player_LoginResp_Login_NonexistentAccount
	//}
	//response.Id = strconv.Itoa(id)
	data, err := proto.Marshal(&response)
	if err != nil {
		panic(err)
	}
	return addprotocolnum(uint16(Service.ModuleId_AvatarResp),data)
}