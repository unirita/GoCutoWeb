package pathutil

import (
	"testing"
)

func TestGetRootPath_正常にRootPathが取得できる(t *testing.T) {
	s := GetRootPath()
	if len(s) == 0 {
		t.Error("取得失敗。")
	}
}
