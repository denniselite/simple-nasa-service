package main

import (
	"github.com/valyala/fasthttp"
	"log"
	"flag"
	"io/ioutil"
	"github.com/denniselite/simple-nasa-service/structs"
	"gopkg.in/yaml.v2"
	"github.com/denniselite/simple-nasa-service/handlers"
	"fmt"
	"github.com/denniselite/simple-nasa-service/service"
	"github.com/denniselite/simple-nasa-service/libs"
	"runtime"
)

func main()  {
	numCPU := runtime.NumCPU()
	fmt.Println("Numbers of CPU:", numCPU)
	runtime.GOMAXPROCS(numCPU)
	log.SetFlags(log.Ldate | log.Lmicroseconds | log.Lshortfile | log.LUTC)

	config := loadConfig()

	pgConnectionString := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=disable",
		config.Db.Host,
		config.Db.Port,
		config.Db.Username,
		config.Db.Password,
		config.Db.Database,
	)

	db, err := libs.ConnectDB("postgres", pgConnectionString)
	if err != nil {
		panic(err)
	}

	ns := new(libs.NasaServerManager)
	ns.Config = config.NSManager

	nasaService := new(service.NasaService)
	nasaService.NumCPU = numCPU
	nasaService.Run(db, ns)

	if err := fasthttp.ListenAndServe(config.Listen, fastHTTPHandler); err != nil {
		log.Fatalf("Error in ListenAndServe: %s", err)
	}
}

func loadConfig() structs.Config {
	var filename string

	// register flags
	flag.StringVar(&filename, "config", "", "config filename")
	flag.StringVar(&filename, "c", "", "config filename (shorthand)")

	flag.Parse()

	config := structs.Config{}

	configData, err := ioutil.ReadFile(filename)
	if err != nil {
		panic(err)
	}

	err = yaml.Unmarshal(configData, &config)
	if err != nil {
		panic(err)
	}

	return config
}


func fastHTTPHandler(ctx *fasthttp.RequestCtx) {
	switch string(ctx.Path()) {
	case "/":
		handlers.MainHandler(ctx)
	default:
		ctx.Error("Unsupported path", fasthttp.StatusNotFound)
	}
}