package api

import (
	"net/http"
)

// Handler returns an instance of httprouter.Router that handle APIs registered here
func (rt *_router) Handler() http.Handler {
	// Login
	rt.router.POST("/session/login", rt.wrap(rt.doLogin))
	// User
	rt.router.GET("/users/:username/profile", rt.wrap(rt.getUserProfile))
	rt.router.PUT("/users/:username", rt.wrap(rt.setMyUserName))
	// Profile picture
	rt.router.GET("/users/:username/picture", rt.wrap(rt.getUserPicture))
	rt.router.PUT("/users/:username/picture", rt.wrap(rt.setMyPhoto))
	// Conversation
	rt.router.GET("/users/:username/conversations", rt.wrap(rt.getMyConversations))
	rt.router.GET("/users/:username/conversations/:conversation_id", rt.wrap(rt.getConversation))
	// Message
	rt.router.POST("/users/:username/conversations/:conversation_id/messages", rt.wrap(rt.sendMessage))
	rt.router.POST("/users/:username/conversations/:conversation_id/messages/:message_id/forward", rt.wrap(rt.forwardMessage))
	rt.router.DELETE("/users/:username/conversations/:conversation_id/messages/:message_id", rt.wrap(rt.deleteMessage))	
	// Comment
	rt.router.POST("/users/:username/conversations/:conversation_id/messages/:message_id/comments", rt.wrap(rt.commentMessage))
	rt.router.DELETE("/users/:username/conversations/:conversation_id/messages/:message_id/comments", rt.wrap(rt.uncommentMessage))
	// Group
	rt.router.PUT("/users/:username/groups/:group_id/photo", rt.wrap(rt.setGroupPhoto))
	rt.router.PUT("/users/:username/groups/:group_id/name", rt.wrap(rt.setGroupName))
	rt.router.POST("/users/:username/groups", rt.wrap(rt.createGroup))
	rt.router.POST("/users/:username/groups/:group_id/members", rt.wrap(rt.addToGroup))
	rt.router.DELETE("/users/:username/groups/:group_id/members/:member_username", rt.wrap(rt.leaveGroup))

	// Special routes
	rt.router.GET("/liveness", rt.liveness)

	return rt.router
}
