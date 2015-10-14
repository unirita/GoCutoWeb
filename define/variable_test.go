package define

import (
	"testing"
)

func TestExpandVariables(t *testing.T) {
	msg := "AAAA[$WAPARAM1]BBB[$WAPARAM2]CCC[$WAPARAM3]DDD[$WAPARAM4]"
	result := ExpandVariables(msg, "ZZZ", "YYY", "XXX")
	veryfy := "AAAA[ZZZ]BBB[YYY]CCC[XXX]DDD[$WAPARAM4]"
	if result != veryfy {
		t.Errorf("Wrong result [%v], want to [%v].", result, veryfy)
	}

	result = ExpandVariables(msg, "ZZZ", "YYY", "XXX", "WWW", "VVV")
	veryfy = "AAAA[ZZZ]BBB[YYY]CCC[XXX]DDD[WWW]"
	if result != veryfy {
		t.Errorf("Wrong result [%v], want to [%v].", result, veryfy)
	}
}
