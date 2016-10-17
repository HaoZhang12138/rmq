####说明
基于**rabbitmq**的消息队列发布/订阅模式的实现
####编译
  - `go build main.go`

####运行
- 作为生产者运行： `./main -p`
- 作为消费者运行： `./main -c`

####依赖库
1. `github.com/streadway/amqp`
2. `gopkg.in/mgo.v2/bson`
3. `github.com/garyburd/redigo/redis`



