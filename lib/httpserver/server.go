package httpserver

import (
	"fmt"
	"github.com/afex/hystrix-go/hystrix"
	"github.com/capell/capell_scan/lib/logger"
	"github.com/fvbock/endless"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"net/http"
	_ "net/http/pprof"
	"time"
)

type HttpServer struct {
	opt *Option
	e   *gin.Engine

	healthStatus int

	nHystrix *hystrix.CircuitBreaker
}

type log struct {
}

func (log) Write(b []byte) (n int, err error) {
	logger.Info("%s", string(b))
	return len(b), nil
}

func NewHttpServer(opt *Option) (s *HttpServer, err error) {
	e := gin.New()
	e.UseH2C = true
	e.Use(cors.Default())
	gin.DisableConsoleColor()
	gin.DefaultWriter = log{}
	e.Use(gin.LoggerWithFormatter(func(param gin.LogFormatterParams) string {
		return fmt.Sprintf("%v %v %v %v %v %.3f %v",
			param.ClientIP,
			param.Method,
			param.Path,
			param.StatusCode,
			param.BodySize,
			float64(param.Latency)/float64(time.Millisecond),
			param.ErrorMessage,
		)
	}))
	e.Use(gin.CustomRecovery(func(c *gin.Context, recovered interface{}) {
		err, ok := recovered.(string)
		if ok {
			c.String(http.StatusInternalServerError, "Internal Service Error")
		}
		c.AbortWithStatus(http.StatusInternalServerError)
		logger.Alert("panic request:%v err:%v", c.Request.URL.String(), err)
		logger.PrintStack(err)
	}))

	s = &HttpServer{
		opt:          opt,
		e:            e,
		healthStatus: 200,
	}
	if opt.EnableDebug {
		s.EnableDebug()
	}
	if opt.EnableHealthCheck {
		s.EnableHealthCheck()
	}
	if !opt.DisableMetrics {
		s.e.GET("/metrics", s.httpDefault)
		s.e.POST("/metrics", s.httpDefault)
		//err = s.EnableMetrics()
		if err != nil {
			s = nil
			return
		}
	}
	return
}

func (s *HttpServer) Run() (err error) {
	return endless.ListenAndServe(s.opt.Addr, s.e)
}

func (s *HttpServer) httpDefault(c *gin.Context) {
	http.DefaultServeMux.ServeHTTP(c.Writer, c.Request)
}

func (s *HttpServer) EnableDebug() {
	s.e.GET("/debug/*p", s.httpDefault)
	s.e.POST("/debug/*p", s.httpDefault)
}

func (s *HttpServer) SetHealthStatus(status int) {
	s.healthStatus = status
}

func (s *HttpServer) EnableHealthCheck() {
	s.e.GET("/healthcheck", func(c *gin.Context) {
		c.JSON(s.healthStatus, map[string]interface{}{
			"timestamp": time.Now().Unix(),
		})
	})
	s.e.POST("/healthcheck/set", func(c *gin.Context) {
		var p struct {
			Status int `json:"status"`
		}
		err := c.BindJSON(&p)
		if err != nil {
			c.JSON(400, "require parameter error")
			return
		}
		s.SetHealthStatus(p.Status)
		c.JSON(200, nil)
	})
}
func (s *HttpServer) Engine() *gin.Engine {
	return s.e
}
