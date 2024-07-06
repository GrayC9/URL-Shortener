package storage

import (
	"context"
	"time"

	"github.com/GrayC9/URL-Shortener/internal/config"
	"github.com/jmoiron/sqlx"
	"github.com/sirupsen/logrus"
)

type MaridDB struct {
	DB   *sqlx.DB
	logg *logrus.Logger
	conf *config.Config
}

func DB(logger *logrus.Logger, cnf *config.Config) *MaridDB {
	return &MaridDB{
		logg: logger,
		conf: cnf,
	}
}

func (db *MaridDB) connect() error {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*3)
	defer cancel()
	d, err := sqlx.ConnectContext(ctx, "mariadb", db.conf.DBconfig)
	if err != nil {
		db.logg.Errorln(err)
		return err
	}
	if err := d.DB.Ping(); err != nil {
		db.logg.Errorln(err)
		return err
	}
	db.DB = d
	db.logg.Infoln("DB configurate")

	return nil
}

func (db *MaridDB) Save(short, original string) {

}

func (db *MaridDB) Get(short string) string {

	db.logg.Errorln("No such shortcode in storage")
	return ""
}
