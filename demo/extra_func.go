package models

import (
	"github.com/civet148/log"
	"github.com/civet148/sqlca/v2"
	"time"
)

func NowTime() string {
	return time.Now().Format("2006-01-02 15:04:05")
}

func (do *InventoryData) BeforeCreate(db *sqlca.Engine) error {
	do.CreateTime = NowTime()
	do.UpdateTime = NowTime()
	log.Infof("%+v", do)
	return nil
}

func (do *InventoryData) BeforeUpdate(db *sqlca.Engine) error {
	do.UpdateTime = NowTime()
	log.Infof("%+v", do)
	return nil
}

func (do *InventoryData) BeforeDelete(db *sqlca.Engine) error {
	log.Infof("%+v", do)
	return nil
}

func (do *InventoryData) AfterCreate(db *sqlca.Engine) error {
	log.Infof("%+v", do)
	return nil
}

func (do *InventoryData) AfterUpdate(db *sqlca.Engine) error {
	log.Infof("%+v", do)
	return nil
}

func (do *InventoryData) AfterDelete(db *sqlca.Engine) error {
	log.Infof("%+v", do)
	return nil
}

