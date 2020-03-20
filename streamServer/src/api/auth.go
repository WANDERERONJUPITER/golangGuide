package main

import (
	"golangGuide/streamServer/src/api/defs"
	"golangGuide/streamServer/src/api/sessions"
	"net/http"
)

//X 开头构成了自定义的header
var HEADER_FIELD_SESSION = "X-Session-Id"
var HEADER_FIELD_UNAME = "X-User-Name"

// Check if the current user has the permission
// Use session id to do the check
func validateUserSession(r *http.Request) bool {
	sid := r.Header.Get(HEADER_FIELD_SESSION)
	if len(sid) == 0 {
		return false
	}

	uname, ok := sessions.IsSessionExpired(sid)
	if ok {
		return false
	}

	r.Header.Add(HEADER_FIELD_UNAME, uname)
	return true
}

func ValidateUser(w http.ResponseWriter, r *http.Request) bool {
	/*TODO   IAM模块  云计算平台，电商平台用来管理用户权限，区分普通用户和后台用户  SSO, Rbac, role based access control  */
	uname := r.Header.Get(HEADER_FIELD_UNAME)
	if len(uname) == 0 {
		sendErrorResponse(w, defs.ErrorNotAuthUser)
		return false
	}

	return true
}
