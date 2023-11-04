package gopiper

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"testing"

	"github.com/admpub/gohttp"
)

func TestAdvance(t *testing.T) {
	lv1URL := `https://www.autohome.com.cn/grade/carhtml/A.html`
	req := gohttp.New()
	resp, _ := req.Get(lv1URL).End()
	defer resp.Body.Close()
	body, _ := ioutil.ReadAll(resp.Body)

	pipe := PipeItem{}

	err := json.Unmarshal([]byte(`{
		"selector": "dl",
		"type": "array",
		"name": "",
		"filter": "",
		"subitem": [
			{
				"name": "brand",
				"selector": "dt div a",
				"type": "map",
				"filter": "string"
			}
		]
	}`), &pipe)

	if err != nil {
		panic(err)
	}

	fmt.Println(string(body))

	v, _ := pipe.PipeBytes(body, "html")
	fmt.Println(`=== [ lv1URL ] ===================================\`)
	showjson(v)
	fmt.Println(`=== [ lv1URL ] ===================================/`)
	fmt.Println()
	// panic(`p`)
}
