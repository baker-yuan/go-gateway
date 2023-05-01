package httprouter

import (
	"fmt"
	"net/http"
	"testing"
)

// curl http://127.0.0.1:8080/articles
// curl http://127.0.0.1:8080/articles/1
// curl http://127.0.0.1:8080/users
// curl http://127.0.0.1:8080/users/1
// curl http://127.0.0.1:8080/product/aAdd
// curl http://127.0.0.1:8080/product/bAdd
func TestHttp(t *testing.T) {
	var router *RouterRepo = New()

	// ------------------------- 初始化 start -------------------------
	// 博客系统
	articleSrv := "baker-blog-article"
	router.GET("/articles", &ServiceInfo{ID: 1, Application: articleSrv})        // 获取所有文章
	router.GET("/articles/:id", &ServiceInfo{ID: 2, Application: articleSrv})    // 获取单个文章
	router.POST("/articles", &ServiceInfo{ID: 3, Application: articleSrv})       // 创建文章
	router.PUT("/articles/:id", &ServiceInfo{ID: 4, Application: articleSrv})    // 更新文章
	router.DELETE("/articles/:id", &ServiceInfo{ID: 5, Application: articleSrv}) // 删除文章

	// 用户服务
	userSrv := "baker-blog-user"
	router.GET("/users", &ServiceInfo{ID: 10, Application: userSrv})             // 获取所有用户信息
	router.GET("/users/:user_id", &ServiceInfo{ID: 11, Application: userSrv})    // 获取用户信息
	router.POST("/users", &ServiceInfo{ID: 12, Application: userSrv})            // 创建用户
	router.PUT("/users/:user_id", &ServiceInfo{ID: 13, Application: userSrv})    // 更新用户信息
	router.DELETE("/users/:user_id", &ServiceInfo{ID: 14, Application: userSrv}) // 删除用户

	// 商品服务
	productSrv := "baker-blog-product"
	router.GET("/product/*Add", &ServiceInfo{ID: 10, Application: productSrv}) // 获取所有用户信息

	// ------------------------- 初始化 end -------------------------

	// 1、绑定回调
	http.HandleFunc("/", func(rsp http.ResponseWriter, req *http.Request) {
		path := req.URL.Path
		method := req.Method
		service, params, redirect := router.Lookup(method, path)
		fmt.Println(service, params, redirect)
	})

	fmt.Println("listen server on http://localhost:8080")

	// 2、注册服务
	http.ListenAndServe(":8080", nil)
}
