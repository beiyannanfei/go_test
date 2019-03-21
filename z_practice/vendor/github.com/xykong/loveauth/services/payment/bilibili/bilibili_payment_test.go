package bilibili

import (
	"encoding/json"
	"fmt"
	"testing"
)

func TestQueryOrderTest(t *testing.T) {
	//resp,err:=QueryBiliOrder("5474762699113183", 396481564)
	//5485075166783503
	//5485853950563071
	//resp,err:=QueryBiliOrder("5485075166783503", 348731357)
	// 问题单号734497494474.23.2019-01-27T18:36:32+08:00.734758094041
	//         734497494474.23.2019-01-27T18:36:32+08:00.734758094041
	//resp,err:=QueryBiliOrder("5485854043773694", 348731357)
	resp, err := QueryBiliOrder("5485853950563071", 348731357)
	data, _ := json.Marshal(resp)
	fmt.Println(string(data), err)
}
