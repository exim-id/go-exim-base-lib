package db

import (
	"context"
	"log"
	"math"
	"runtime/debug"
	"sync"

	"github.com/exim-id/go-exim-base-lib/dto"
	"github.com/exim-id/go-exim-base-lib/env"
	"github.com/exim-id/go-exim-base-lib/errors"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type MongoTransaction[T interface{}] struct {
	sync.Mutex
	pesimisticLock bool
	db             *mongo.Database
	cli            *mongo.Client
	session        mongo.Session
}

func (m *MongoTransaction[T]) Transaction(doing func(*mongo.Database) T) T {
	if m.pesimisticLock {
		m.Lock()
	}
	cli, err := MongoDb()
	m.cli = cli
	defer errors.MongoClose(cli)
	if err != nil {
		panic(err)
	}
	m.db = cli.Database(env.GetMongoDbName())
	m.session, err = cli.StartSession()
	defer m.deferTransaction()
	if err != nil {
		panic(err)
	}
	return doing(m.db)
}

func (m *MongoTransaction[T]) deferTransaction() {
	defer func() {
		if m.pesimisticLock {
			m.Unlock()
		}
	}()
	if r := recover(); r != nil {
		log.Println("Catch", r)
		log.Println("Stack", string(debug.Stack()))
		if err := m.session.AbortTransaction(context.Background()); err != nil {
			panic(err)
		}
	} else {
		if err := m.session.CommitTransaction(context.Background()); err != nil {
			panic(err)
		}
	}
	m.session.EndSession(context.Background())
}

func NewMongoTransaction[T interface{}](pesimisticLock bool) *MongoTransaction[T] {
	return &MongoTransaction[T]{pesimisticLock: pesimisticLock}
}

func OnMongoConnection[T interface{}](doing func(*mongo.Database) T) T {
	cli, err := MongoDb()
	defer errors.MongoClose(cli)
	if err != nil {
		panic(err)
	}
	db := cli.Database(env.GetMongoDbName())
	return doing(db)
}

func MongoDbPage[T interface{}, W interface{}](req dto.PageReq[W], doing func(*mongo.Database) *mongo.Collection) dto.PageRes[T] {
	return OnMongoConnection[dto.PageRes[T]](func(d *mongo.Database) dto.PageRes[T] { return mongoDbPage[T](doing(d), req) })
}

func mongoDbPage[T interface{}, W interface{}](coll *mongo.Collection, req dto.PageReq[W]) dto.PageRes[T] {
	var result dto.PageRes[T]
	result.Page = req.Page
	result.Size = req.Size
	opt := options.Find()
	size := int64(req.Size)
	opt.Limit = &size
	skip := int64((req.Page - 1) * req.Size)
	opt.Skip = &skip
	cur, err := coll.Find(context.Background(), req.Where, opt)
	defer errors.MongoCursorClose(cur)
	if err != nil {
		panic(err)
	}
	var response []T
	if err = cur.All(context.Background(), &response); err != nil {
		panic(err)
	}
	result.Datas = response
	count, err := coll.CountDocuments(context.Background(), req.Where)
	if err != nil {
		panic(err)
	}
	result.PageCount = uint64(math.Ceil(float64(count) / float64(req.Size)))
	return result
}

func MongoDbFindAllCollection[T interface{}, W interface{}](filter W, doing func(*mongo.Database) *mongo.Collection) []T {
	return OnMongoConnection[[]T](func(d *mongo.Database) []T { return MongoDbFindAllFromCollection[T](doing(d), filter) })
}

func MongoDbFindAllFromCollection[T interface{}, W interface{}](coll *mongo.Collection, filter W) []T {
	cur, err := coll.Find(context.Background(), filter, options.Find())
	defer errors.MongoCursorClose(cur)
	if err != nil {
		panic(err)
	}
	var result []T
	if err = cur.All(context.Background(), &result); err != nil {
		panic(err)
	}
	return result
}

func MongoDb() (*mongo.Client, error) {
	client, err := mongo.Connect(context.Background(), options.Client().ApplyURI(env.GetMongoDBUrl()))
	return client, err
}
