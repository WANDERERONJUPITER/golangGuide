package dbops

import (
	"log"
	_"github.com/go-sql-driver/mysql"
)

/*
ap-->video id --> mysql
dispatcher --> mysql -->video id --> datachannel
executor --> datachannel --> video id -->delete videos
*/

func ReadVideoDeletionRecord(count int) ([]string,error)  {
	stmtOut, err := dbConn.Prepare("SELECT video_id FROM video_del_rec LIMIT ?")

	var ids []string

	if err != nil {
		return ids, err
	}

	rows, err := stmtOut.Query(count)
	if err != nil {
		log.Printf("Query VideoDeletionRecord error: %v", err)
		return ids, err
	}

	for rows.Next() {
		var id string
		if err := rows.Scan(&id); err != nil {
			return ids, err
		}

		ids = append(ids, id)
	}

	defer stmtOut.Close()
	return ids, nil
}

func DelVideoDeletionRecord(vid string) error {
	stmtDel, err := dbConn.Prepare("DELETE FROM video_del_rec WHERE video_id=?")
	if err != nil {
		return err
	}

	_, err = stmtDel.Exec(vid)
	if err != nil {
		//TODO  log 做切割，压缩，分级
		log.Printf("Deleting VideoDeletionRecord error: %v", err)
		return err
	}

	defer stmtDel.Close()
	return nil
}