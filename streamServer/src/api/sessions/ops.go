package sessions

import (
	dbops "golangGuide/streamServer/src/api/dbOps"
	"golangGuide/streamServer/src/api/defs"
	"golangGuide/streamServer/src/api/utils"
	"sync"
	"time"
)

// 考虑当我们增加一项服务，如redis来解决此处问题，势必会增加系统的复杂度，在收益和支持之间平衡
// sync.Map  是在1.9之后才加入的，支持并发的读写，读的时候1w并发和10w并发趋于常数，没有太大的区别，但在写的时候，每次都要加一个全局锁，是非常耗时的
var sessionMap *sync.Map

func init() {
	sessionMap = &sync.Map{}
}

func LoadSessionsFromDB() {
	r, err := dbops.RetrieveAllSessions()
	if err!=nil {
		return
	}

	//将读取到的数据存储在cache
	r.Range(func(key, value interface{}) bool {
		ss := value.(*defs.SimpleSession)
		sessionMap.Store(key, ss)
		return true
	})

}

func GenerateNewSessionId(userName string) string {
	id, _ := utils.NewUUID()
	ct := nowInMilli()
	ttl := ct + 30 * 60 * 1000// Severs session valid time: 30 min

	ss := &defs.SimpleSession{Username: userName, TTL: ttl}
	sessionMap.Store(id, ss)
	dbops.InsertSession(id, ttl, userName)

	return id
}

func IsSessionExpired(sid string) (string, bool) {
	ss, ok := sessionMap.Load(sid)
	ct := nowInMilli()
	if ok {
		if ss.(*defs.SimpleSession).TTL < ct {
			deleteExpiredSession(sid)
			return "", true
		}

		return ss.(*defs.SimpleSession).Username, false
	} else {//  保证数据逻辑一致
		ss, err := dbops.RetrieveSession(sid)
		if err != nil || ss == nil {
			return "", true
		}

		if ss.TTL < ct {
			deleteExpiredSession(sid)
			return "", true
		}

		sessionMap.Store(sid, ss)
		return ss.Username, false
	}
	return "", true
}

func nowInMilli() int64{
	return time.Now().UnixNano()/1000000
}

func deleteExpiredSession(sid string) {
	sessionMap.Delete(sid)
	dbops.DeleteSession(sid)
}