package web

import (
	"net/http"
	"strings"
	"text/template"
)

// HandlerFunc defines the request handler used by gee
type HandlerFunc func(*Context)

// Engine implement the interface of ServeHTTP
type Engine struct {
	*RouterGroup
	trees         map[string]*node   // 每种请求方式对应一颗树
	htmlTemplates *template.Template // for html render
	funcMap       template.FuncMap   // for html render
}

// New is the constructor of gee.Engine
func New() *Engine {
	engine := &Engine{
		trees: make(map[string]*node),
	}
	engine.RouterGroup = &RouterGroup{engine: engine}
	return engine
}

// Run defines the method to start a http server
func (engine *Engine) Run(addr string) (err error) {
	return http.ListenAndServe(addr, engine)
}

func (engine *Engine) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	c := newContext(w, req)
	c.engine = engine
	engine.handleHTTPRequest(c)
}

func (engine *Engine) handleHTTPRequest(c *Context) {
	if _, ok := engine.trees[c.Method]; !ok {
		c.String(http.StatusNotFound, "404 NOT FOUND: %s\n", c.Path)
		return
	}

	// 方法对应的根节点
	root := engine.trees[c.Method]
	n, params := root.getRoute(c.Path)
	if n != nil {
		c.Params = params
		c.handlers = n.handlers
	} else {
		c.handlers = nil
	}
	c.Next()
}

func (engine *Engine) SetFuncMap(funcMap template.FuncMap) {
	engine.funcMap = funcMap
}

func (engine *Engine) LoadHTMLGlob(pattern string) {
	engine.htmlTemplates = template.Must(template.New("").Funcs(engine.funcMap).ParseGlob(pattern))
}

func (engine *Engine) addRoute(method string, pattern string, handlers ...HandlerFunc) {
	assert1(pattern[0] == '/', "path must begin with '/'")
	assert1(method != "", "HTTP method can not be empty")
	assert1(len(handlers) > 0, "there must be at least one handler")

	if _, ok := engine.trees[method]; !ok {
		engine.trees[method] = &node{}
	}
	parts := parsePattern(pattern)
	engine.trees[method].insert(pattern, parts, 0, handlers...)
}

// parsePattern 将 pattern 解析为字符串数组（Only one * is allowed）
func parsePattern(pattern string) []string {
	vs := strings.Split(pattern, "/")

	parts := make([]string, 0)
	for i, item := range vs {
		if item != "" {
			if item[0] == '*' {
				assert1(i == len(vs)-1, "pattern can only have one * and it must be the last part")
			}
			parts = append(parts, item)
		}
	}
	return parts
}

func assert1(guard bool, text string) {
	if !guard {
		panic(text)
	}
}
