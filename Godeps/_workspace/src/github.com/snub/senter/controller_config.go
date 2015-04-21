package senter

import (
	"database/sql"
	"github.com/jinzhu/gorm"
	"time"
)

const controllerConfigTableName string = "sensor_controller_config"

type ControllerConfig struct {
	Id             int64
	ControllerId   int64
	IpAddress      sql.NullString
	UpdateInterval sql.NullInt64
	NtpIpAddress   sql.NullString
	CreatedAt      time.Time
	UpdatedAt      time.Time
}

func NewControllerConfig(controllerId int64) *ControllerConfig {
	return &ControllerConfig{ControllerId: controllerId}
}

//func NewControllerConfig(controllerId int64, ipAddress string, updateInterval int64, ntpIpAddress string) *ControllerConfig {
//	return &ControllerConfig{ControllerId: controllerId, Id: 0, IpAddress: ipAddress, UpdateInterval: updateInterval, NtpIpAddress: ntpIpAddress}
//}

func LoadControllerConfigByControllerId(controllerId int64) *ControllerConfig {
	logger.Printf("load controller config by controller id: %d\n", controllerId)
	db := getDb()
	var cs []ControllerConfig
	if err := db.Where("controller_id = ?", controllerId).Find(&cs).Error; err != nil {
		if err != gorm.RecordNotFound {
			logger.Printf("unable to load controller config bycontroller id: %s\n", err)
		} else {
			logger.Printf("no record found: controller id = %d", controllerId)
		}
		return NewControllerConfig(controllerId)
	}
	logger.Printf("cs: %v\n", cs)
	if len(cs) == 0 {
		return NewControllerConfig(controllerId)
	}
	if len(cs) > 1 {
		logger.Println("more than one result by controller id: %d\n", controllerId)
		return nil
	}
	return &(cs[0])
}

func (c ControllerConfig) TableName() string {
	return controllerConfigTableName
}

func (c *ControllerConfig) New() bool {
	return getDb().NewRecord(c)
}

// TODO on update check rows affected
func (c *ControllerConfig) Save() {
	db := getDb()
	if db.NewRecord(c) {
		if err := db.Create(c).Error; err != nil {
			logger.Printf("unable to create controller config: %s\n", err)
		}
	} else {
		if err := db.Save(c).Error; err != nil {
			logger.Printf("unable to save controller config: %s\n", err)
		}
	}
}
