[toc]

# Redis
## 参考网址
[redis连环炮面试](https://www.cnblogs.com/myseries/p/12071608.html)

## 安装及启动配置
[centos7安装redis6.0](https://www.cnblogs.com/kurtye/p/14503897.html)

### docker启动redis
[docker操作redis](https://blog.csdn.net/occultskyrong/article/details/85199926)
```
启动redis
docker run -p 6379:6379 --name redis -d redis:latest redis-server

连接redis
docker exec -ti d0b86 redis-cli
docker exec -ti d0b86 redis-cli -h 172.17.0.3 -p 6379 

docker exec -it redis_name redis-cli -h 192.168.1.100 -p 6379 -a your_password

查看ip
docker inspect redis_s | grep IPAddress 查看ip
```

## redis学习路线
**两大维度**：   
    系统维度和应用维度        
**三大主线**：   
   1. 高性能：线程模型、数据结构、持久化、网络框架；

   2. 高可靠性：主从复制、哨兵机制；

   3. 高扩展性：数据分片、负载均衡。

      ![image-20220818164208217](.\image\Redis\image-20220818164208217.png)

      ![image-20220818164224052](E:\calmkNote\image\Redis\image-20220818164224052.png)

   4. ![image-20220818165646917](.\Redis.assets\image-20220818165646917.png)

### 数据结构及操作
**五种基本数据结构**

- 字符串：sds
    SDS：保存了字符串长度及free的空间
- 列表：hash表 + 链表
    - 渐进式hash
    rehashidx初始为0，每次对字典增删改查就会将一个hash桶的数据移动到新的hash表，index++，全部rehash结束后，设为-1
- 集合
- hash表
- 有序集合

**底层实现**
- sds动态字符串
- 整数集合
- 链表
- 跳跃表
- 字典
- 压缩列表
- 未完待续

### 高性能高可靠性技术
#### 数据持久化
持久化机制包括快照和写文件两种
- RDB快照
    - 根据设定规则自动快照、save|bgsave命令、执行flushall命令【该命令会清空数据库所有数据，只要设定了自动规则快照。如果未设定自动快照规则，则不会进行快照】、主从复制时。
    - 步骤：
        1. 创建子进程，父进程继续对外处理请求，子进程复制；
        2. 当写完所有数据后，用该文件替换旧的快照
    - 保存的数据为执行frok那一瞬间的数据，因为克隆进程是写时复制，frok后修改的数据不会被保存，可能会丢
- AOF写文件
    - frok后，父进程的读写操作还是会追加到旧的文件中，同时会保存在缓冲区一份。重写完成后把缓冲区的一次性写入文件中，之后替换旧文件

#### 主从复制 master-slave
1.  常用配置   
```
主从配置
    info replication        //查看主从配置
    slaveof [ip] [port]     //salve配置，成为某个master的slave
```
#### 哨兵机制 sentinel
- 关键数据结构
```C
//每一个master实例，有一个结构体
struct sentinelRedisInstance{
    int flag;       //(SRI_MASTER|SRI_SLAVE|SRI_SENTINEL)  (SRI_S_DOWN 主观下线) (SRI_O_DOWN 客观下线)
    char *name;
    struct sentinelAddr {
        char *ip;
        int port;
    };
    
    dict *sentinels;        //记录master的slaves   
    dict *slaves;           //记录监控同一个master的sentinel
}
```
- 哨兵启动流程
    1.  启动并初始化sentinel     
        初始化master资源，并和master建立网络连接（包括发送消息和接受订阅消息两条连接）
        ![image](https://note.youdao.com/yws/res/5963/6B7818B648D246429B4224FAE0FA21DB)
    2.  向master发送INFO命令，获取主服务器信息；
    3.  根据返回INFO命令，将slaver服务器信息更新在sentinel中；
    4.  对刚刚获取到的slaver服务器建立两个网络连接，并发送INFO更新信息；
    5.  向master发送订阅消息，之后接收从主服务返回的订阅消息     
        ![image](https://note.youdao.com/yws/res/5971/BB7CFAA38CCA4E27A9470BBE267D24A5)
    6.  会有多个sentinel对redis服务器进行监听，当收到多个订阅信息后，会进行判断：
        - 如果是自己发送的消息获得的订阅信息，则进行丢弃
        - 如果是其他sentinel发送消息后收到的订阅信息，则根据订阅信息更新相关结构体
        ![image](https://note.youdao.com/yws/res/5969/B47B70B2F96B4157ACF2A34DF50B77FA)
    7.  根据订阅信息，更新sentinels字典，每个sentinel会与其他sentinel建立命令连接；    
        ![image](https://note.youdao.com/yws/res/5964/59C54416FDBF4A1285D8CC267A4B4C1C) 

- 哨兵监控流程   
    1. 主观下线
        每秒会向自己监控的所有实例(master slave sentinel)发送PING命令，超过配置时间仍未收到有效命令后，就会判定服务已经下线；
        (有效回复： +PANG   -LOADING    -MASTERDOWN)
    2. 客观下线
        当判定为主观下线后，会询问其他sentinel，看看他们是否判定主观下线；    
        监视同一master的不同sentinel会有不同的配置，其中一个判断客观下线，其他的并不一定会判断客观下线。
        询问流程：  
        - 发送消息：sentinel is-master-down-by-addr
        - 接收消息：sentinel is-master-down-by-addr命令回复： multi bulk
        - 发送消息的sentinel会统计接收到的消息判断是否客观下线      
    3.  选举领头sentinel
        当有sentinel判定客观下线时，会进行选举，选举出领头sentinel来进行后续的故障恢复
    4.  故障转移
        1.  从所有的slavers中选择一个作为master;     
            挑选状态良好、数据完整的服务器，发送 slaver no one命令，将其升级为master
        2.  让剩下的slavers作为新的master的slavers，进行复制操作；
        3.  将已经下线的master作为新的master的slaver,当其重新上线时，会重新作为slaver

#### 故障自动恢复
#### 切片集群
---

## 常见面试题
### redis为什么这么快
（操作内存、数据结构设计合理 hash表操作快、单线程无锁，不需要切换上下文、io多路复用）

### Redis性能影响 异步模式：阻塞点
- **和客户端交互的阻塞点**
1.  集合全量查询和聚合操作
        操作耗时过长，会阻塞主线程；
2.  bigkey删除操作
        删除操作会将数据清除，并在原有内存处添加链表保证操作系统对空闲内存的掌控，如果删除大量内容，此操作会耗时，导致主线程阻塞；
3.  清除数据库
        清楚数据库会删除所有的键值对；

- **和磁盘交互的阻塞点**    
IO操作一般都比较费时费力，因此在进行快照(RDB)或者写文件(AOF)时，利用子进程进行负责写入，这样磁盘IO不会阻塞主线程。
1.  AOF日志文件写入
    但是如有大量文件需要写入AOF，并需要同步写回，则会阻塞主线程
- **主从节点交互时的阻塞点**

---
## redis缓存

### 缓存类型
- **只读缓存**    
    - 缓存只做读模式，如果缓存中有数据，则从缓存中读取；如果缓存没有命中，则从后台数据库中读取数据并写入缓存。    
    - 如果需要修改后台数据，则直接操作数据库，操作完成后删除缓存中的数据。
- **读写缓存**
    - 同步独写模式    
        如果存在读操作，则修改缓存，并同步修改数据库，操作完成后返回
    - 异步独写模式        
        当发生写操作时，只修改缓存，等到缓存失效后，再将数据写回数据库

### 缓存策略

- **cache aside 旁路缓存**
策略以数据库中的数据为准，缓存中的数据是按需加载    
- 操作流程：
读缓存时先查缓存，miss后查db并写缓存。写操作时直接修改数据库，后删除缓存。因为缓存是内存操作，速度快
不可以先删缓存，再更新数据库。（事务A删除缓存后，事务B读缓存，miss后读db并更新缓存，此时事务A更新数据库，出现一致性问题）
更新数据库再删缓存也可能出现问题，很难出现请求B更新数据库并删除缓存，A才更新完缓存的情况。如果说A早于B更新缓存，而B删除缓存，此时数据还是一致的
- 缺陷、
    首次数据一定不在缓存，需要预加载；写操作多会导致频繁删缓存，最好用于读多写少的场景
- **write/read through**
缓存与后端数据库的交互细节对应用层服务隐藏，应用层直接操作缓存，缓存向DB请求数据并更新
写缓存时，cache不存在，直接更新db，缓存存在，先更新缓存，缓存服务自己更新DB
读缓存时先读cache，miss后从数据库加载
- **write back**
Read/Write Through 是同步更新 cache 和 DB，而 Write Behind Caching 则是只更新缓存，不直接更新 DB，而是改为异步批量的方式来更新 DB


### 缓存异常
- **缓存一致性**
    - 缓存和后台数据库中都有数据，且数据一致；
    - 缓存中无数据，数据库中数据为最新值
    ![image](https://note.youdao.com/yws/res/5966/E888D8532B0A4323B5339D2BD53EACFA)

- **缓存异常情况**
  ![image](https://note.youdao.com/yws/res/5967/AA21B9124A414FA0B22EEE75E3B503E1)    
    - 缓存雪崩
        现象：大量缓存同时失效，缓存服务器宕机
        解决办法：过期时间加上随机值，不要同时失效
    
    - 缓存击穿
        现象：热点数据到期，大量用户同时访问热点数据
        解决方法：热点数据不设置过期时间
    - 缓存穿透
        现象：访问不存在的数据，频繁访问增加压力
        解决方法：设置null数据
---

## 实践中遇到的问题
### 数据倾斜
![image](https://note.youdao.com/yws/res/5970/3CD3C067E4AC4070ABBA4E47C23D01D1)


## 源码
[redis源码结构](https://www.cnblogs.com/breka/articles/9914787.html)