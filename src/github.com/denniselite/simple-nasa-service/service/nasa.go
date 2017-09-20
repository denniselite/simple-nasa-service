package service

import (
	"github.com/denniselite/simple-nasa-service/libs"
	"github.com/denniselite/simple-nasa-service/structs"
	"log"
	"fmt"
	"time"
	"strconv"
	"errors"
)

type NasaService struct {
	Db *libs.DB
	NSManager *libs.NasaServerManager
	pagesProcessed int
	NumCPU int
}

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

	var chunk int
	var chunkRest int

	chunk = response.Page.TotalPages / s.NumCPU
	chunkRest = response.Page.TotalPages % s.NumCPU

	var routinesCount int
	if s.NumCPU > 1 {
		routinesCount = s.NumCPU - 1
	} else {
		routinesCount = 1
	}

	for i := 0; i < routinesCount; i++ {
		go func(i int) {
			var endCount int
			if i != routinesCount {
				endCount = (i + 1) * chunk
			} else {
				endCount = (i + 1) * chunk + chunkRest
			}
			for endCount > response.Page.Number {
				response, err = s.NSManager.GetNEOInfo(i * chunk + 1)
				if err != nil {
					log.Printf("Error to process NASA API response: %s", err.Error())
					return
				}
				err = s.saveNeoFromResponse(&response)
				if err != nil {
					log.Println(err)
				}
			}
		}(i)
	}
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

	s.pagesProcessed++
	log.Println(fmt.Sprintf("Pages processed: %d of %d", s.pagesProcessed, response.Page.TotalPages))
	if (response.Page.Number > 0) && (response.Page.Number == response.Page.TotalPages) {
		log.Println("Processed!")
	}
	return
}