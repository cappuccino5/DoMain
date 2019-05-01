package utils

import (
	"crypto/md5"
	"fmt"
	"io"
	"math/rand"
	"regexp"
	"strconv"
	"time"
)

// 登陆密码实现，先使用时间戳生成d5加盐，然后使用密码+盐生成Md5
//md5方法,密码加盐
func Md5V(str string) string {
	w := md5.New()
	io.WriteString(w, str)
	md5str := fmt.Sprintf("%x", w.Sum(nil))
	return md5str
}

//生成盐
func CreateSalt() string {
	now := time.Now().Unix()
	now_str := strconv.FormatInt(now, 10)
	time_str := Md5V(now_str)
	rand_int := RandInt64(100, 999)
	rand_str := strconv.FormatInt(rand_int, 10)
	return time_str + rand_str
}

//随机数(100,999)
func RandInt64(min, max int64) int64 {
	if min >= max || min == 0 || max == 0 {
		return max
	}
	return rand.Int63n(max-min) + min
}

//正则匹配,匹配规则可以一直加 例子:RegexpMatch("email","nimin@qq.com" )
func RegexpMatch(pattern_type string, source string) bool {
	pattern_list := map[string]string{}
	pattern_list["ip"] = "(2(5[0-5]{1}|[0-4]\\d{1})|[0-1]?\\d{1,2})(\\.(2(5[0-5]{1}|[0-4]\\d{1})|[0-1]?\\d{1,2})){3}"
	pattern_list["email"] = "^[a-zA-Z0-9_.-]+@[a-zA-Z0-9-]+(\\.[a-zA-Z0-9-]+)*\\.[a-zA-Z0-9]{2,6}$"
	pattern_list["qq"] = "^[1-9]\\d{4,10}$"
	pattern := pattern_list[pattern_type]
	reg := regexp.MustCompile(pattern)
	if res := reg.FindAllString(source, -1); res == nil {
		return false
	} else {
		return true
	}
}

// 测试数组
var general = []string{
	"url",
}
// 普通日志和警告日志存在一个map里面
var operation = map[string]string{
	"url1": "general",
	"url2": "warn",
}
// 日志等级划分
func GeneralLog(url string) error {
	name, ok := operation[url]
	if !ok {
		return nil
	}
	fmt.Println(name)
	//l, err := lo.NewLog(name)
	//if err != nil {
	//	return err
	//}
	//l.InsertDb(v, param)
	return nil
}
