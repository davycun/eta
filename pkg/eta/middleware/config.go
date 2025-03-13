package middleware

import "github.com/gin-gonic/gin"

type MidOption struct {
	Order       int    //middleware的排序号
	Name        string //名字，相同名字后面注册的会覆盖之前
	HandlerFunc gin.HandlerFunc
}
