package wmongo

import (
	"context"
	"fmt"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"reflect"
	"strings"
	"time"
)

type ColOperator struct {
	cli      *mongo.Client
	database string
	table    string
	col      *mongo.Collection
	timeout  time.Duration
}

func NewColOperator(cli *mongo.Client, db, col string) *ColOperator {
	c := &ColOperator{
		cli:      cli,
		database: db,
		table:    col,
		col:      cli.Database(db).Collection(col),
		timeout:  time.Duration(10) * time.Second,
	}
	return c
}

func (c *ColOperator) ScanCursor(ctx context.Context, cursor *mongo.Cursor, result interface{}) (err error) {
	defer func() {
		err1 := cursor.Close(ctx)
		if err == nil && err1 != nil {
			err = err1
		}
	}()

	resultv := reflect.ValueOf(result)
	if resultv.Kind() != reflect.Ptr || resultv.Elem().Kind() != reflect.Slice {
		err = fmt.Errorf("result argument must be a slice address")
		return
	}
	slicev := resultv.Elem()
	slicev = slicev.Slice(0, slicev.Cap())
	elemt := slicev.Type().Elem()

	i := 0
	for ; cursor.Next(ctx); i++ {
		if slicev.Len() == i {
			// slice长度耗尽时，通过append触发slice_grow
			// 并将slice的len扩大到与cap相同，保证新增长的空间可用index索引
			elemp := reflect.New(elemt)
			err = cursor.Decode(elemp.Interface())
			if err != nil {
				return
			}
			slicev = reflect.Append(slicev, elemp.Elem())
			slicev = slicev.Slice(0, slicev.Cap())
		} else {
			// slice长度未耗尽时通过index索引使用剩余空间
			err = cursor.Decode(slicev.Index(i).Addr().Interface())
			if err != nil {
				return
			}
		}
	}
	// 将slice恢复为真正的长度
	resultv.Elem().Set(slicev.Slice(0, i))

	err = cursor.Err()
	if err != nil {
		return
	}
	return
}

func (c *ColOperator) GenSortBson(sort []string) (result bson.D) {
	result = bson.D{}
	for _, v := range sort {
		if strings.HasPrefix(v, "-") {
			v = strings.TrimPrefix(v, "-")
			result = append(result, bson.E{Key: v, Value: -1})
		} else {
			v = strings.TrimPrefix(v, "+")
			result = append(result, bson.E{Key: v, Value: 1})
		}
	}
	return
}

func (c *ColOperator) FindOneAndUpdate(ctx context.Context, filter, update, result interface{},
	sort []string, upsert, returnNew bool, opts ...*options.FindOneAndUpdateOptions) (has bool, err error) {
	opt := options.FindOneAndUpdate()
	if len(sort) > 0 {
		opt.SetSort(c.GenSortBson(sort))
	}
	opt.SetUpsert(upsert)
	if returnNew {
		opt.SetReturnDocument(options.After)
	} else {
		opt.SetReturnDocument(options.Before)
	}
	opts = append(opts, opt)
	ctx, cancel := context.WithTimeout(ctx, c.timeout)
	defer cancel()
	err = c.col.FindOneAndUpdate(ctx, filter, update, opts...).Decode(result)
	if err == mongo.ErrNoDocuments {
		has, err = false, nil
		return
	}
	if err != nil {
		return
	}
	has = true
	return
}

func (c *ColOperator) Find(ctx context.Context, filter interface{}, result interface{},
	sort []string, skip, limit int64, opts ...*options.FindOptions) (err error) {
	opt := options.Find().SetSkip(skip).SetLimit(limit)
	if len(sort) > 0 {
		opt.SetSort(c.GenSortBson(sort))
	}
	opts = append(opts, opt)

	ctx, cancel := context.WithTimeout(ctx, c.timeout)
	defer cancel()
	cursor, err := c.col.Find(ctx, filter, opts...)
	if err != nil {
		if err == mongo.ErrNilDocument {
			err = nil
		}
		return
	}
	err = c.ScanCursor(ctx, cursor, result)
	if err != nil {
		return
	}
	return
}

func (c *ColOperator) FindOne(ctx context.Context, filter interface{}, result interface{},
	sort []string, skip int64, opts ...*options.FindOneOptions) (has bool, err error) {
	opt := options.FindOne().SetSkip(skip)
	if len(sort) > 0 {
		opt.SetSort(c.GenSortBson(sort))
	}
	opts = append(opts, opt)

	ctx, cancel := context.WithTimeout(ctx, c.timeout)
	defer cancel()
	err = c.col.FindOne(ctx, filter, opts...).Decode(result)
	if err != nil {
		return
	}
	has = true
	return
}

