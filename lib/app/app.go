package app

import (
	"context"
	"github.com/capell/capell_scan/lib/discovery"
	"google.golang.org/grpc"
	"log"
	"math/rand"
	"net"
	"os"

	"github.com/BurntSushi/toml"
	"github.com/go-redis/redis/v8"
	"go.mongodb.org/mongo-driver/mongo"
	mongoOptions "go.mongodb.org/mongo-driver/mongo/options"
	"gorm.io/gorm"

	"github.com/capell/capell_scan/lib/db"
	"github.com/capell/capell_scan/lib/httpclient"
	"github.com/capell/capell_scan/lib/httpserver"
	"github.com/capell/capell_scan/lib/logger"
	stime "github.com/capell/capell_scan/lib/time"
)

func Assert(tips string, err error) {
	if err != nil {
		logger.Error("%s err:%v", tips, err)
		log.Fatalf("%s err:%v", tips, err)
	} else {
		logger.Info("%s", tips)
	}
}

type App struct {
	opt        *Option
	httpServ   *httpserver.HttpServer
	grpcServer *grpc.Server
	redis      map[string]*redis.Client
	mongos     map[string]*mongo.Client
	etcd       *discovery.Register
}

func (a *App) Init(opt *Option, allopt interface{}) (err error) {
	rand.Seed(stime.UnixNano())
	if allopt != nil {
		if len(os.Args) != 2 {
			log.Fatalf("Usage:%s <conf.toml>", os.Args[0])
			return
		}

		cfg := os.Args[1]
		_, err = toml.DecodeFile(cfg, allopt)
		if err != nil {
			log.Fatalf("parse option file:%v err:%v", cfg, err)
			return
		}
	}

	a.opt = opt

	err = logger.Init(&opt.Log)
	if err != nil {
		log.Fatalf("init logger err:%v", err)
		return
	}

	if opt.HttpClient != nil {
		err := httpclient.Init(opt.HttpClient)
		if err != nil {
			logger.Error("init http client with option:%#v err:%v", opt.HttpClient, err)
			return err
		}
	}

	if opt.HttpServer != nil {
		s, err := httpserver.NewHttpServer(opt.HttpServer)
		if err != nil {
			logger.Error("new http server with option:%#v err:%v", *opt.HttpServer, err)
			return err
		}
		a.httpServ = s
	}

	if opt.GrpcServer != nil {
		server := grpc.NewServer()
		a.grpcServer = server
	}

	if opt.Databases != nil {
		err = a.loadDbs(opt.Databases)
		if err != nil {
			return err
		}
	}

	if opt.Mongos != nil {
		err = a.loadMongos(opt.Mongos)
		if err != nil {
			return err
		}
	}
	return
}

func (a *App) Run() {
	if a.grpcServer != nil {
		go func() {
			addr := a.opt.GrpcServer.Addr
			lis, err := net.Listen("tcp", addr)
			if err != nil {
				panic(err)
			}
			if err := a.grpcServer.Serve(lis); err != nil {
				panic(err)
			}
		}()
	}
	if a.httpServ != nil {
		a.httpServ.Run()
	}
}

func (a *App) HttpServer() *httpserver.HttpServer {
	return a.httpServ
}

func (a *App) loadDbs(opt map[string]db.Option) (err error) {
	err = db.Init(opt)
	return
}

func (a *App) GetDB(name string) (gdb *gorm.DB) {
	gdb, _ = db.GetDB(name)
	return
}

func (a *App) loadRedis(opt map[string]RedisOption) (err error) {
	rdb := make(map[string]*redis.Client, len(opt))
	for k, v := range opt {
		logger.Info("open redis.%v:%v/%v", k, v.Addr, v.DB)
		opts := &redis.Options{
			Addr:     v.Addr,
			Password: v.Password,
			DB:       v.DB,
		}
		c := redis.NewClient(opts)
		rdb[k] = c
	}
	a.redis = rdb
	return
}

func (a *App) GetRedis(name string) (c *redis.Client) {
	if a.redis == nil {
		return
	}
	c = a.redis[name]
	return
}

func (a *App) loadMongos(opt map[string]MongoOption) (err error) {
	ms := make(map[string]*mongo.Client, len(opt))
	for k, v := range opt {
		logger.Info("connect mongo.%v:%v", k, v.Url)
		clientOpts := mongoOptions.Client().ApplyURI(v.Url)
		client, err := mongo.Connect(context.TODO(), clientOpts)
		if err != nil {
			logger.Error("connect mongo.%v:%v err:%v", k, v.Url, err)
			return err
		}
		ms[k] = client
	}
	a.mongos = ms
	return
}

func (a *App) GetMongo(name string) (c *mongo.Client) {
	if a.mongos == nil {
		return
	}
	c = a.mongos[name]
	return
}

func (a *App) InitEtcd(addr string) (etcdRegister *discovery.Register) {
	etcdAddress := []string{addr}
	// 服务注册
	etcdRegister = discovery.NewRegister(etcdAddress)
	//grpcAddress := "127.0.0.1:10002"
	//defer etcdRegister.Stop()
	//taskNode := discovery.Server{
	//	Name: "swap",
	//	Addr: grpcAddress,
	//}
	//if _, err := etcdRegister.Register(taskNode, 10); err != nil {
	//	panic(fmt.Sprintf("start server failed, err: %v", err))
	//}
	return
}

func (a *App) GetGrpcServer() (s *grpc.Server) {
	return a.grpcServer
}
