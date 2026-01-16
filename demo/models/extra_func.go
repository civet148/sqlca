package models

import (
	"time"

	"github.com/civet148/sqlca/v3"
)

func NowTime() string {
	return time.Now().Format("2006-01-02 15:04:05")
}

func (do *InventoryData) BeforeQueryData(db *sqlca.Engine) error {
	//log.Debugf("%+v", do)
	return nil
}

func (do *InventoryData) BeforeCreateData(db *sqlca.Engine) error {
	do.CreateTime = NowTime()
	do.UpdateTime = NowTime()
	//log.Debugf("%+v", do)
	return nil
}

func (do *InventoryData) BeforeUpdateData(db *sqlca.Engine) error {
	do.UpdateTime = NowTime()
	//log.Debugf("%+v", do)
	return nil
}

func (do *InventoryData) BeforeDeleteData(db *sqlca.Engine) error {
	//log.Debugf("%+v", do)
	return nil
}

func (do *InventoryData) AfterQueryData(db *sqlca.Engine) error {
	//log.Debugf("%+v", do)
	return nil
}

func (do *InventoryData) AfterCreateData(db *sqlca.Engine) error {
	//log.Debugf("%+v", do)
	return nil
}

func (do *InventoryData) AfterUpdateData(db *sqlca.Engine) error {
	//log.Debugf("%+v", do)
	return nil
}

func (do *InventoryData) AfterDeleteData(db *sqlca.Engine) error {
	//log.Debugf("%+v", do)
	return nil
}
