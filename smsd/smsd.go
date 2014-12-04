package smsd

import (
	"encoding/json"
	"fmt"
	"github.com/gilankpam/earthquapps-smsd/api"
	_ "github.com/go-sql-driver/mysql"
	"github.com/jmoiron/sqlx"
	"github.com/kr/beanstalk"
	"log"
	"sync"
	"time"
)

var (
	wg        sync.WaitGroup
	apiClient *api.Client
	db        *sqlx.DB
	config    *Config
)

type sms struct {
	Number  string `json:"number"`
	Message string `json:"message"`
}

type phone struct {
	Id     int64
	Number string
}

func Serve(configFile string) {
	config, err := GetConfig(configFile)
	if err != nil {
		panic(err)
	}
	db, err := sqlx.Open("mysql", config.DBUrl)
	if err != nil {
		log.Fatalln(err)
	}
	apiClient = api.NewClient(
		config.Api.CustomerKey,
		config.Api.CustomerKeySecret,
		config.Api.AccessToken,
		config.Api.AccessTokenSecret,
	)
	c, err := beanstalk.Dial("tcp", fmt.Sprintf("%s:%s", config.Beanstalk.Host, config.Beanstalk.Port))
	if err != nil {
		log.Fatalf("Error conecting to beanstalk: %v", err)
	}
	news := beanstalk.NewTubeSet(c, "news")
	verf := beanstalk.NewTubeSet(c, "verification")

	wg.Add(1)
	go listenNewsLoop(news, db)
	wg.Add(1)
	go listenVerLoop(verf)
	log.Println("Go SMS Daemon is running . . .")
	wg.Wait()

}

func listenNewsLoop(news *beanstalk.TubeSet, db *sqlx.DB) {
	for {
		id, body, err := news.Reserve(time.Duration(time.Second))
		if cerr, ok := err.(beanstalk.ConnError); ok && cerr.Err == beanstalk.ErrTimeout {
			//do nothing?
		} else if err != nil {
			panic(err)
		} else {
			log.Println("news : " + string(body))
			go broadcast(string(body), db)
			news.Conn.Delete(id)
		}

	}
	wg.Done()
}

//Send news message to all subs
func broadcast(msg string, db *sqlx.DB) {
	phones := []phone{}
	db.Select(&phones, "SELECT id, number FROM phone WHERE active=true")
	for _, phone := range phones {
		log.Printf("Sending news to : %s", phone.Number)
		go apiClient.SendSMSBulk(phone.Number, msg)
	}
}

//Listen for verivication queue
func listenVerLoop(verf *beanstalk.TubeSet) {
	for {
		id, body, err := verf.Reserve(time.Duration(time.Second))
		if cerr, ok := err.(beanstalk.ConnError); ok && cerr.Err == beanstalk.ErrTimeout {
			//do nothing?
		} else if err != nil {
			panic(err)
		} else {
			sms := new(sms)
			json.Unmarshal(body, sms)
			log.Printf("Sending ver code : %s/%s\n", sms.Number, sms.Message)
			go apiClient.SendSMSBulk(sms.Number, sms.Message)
			verf.Conn.Delete(id)
		}

	}
	wg.Done()
}
