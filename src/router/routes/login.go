package routes

import (
	"go-book-api/src/controllers"
	"net/http"
)

var loginRoutes = Route{
	URI:          "/login",
	Method:       http.MethodPost,
	Function:     controllers.Login,
	AuthRequired: false,
}
