# mongo-mload 工具使用介绍
===


## 功能
 MongoDB 压力测试工具

##2015.02.11更新：
		新增update模式，在--operation后使用update参数


## 参数
	  --host   	       压测目标（如127.0.0.1）
	  --port    	       压测端口（default 27017）
	  --db             操作数据库名称（默认mongobench）
	  --collection             操作数据库表（默认data_test）
	  --userName         用户名（如果有）
	  --passWord         密码（如果有）
	  --cpunum	       多核设定（默认1）
	  --procnum	       压测并发个数（默认4）
	  --datanum	       每个线程插入数据条数（默认10000）
	  --logpath	       日志路径（默认./log.log）
	  --jsonfile	       希望插入document路径（不选用该参数则使用默认的插入格式）
	  --operation	       压测模式（insert,prepare,query,tps,update）prepare模式会在插入完成后为查询会用的项添加索引
	  --queryall	       压测模式为query的时候，是否返回所有查询到的结果（默认false，即db.xx.findOne()）
	  --clean		       是否清理数据(默认false，如果为true将drop数据库mongobench)
	  --geo          是否进行空间地理数据的测试（默认false, 即普通查询和索引；true 则使用经纬度类型数据进行查询）
    ---geofield          空间地理查询测试使用的2d sphere字段名称（默认 loc）

##测试实例
 // --logpath /tmp/log.log
###插入测试
		首先清理数据库：
		go run mload.go --host 127.0.0.1 --clean true

		再来进行插入测试：
		使用8核cpu，8个并发，每个并发插入100000条数据，日志输入为/tmp/log.log，插入的每条数据为./test_data.json中的内容

		go run mload.go --host 127.0.0.1 --datanum 100000 --procnum 8 --cpunum 8 --jsonfile ./test_data.json --operation insert

###查询测试
		首先清理数据库：
		go run mload.go --host 127.0.0.1 --clean true

		再来为查询准备数据（比如准备1000000条）：
		go run mload.go --host 127.0.0.1 --datanum 1000000 --procnum 1 --operation prepare

		接下来进行测试（limit one的）：
		使用8核cpu，8个并发
		go run mload.go --host 127.0.0.1 --datanum 1000000 --procnum 8 --cpunum 8 --operation query

		在进行非limit one的：
		使用8核cpu，8个并发
		go run mload.go --host 127.0.0.1 --datanum 1000000 --procnum 8 --cpunum 8 --operation query  --queryall true
###读写测试
		首先清理数据库：
		go run mload.go --host 127.0.0.1 --clean true

		再来为查询准备数据（比如准备1000000条）：
		go run mload.go --host 127.0.0.1 --datanum 1000000 --procnum 1 --logpath /tmp/log.log --operation prepare
		再来进行测试		
		go run mload.go --host 127.0.0.1 --datanum 1000000 --procnum 1 --logpath /tmp/log.log --operation tps
###更新测试
		首先清理数据库：
		go run mload.go --host 127.0.0.1 --clean true

		再来为查询准备数据（比如准备1000000条）：
		go run mload.go --host 127.0.0.1 --datanum 10 --procnum 1 --operation prepare
		再来进行update压测
		go run mload.go --host 127.0.0.1 --datanum 1 --procnum 10 --operation update

### Geo查询测试
    首先清理数据库：
		go run mload.go --host 127.0.0.1 --clean true

		再来为查询准备数据（比如准备1000000条）：
    go run mload.go --host 127.0.0.1 --datanum 100000 --procnum 2 --operation prepare
    go run mload.go --host 127.0.0.1 --datanum 100000 --procnum 2 --operation prepare --db test --collection testccc --geofield gps --geo

    接下来进行测试（limit one的）：
    使用8核cpu，8个并发
    go run mload.go --host 127.0.0.1 --datanum 100000 --procnum 8 --cpunum 4 --operation query  --geofield loc --geo true
    go run mload.go --host 127.0.0.1 --datanum 100000 --procnum 8 --cpunum 4 --operation query --db metok_core --collection cell_position --geofield loc --geo

    在进行非limit one的：
    使用8核cpu，8个并发
    go run mload.go --host 127.0.0.1 --datanum 100000 --procnum 8 --cpunum 4 --operation query  --queryall true --geofield loc --geo true
    go run mload.go --host 127.0.0.1 --datanum 100000 --procnum 8 --cpunum 4 --operation query  --queryall true --db metok_core --collection cell_position --geofield loc --geo

## 待完成
		1.日志完善
		3.update压测
		4.单实例、多数据库 的压测
