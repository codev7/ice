package ice

import ()

type Middleware func(interface{}, HttpConn, func() interface{}) interface{}

type MiddlewareProvider interface {
	Middlewares() []Middleware
}

func Anonymous(r interface{}, conn HttpConn, next func() interface{}) interface{} {
	if conn.User() == nil || conn.User().CheckRole("?") {
		return next()
	}
	return ForbiddenError("Only unauthorized users are allowed.")
}

func Authenticated(r interface{}, conn HttpConn, next func() interface{}) interface{} {
	if conn.User() != nil && conn.User().CheckRole("*") {
		return next()
	}
	return ForbiddenError("Only authorized users are allowed.")
}
