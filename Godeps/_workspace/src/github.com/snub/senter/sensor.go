package senter

import (
	"database/sql"
	"github.com/jinzhu/gorm"
	"time"
)

const sensorTableName string = "sensor"

type Sensor struct {
	Id            int64
	DeviceAddress string
	ControllerId  int64
	Name          sql.NullString
	Description   sql.NullString
	CreatedAt     time.Time
	UpdatedAt     time.Time
}

func NewSensor(deviceAddress string) *Sensor {
	return &Sensor{Id: 0, ControllerId: 0, DeviceAddress: deviceAddress}
}

func LoadSensors() []*Sensor {
	logger.Println("load sensors")
	db := getDb()
	var ss []*Sensor
	query := db.Find(&ss)
	if query.Error != nil {
		logger.Println("unable to load sensors")
		return nil
	}
	return ss
}

func LoadSensorById(sensorId int64) *Sensor {
	logger.Printf("load sensors by id: %d\n", sensorId)
	db := getDb()
	var sensor Sensor
	if err := db.First(&sensor, sensorId).Error; err != nil {
		logger.Printf("unable to load senspr by id: %s\n", err)
		return NewSensor("")
	}
	logger.Printf("sensor: %v\n", sensor)
	return &sensor
}

func LoadSensorsByControllerId(controllerId int64) []*Sensor {
	logger.Printf("load sensors by controller id: %d\n", controllerId)
	db := getDb()
	var ss []*Sensor
	if err := db.Where("controller_id = ?", controllerId).Find(&ss).Error; err != nil {
		if err != gorm.RecordNotFound {
			logger.Printf("unable to load sensors by controller id: %s\n", err)
		} else {
			logger.Printf("no record found: controller id = %d\n", controllerId)
		}
	}
	return ss
}

// TODO better error handling
func LoadSensorByDeviceAddress(deviceAddress string) *Sensor {
	logger.Printf("load sensor by device address: %s\n", deviceAddress)
	db := getDb()
	var ss []Sensor
	if err := db.Where("device_address = ?", deviceAddress).Find(&ss).Error; err != nil {
		if err != gorm.RecordNotFound {
			logger.Printf("unable to load sensor by device address: %s\n", err)
		} else {
			logger.Printf("no record found: device address = %s\n", deviceAddress)
		}
		return NewSensor(deviceAddress)
	}
	logger.Printf("ss: %v\n", ss)
	if len(ss) == 0 {
		return NewSensor(deviceAddress)
	}
	if len(ss) > 1 {
		logger.Printf("more than one result by device address: %s\n", deviceAddress)
		return nil

	}
	return &(ss[0])
}

func (s Sensor) TableName() string {
	return sensorTableName
}

func (s *Sensor) New() bool {
	return getDb().NewRecord(s)
}

func (s *Sensor) SetControllerId(controllerId int64) {
	s.ControllerId = controllerId
}

func (s *Sensor) Create() {
	db := getDb()
	if db.NewRecord(s) {
		if err := db.Create(s).Error; err != nil {
			logger.Printf("unable to create sensor: %s\n", err)
		}
	} else {
		logger.Printf("cannot create, sensor already exists with id: %d\n", s.Id)
	}
}

// TODO on update check rows affected
func (s *Sensor) Save() {
	db := getDb()
	if db.NewRecord(s) {
		if err := db.Create(s).Error; err != nil {
			logger.Printf("unable to create sensor: %s\n", err)
		}
	} else {
		if err := db.Save(s).Error; err != nil {
			logger.Printf("unable to save sensor: %s\n", err)
		}
	}
}
