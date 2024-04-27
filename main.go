package main

import (
	"database/sql"
	"flag"
	"fmt"
	"log"
	"os"

	sql2 "github.com/minitauros/swagen/sql"
	"github.com/minitauros/swagen/swagger"

	"gopkg.in/yaml.v2"
)

type Config struct {
	DB struct {
		Host     string
		Port     string
		User     string
		Password string
		Dbname   string
		Schema   string
	}
	Service   swagger.ServiceInfo
	Resources map[string]swagger.Resource // Table name => resource
}

func main() {
	configFilePath := flag.String("conf", "", "config file path")

	flag.Parse()

	if *configFilePath == "" {
		log.Fatal("No config file path given. Use the -conf flag.")
	}

	configBytes, err := os.ReadFile(*configFilePath)
	if err != nil {
		log.Fatal(err)
	}

	conf := Config{}
	err = yaml.Unmarshal(configBytes, &conf)
	if err != nil {
		log.Fatal(err)
	}

	psqlInfo := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
		conf.DB.Host, conf.DB.Port, conf.DB.User, conf.DB.Password, conf.DB.Dbname)

	db, err := sql.Open("postgres", psqlInfo)
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()
	db.Exec("set search_path='" + conf.DB.Schema + "'")

	generator := swagger.Generator{
		TableService: sql2.NewTableService(db),
		Resources:    conf.Resources,
		ServiceInfo:  conf.Service,
	}

	swag, err := generator.Generate()
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println(swag)
}
