package utils

import "fmt"

// 权限管理操作日志
type OperationLog interface {
	InsertDb(url, param interface{}) error
}

type Instance func() OperationLog

var adapter = make(map[string]Instance)

func Register(name string, log Instance) {
	if _, ok := adapter[name]; ok {
		panic("warn: Register called twice for adapter " + name)
	}
	adapter[name] = log
}

func NewLog(name string) (opt OperationLog, err error) {
	instanceFunc, ok := adapter[name]
	if !ok {
		err = fmt.Errorf("log: unknown adapter name %v (forgot to import?)", name)
		return
	}

	return instanceFunc(), nil
}