func (c *ColOperator) FindById(ctx context.Context, Id interface{}, result interface{},
	opts ...*options.FindOneOptions) (has bool, err error) {

	filter := bson.M{"_id": Id}
	return c.FindOne(ctx, filter, result, nil, 0, opts...)
}

func (c *ColOperator) InsertOne(ctx context.Context, document interface{},
	opts ...*options.InsertOneOptions) (insertedID interface{}, err error) {

	ctx, cancel := context.WithTimeout(ctx, c.timeout)
	defer cancel()
	result, err := c.col.InsertOne(ctx, document, opts...)
	if err != nil {
		return
	}
	insertedID = result.InsertedID
	return
}

func (c *ColOperator) UpdateOne(ctx context.Context, filter, update interface{}, upsert bool, opts ...*options.UpdateOptions) (result *mongo.UpdateResult, err error) {
	opt := options.Update()
	opt.SetUpsert(upsert)
	opts = append(opts, opt)
	ctx, cancel := context.WithTimeout(ctx, c.timeout)
	defer cancel()

	resultOri, err := c.col.UpdateOne(ctx, filter, update, opts...)
	if err != nil {
		return
	}
	result = &mongo.UpdateResult{
		MatchedCount:  resultOri.MatchedCount,
		ModifiedCount: resultOri.ModifiedCount,
		UpsertedCount: resultOri.UpsertedCount,
		UpsertedID:    resultOri.UpsertedID,
	}
	return
}

func (c *ColOperator) Upsert(ctx context.Context, filter, update interface{}, opts ...*options.UpdateOptions) (result *mongo.UpdateResult, err error) {
	return c.UpdateOne(ctx, filter, update, true, opts...)
}

func (c *ColOperator) UpsertID(ctx context.Context, ID interface{}, update interface{}) (result *mongo.UpdateResult, err error) {
	filter := bson.M{"_id": ID}
	return c.UpdateOne(ctx, filter, update, true)
}

func (c *ColOperator) DeleteID(ctx context.Context, ID interface{},
	opts ...*options.DeleteOptions) (has bool, err error) {
	filter := bson.M{"_id": ID}
	return c.DeleteOne(ctx, filter, opts...)
}

func (c *ColOperator) DeleteOne(ctx context.Context, filter interface{},
	opts ...*options.DeleteOptions) (has bool, err error) {
	ctx, cancel := context.WithTimeout(ctx, c.timeout)
	defer cancel()
	result, err := c.col.DeleteOne(ctx, filter, opts...)
	if err != nil {
		return
	}
	has = result.DeletedCount > 0
	return
}

func (c *ColOperator) DeleteMany(ctx context.Context, filter interface{},
	opts ...*options.DeleteOptions) (deletedCnt int64, err error) {
	ctx, cancel := context.WithTimeout(ctx, c.timeout)
	defer cancel()
	result, err := c.col.DeleteMany(ctx, filter, opts...)
	if err != nil {
		return
	}
	deletedCnt = result.DeletedCount
	return
}

//func (c *ColOperator) FindByFilter(ctx context.Context, filter, sort, result interface{}, skip, limit int64) (err error) {
//	opt := options.Find()
//	if skip != 0 {
//		opt.SetSkip(skip)
//	}
//	if limit != 0 {
//		opt.SetLimit(limit)
//	}
//	opt.SetSort(sort)
//	cur, err1 := c.col.Find(ctx, filter, opt)
//	if err1 != nil {
//		err = err1
//		return
//	}
//	defer cur.Close(ctx)
//	err = cur.All(ctx, result)
//	return
//}

//------------------------------------------------------------------------------------------------------------------------------------------------------------

