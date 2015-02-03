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
var r *rand.Rand

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

func (mongobench *Mongobench) QueryData(all bool, ch chan int) {

	datainfo := Datainfo{"Edison", 1}
	if all == true {
		result1 := []int{}
		for i := 0; i < mongobench.datanum; i++ {

			b := datainfo.Num * r.Intn(mongobench.datanum)
			query := bson.M{"Num": b}
			mongobench.mongoCollection.Find(query).All(&result1)

		}
	} else {

		var result1 interface{}
		for i := 0; i < mongobench.datanum; i++ {

			b := datainfo.Num * r.Intn(mongobench.datanum)
			query := bson.M{"Num": b}
			mongobench.mongoCollection.Find(query).One(&result1)

		}
	}
	ch <- 1

}
func (mongobench *Mongobench) CleanJob() {
	logger.Println("Start clean database :mongobench")
	mongobench.mongoDatabase.DropDatabase()

}
func (mongobench *Mongobench) AddIndex() {
	logger.Println("Start build index Num")

	mongobench.mongoCollection.EnsureIndexKey("Num")
}

func main() {
	var host, userName, passWord, port, logpath, jsonfile string
	var operation string
	var queryall, clean bool
	var cpunum, datanum, procnum int
	var err1 error
	var multi_logfile []io.Writer

	r = rand.New(rand.NewSource(time.Now().UnixNano()))

	flag.StringVar(&host, "host", "", "The mongodb host")
	flag.StringVar(&userName, "userName", "", "The mongodb username")
	flag.StringVar(&passWord, "passWord", "", "The mongodb password")
	flag.StringVar(&port, "port", "27017", "The mongodb port")
	flag.IntVar(&cpunum, "cpunum", 1, "The cpu number wanna use")
	flag.IntVar(&datanum, "datanum", 10000, "The data count per proc")
	flag.IntVar(&procnum, "procnum", 4, "The proc num ")
	flag.StringVar(&logpath, "logpath", "./log.log", "the log path ")
	flag.StringVar(&jsonfile, "jsonfile", "", "the json file u wanna insert(only one json )")
	flag.StringVar(&operation, "operation", "", "the operation ")
	flag.BoolVar(&queryall, "queryall", false, "query all or limit one")
	flag.BoolVar(&clean, "clean", false, "Drop the Database mongobench")
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

	if host != "" && operation != "" && clean == false {
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

		if operation == "insert" {
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
		} else if operation == "prepare" {
			chs := make([]chan int, mongobench.procnum)
			runtime.GOMAXPROCS(mongobench.cpunum)
			for i := 0; i < mongobench.procnum; i++ {
				fmt.Println(i)

				chs[i] = make(chan int)

				go mongobench.InsertData("no", chs[i])

			}

			for _, cha := range chs {
				<-cha

			}
			mongobench.AddIndex()

		} else if operation == "query" {

			chs := make([]chan int, mongobench.procnum)
			runtime.GOMAXPROCS(mongobench.cpunum)
			for i := 0; i < mongobench.procnum; i++ {
				fmt.Println(i)

				chs[i] = make(chan int)

				go mongobench.QueryData(queryall, chs[i])

			}

			for _, cha := range chs {
				<-cha

			}

		} else if operation == "tps" {
			//query
			chs := make([]chan int, mongobench.procnum)
			runtime.GOMAXPROCS(mongobench.cpunum)
			for i := 0; i < mongobench.procnum; i++ {
				fmt.Println(i)

				chs[i] = make(chan int)

				go mongobench.QueryData(queryall, chs[i])

			}
			//insert
			chs1 := make([]chan int, mongobench.procnum)
			runtime.GOMAXPROCS(mongobench.cpunum)
			for i := 0; i < mongobench.procnum; i++ {
				fmt.Println(i)

				chs1[i] = make(chan int)

				go mongobench.InsertData("no", chs1[i])

			}
			for _, chb := range chs1 {
				<-chb

			}
			for _, cha := range chs {
				<-cha

			}

		} else {
			fmt.Println("Only support operation prepare/insert/query!")

		}

		logger.Println("=====Done.=====")
	} else if host != "" && clean == true {

		logger.Println("=====job start.=====")
		logger.Println("start init colletion")

		mongobench := Newmongobench(host, userName, passWord, port, cpunum, datanum, procnum)

		mongobench.Conn(mongourl)
		defer mongobench.mongoClient.Close()
		mongobench.CleanJob()
	} else {
		fmt.Println("Please use -help to check the usage")
		fmt.Println("At least need host parameter")

	}

}
