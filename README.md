# mongo-mload 工具使用介绍
===


## 功能
 MongoDB 压力测试工具

##2015.01.20更新：
		第一版本完成
		实现了：insert压力测试，日志还未完善。暂时还是通过mongostat查看测试情况。

## 参数
		--host   	 压测目标（如127.0.0.1）
		--port    	 压测端口（default 27017）
		--userName   用户名（如果有）
		--passWord   密码（如果有）
		--cpunum	 多核设定（默认1）
		--procnum	 压测并发个数（默认4） 
		--datanum	 每个线程插入数据条数（默认10000）
		--logpath	 日志路径（默认./log.log）
		--jsonfile	 希望插入document路径（不选用该参数则使用默认的插入格式）
		--insert	 插入模式（默认false，使用true时会进行插入压测）


##测试实例

		使用8核cpu，8个并发，每个并发插入100000条数据，日志输入为/tmp/log.log，插入的每条数据为./test_data.json中的内容

		go run mload.go --host 127.0.0.1 --datanum 100000 --procnum 8 --cpunum 8 --logpath /tmp/log.log --jsonfile ./test_data.json --insert true



## 待完成
		1.日志完善
		2.查询压测
		3.update压测
		4.单实例、多数据库 的压测
