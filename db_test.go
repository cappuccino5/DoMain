package mysql

import (
	"github.com/go-sql-driver/mysql"
	"github.com/jinzhu/gorm"
	"github.com/sirupsen/logrus"
	"go-hisens/chat/model"
	"testing"
	"time"
)

func TestConnectDb(t *testing.T) {
	type DbLogger struct {
		gorm.Logger
	}
	qdb, err := gorm.Open("mysql", "root:zero_lee@tcp(127.0.0.1:3306)/hisens?charset=utf8&parseTime=True&loc=Local")
	if err != nil {
		panic(err)
	}
	qdb.SingularTable(true)
	qdb.SetLogger(DbLogger{})
	mysql.SetLogger(DbLogger{})
	qdb.LogMode(true) // 显示sql语句
	if err = qdb.DB().Ping(); err != nil {
		panic(err)
	}
	mdb := connectDb
	if qdb == nil || mdb == nil {
		t.Fatal("connect db error")
	}

	t.Log("connect db success")
}

func testGormUpdate(t *testing.T) {
	mdb := connectDb()
	mdb.Debug().Table(`user`).Where(`id=?`, 50008).Updates(map[string]interface{}{
		"user_name": "lisi",
		"city":      "shenzhen",
	})
}

// 数据库
type DbTestLogger struct {
	gorm.Logger
}

func (l DbTestLogger) Print(values ...interface{}) {
	logrus.Info(values...)
}
func connectDb() *gorm.DB {

	/*
	用户账号：root
	用户密码：zero_lee
	db_Name：hisens
	服务器地址和端口：127.0.0.1:3306
	*/
	qdb, err := gorm.Open("mysql", "root:zero_lee@tcp(127.0.0.1:3306)/hisens?charset=utf8&parseTime=True&loc=Local")
	if err != nil {
		panic(err)
	}
	qdb.SingularTable(true)
	qdb.SetLogger(DbTestLogger{})
	mysql.SetLogger(DbTestLogger{})
	qdb.LogMode(true) // 显示sql语句
	if err = qdb.DB().Ping(); err != nil {
		panic(err)
	}
	return qdb
}

/* go test -bench=MessageInset db_test.go -v
 500           2822265 ns/op
 500           2990234 ns/op
 500           2748047 ns/op
 500           2863281 ns/op
 500           2693359 ns/op
 500           2921875 ns/op
 500           2697265 ns/op
 500           2634765 ns/op
 500           2861328 ns/op
 500           2707031 ns/op
*/
func BenchmarkMessageInset(b *testing.B) {
	mdb := connectDb()
	num := uint64(1000)
	nowTime := uint64(time.Now().Unix())
	for i := 0; i < b.N; i++ {
		num++
		offlineMsg := model.OfflineMsg{
			From:       num,
			To:         num,
			Type:       1,
			CreateTime: nowTime,
			UpdateTime: nowTime,
			Event:      1,
			Url:        "http://",
			Text:       "this is test messages",
			Title:      "messages info ",
			Intro:      "bak",
			ThumbUrl:   "http://xx",
			Accessory:  1,
		}
		if err := mdb.Create(&offlineMsg).Error; err != nil {
			b.Error(err)
			return
		}
		b.Log(i)
	}
	b.Log("test message success")
}

/* go test -bench=MessageQuery db_test.go -v
 200           9677734 ns/op
 200           9914546 ns/op
 200           9931640 ns/op
 100          10034174 ns/op
 200           9775390 ns/op
 100          10087885 ns/op
 200           9848633 ns/op
 100          10209956 ns/op
 200           9731445 ns/op
 100          10029297 ns/op

*/
func BenchmarkMessageQuery(b *testing.B) {
	mdb := connectDb()
	num := uint64(1000)
	var offlineMsg []model.OfflineMsg
	for i := 0; i < b.N; i++ {
		num++
		if err := mdb.Model("offline_msg").Where("offline_msg.from = ?", 1001).Find(&offlineMsg).Error; err != nil {
			b.Error(err)
			return
		}
		b.Log(i)
	}
	b.Log("test message success")
}