//func (c *ColOperator) UpdateOne(ctx context.Context, filter, obj interface{}, opts ...*options.UpdateOptions) (result *mongo.UpdateResult, err error) {
//	result, err = c.col.UpdateOne(ctx, filter, obj, opts...)
//	return
//}
//
//func (c *ColOperator) UpdateMany(ctx context.Context, filter, obj interface{}, opts ...*options.UpdateOptions) (result *mongo.UpdateResult, err error) {
//	result, err = c.col.UpdateMany(ctx, filter, obj, opts...)
//	return
//}
//
//func (c *ColOperator) UpsertOneWthFilter(ctx context.Context, query bson.D, obj interface{}) (result *mongo.UpdateResult, err error) {
//	opts := options.Update().SetUpsert(true)
//	result, err = c.col.UpdateOne(ctx,
//		query,
//		bson.D{
//			{"$set", obj},
//		}, opts)
//	return
//}
//
//func setIntFieldValue(obj interface{}, fieldName string, val int64, err *error) func() {
//	v := reflect.ValueOf(obj)
//	ignore := false
//	switch v.Kind() {
//	case reflect.Pointer, reflect.Interface:
//		v = v.Elem()
//	case reflect.Struct:
//	default:
//		ignore = true
//	}
//	if ignore {
//		return func() {}
//	}
//	var oldVal int64
//	var pf *reflect.Value
//	if !ignore && v.IsValid() {
//		f := v.FieldByName(fieldName)
//		switch f.Kind() {
//		case reflect.Int64:
//			oldVal = f.Int()
//			f.SetInt(val)
//			pf = &f
//		case reflect.Uint64:
//			oldVal = int64(f.Uint())
//			f.SetUint(uint64(val))
//			pf = &f
//		}
//	}
//	return func() {
//		if pf != nil {
//			if err == nil || *err != nil {
//				switch pf.Kind() {
//				case reflect.Int64:
//					pf.SetInt(oldVal)
//				case reflect.Uint64:
//					pf.SetUint(uint64(oldVal))
//				}
//			}
//		}
//	}
//}
//
//func (c *ColOperator) UpsertOneWithCtUt(ctx context.Context, id, obj interface{}) (result *mongo.UpdateResult, err error) {
//	now := time.Now().UnixMilli()
//	defer setIntFieldValue(obj, "Ut", now, &err)()
//	defer setIntFieldValue(obj, "Ct", 0, nil)()
//	opts := options.Update().SetUpsert(true)
//	result, err = c.col.UpdateOne(ctx, bson.D{{"_id", id}},
//		bson.D{
//			{"$set", obj},
//			{"$setOnInsert", bson.D{{"ct", now}}},
//		}, opts)
//	return
//}
//
//func (c *ColOperator) UpsertOneWithShardKey(ctx context.Context, shardKey string, shardVal, id, obj interface{}) (result *mongo.UpdateResult, err error) {
//	opts := options.Update().SetUpsert(true)
//	result, err = c.col.UpdateOne(ctx,
//		bson.D{
//			{"_id", id},
//			{shardKey, shardVal},
//		},
//		bson.D{
//			{"$set", obj},
//		}, opts)
//	return
//}
//
//func (c *ColOperator) UpsertOneWithShardKeyCtUt(ctx context.Context, shardKey string, shardVal, id, obj interface{}) (result *mongo.UpdateResult, err error) {
//	now := time.Now().UnixMilli()
//	defer setIntFieldValue(obj, "Ut", now, &err)()
//	defer setIntFieldValue(obj, "Ct", 0, nil)()
//	opts := options.Update().SetUpsert(true)
//	result, err = c.col.UpdateOne(ctx,
//		bson.D{
//			{"_id", id},
//			{shardKey, shardVal},
//		},
//		bson.D{
//			{"$set", obj},
//			{"$setOnInsert", bson.D{{"ct", now}}},
//		}, opts)
//	return
//
//}
//
//func (c *ColOperator) FindOneById(ctx context.Context, id, obj interface{}) (err error) {
//	err = c.col.FindOne(ctx, bson.D{{"_id", id}}).Decode(obj)
//	if err != nil && err != mongo.ErrNoDocuments {
//		return
//	}
//	return
//}
//
//func (c *ColOperator) FindOneByFilter(ctx context.Context, filter, sort, obj interface{}) (err error) {
//	opt := &options.FindOneOptions{}
//	opt.SetSort(sort)
//	err = c.col.FindOne(ctx, filter, opt).Decode(obj)
//	return
//}
//
//func (c *ColOperator) FindByFilter(ctx context.Context, filter, sort, result interface{}, skip, limit int64) (err error) {
//	opt := options.Find()
//	if skip != 0 {
//		opt.SetSkip(skip)
//	}
//	if limit != 0 {
//		opt.SetLimit(limit)
//	}
//	opt.SetSort(sort)
//	cur, err1 := c.col.Find(ctx, filter, opt)
//	if err1 != nil {
//		err = err1
//		return
//	}
//	defer cur.Close(ctx)
//	err = cur.All(ctx, result)
//	return
//}
//
//func (c *ColOperator) CountByFilter(ctx context.Context, filter interface{}) (count int64, err error) {
//	opt := &options.CountOptions{}
//	count, err = c.col.CountDocuments(ctx, filter, opt)
//	return
//}
//
//func (c *ColOperator) Aggregate(ctx context.Context, pipeline mongo.Pipeline, result any) (err error) {
//	opt := options.Aggregate()
//	cur, aErr := c.col.Aggregate(ctx, pipeline, opt)
//	if aErr != nil {
//		return aErr
//	}
//	err = cur.All(ctx, result)
//	return
//}
//
//func (c *ColOperator) InsertMany(ctx context.Context, documents []any) (result *mongo.InsertManyResult, err error) {
//	opt := &options.InsertManyOptions{}
//	result, err = c.col.InsertMany(ctx, documents, opt)
//	return
//}
