package sqlca

import (
	"context"
	"github.com/civet148/log"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"
	"time"
)

const (
	defaultTimeoutSeconds = 3
	defaultReadWriteTimeoutSeconds = 30*60
)

type MgoExecutor struct {
	client *mongo.Client
	db     *mongo.Database
}

func newMgoExecutor(strDatabaseName, strDSN string) (executor, error) {

	if strDatabaseName == "" {
		return nil, log.Errorf("database name require")
	}
	ctx, cancel := newContextTimeout(defaultTimeoutSeconds)
	defer cancel()
	client, err := mongo.Connect(ctx, options.Client().ApplyURI(strDSN))
	if err != nil {
		log.Errorf("connect %s error [%s]", strDSN, err)
		return nil, err
	}
	err = client.Ping(ctx, readpref.Primary())
	if err != nil {
		log.Errorf("[%s] ping error [%s]", strDSN, err)
		return nil, err
	}

	return &MgoExecutor{
		client: client,
		db:     client.Database(strDatabaseName),
	}, nil
}

func newContextTimeout(timeout int) (context.Context, context.CancelFunc) {
	return context.WithTimeout(context.Background(), time.Duration(timeout)*time.Second)
}

func (m *MgoExecutor) Ping() (err error) {
	return nil
}

func (m *MgoExecutor) SetMaxOpenConns(n int) {

}

func (m *MgoExecutor) SetMaxIdleConns(n int) {

}

func (m *MgoExecutor) SetConnMaxLifetime(d time.Duration) {

}

func (m *MgoExecutor) SetConnMaxIdleTime(d time.Duration) {

}

func (m *MgoExecutor) Exec(e *Engine, strSQL string) (rowsAffected, lastInsertId int64, err error) {
	return
}

func (m *MgoExecutor) Query(e *Engine, strSQL string) (count int64, err error) {
	ctx, cancel := newContextTimeout(defaultReadWriteTimeoutSeconds)
	defer cancel()
	var cursor *mongo.Cursor
	type User struct {
		Name string `bson:"name"`
		Sex string `bson:"sex"`
		Age uint8 `bson:"age"`
	}
	cursor, err = m.db.Collection("user").Find(ctx, nil)
	if err != nil {
		return 0, log.Errorf("find error [%s]", err)
	}
	var user []User
 	if err = cursor.All(ctx, &user); err != nil {
		return 0, log.Errorf("cursor all error [%s]", err)
	}
	return
}

func (m *MgoExecutor) QueryEx(e *Engine, strSQL string) (rowsAffected, total int64, err error) {
	return
}

func (m *MgoExecutor) QueryRaw(e *Engine, strSQL string) (rowsAffected int64, err error) {
	return
}

func (m *MgoExecutor) QueryMap(e *Engine, strSQL string) (rowsAffected int64, err error) {

	return
}

func (m *MgoExecutor) Update(e *Engine, strSQL string) (rowsAffected int64, err error) {

	return
}

func (m *MgoExecutor) Insert(e *Engine, strSQL string) (lastInsertId int64, err error) {

	return
}

func (m *MgoExecutor) Upsert(e *Engine, strSQL string) (lastInsertId int64, err error) {

	return
}

func (m *MgoExecutor) Delete(e *Engine, strSQL string) (rowsAffected int64, err error) {

	return
}

func (m *MgoExecutor) Close() error {
	ctx, cancel := newContextTimeout(3)
	defer cancel()
	return m.client.Disconnect(ctx)
}

func (m *MgoExecutor) txBegin() (tx executor, err error) {
	return
}

func (m *MgoExecutor) txGet(e *Engine, dest interface{}, strQuery string, args ...interface{}) (count int64, err error) {
	return
}

func (m *MgoExecutor) txExec(e *Engine, strQuery string, args ...interface{}) (lastInsertId, rowsAffected int64, err error) {
	return
}

func (m *MgoExecutor) txRollback() (err error) {
	return
}

func (m *MgoExecutor) txCommit() (err error) {
	return
}

//------------------------------------------------------------------------------------------------------

func (m *MgoExecutor) collection(strName string, opts...*options.CollectionOptions) *mongo.Collection {
	return m.db.Collection(strName, opts...)
}
