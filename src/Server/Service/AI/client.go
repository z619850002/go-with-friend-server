package AI

import (
	"context"
	"flag"
	"fmt"
	"log"
	"strconv"
	"strings"
	"time"
	"github.com/fatih/color"
	"github.com/google/uuid"
	"google.golang.org/grpc"
	"google.golang.org/genproto/googleapis/rpc/code"
	game "go-with-friend-server/src/Server/Service/AI/proto"
	"go-with-friend-server/src/Server/Service"
	"github.com/golang/protobuf/proto"
	"encoding/binary"
)

type logger struct {
	// Registered clients.
	clients map[string]*aiclient
	// Register requests from the clients.
	register chan *aiclient
	// Unregister requests from clients.
	unregister chan *aiclient
}

func newlogger() *logger {
	return &logger{
		register:   make(chan *aiclient),
		unregister: make(chan *aiclient),
		clients:    make(map[string]*aiclient),
	}
}

func (h *logger) Run() {
	for {
		select {
		case aiclient := <-h.register:
			h.clients[aiclient.id] = aiclient
		case aiclient := <-h.unregister:
			if _, ok := h.clients[aiclient.id]; ok {
				delete(h.clients, aiclient.id)
			}
		}
	}
}

type aiclient struct {
	//玩家的id
	id string
	//playerid
	playerid string
	//注册器
	logger *logger
	// 与ai端连接的client
	aiserver_c  game.GameClient
	// 传输玩家坐标的chan
	 playerposition chan string
}


var (
	address  string
	conn 	 *grpc.ClientConn
	hub      *logger
)

func init(){
	flag.StringVar(&address, "address", "172.26.160.19:50052", "server address")
	flag.Parse()
	//获得grpc连接conn
	var err error
	conn, err = grpc.Dial(address, grpc.WithInsecure())
	if err != nil {
		log.Fatalf("connect error: %v", err)
	}
	//defer conn.Close()
	//新建一个注册器hub管理所有的ai连接
	hub = newlogger()
	go hub.Run()
}

func Startaigame(replay []byte, send chan []byte){
	//首先需要对数据进行解析
	airequest := &Service.Ai_Ainewgamereq{}
	err := proto.Unmarshal(replay, airequest)
	if err != nil {
		log.Fatal("unmarshaling error: ", err)
	}
	//define a channel for ai to recieve position and for player to set their position
	var playerid string
	flag.StringVar(&playerid, "player", uuid.New().String(), "player id")
	//使用conn获得client实例
	ai_client := game.NewGameClient(conn)
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	//新建一个aiclient，存储玩家相关的数据
	aiclient := &aiclient{id:airequest.Id, playerid:playerid, logger:hub,  aiserver_c:ai_client, playerposition:make(chan string)}
	aiclient.logger.register <- aiclient

	player := &game.Player{Id: aiclient.playerid}
	r, err := aiclient.aiserver_c.Login(ctx, player)
	if err != nil {
		log.Fatalf("login error: %v", err)
	}
	if r.Status.Code != int32(code.Code_OK) {
		log.Fatalf("login error: %s", r.Status.Message)
	}
	fmt.Println(color.GreenString("player %s ：login success!", aiclient.playerid))

	if airequest.Flag == true {
		//start a new game
		//first step ： according to the var playercolor: true means white and false means black, playercolor deault black
		if airequest.Color == true {
			r, err := aiclient.aiserver_c.SetWhite(context.Background(), &game.Player{Id: aiclient.playerid})
			if err != nil {
				log.Fatalf("SetWhite error: %v", err)
			}
			fmt.Printf("OpenGo placed a stone at coordinates %d, %d\n", r.X, r.Y)
			//todo : send the posotion to client
			send <- sendaistep(r.X, r.Y)
		}
		for{
			message, ok := <- aiclient.playerposition
			if !ok {
				return
			}
			a := strings.Split(message, ",")
			temp, _ := strconv.ParseInt(a[0], 10, 64)
			coordinateX := int32(temp)
			temp, _ = strconv.ParseInt(a[1], 10, 64)
			coordinateY := int32(temp)
			request := &game.Step{
				X:      coordinateX,
				Y:      coordinateY,
				Player: &game.Player{Id: aiclient.playerid},
			}
			r, err := aiclient.aiserver_c.Play(context.Background(), request)
			if err != nil {
				log.Fatalf("play error: %v", err)
			}
			fmt.Printf("OpenGo placed a stone at coordinates %d, %d\n", r.X, r.Y)
			//todo : send the posotion to client
			send <- sendaistep(r.X, r.Y)
		}
	} else {
		//start a provious game

	}
}
func  sendaistep(x,y int32)[]byte{
	position := Service.Ai_Aistep{}
	position.Position.X = x
	position.Position.Y = y
	data, err := proto.Marshal(&position)
	if err != nil {
		panic(err)
	}
	// add the protocol to the header
	m := make([]byte, 2+len(data))
	// 默认使用大端序
	binary.BigEndian.PutUint16(m, uint16(Service.ModuleId_PendingResp))
	copy(m[2:], data)
	fmt.Println(m)
	//return the response
	return m
}



func Setplayerposition(replay []byte){
	//首先需要对数据进行解析
	airequest := &Service.Ai_Playerstep{}
	err := proto.Unmarshal(replay, airequest)
	if err != nil {
		log.Fatal("unmarshaling error: ", err)
	}
	//解析出坐标，变成x，y的格式（以，分割开来）
	x := turnString(airequest.Position.X)
	y := turnString(airequest.Position.Y)
	str := x+","+y
	//放入chan当中
	aiclient := hub.clients[airequest.Id]
	aiclient.playerposition <- str
}

func turnString(n int32) string {
	buf := [11]byte{}
	pos := len(buf)
	i := int64(n)
	signed := i < 0
	if signed {
		i = -i
	}
	for {
		pos--
		buf[pos], i = '0'+byte(i%10), i/10
		if i == 0 {
			if signed {
				pos--
				buf[pos] = '-'
			}
			return string(buf[pos:])
		}
	}
}