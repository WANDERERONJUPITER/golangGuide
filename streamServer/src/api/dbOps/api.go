package dbops

import (
	"database/sql"
	"golangGuide/streamServer/src/api/defs"
	"golangGuide/streamServer/src/api/utils"

	"log"
	"time"
)

func AddUserCredential(loginName, pwd string) error {
	//不要用加号来连接各部分操作语句    1=1
	qs := `insert into users (login_name ,pwd ) values (?,?)`
	stmtIns, err := dbConn.Prepare(qs)
	if err != nil {
		return err
	}

	_, err = stmtIns.Exec(loginName, pwd)
	if err != nil {
		return err
	}
	//注意 defer在栈退出的时候才会调用，对性能会有些许损耗,但考虑到上面出现错误之后下面语句将不再执行，将会存在大量未关闭链接
	defer stmtIns.Close()
	return nil
}

func GetUserCredential(loginName string) (string, error) {
	qs := "select pwd from users where login_name = ?"
	stmtOut, err := dbConn.Prepare(qs)
	if err != nil {
		log.Printf("%s", err)
		return "", err
	}

	// 声明变量接收
	var pwd string

	err = stmtOut.QueryRow(loginName).Scan(&pwd)
	// 这里如果没有 QueryRow(loginName) 没有拿到数据会返回对象 通过后面的scan()  将sql.ErrNoRows 带出
	if err != nil && err != sql.ErrNoRows {
		return "", nil
	}

	defer stmtOut.Close()
	return pwd, nil
}

func GetUser(loginName string) (*defs.User, error) {
	stmtOut, err := dbConn.Prepare("SELECT id, pwd FROM users WHERE login_name = ?")
	if err != nil {
		log.Printf("%s", err)
		return nil, err
	}

	var id int
	var pwd string

	err = stmtOut.QueryRow(loginName).Scan(&id, &pwd)
	if err != nil && err != sql.ErrNoRows {
		return nil, err
	}

	if err == sql.ErrNoRows {
		return nil, nil
	}

	res := &defs.User{Id: id, LoginName: loginName, Pwd: pwd}

	defer stmtOut.Close()

	return res, nil
}

func DeleteUser(loginName, pwd string) error {
	qs := "delete from users where login_name=? and pwd= ?"
	stmtDel, err := dbConn.Prepare(qs)
	if err != nil {
		log.Printf("delete user failed: %s", err)
		return err
	}

	_, err = stmtDel.Exec(loginName, pwd)
	if err != nil {
		return err
	}
	defer stmtDel.Close()
	return nil
}

//  视频部分     TODO   考虑后面重写这一块  分包管理用户 视频 等模块
func AddNewVideo(aid int, name string) (*defs.VideoInfo, error) {
	// create uuid
	vid, err := utils.NewUUID()
	if err != nil {
		return nil, err
	}

	// create time  -->  db 注意这里写库的时间和显示的时间通常会存在误差（时延等）但相对顺序不会改变，而且格式（string  和  time）不同，

	t := time.Now()
	ctime := t.Format("Jan 02 2006, 15:04:05")

	qs := "insert into video_info (id, author_id, name, display_ctime) values (?, ?, ?,  ?)"
	stmtIns, err := dbConn.Prepare(qs)
	if err != nil {
		return nil, err
	}

	_, err = stmtIns.Exec(vid, aid, name, ctime)
	if err != nil {
		return nil, err
	}

	res := &defs.VideoInfo{
		Id:           vid,
		AuthorId:     aid,
		Name:         name,
		DisplayCtime: ctime,
	}
	defer stmtIns.Close()
	return res, nil
}

func GetVideoInfo(vid string) (*defs.VideoInfo, error) {
	qs := "select author_id,name,display_ctime from video_info where id = ?"
	stmtOut, err := dbConn.Prepare(qs)
	if err != nil {
		return nil, err
	}

	var aid int
	var dct, name string
	err = stmtOut.QueryRow(vid).Scan(&aid, &name, &dct)
	if err != nil && err != sql.ErrNoRows {
		return nil, err
	}

	if err == sql.ErrNoRows {
		return nil, nil
	}

	defer stmtOut.Close()

	res := &defs.VideoInfo{Id: vid, AuthorId: aid, Name: name, DisplayCtime: dct}

	return res, nil
}

func ListVideoInfo(uname string, from, to int) ([]*defs.VideoInfo, error) {
	qs := `SELECT video_info.id, video_info.author_id, video_info.name, video_info.display_ctime FROM video_info 
		INNER JOIN users ON video_info.author_id = users.id
		WHERE users.login_name = ? AND video_info.create_time > FROM_UNIXTIME(?) AND video_info.create_time <= FROM_UNIXTIME(?) 
		ORDER BY video_info.create_time DESC`
	stmtOut, err := dbConn.Prepare(qs)
	if err != nil {
		panic(err)
	}

	var res []*defs.VideoInfo

	rows, err := stmtOut.Query(uname, from, to)
	if err != nil {
		log.Printf("%s", err)
		return res, err
	}

	for rows.Next() {
		var id, name, ctime string
		var aid int
		if err := rows.Scan(&id, &aid, &name, &ctime); err != nil {
			return res, err
		}

		vi := &defs.VideoInfo{Id: id, AuthorId: aid, Name: name, DisplayCtime: ctime}
		res = append(res, vi)
	}

	defer stmtOut.Close()

	return res, nil
}

func DeleteVideoInfo(vid string) error {
	qs := "DELETE FROM video_info WHERE id=?"
	stmtDel, err := dbConn.Prepare(qs)
	if err != nil {
		return err
	}

	_, err = stmtDel.Exec(vid)
	if err != nil {
		return err
	}

	defer stmtDel.Close()
	return nil
}

// 留言模块
func AddNewComments(vid string, aid int, content string) error {

	id, err := utils.NewUUID()
	if err != nil {
		return err
	}

	qs := "INSERT INTO comments (id, video_id, author_id, content) values (?, ?, ?, ?)"
	stmtIns, err := dbConn.Prepare(qs)
	if err != nil {
		return err
	}

	_, err = stmtIns.Exec(id, vid, aid, content)
	if err != nil {
		return err
	}

	defer stmtIns.Close()
	return nil
}

func ListComments(vid string, from, to int) ([]*defs.Comment, error) {
	qs := ` SELECT comments.id, users.Login_name, comments.content FROM comments
		INNER JOIN users ON comments.author_id = users.id
		WHERE comments.video_id = ? AND comments.time > FROM_UNIXTIME(?) AND comments.time <= FROM_UNIXTIME(?)
		ORDER BY comments.time DESC`
	stmtOut, err := dbConn.Prepare(qs)

	var res []*defs.Comment

	rows, err := stmtOut.Query(vid, from, to)
	if err != nil {
		return res, err
	}

	for rows.Next() {
		var id, name, content string
		if err := rows.Scan(&id, &name, &content); err != nil {
			return res, err
		}

		c := &defs.Comment{Id: id, VideoId: vid, Author: name, Content: content}
		res = append(res, c)
	}
	defer stmtOut.Close()

	return res, nil
}
