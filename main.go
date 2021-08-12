package main

import (
	"file/controller"
	"file/lib/config"
	"file/middleware"
	"fmt"
	"github.com/julienschmidt/httprouter"
	"log"
	"net/http"
)

// 定义中间件结构体
type middleWareHandler struct {
	r *httprouter.Router
}

// NewMiddleWareHandler 工厂函数，创建中间件 实例
func NewMiddleWareHandler(r *httprouter.Router) http.Handler {
	m := middleWareHandler{}
	m.r = r
	return m
}

// 实现go语言接口，这里就可以拦截到每个请求了
func (m middleWareHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	middleware.ReqLog(w, r)
	m.r.ServeHTTP(w, r)
}

// RegisterHandlers 路由配置
func RegisterHandlers() *httprouter.Router {
	router := httprouter.New()

	// 处理静态文件
	// http://localhost:8001/file-assets/libs/ajax_upload.js
	router.ServeFiles(config.BaseAssets+"/*filepath", http.Dir("./resources/assets"))

	// 异常捕获
	router.PanicHandler = func(w http.ResponseWriter, r *http.Request, v interface{}) {
		w.WriteHeader(http.StatusInternalServerError)
		_, _ = fmt.Fprintf(w, "error:%s", v)
	}

	// 捕获访问的页面不存在时，显示的页面
	router.NotFound = http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// core.WriteError(w, "PAGE_NOT_FOUND", "Requested resource could not be found")
		_, _ = fmt.Fprint(w, "暂无路由!\n")
	})

	controller.NewFileController(router)
	controller.NewEnumController(router)
	return router
}

func main() {
	// 创建并注册路由
	r := RegisterHandlers()
	// 路由通过中间件返回 Handler 方法
	mh := NewMiddleWareHandler(r)

	addr := fmt.Sprintf(":%d", config.AppPort)
	log.Println(fmt.Sprintf("http://localhost%s", addr))
	log.Fatal(http.ListenAndServe(addr, mh))
}
