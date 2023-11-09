package rest

import (
	"github.com/aeroideaservices/focus/menu/rest/handlers"
	"github.com/aeroideaservices/focus/menu/rest/services"
	"github.com/gin-gonic/gin"
)

type Router struct {
	menuHandler     *handlers.MenuHandler
	menuItemHandler *handlers.MenuItemHandler
	domainsHandler  *handlers.DomainsHandler
	errorHandler    services.ErrorHandler
}

func NewRouter(
	menuHandler *handlers.MenuHandler,
	menuItemHandler *handlers.MenuItemHandler,
	domainsHandler *handlers.DomainsHandler,
	errorHandler services.ErrorHandler,
) *Router {
	return &Router{
		menuHandler:     menuHandler,
		menuItemHandler: menuItemHandler,
		domainsHandler:  domainsHandler,
		errorHandler:    errorHandler,
	}
}

func (r *Router) SetRoutes(routerGroup *gin.RouterGroup) {
	menus := routerGroup.Group("menus")
	menus.Use(r.errorHandler.Handle) // отлов ошибок

	menus.GET("", r.menuHandler.List)
	menus.POST("", r.menuHandler.Create)

	menu := menus.Group(":" + handlers.MenuIdParam)
	menu.GET("", r.menuHandler.Get)
	menu.PUT("", r.menuHandler.Update)
	menu.DELETE("", r.menuHandler.Delete)

	menuItems := menu.Group("items")
	menuItems.GET("", r.menuItemHandler.List)
	menuItems.POST("", r.menuItemHandler.Create)

	menuItem := menuItems.Group(":" + handlers.MenuItemIdParam)
	menuItem.GET("", r.menuItemHandler.Get)
	menuItem.PUT("", r.menuItemHandler.Update)
	menuItem.DELETE("", r.menuItemHandler.Delete)
	menuItem.POST("move", r.menuItemHandler.Move)

	domain := menus.Group("domains")
	domain.GET("", r.domainsHandler.List)
	domain.POST("", r.domainsHandler.Create)
}
