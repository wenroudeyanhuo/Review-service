package data

import (
	"errors"
	"fmt"
	"gorm.io/driver/mysql"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"review-service/internal/conf"
	"review-service/internal/data/query"
	"strings"

	"github.com/go-kratos/kratos/v2/log"
	"github.com/google/wire"
)

// ProviderSet is data providers.
var ProviderSet = wire.NewSet(NewData, NewReviewerRepo, NewDB)

// Data .
type Data struct {
	// TODO wrapped database client
	query *query.Query
	log   *log.Helper
}

// NewData .
func NewData(db *gorm.DB, logger log.Logger) (*Data, func(), error) {
	cleanup := func() {
		log.NewHelper(logger).Info("closing the data resources")
	}
	//非常重要!  为gen生成的代码指定了数据库对象
	query.SetDefault(db)
	return &Data{query: query.Q, log: log.NewHelper(logger)}, cleanup, nil
}
func NewDB(cfg *conf.Data) (*gorm.DB, error) {
	switch strings.ToLower(cfg.Database.GetDriver()) {
	case "mysql":
		db, err := gorm.Open(mysql.Open(cfg.Database.GetSource()))
		if err != nil {
			panic(fmt.Errorf("connect db fail: %w", err))
		}
		return db, err
	case "sqlite":
		db, err := gorm.Open(sqlite.Open(cfg.Database.GetSource()))
		if err != nil {
			panic(fmt.Errorf("connect db fail: %w", err))
		}
		return db, err
	}
	return nil, errors.New("unknown database driver")
}
