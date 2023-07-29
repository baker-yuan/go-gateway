package main

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"

	"github.com/xeipuuv/gojsonschema"
)

// https://vue-json-schema-form.lljj.me/
// https://github.com/zyqwst/json-schema-editor-vue
// https://github.com/lljj-x/vue-json-schema-form/tree/master/packages/lib/vue3/vue3-form-element

//
// curl --location --request POST 'http://127.0.0.1:8080/submit' \
// --header 'Content-Type: application/json' \
// --header 'Cookie: uid=1' \
// --data-raw '{
//    "name": "baker",
//    "age": 18
// }'
func main() {
	http.HandleFunc("/submit", submitHandler)
	log.Fatal(http.ListenAndServe(":8080", nil))
}

func submitHandler(w http.ResponseWriter, r *http.Request) {
	// 读取请求体中的JSON数据
	body, err := io.ReadAll(r.Body)
	if err != nil {
		http.Error(w, "Failed to read request body", http.StatusBadRequest)
		return
	}

	// 解析JSON数据到一个map或结构体中
	var data map[string]interface{}
	err = json.Unmarshal(body, &data)
	if err != nil {
		http.Error(w, "Failed to parse JSON data", http.StatusBadRequest)
		return
	}

	// 加载JSON Schema文件
	schemaLoader := gojsonschema.NewStringLoader(`
{
  type: 'object',
  required: [
    'userName',
    'age',
  ],
  properties: {
    userName: {
      type: 'string',
      title: '用户名',
      default: 'Liu.Jun',
    },
    age: {
      type: 'number',
      title: '年龄'
    },
    bio: {
      type: 'string',
      title: '签名',
      minLength: 10,
      default: '知道的越多、就知道的越少',
      'ui:options': {
        placeholder: '请输入你的签名',
        type: 'textarea',
        rows: 1
      }
    }
  }
}
`)

	// 加载要验证的JSON数据
	documentLoader := gojsonschema.NewGoLoader(data)

	// 验证JSON数据是否符合JSON Schema
	result, err := gojsonschema.Validate(schemaLoader, documentLoader)
	if err != nil {
		http.Error(w, "Failed to validate JSON data", http.StatusInternalServerError)
		return
	}

	// 检查验证结果
	if result.Valid() {
		fmt.Fprint(w, "JSON data is valid")
	} else {
		fmt.Fprint(w, "JSON data is not valid")
		for _, err := range result.Errors() {
			fmt.Fprintln(w, "- "+err.String())
		}
	}
}
