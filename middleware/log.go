package middleware

import (
	"bytes"
	"io/ioutil"
	"log"
	"net/http"
)

func ReqLog(_ http.ResponseWriter, r *http.Request) {
	_ = r.ParseForm()

	log.Println(r.Method, r.URL.Path)

	// 解析 queryString 参数
	var queryString = ""
	for key := range r.Form {
		queryString += `"` + key + `":"` + r.Form.Get(key) + `",`
	}
	if len(queryString) > 0 {
		queryString = string([]byte(queryString)[:len(queryString)-1])
		log.Println(`{` + queryString + `}`)
	}

	// 解析body参数
	buff, _ := ioutil.ReadAll(r.Body)
	_ = r.Body.Close()                               // 必须关闭后在给body 赋值
	r.Body = ioutil.NopCloser(bytes.NewBuffer(buff)) // 不赋值后面使用会报错
	bodyString := string(buff)
	if len(bodyString) > 0 {
		log.Println(bodyString)
	}
}
