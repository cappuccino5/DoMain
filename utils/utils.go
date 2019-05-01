package utils

import (
	"crypto/md5"
	"fmt"
	"io"
)
// 登陆密码实现，先使用时间戳生成d5加盐，然后使用密码+盐生成Md5
//md5方法,密码加盐
func Md5V(str string) string {
	w := md5.New()
	io.WriteString(w, str)
	md5str := fmt.Sprintf("%x", w.Sum(nil))
	return md5str
}

