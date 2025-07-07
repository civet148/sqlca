package models

import (
	"github.com/civet148/sqlca/v3"
	"time"
)

func NowTime() string {
	return time.Now().Format("2006-01-02 15:04:05")
}

func (do *InventoryData) BeforeQuery(db *sqlca.Engine) error {
	//log.Debugf("%+v", do)
	return nil
}

func (do *InventoryData) BeforeCreate(db *sqlca.Engine) error {
	do.CreateTime = NowTime()
	do.UpdateTime = NowTime()
	//log.Debugf("%+v", do)
	return nil
}

func (do *InventoryData) BeforeUpdate(db *sqlca.Engine) error {
	do.UpdateTime = NowTime()
	//log.Debugf("%+v", do)
	return nil
}

func (do *InventoryData) BeforeDelete(db *sqlca.Engine) error {
	//log.Debugf("%+v", do)
	return nil
}

func (do *InventoryData) AfterQuery(db *sqlca.Engine) error {
	//log.Debugf("%+v", do)
	return nil
}

func (do *InventoryData) AfterCreate(db *sqlca.Engine) error {
	//log.Debugf("%+v", do)
	return nil
}

func (do *InventoryData) AfterUpdate(db *sqlca.Engine) error {
	//log.Debugf("%+v", do)
	return nil
}

func (do *InventoryData) AfterDelete(db *sqlca.Engine) error {
	//log.Debugf("%+v", do)
	return nil
}
