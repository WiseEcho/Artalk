package db

import (
	"fmt"
	"time"

	"github.com/artalkjs/artalk/v2/internal/config"
	"github.com/artalkjs/artalk/v2/internal/log"
	"github.com/artalkjs/artalk/v2/internal/log/zapgorm2"
	"gorm.io/gorm"
	gormlogger "gorm.io/gorm/logger"
	"gorm.io/gorm/schema"
)

func NewDB(conf config.DBConf) (*gorm.DB, error) {
	gormConfig := &gorm.Config{
		Logger: zapgorm2.New(log.GetLogger(), gormlogger.Config{
			Colorful:                  true,
			IgnoreRecordNotFoundError: true,
			SlowThreshold:             200 * time.Millisecond,
			LogLevel:                  gormlogger.Info,
		}),
	}

	// Enable Prepared Statement by default
	if prepareStmt := conf.PrepareStmt; prepareStmt != nil {
		gormConfig.PrepareStmt = *prepareStmt
	} else {
		gormConfig.PrepareStmt = true
	}

	var dsn string
	if conf.Dsn != "" {
		dsn = conf.Dsn
	} else {
		dsn = getDsnByConf(conf)
	}

	switch conf.Type {
	case config.TypeSQLite:
		return OpenSQLite(dsn, gormConfig)
	case config.TypeMySql:
		return OpenMySql(dsn, gormConfig)
	case config.TypePostgreSQL:
		return OpenPostgreSQL(dsn, gormConfig)
	case config.TypeMSSQL:
		return OpenSqlServer(dsn, gormConfig)
	}

	return nil, fmt.Errorf("unsupported database type %s", conf.Type)
}

func NewTestDB() (*gorm.DB, error) {
	return OpenSQLite("file::memory:?cache=shared", &gorm.Config{
		DisableForeignKeyConstraintWhenMigrating: true,
		NamingStrategy: schema.NamingStrategy{
			TablePrefix: "atk_",
		},
	})
}

func CloseDB(db *gorm.DB) error {
	ddb, err := db.DB()
	if err != nil {
		return err
	}
	return ddb.Close()
}
