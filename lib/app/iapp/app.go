package iapp

import (
	"github.com/go-redis/redis/v8"
	"go.mongodb.org/mongo-driver/mongo"
	"google.golang.org/grpc"
	"gorm.io/gorm"
)

type IApp interface {
	GetDB(name string) *gorm.DB
	GetRedis(name string) *redis.Client
	GetMongo(name string) (c *mongo.Client)
	GetGrpcServer() (s *grpc.Server)
}
