package service

import (
	"github.com/denniselite/simple-nasa-service/libs"
	"github.com/denniselite/simple-nasa-service/structs"
	"log"
	"fmt"
	"time"
	"strconv"
	"errors"
	"github.com/jinzhu/gorm"
)

type NasaService struct {
	Db *libs.DB
	NSManager *libs.NasaServerManager
	NumCPU int
}

var nasaService *NasaService

func (s *NasaService) Run(db *libs.DB, ns *libs.NasaServerManager) {
	s.Db = db
	s.NSManager = ns

	if s.isDbEmpty() {
		s.initDB()
		err := s.processNEOs()
		if err != nil {
			panic(err)
		}
	}

	nasaService = s
}

func GetNasaService() (s *NasaService) {
	s = nasaService
	return
}

func (s *NasaService) isDbEmpty() bool {
	return !s.Db.HasTable(&structs.NEO{}) && !s.Db.HasTable(&structs.NEOData{})
}

func (s *NasaService) initDB()  {
	if !s.Db.HasTable(&structs.NEO{}) {
		s.Db.AutoMigrate(&structs.NEO{})
	}
	if !s.Db.HasTable(&structs.NEOData{}) {
		s.Db.AutoMigrate(&structs.NEOData{})
	}
}

func (s *NasaService) processNEOs() (err error) {
	var response structs.NasaResponse
	response.Page.TotalPages = 1

	log.Println("NEO data aggregation has been started")
	response, err = s.NSManager.GetNEOInfo(0)
	if err != nil {
		log.Printf("Error to process NASA API response: %s", err.Error())
		return
	}
	err = s.saveNeoFromResponse(&response)
	if err != nil {
		return
	}

	for i := 1; i <= response.Page.TotalPages; i++ {
		response, err = s.NSManager.GetNEOInfo(i)
		if err != nil {
			log.Printf("Error to process NASA API response: %s", err.Error())
			return
		}
		err = s.saveNeoFromResponse(&response)
		if err != nil {
			log.Println(err)
			break;
		}
		fmt.Printf("\rPages processed: %d of %d", i, response.Page.TotalPages)
	}
	log.Println("Processed!")
	return
}

func (s *NasaService) saveNeoFromResponse(response *structs.NasaResponse) (err error) {
	for _, v := range response.NearEarthObjects {
		neo := new(structs.NEO)
		neo.Name = v.Name
		neo.Reference = v.NeoReferenceId
		neo.IsHazardous = v.IsPotentiallyHazardousAsteroid
		err = s.Db.Save(neo).Error
		if err != nil {
			err = errors.New(fmt.Sprintf("Error to save NEO to DB: %s", err.Error()))
			return
		}
		for _, vData := range v.NearEarthObjectData {
			neoData := new(structs.NEOData)
			neoData.Reference = v.NeoReferenceId
			neoData.Date, err = time.Parse("2006-01-02", vData.Date)
			if err != nil {
				err = errors.New(fmt.Sprintf("Error to parse the date value: %s", err.Error()))
				return
			}
			neoData.Speed, err = strconv.ParseFloat(vData.RelativeVelocity.KilometersPerHour, 32)
			if err != nil {
				err = errors.New(fmt.Sprintf("Error to parse the NEO speed value: %s", err.Error()))
				return
			}
			err = s.Db.Save(neoData).Error
			if err != nil {
				err = errors.New(fmt.Sprintf("Error to save NEO to DB: %s", err.Error()))
				return
			}
		}
	}
	return
}

func (s *NasaService) GetHazardous() (res []structs.NEO, err error) {
	err = s.Db.Exec("SELECT neos.").
		Where("is_hazardous = TRUE").
		Find(&res).Error
	return
}

func (s *NasaService) GetFastest(isHazardous bool) (res structs.NEO, err error) {
	err = s.Db.Debug().
		Where("is_hazardous = ?", isHazardous).
		Preload("NEOData", func(db *gorm.DB) *gorm.DB {
		return db.Where("neo_data.speed > 0").Order("neo_data.speed DESC")
	}).
		First(&res).Error
	return
}