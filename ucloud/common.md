[toc]

# 常用笔记

## topid-oid

"request_uuid": "{{$guid}}",

Test03: top-60036,     oid:1260006

现网：topid-56028911    oid-64203518

打包rpm: 10.69.52.13

chaos: 10.69.202.211

## urouter-api灰度节点

```
# 内部
172.18.211.253:5100
10.68.132.250:5100

# 外部
172.18.211.252:4200
10.68.132.249:4200
```

## 公网镜像向test03发镜像

```
1. 你提供一个镜像所在区域的us3的bucket地址、公钥/私钥，我把镜像拷贝到你的us3
2. 你把镜像从us3下载到你本地
3. 你把镜像scp到test03的跳板机
4. 你把镜像scp到root@10.72.164.118:/data/zhaokai下
5. 以上都做完后，给我提供需要上传的镜像的公司id和项目id，我把镜像给你传到test环境
```

## 服务部署

### uimage-access



#### docker部署

```
version=udc-ant-international_v1_cmp
instance=10.72.136.169
zk=10.72.137.139:2181,10.72.137.140:2181,10.72.137.85:2181
zone=666003001

docker login hub.ucloudadmin.com -u uhost-api -p 2h8e9rqpokc9d5ar27i60davw38s2q1r

if [ "$access" != "0" ]; then
  docker stop uimage-access
  docker rm -f uimage-access
fi
docker run -d --restart=always \
  --name uimage-access \
  --network host \
  -v /data/image-go/access/:/data/image-go/access/ \
  -v /etc/localtime:/etc/localtime:ro \
  hub.ucloudadmin.com/sre-uhost/uimage:$version uimage_access --zk_address=$zk --zone_id=$zone

if [ "$proxy" != "0" ]; then
  docker stop uimage-proxy
  docker rm -f uimage-proxy
fi
docker run -d --restart=always \
  --name uimage-proxy \
  --network host \
  -v /data/image-go/proxy/:/data/image-go/proxy/ \
  -v /etc/localtime:/etc/localtime:ro \
  hub.ucloudadmin.com/sre-uhost/uimage:$version uimage_proxy --zk_address=$zk --instance_id=$instance
```

#### ally部署

```
version=udc-ant-international_v1_cmp
instance=10.72.136.169
zk=10.72.137.139:2181,10.72.137.140:2181,10.72.137.85:2181
zone=666003001

ally invite uimage-access-go \
	--app-bin /data/image-go/bin/access/access \
	-- --zk_address=$zk \
	--zone_id=$zone


ally invite uimage-proxy-go \
	--app-bin /data/image-go/bin/proxy/proxy \
	-- --zk_address=$zk \
	--instance_id=$instance
```

## humming发布

```
md5sum    humming-1.2.26-1.ad3158e.x86_64.rpm
99ca67102a313a40bbd7444cdbde1ddc  humming-1.2.26-1.ad3158e.x86_64.rpm

# 更改tag 和 md5

curl --location --request POST 'http://gw.iaas.ucloudadmin.com/terra/api/version' \
--header 'auth-key: eyJ0eXAiOiJKV1QiLCJhbGciOiJIUzI1NiJ9.eyJleHAiOjE3NzM4MTg4NjcsIm5iZiI6MTY3ODc3ODg0NywiaWF0IjoxNjc4Nzc4ODU3LCJzdWIiOiJteSB0b2tlbiIsImRhdGEiOnsiaWQiOjE0OTgsInVzZXJuYW1lIjoidWhvc3QtZGV2Iiwibmlja25hbWUiOiJ1aG9zdFx1NzgxNFx1NTNkMVx1NGUxM1x1NzUyOFx1OGQyNlx1NTNmNyIsInN1cGVydXNlciI6ZmFsc2V9fQ.puCiomnurO7FXX9ZDomrk9iayc07-iZByGFgZ30C_Pk' \
--header 'Content-Type: application/json' \
--data '    {
    "service_id": "service-d7b2171b",
    "version_name": "v1.2.26",
    "tag": "1.2.26",
    "md5sum": "99ca67102a313a40bbd7444cdbde1ddc",
    "is_default": 0
}'
```



## nbd挂载

```
# humming
koala_ip=10.66.170.245
koala_port=6000
nbd=/dev/nbd10
ubs_id=bsi-7c83832
/opt/humming/nbd-client-x86 $koala_ip $koala_port $nbd -b 512 -p -v $ubs_id -t 100


# 宿主机

```

## test03部署humming

```
# 部署master
ally invite humming-master-instance-1 --app-bin /data/zhaokai/humming-1/humming -- master --addr 10.72.164.118 --mongo mongodb://10.72.136.142:27017,10.72.137.78:27017,10.72.137.90:27017 --zookeeper 10.72.137.139:2181,10.72.137.140:2181,10.72.137.85:2181 --log /data/zhaokai/humming-1/master.log -i 1 -p 10612 -g 10613


# 部署agent
ally invite humming-agent-instance-1 --app-bin /data/zhaokai/humming-1/humming -- agent --addr 10.72.164.118 --zone 666003001 --zookeeper 10.72.137.139:2181,10.72.137.140:2181,10.72.137.85:2181 --log /data/zhaokai/humming-1/agent.log  -p 10614 -g 10615 -i 1 

```

## ansible使用

```

ansible all -i "2003:da8:2004:1000:0a3b:b013:7f00:52de," -e "ansible_user=ubuntu ansible_password=SvxrU4HPBO" -m raw -a "sudo apt install -y getenforce"


cmd="free"
ansible all -i "106.75.163.157," -e "ansible_user=root ansible_password=calmxkk111" -m raw -a free


2003:da8:2004:1000:0a3b:1718:7f00:52de

# test03
cmd="getenforce"
ansible all -i "2003:da8:2004:1000:0a3b:1718:7f00:52de," -e "ansible_user=root ansible_password=calmxkk111" -m raw -a "$cmd" -vvv

117.50.91.169
```
