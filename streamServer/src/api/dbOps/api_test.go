package dbops

//  TODO 每次提交代码之前先跑 test
import (
	"fmt"
	"strconv"
	"testing"
	"time"
)

// 保证数据之间不会循环依赖和被干扰  每次我们创建数据完成测试后直接将数据清除
// init (DB login, truncate tables) --> run test --> clear data (truncate tables)

var tempvid string

func clearTables() {
	dbConn.Exec("truncate users")
	dbConn.Exec("truncate video_info")
	dbConn.Exec("truncate comments")
	dbConn.Exec("truncate sessions")
}

func TestMain(m *testing.M) {
	clearTables()
	m.Run()
	clearTables()
}

//为保证顺序，我们使用subTest
func TestUserWorkFlow(t *testing.T) {
	t.Run("Add", testAddUser)
	t.Run("Get", testGetUser)
	t.Run("Delete", testDeleteUser)
	t.Run("Reget", testReGetUser)
}

func testAddUser(t *testing.T) {
	//TODO  这里使用mock data
	err := AddUserCredential("aaron", "anderson")
	if err != nil {
		t.Errorf("error of addUser: %v", err)
	}
}

func testGetUser(t *testing.T) {
	pwd, err := GetUserCredential("aaron")
	if err != nil || pwd != "anderson" {
		t.Errorf("error of getUser: %v", err)
	}
}

func testDeleteUser(t *testing.T) {
	err := DeleteUser("aaron", "anderson")
	if err != nil {
		t.Errorf("error of deleteUser: %v", err)
	}
}

//TODO  考虑这里实施的必要性，是否说测试删除成功之后我们根本没有必要再获取一次
func testReGetUser(t *testing.T) {
	pwd, err := GetUserCredential("aaron")
	if err != nil {
		t.Errorf("error of reGetUser: %v", err)
	}

	//TODO  注意这里删除后，查询不存在的键返回的数据是""
	if pwd != "" {
		t.Errorf("delete user test failed: %s", err)
	}
}

// ---------------- video test  begin--------------
func TestVideoWorkFlow(t *testing.T) {
	clearTables()
	//清除数据后 要先有个用户
	t.Run("PrepareUser", testAddUser)
	t.Run("AddVideo", testAddVideoInfo)
	t.Run("GetVideo", testGetVideoInfo)
	t.Run("DelVideo", testDeleteVideoInfo)
	t.Run("RegetVideo", testRegetVideoInfo)
}

func testAddVideoInfo(t *testing.T) {
	vi, err := AddNewVideo(1, "my-video")
	if err != nil {
		t.Errorf("Error of AddVideoInfo: %v", err)
	}
	tempvid = vi.Id
}

func testGetVideoInfo(t *testing.T) {
	_, err := GetVideoInfo(tempvid)
	if err != nil {
		t.Errorf("Error of GetVideoInfo: %v", err)
	}
}

func testDeleteVideoInfo(t *testing.T) {
	err := DeleteVideoInfo(tempvid)
	if err != nil {
		t.Errorf("Error of DeleteVideoInfo: %v", err)
	}
}

func testRegetVideoInfo(t *testing.T) {
	vi, err := GetVideoInfo(tempvid)
	if err != nil || vi != nil {
		t.Errorf("Error of RegetVideoInfo: %v", err)
	}
}

//comments
func TestComments(t *testing.T) {
	clearTables()
	t.Run("AddUser", testAddUser)
	t.Run("AddCommnets", testAddComments)
	t.Run("ListComments", testListComments)
}

func testAddComments(t *testing.T) {
	vid := "12345"
	aid := 1
	content := "I like this video"

	err := AddNewComments(vid, aid, content)

	if err != nil {
		t.Errorf("Error of AddComments: %v", err)
	}
}

func testListComments(t *testing.T) {
	vid := "12345"
	from := 1514764800
	to, _ := strconv.Atoi(strconv.FormatInt(time.Now().UnixNano()/1000000000, 10))

	res, err := ListComments(vid, from, to)
	if err != nil {
		t.Errorf("Error of ListComments: %v", err)
	}

	for i, ele := range res {
		fmt.Printf("comment: %d, %v \n", i, ele)
	}
}
