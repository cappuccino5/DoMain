package utils
// 普通日志(查询)
type General struct {
	Router map[string]string
}

var GeneralUrl = []string{
	"url",
}

// 存数据库操作
func (l *General) InsertDb(url, param interface{}) error {
	// todo 例子，伪代码
	err := model.InsertLog(constantlog{
		UserName: "用户名",
		MainMenu: "url",
		Content:  "请求的参数"，
		CreateTime "修改时间",
		OperatorId:"操作人id",
		Ip:"操作人ip",
	})
	return err
}

func newGeneral() OperationLog {
	return &General{}
}

func init() {
	Register("general", newGeneral)
}
