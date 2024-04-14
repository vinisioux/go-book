package routes

import (
	"go-book-api/src/controllers"
	"net/http"
)

var postsRoutes = []Route{
	{
		URI:          "/posts",
		Method:       http.MethodPost,
		Function:     controllers.CreatePost,
		AuthRequired: true,
	},
	{
		URI:          "/posts",
		Method:       http.MethodGet,
		Function:     controllers.GetPosts,
		AuthRequired: true,
	},
	{
		URI:          "/posts/{postId}",
		Method:       http.MethodGet,
		Function:     controllers.GetPost,
		AuthRequired: true,
	},
	{
		URI:          "/posts/{postId}",
		Method:       http.MethodPut,
		Function:     controllers.UpdatePost,
		AuthRequired: true,
	},
	{
		URI:          "/posts/{postId}",
		Method:       http.MethodDelete,
		Function:     controllers.DeletePost,
		AuthRequired: true,
	},
	{
		URI:          "/users/{userId}/posts",
		Method:       http.MethodGet,
		Function:     controllers.GetUserPosts,
		AuthRequired: true,
	},
	{
		URI:          "/posts/{postId}/like",
		Method:       http.MethodPost,
		Function:     controllers.LikePost,
		AuthRequired: true,
	},
	{
		URI:          "/posts/{postId}/unlike",
		Method:       http.MethodPost,
		Function:     controllers.UnlikePost,
		AuthRequired: true,
	},
}
