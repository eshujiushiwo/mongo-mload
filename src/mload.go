package main

import (
	"flag"
	"fmt"
	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
	"io"
	"log"
	"os"
	"runtime" //for goroutine
	"time"
	//"reflect" //for test
	"encoding/json"
	"io/ioutil"
	"math/rand"
)

var logfile *os.File
var logger *log.Logger
var datainfo *Datainfo
var jsondata = make(map[string]interface{})
var jsonMap = make(map[string]interface{})

type Mongobench struct {
	host            string
	port            string
	userName        string
	passWord        string
	cpunum          int
	datanum         int
	procnum         int
	mongoClient     *mgo.Session
	mongoDatabase   *mgo.Database
	mongoCollection *mgo.Collection
}

type Datainfo struct {
	Name string
	Num  int
}

func GetMongoDBUrl(addr, userName, passWord string, port string) string {
	var mongoDBUrl string

	if port == "no" {
		if userName == "" || passWord == "" {

			mongoDBUrl = "mongodb://" + addr

		} else {

			mongoDBUrl = "mongodb://" + userName + ":" + passWord + "@" + addr

		}

	} else {

		if userName == "" || passWord == "" {

			mongoDBUrl = "mongodb://" + addr + ":" + port

		} else {

			mongoDBUrl = "mongodb://" + userName + ":" + passWord + "@" + addr + ":" + port

		}

	}
	return mongoDBUrl
}

func Newmongobench(host, userName, passWord string, port string, cpunum, datanum, procnum int) *Mongobench {

	mongobench := &Mongobench{host, port, userName, passWord, cpunum, datanum, procnum, nil, nil, nil}

	return mongobench
}

func ReadJson(filename string) (map[string]interface{}, error) {
	bytes, err := ioutil.ReadFile(filename)
	if err != nil {
		logger.Println("ReadFile:", err.Error())
		return nil, err
	}
	if err := json.Unmarshal(bytes, &jsondata); err != nil {
		logger.Println("unmarshal:", err.Error())
		return nil, err
	}
	return jsondata, nil

}

func (mongobench *Mongobench) Conn(mongoDBUrl string) {
	var err error
	mongobench.mongoClient, err = mgo.Dial(mongoDBUrl)
	if err != nil {
		logger.Println("Connect to", mongobench.host, "UserName is ", mongobench.userName, "PassWord is", mongobench.passWord, "Failed")
	}
	logger.Println("Connect to", mongobench.host, "UserName is ", mongobench.userName, "PassWord is", mongobench.passWord, "Success")
	mongobench.mongoDatabase = mongobench.mongoClient.DB("mongobench")
	mongobench.mongoCollection = mongobench.mongoDatabase.C("data_test")

}

func (mongobench *Mongobench) InsertData(jsoninfo string, ch chan int) {
	r := rand.New(rand.NewSource(time.Now().UnixNano()))
	datainfo := Datainfo{"Edison", 1}
	for i := 0; i < mongobench.datanum; i++ {

		a := datainfo.Num * r.Intn(mongobench.datanum)
		if jsoninfo == "no" {
			err := mongobench.mongoCollection.Insert(bson.M{"Name": datainfo.Name, "Num": a})

			if err != nil {
				logger.Println("insert failed:", err)
				os.Exit(1)
			}
		} else if jsoninfo == "yes" {

			err := mongobench.mongoCollection.Insert(jsonMap)
			if err != nil {
				logger.Println("insert failed:", err)
				os.Exit(1)
			}
		}

	}
	ch <- 1

}

func main() {
	var host, userName, passWord, port, logpath, jsonfile string
	var insert bool
	var cpunum, datanum, procnum int
	var err1 error
	var multi_logfile []io.Writer
	//The DataInfo
	//datainfo := Datainfo{"Edison", 1}
	flag.StringVar(&host, "host", "", "The mongodb host")
	flag.StringVar(&userName, "userName", "", "The mongodb username")
	flag.StringVar(&passWord, "passWord", "", "The mongodb password")
	flag.StringVar(&port, "port", "27017", "The mongodb port")
	flag.IntVar(&cpunum, "cpunum", 1, "The cpu number wanna use")
	flag.IntVar(&datanum, "datanum", 10000, "The data count per proc")
	flag.IntVar(&procnum, "procnum", 4, "The proc num ")
	flag.StringVar(&logpath, "logpath", "./log.log", "the log path ")
	flag.StringVar(&jsonfile, "jsonfile", "", "the json file u wanna insert(only one json )")
	flag.BoolVar(&insert, "insert", false, "set true to do insert ")
	flag.Parse()

	logfile, err1 = os.OpenFile(logpath, os.O_RDWR|os.O_CREATE, 0666)
	defer logfile.Close()
	if err1 != nil {
		fmt.Println(err1)
		os.Exit(-1)
	}
	multi_logfile = []io.Writer{
		logfile,
		os.Stdout,
	}
	logfiles := io.MultiWriter(multi_logfile...)
	logger = log.New(logfiles, "\r\n", log.Ldate|log.Ltime|log.Lshortfile)

	mongourl := GetMongoDBUrl(host, userName, passWord, port)

	if host != "" {
		logger.Println("=====job start.=====")
		logger.Println("start init colletion")
		mongobench := Newmongobench(host, userName, passWord, port, cpunum, datanum, procnum)

		mongobench.Conn(mongourl)
		defer mongobench.mongoClient.Close()

		if jsonfile != "" {
			var err error
			jsonMap, err = ReadJson(jsonfile)
			logger.Println(jsonMap)
			if err != nil {
				logger.Println(err)
			}
		}

		if insert == true {
			chs := make([]chan int, mongobench.procnum)
			runtime.GOMAXPROCS(mongobench.cpunum)
			for i := 0; i < mongobench.procnum; i++ {
				fmt.Println(i)

				chs[i] = make(chan int)

				if jsonfile == "" {
					go mongobench.InsertData("no", chs[i])
				} else {
					go mongobench.InsertData("yes", chs[i])
				}
			}

			for _, cha := range chs {
				<-cha

			}
		} else {
			logger.Println("insert is set to ", insert)
		}
		logger.Println("=====Done.=====")
	} else {
		fmt.Println("Please use -help to check the usage")
		fmt.Println("At least need host parameter")

	}

}
