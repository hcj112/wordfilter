package http

import (
	"github.com/gin-gonic/gin"
	"github.com/hcj112/wordfilter/internal/filter"
	"github.com/hcj112/wordfilter/internal/conf"
)

type Server struct {
	engine *gin.Engine
	filter *filter.Filter
}

func New(c *conf.HTTPServer, f *filter.Filter) *Server {
	gin.SetMode(gin.ReleaseMode)
	engine := gin.New()
	engine.Use(loggerHandler, recoverHandler)
	go func() {
		if err := engine.Run(c.Addr); err != nil {
			panic(err)
		}
	}()
	s := &Server{
		engine: engine,
		filter: f,
	}
	s.initRouter()
	return s
}

func (s *Server) initRouter() {
	group := s.engine.Group("/filter")
	group.GET("/keyword/add", s.add)
	group.GET("/keyword/remove", s.remove)
	group.GET("/keyword/replace", s.replace)
	group.GET("/keyword/list", s.list)
}

func (s *Server) Close() {
	s.filter.Close()
}

func (s *Server) add(c *gin.Context) {
	var arg struct {
		Keyword string `form:"keyword" binding:"required"`
	}
	if err := c.BindQuery(&arg); err != nil {
		errors(c, RequestErr, err.Error())
		return
	}
	if err := s.filter.Add(arg.Keyword); err != nil {
		result(c, nil, RequestErr)
		return
	}
	result(c, nil, OK)
	return
}

func (s *Server) remove(c *gin.Context) {
	var arg struct {
		Keyword string `form:"keyword" binding:"required"`
	}
	if err := c.BindQuery(&arg); err != nil {
		errors(c, RequestErr, err.Error())
		return
	}
	if err := s.filter.Remove(arg.Keyword); err != nil {
		result(c, nil, RequestErr)
		return
	}
	result(c, nil, OK)
	return
}

func (s *Server) replace(c *gin.Context) {
	var arg struct {
		Keyword string `form:"keyword" binding:"required"`
	}
	if err := c.BindQuery(&arg); err != nil {
		errors(c, RequestErr, err.Error())
		return
	}
	keyword := s.filter.Filter(arg.Keyword)
	result(c, keyword, OK)
}

func (s *Server) list(c *gin.Context) {
	keywrods,err := s.filter.List()
	if err != nil {
		result(c, nil, RequestErr)
		return
	}
	result(c, keywrods, OK)
}

