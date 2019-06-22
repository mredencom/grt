package qhttp

import "net/http"

// 创建一个请求结构
type Request struct {
	*http.Request 	 					 // 内置的http请求方法
	parsedGet     bool                   // GET参数是否已经解析
	parsedPost    bool                   // POST参数是否已经解析
	queryVars     map[string][]string    // GET参数
	routerVars    map[string][]string    // 路由解析参数
	exit          bool                   // 是否退出当前请求流程执行
	Id            int                    // 请求id(唯一)
	//Server        *Server                // 请求关联的服务器对象
	//Cookie        *Cookie                // 与当前请求绑定的Cookie对象(并发安全)
	//Session       *Session               // 与当前请求绑定的Session对象(并发安全)
	//Response      *Response              // 对应请求的返回数据操作对象
	//Router        *Router                // 匹配到的路由对象
	EnterTime     int64                  // 请求进入时间(微秒)
	LeaveTime     int64                  // 请求完成时间(微秒)
	params        map[string]interface{} // 开发者自定义参数(请求流程中有效)
	parsedHost    string                 // 解析过后不带端口号的服务器域名名称
	clientIp      string                 // 解析过后的客户端IP地址
	rawContent    []byte                 // 客户端提交的原始参数
	isFileRequest bool                   // 是否为静态文件请求(非服务请求，当静态文件存在时，优先级会被服务请求高，被识别为文件请求)
}