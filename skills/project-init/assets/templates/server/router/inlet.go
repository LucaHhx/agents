package router

import (
	"github.com/gin-gonic/gin"
	// HZ:ROUTER:PACKAGE_IMPORTS
)

var (
	V1RouterGroupApp = new(V1RouterGroup)
)

type V1RouterGroup struct {
	// HZ:ROUTER:PACKAGE_FIELDS
}

func (routerGroup V1RouterGroup) InitRouter(Router *gin.RouterGroup) {
	r := Router.Group("v1")
	_ = r
	// HZ:ROUTER:PACKAGE_INIT
}
