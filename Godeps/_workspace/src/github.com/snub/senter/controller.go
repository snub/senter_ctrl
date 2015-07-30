package senter

import (
	"github.com/jinzhu/gorm"
	"time"
)

const controllerTableName string = "sensor_controller"

type Controller struct {
	Id          int64
	MacAddress  string
	LastStartup time.Time
	CreatedAt   time.Time
	UpdatedAt   time.Time
}

func NewController(macAddress string, lastUpdate int64) *Controller {
	return &Controller{Id: 0, MacAddress: macAddress, LastStartup: time.Unix(lastUpdate, 0).UTC()}
}

func LoadControllers() []*Controller {
	logger.Println("load controllers")
	db := getDb()
	var cs []*Controller
	query := db.Find(&cs)
	if query.Error != nil {
		logger.Println("unable to load controllers")
		return nil
	}
	return cs
}

func LoadControllerById(controllerId int64) *Controller {
	logger.Printf("load controller by id: %d\n", controllerId)
	db := getDb()
	var controller Controller
	if err := db.First(&controller, controllerId).Error; err != nil {
		logger.Printf("unable to load controller by id: %s\n", err)
		return NewController("", 0)
	}
	logger.Printf("controller: %v\n", controller)
	return &controller
}

// TODO better error handling
func LoadControllerByMacAddress(macAddress string) *Controller {
	logger.Printf("load controller by mac address: %s\n", macAddress)
	db := getDb()
	var cs []Controller
	if err := db.Where("mac_address = ?", macAddress).Find(&cs).Error; err != nil {
		if err != gorm.RecordNotFound {
			logger.Printf("unable to load controller by mac address: %s\n", err)
		} else {
			logger.Printf("no record found: mac address = %s", macAddress)
		}
		return NewController(macAddress, 0)
	}
	logger.Printf("cs: %v\n", cs)
	if len(cs) == 0 {
		return NewController(macAddress, 0)
	}
	if len(cs) > 1 {
		logger.Println("more than one result by mac address: %s\n", macAddress)
		return nil
	}
	return &(cs[0])
}

func (c Controller) TableName() string {
	return controllerTableName
}

func (c *Controller) New() bool {
	return getDb().NewRecord(c)
}

func (c *Controller) SetLastStartup(lastUpdate int64) {
	c.LastStartup = time.Unix(lastUpdate, 0).UTC()
}

// TODO on update check rows affected
func (c *Controller) Save() {
	db := getDb()
	if db.NewRecord(c) {
		if err := db.Create(c).Error; err != nil {
			logger.Printf("unable to create controller: %s\n", err)
		}
	} else {
		if err := db.Save(c).Error; err != nil {
			logger.Printf("unable to save controller: %s\n", err)
		}
	}
}
