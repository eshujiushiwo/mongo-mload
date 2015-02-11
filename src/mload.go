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

func (mongobench *Mongobench) UpdateData(ch chan int) {
	datainfo := Datainfo{"Edison", 1}
	for i := 0; i < mongobench.datanum; i++ {
		b := datainfo.Num * r.Intn(mongobench.datanum)
		query := bson.M{"Num": b}
		//update := bson.M{"$set": bson.M{"ii": []bson.M{{"i": 58506, "a": 1}, {"i": 59015, "a": 1}, {"i": 58415, "a": 2}, {"i": 58414, "a": 1}, {"i": 59014, "a": 1}, {"i": 59215, "a": 1}, {"i": 58015, "a": 2}, {"i": 58014, "a": 2}, {"i": 32001, "a": 0}, {"i": 112010, "a": 2}, {"i": 59106, "a": 4}, {"i": 58013, "a": 1}, {"i": 58212, "a": 2}, {"i": 58213, "a": 2}, {"i": 59412, "a": 1}, {"i": 50009, "a": 7}, {"i": 59306, "a": 2}, {"i": 112073, "a": 4}, {"i": 58208, "a": 2}, {"i": 40016, "a": 0}, {"i": 41003, "a": 0}, {"i": 58507, "a": 1}, {"i": 58011, "a": 1}, {"i": 59414, "a": 1}, {"i": 58106, "a": 3}, {"i": 40018, "a": 0}, {"i": 80208, "a": 0}, {"i": 59505, "a": 5}, {"i": 58305, "a": 6}, {"i": 58505, "a": 6}, {"i": 112013, "a": 2}, {"i": 59407, "a": 0}, {"i": 58204, "a": 0}, {"i": 112011, "a": 8}, {"i": 58207, "a": 13}, {"i": 59009, "a": 0}, {"i": 80307, "a": 0}, {"i": 58304, "a": 10}, {"i": 80109, "a": 2}, {"i": 80111, "a": 11}, {"i": 80102, "a": 2}, {"i": 41061, "a": 0}, {"i": 80110, "a": 7}, {"i": 80101, "a": 2}, {"i": 80103, "a": 2}, {"i": 80107, "a": 1}, {"i": 80108, "a": 1}, {"i": 80112, "a": 14}, {"i": 80106, "a": 1}, {"i": 80104, "a": 2}, {"i": 80105, "a": 0}, {"i": 58302, "a": 1}, {"i": 58006, "a": 1}, {"i": 59206, "a": 3}, {"i": 58012, "a": 3}, {"i": 58412, "a": 5}, {"i": 59013, "a": 6}, {"i": 58504, "a": 24}, {"i": 59104, "a": 15}, {"i": 59011, "a": 2}, {"i": 59010, "a": 1}, {"i": 58411, "a": 7}, {"i": 58410, "a": 3}, {"i": 58010, "a": 2}, {"i": 59210, "a": 2}, {"i": 58211, "a": 4}, {"i": 59411, "a": 3}, {"i": 58210, "a": 4}, {"i": 59410, "a": 3}, {"i": 58409, "a": 0}, {"i": 58408, "a": 2}, {"i": 59008, "a": 0}, {"i": 59408, "a": 1}, {"i": 59207, "a": 3}, {"i": 58104, "a": 5}, {"i": 59504, "a": 12}, {"i": 58503, "a": 2}, {"i": 59103, "a": 0}, {"i": 58303, "a": 3}, {"i": 59503, "a": 6}, {"i": 50003, "a": 68}, {"i": 58206, "a": 2}, {"i": 59502, "a": 2}, {"i": 59409, "a": 3}, {"i": 59209, "a": 3}, {"i": 58009, "a": 2}, {"i": 40008, "a": 0}, {"i": 58103, "a": 3}, {"i": 59303, "a": 2}, {"i": 59012, "a": 4}, {"i": 59212, "a": 3}, {"i": 59305, "a": 12}, {"i": 59105, "a": 9}, {"i": 58413, "a": 7}, {"i": 58404, "a": 2}, {"i": 58105, "a": 8}, {"i": 58008, "a": 5}, {"i": 59208, "a": 6}, {"i": 58407, "a": 0}, {"i": 59007, "a": 4}, {"i": 58406, "a": 0}, {"i": 58007, "a": 0}, {"i": 59404, "a": 0}, {"i": 59005, "a": 1}, {"i": 58405, "a": 0}, {"i": 50008, "a": 13}, {"i": 58102, "a": 5}, {"i": 59302, "a": 4}, {"i": 58502, "a": 0}, {"i": 59102, "a": 7}, {"i": 58501, "a": 2}, {"i": 59101, "a": 2}, {"i": 40003, "a": 14}, {"i": 41001, "a": 0}, {"i": 50004, "a": 0}, {"i": 40001, "a": 666}, {"i": 59205, "a": 5}, {"i": 59004, "a": 1}, {"i": 58205, "a": 1}, {"i": 59405, "a": 2}, {"i": 112001, "a": 2}, {"i": 58301, "a": 18}, {"i": 59501, "a": 14}, {"i": 58101, "a": 2}, {"i": 59301, "a": 0}, {"i": 58005, "a": 0}, {"i": 59204, "a": 2}, {"i": 58004, "a": 1}, {"i": 58203, "a": 1}, {"i": 58202, "a": 0}, {"i": 59403, "a": 1}, {"i": 30504, "a": 0}, {"i": 30505, "a": 0}, {"i": 58403, "a": 5}, {"i": 59003, "a": 6}, {"i": 30503, "a": 0}, {"i": 31006, "a": 1}, {"i": 58209, "a": 5}, {"i": 58402, "a": 0}, {"i": 59002, "a": 8}, {"i": 30502, "a": 0}, {"i": 30506, "a": 2}, {"i": 58003, "a": 7}, {"i": 59402, "a": 0}, {"i": 50006, "a": 0}, {"i": 41032, "a": 0}, {"i": 59406, "a": 2}, {"i": 30006, "a": 2}, {"i": 30501, "a": 0}, {"i": 58002, "a": 11}, {"i": 59203, "a": 0}, {"i": 30005, "a": 1}, {"i": 31005, "a": 1}, {"i": 58401, "a": 4}, {"i": 59202, "a": 0}, {"i": 30004, "a": 1}, {"i": 31004, "a": 1}, {"i": 58201, "a": 0}, {"i": 40002, "a": 98}, {"i": 59006, "a": 11}, {"i": 30003, "a": 1}, {"i": 31003, "a": 1}, {"i": 59201, "a": 5}, {"i": 40017, "a": 0}, {"i": 59401, "a": 0}, {"i": 50007, "a": 0}, {"i": 40009, "a": 0}, {"i": 50005, "a": 12}, {"i": 30002, "a": 1}, {"i": 31002, "a": 1}, {"i": 59001, "a": 26}, {"i": 30001, "a": 11}, {"i": 31001, "a": 1}, {"i": 58001, "a": 14}, {"i": 90001, "a": 6721}, {"i": 90005, "a": 39, "t": 1423470994}, {"i": 90006, "a": 2, "t": 1423470694}, {"i": 90003, "a": 980639}, {"i": 90004, "a": 99}, {"i": 90008, "a": 4649}, {"i": 80212, "a": 1}, {"i": 80203, "a": 0}, {"i": 80303, "a": 0}, {"i": 80201, "a": 1}, {"i": 80301, "a": 0}, {"i": 80207, "a": 4}, {"i": 80308, "a": 0}, {"i": 80204, "a": 1}, {"i": 80206, "a": 0}, {"i": 80304, "a": 0}, {"i": 80306, "a": 0}, {"i": 80202, "a": 1}, {"i": 80205, "a": 0}, {"i": 80209, "a": 4}, {"i": 80305, "a": 0}}}}
		update := bson.M{"$set": bson.M{"ii": []bson.M{{"i": 58506, "a": 1}}}}
		mongobench.mongoCollection.Update(query, update)
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

		} else if operation == "update" {
			ch := make([]chan int, mongobench.procnum)
			runtime.GOMAXPROCS(mongobench.cpunum)
			for i := 0; i < mongobench.procnum; i++ {
				fmt.Println(i)

				ch[i] = make(chan int)

				go mongobench.UpdateData(ch[i])

			}

			for _, cha := range ch {
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
