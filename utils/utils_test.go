package utils

import "testing"

func TestLog(t *testing.T) {
	res := GeneralLog("url")
	t.Log(res)
}
