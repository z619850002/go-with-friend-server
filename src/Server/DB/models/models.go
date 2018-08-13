package models

import "time"
//construct all models here
//then we can use the orm to change these structs to tables in database

//Entities

type Account struct {
	Id int
	//with player
	Player *Player `orm:"rel(fk)"`
	/*-----------------------------attribute----------------------------------*/
	IdentityType string `orm:"size(50)"`
	Identifier   string `orm:"size(100)"`
	Credential   string `orm:"size(100)"`
}


type Level struct {
	Id int

	//with player
	Player *Player	`orm:"rel(one)"`
	/*-----------------------------attribute----------------------------------*/
	Experience  int64
	FallNumber	int64
	TakeNumber	int64
	PracticeInfo	string	`orm:"size(500)"`
	//fight with friends
	FriendsBattle	int
	//smart battle
	SmartBattle		int
	Circumstance 	int
}

type Player struct {
	Id int
	//with Account
	Account []*Account `orm:"reverse(many)"`
	//with Level
	Level *Level	`orm:"reverse(one)"`
	//with History
	CollectedHistory []*History `orm:"rel(m2m); rel_table(collected_history)"`
	PlayerHistory    []*History `orm:"reverse(many)"`
	//with GoodChoice
	GoodChoice []*GoodChoice `orm:"rel(m2m); rel_table(praised_goodchoice)"`
	//with Message
	Messages []*Message `orm:"reverse(many)"`
	/*-----------------------------attribute----------------------------------*/
	NickName   string `orm:"size(50)"`
	Avatar     string
	AvatarEdge string
	Isonline  bool			//是否在线
	dan 	int				//段位
}

type History struct {
	Id int
	//with Player
	CollectedByPlayers []*Player `orm:"reverse(many)"`
	Player1            *Player   `orm:"rel(fk)"`
	Player2            *Player   `orm:"rel(fk)"`
	//Players		[]*Player		`orm:"reverse(many)"`
	//with message
	Messages []*Message `orm:"reverse(many)"`
	//with GoodChoice
	GoodChoices []*GoodChoice `orm:"reverse(many)"`
	/*-----------------------------attribute----------------------------------*/
	Url        string `orm:"size(100)"`
	Commentsnumber	int64
	Bystanders int64
	Like       int64
	Step       int64
	Winner     int64
	//size of the chess board
	Size int64
	Type int64
	Time time.Time
	ChessBoard   string `orm:"size(500)"`
	Turn         bool
	LastTurnTime time.Time
	Player1Score	int
	Player2Score	int
}

type GoodChoice struct {
	Id int
	//with Player
	Players []*Player `orm:"reverse(many)"`
	//with History
	Histories *History `orm:"rel(fk)"`
	/*-----------------------------attribute----------------------------------*/
	Pos          int
	PraiseNumber int
}


type Message struct {
	Id int
	//with battle
	History *History `orm:"rel(fk)"`
	//with player
	Player *Player `orm:"rel(fk)"`
	/*-----------------------------attribute----------------------------------*/
	Data string `orm:"size(100)"`
}


type Invitation struct {
	Id      int
	Inviter *Player  `orm:"rel(fk)"`
	Invitee *Player  `orm:"rel(fk)"`
	Turn    bool
	Time    time.Time
	Size    int //棋盘大小
	Firststep string
	Type string		//邀请方式
}
