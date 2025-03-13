package auth

const (
	Read   Type = 1 << iota
	Edit        //编辑
	Delete      //删除
	Create      //创建
	Usage       //使用

	UsageWithRead Type = 17 // 包含了Usage 和 Read
	Admin         Type = 31 //包含了Read、Edit、Delete、Create、Usage

	Unknown Type = 0
)

type Type int
