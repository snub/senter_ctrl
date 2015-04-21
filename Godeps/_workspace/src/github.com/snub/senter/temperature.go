package senter

import (
	"github.com/jinzhu/gorm"
	"time"
)

const temperatureTableName string = "sensor_temperature"

type Temperature struct {
	Id        int64
	SensorId  int64
	Timestamp time.Time
	Value     float32
}

func NewTemperature(sensor *Sensor, timestamp int64, value float32) *Temperature {
	return &Temperature{0, sensor.Id, time.Unix(timestamp, 0).UTC(), value}
}

func LoadTemperaturesBySensorId(sensorId int64) []*Temperature {
	logger.Printf("load temperatures by sensor id: %d\n", sensorId)
	db := getDb()
	var ts []*Temperature
	if err := db.Where("sensor_id = ?", sensorId).Order("timestamp desc").Find(&ts).Error; err != nil {
		if err != gorm.RecordNotFound {
			logger.Printf("unable to load temperatures by sensor id: %s\n", err)
		} else {
			logger.Printf("no record found: sensor id = %d\n", sensorId)
		}
	}
	return ts
}

func (t Temperature) TableName() string {
	return temperatureTableName
}

func (t *Temperature) Create() {
	db := getDb()
	if db.NewRecord(t) {
		if err := db.Create(t).Error; err != nil {
			logger.Printf("unable to create temperature: %s\n", err)
		}
	} else {
		logger.Printf("cannot create, temperature already exists with id: %d\n", t.Id)
	}
}

func (t *Temperature) Save() {
	db := getDb()
	if db.NewRecord(t) {
		if err := db.Create(t).Error; err != nil {
			logger.Printf("unable to create temperature: %s\n", err)
		}
	} else {
		logger.Println("temperature does not suppor updating")
	}
}
