package database

import (
	"log"

	"github.com/Thales-s-Orgs/VideoEncoder-GoLang/domain"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/sqlite"
	_ "github.com/lib/pq"
)

type Database struct {
	DB            *gorm.DB
	DSN           string
	DSNTest       string
	DBType        string
	DBTypeTest    string
	Debug         bool
	AutomigrateDB bool
	Env           string
}

func NewDB() *Database {
	return &Database{}
}

func NewDBTest() *gorm.DB {
	db := NewDB()
	db.Env = "Test"
	db.DBTypeTest = "sqlite3"
	db.DSNTest = ":memory:"
	db.AutomigrateDB = true
	db.Debug = true

	connection, err := db.Connect()

	if err != nil {
		log.Fatalf("Database could not been started: %v", err)
	}

	return connection
}

func (d *Database) Connect() (*gorm.DB, error) {

	var err error

	if d.Env != "test" {
		d.DB, err = gorm.Open(d.DBType, d.DSN)
	} else {
		d.DB, err = gorm.Open(d.DBTypeTest, d.DSNTest)
	}

	if err != nil {
		return nil, err
	}

	if d.Debug {
		d.DB.LogMode(true)
	}

	if d.AutomigrateDB {
		d.DB.AutoMigrate(&domain.Video{}, &domain.Job{})
		d.DB.Model(domain.Job{}).AddForeignKey("video_id", "videos (id)", "CASCADE", "CASCADE")
	}

	return d.DB, nil

}
