package service

import (
	"github.com/denniselite/simple-nasa-service/libs"
	"github.com/denniselite/simple-nasa-service/structs"
	"log"
	"fmt"
)

type NasaService struct {
	Db *libs.DB
	NSManager *libs.NasaServerManager
}

func (s *NasaService) Run(db *libs.DB, ns *libs.NasaServerManager) {
	s.Db = db
	s.NSManager = ns

	err := s.processNEOs()
	if err != nil {
		panic(err)
	}
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

	var neos []structs.NearEarthObject
	var response structs.NasaResponse
	response.Page.TotalPages = 1

	response, err = s.NSManager.GetNEOInfo(0)
	if err != nil {
		log.Printf("Error to process NASA API response: %s", err.Error())
		return
	}

	var chunk int
	chunk = response.Page.TotalPages / 10
	log.Println(chunk)

	progress := make([]int, 10)
	for i := 0; i < 10; i++ {
		go func(i int) {
			for (i + 1) * chunk > response.Page.Number {
				response, err = s.NSManager.GetNEOInfo(i * chunk + 1)
				if err != nil {
					log.Printf("Error to process NASA API response: %s", err.Error())
					return
				}
				for k, v := range response.NearEarthObjects {
					neos = append(neos, v)
					//res := fmt.Sprintf("\rProcessed %d of %d", k+(response.Page.Number)*response.Page.Size, response.Page.TotalElements)
					res := k+(response.Page.Number)*response.Page.Size
					progress[i] = res
					fmt.Printf("\r%v", progress)
				}
			}
		}(i)
	}
	return
}