package dbops

import (
	"database/sql"
	"golangGuide/streamServer/src/api/defs"
	"log"
	"strconv"
	"sync"
)

func InsertSession(sid string, ttl int64, userName string) error {
	ttlStr := strconv.FormatInt(ttl, 10)
	qs := "INSERT INTO sessions (session_id, TTL, login_name) VALUES (?, ?, ?)"
	stmtIns, err := dbConn.Prepare(qs)
	if err != nil {
		return err
	}

	_, err = stmtIns.Exec(sid, ttlStr, userName)
	if err != nil {
		return err
	}

	defer stmtIns.Close()
	return nil
}

func RetrieveSession(sid string) (*defs.SimpleSession, error) {
	ss := &defs.SimpleSession{}
	qs := "select ttl,login_name from sessions where session_id =?"
	stmtOut, err := dbConn.Prepare(qs)
	if err != nil {
		return nil, err
	}

	var ttl,userName string

	err = stmtOut.QueryRow(sid).Scan(&ttl, &userName)
	if err != nil && err != sql.ErrNoRows {
		return nil, err
	}

	if res, err := strconv.ParseInt(ttl, 10, 64); err == nil {
		ss.TTL = res
		ss.Username = userName
	} else {
		return nil, err
	}

	defer stmtOut.Close()
	return ss, nil
}

func DeleteSession(sid string) error {
	stmtOut, err := dbConn.Prepare("DELETE FROM sessions WHERE session_id = ?")
	if err != nil {
		log.Printf("%s", err)
		return err
	}

	if _, err := stmtOut.Query(sid); err != nil {
		return err
	}

	return nil
}

func RetrieveAllSessions() (*sync.Map, error) {
	m := &sync.Map{}
	stmtOut, err := dbConn.Prepare("SELECT * FROM sessions")
	if err != nil {
		log.Printf("%s", err)
		return nil, err
	}

	rows, err := stmtOut.Query()
	if err != nil {
		log.Printf("%s", err)
		return nil, err
	}

	for rows.Next() {
		var id,ttlStr,loginName string
		if er := rows.Scan(&id, &ttlStr, &loginName); er != nil {
			log.Printf("retrive sessions error: %s", er)
			break
		}

		if ttl, err1 := strconv.ParseInt(ttlStr, 10, 64); err1 == nil{
			ss := &defs.SimpleSession{
				Username: loginName,
				TTL: ttl,
			}
			m.Store(id, ss)
			log.Printf(" session id: %s, ttl: %d", id, ss.TTL)
		}
	}

	return m, nil
}
