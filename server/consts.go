package main

const (
	NeedRegister      = "need register"        //需要先注册
	UserNameRepeated  = "user name repeated"   //玩家名称重复
	RoomNameRepeated  = "room name repeated"   //聊天室名称重复
	AlreadyRegistered = "already registered"   //已经完成注册
	UnknownOrder      = "unknown order"        //未知命令
	RoomNotfound      = "room not found"       //找不到房间
	ParamError        = "param error"          //参数错误
	ParamMiss         = "param miss"           //缺少必要参数
	NeedJoinRoom      = "please join one room" //聊天前要先加入一个房间

	Help       = "/help"       //帮助
	Register   = "/register"   //注册一个用户
	Rooms      = "/rooms"      //显示所有的聊天室
	CreateRoom = "/createRoom" //创建一个聊天室
	LeaveRoom  = "/leaveRoom"  //离开聊天室
	JoinRoom   = "/joinRoom"   //加入聊天室
	CloseRoom  = "/closeRoom"  //关闭聊天室（聊天室创建者）
	Users      = "/users"      //所有的玩家姓名
	Popular    = "/popular"    //最近时间段内，出现最多次的词语
	Stats      = "/stats"      //查看某个玩家的在线时长
)
