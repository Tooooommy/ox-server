package server

import (
	"html/template"
	"log"
	"net/http"
	"path"
	"strings"
)

type HandlerFunc func(*Context)

type Application struct {
	*RouterGroup
	router    *router
	groups    []*RouterGroup
	templates *template.Template
	funcMap   template.FuncMap
}

func New() *Application {
	app := &Application{router: newRouter()}
	app.RouterGroup = &RouterGroup{app: app}
	app.groups = []*RouterGroup{app.RouterGroup}
	return app
}

func (app *Application) addRoute(method string, pattern string, handler HandlerFunc) {
	app.router.addRoute(method, pattern, handler)
}

// GET method
func (app *Application) GET(pattern string, handler HandlerFunc) {
	app.addRoute(http.MethodGet, pattern, handler)
}

// POST method
func (app *Application) POST(pattern string, handler HandlerFunc) {
	app.addRoute(http.MethodPost, pattern, handler)
}

// run on port
func (app *Application) Run(addr string) error {
	return http.ListenAndServe(addr, app)
}

// implement http
func (app *Application) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	var middlewares []HandlerFunc
	for _, group := range app.groups {
		if strings.HasPrefix(req.URL.Path, group.prefix) {
			middlewares = append(middlewares, group.middlewares...)
		}
	}
	c := newContext(w, req)
	c.handlers = middlewares
	c.app = app
	app.router.handle(c)
}

// template
func (app *Application) SetFuncMap(funcMap template.FuncMap) {
	app.funcMap = funcMap
}
func (app *Application) LoadHTMLGlob(pattern string) {
	app.templates = template.Must(template.New("").Funcs(app.funcMap).ParseGlob(pattern))
}

// router group for app
type RouterGroup struct {
	prefix      string
	middlewares []HandlerFunc
	parent      *RouterGroup
	app         *Application
}

func (group *RouterGroup) Group(prefix string) *RouterGroup {
	app := group.app
	newGroup := &RouterGroup{
		prefix: prefix,
		parent: group,
		app:    app,
	}
	app.groups = append(app.groups, newGroup)
	return newGroup
}

func (group *RouterGroup) addRoute(method string, comp string, handler HandlerFunc) {
	pattern := group.prefix + comp
	log.Printf("Route %4s - %s", method, pattern)
	group.app.router.addRoute(method, pattern, handler)
}

func (group *RouterGroup) Use(middlewares ...HandlerFunc) {
	group.middlewares = append(group.middlewares, middlewares...)
}

func (group *RouterGroup) GET(pattern string, handler HandlerFunc) {
	group.addRoute(http.MethodGet, pattern, handler)
}

func (group *RouterGroup) POST(pattern string, handler HandlerFunc) {
	group.addRoute(http.MethodPost, pattern, handler)
}

func (group *RouterGroup) createStaticHandler(relativePath string, fs http.FileSystem) HandlerFunc {
	absolutePath := path.Join(group.prefix, relativePath)
	fileServer := http.StripPrefix(absolutePath, http.FileServer(fs))
	return func(c *Context) {
		file := c.Param("filepath")
		if _, err := fs.Open(file); err != nil {
			c.Status(http.StatusNotFound)
			return
		}
		fileServer.ServeHTTP(c.Writer, c.Req)
	}
}

func (group *RouterGroup) Static(relativePath string, root string) {
	handler := group.createStaticHandler(relativePath, http.Dir(root))
	urlPattern := path.Join(relativePath, "/*filepath")
	group.GET(urlPattern, handler)
}
