[toc]
# Kubernetes




为什么要存在pod（不可以直接用container吗）

控制器模型（保证集群中的容器始终处于一个期望的状态，可以实现自动）
Deployment
StatefulSet 
DaemonSet
创建并运行在每个Node上，且每个Node只有一个这样的Pod，节点加入集群后会自动创建pod，离开集群后自动删除。
控制器从etcd中遍历Node，来维持节点上Pod的数目
Service 将一组pod暴露给外界访问的一种机制，实现主要是通过kube-procxy and iptables
为什么需要service：pod的ip不固定、pod实例之间需要负载均衡
有一个cluster-ip，通过访问这个固定的ip，会将请求负载均衡的访问到自己endport列表上的pod。只能保证集群内部访问
外部访问需要nodeport
Job&CronJob【离线业务】
RABC【基于角色的权限控制】
operator


## 学习文档
https://www.orchome.com/6766

### 框架结构
#### master节点组件
master为集群的控制节点，负责整个集群的管理和控制，kubectl将命令发给master节点，之后master节点管理node节点，可以配置HA，master节点主要包括以下组件：     
- kube-apiserver
提供了http restful接口服务，k8s所有资源的增删改查的唯一入口，也是集群控制的入口进程

- control-manager
资源对象的自动化控制中心

- shcedule
负责调度资源到合适的Node

- etcd
保存资源对象的所有配置数据

- kubectl
命令行管理工具


#### node节点组件
node节点为工作负载节点，运行master分配过来的资源，当某个node宕机时，其资源会呗分配给其他node，node节点主要包括以下组件：    
- kubelet
负责pod对应的容易创建、启停等任务，同时和master合作，实现集群的基本管理功能。Node节点可以动态增加至k8s集群，kubelet会定时向master汇报自己的状况，如资源运行情况、cpu和内存等，保障master清楚node的资源状况，方便资源调度。
如果node不可用，master会触发**工作负载大转移**工作。

- kube-proxy
实现k8s service的通信和负载均衡

- dokeer engine
负责本机的容器创建和管理工作    

-----
### 功能组件
#### pod
- 基础知识
- 技术点：
    - 为什么要存在pod，不能直接使用container?
        - 一组容器难以管理，如果其中一个容器挂了，无法判断整个服务的状态。抽象为pod后，一个容器挂了，则整个服务都挂了；
        - pod 里的所有容器，共享的是同一个 Network Namespace，并且可以声明共享同一个Volume，方便密切关联的容器之间的内部通信；
        - 创建建pod时首先创建基础容pause，该容器创建namespace，后续pod中的container加入该namespace
    
    - pod资源配置
        每个资源都需要配置两个参数：request(最小申请量，系统必须满足)、limits(最大使用量，超过后容器会被kill)
        - cpu
        定义100m，表示0.1个CPU
        - memory
        定义64Mi，表示64MiB

#### Replication Controller
会指定期待的pod数量，当实际pod数量不一致时，会动态调整pod的数量     
提交一个RC资源到集群后，controller manager得到通知，定期巡检pod数目，保证和期望值相同；    
1.2版本后升级为replic set，增加了基于集合的lable selector功能， RC只支持基于等式的lable selector

#### replicSet

#### Deployment
- 无状态管理

- 对replicSet的一种升级，可以知道当前pod的部署进度
    因为一个RS的创建包括pod的创建、调度、绑定节点、node上运行等漫长流程，需要清楚所有pod的部署进度。

- 工作原理
    控制replicaSet对象来实现版本管理，通过replicaSet控制pod数目，每一个版本都有一个replicaSet，实现版本升级和回退，水平收缩和扩容，滚动更新

#### Horizontal Pod Autoscaler
制定规则，让k8s自动管理pod数量，实现水平扩容，指标如下：     
- CPU Utilization Percentage
- 应用程序自定义的度量指标，比如服务在每秒内的相应请求数


#### statefulSet
- 有状态集群特点
    - 每个节点有固定身份id，集群中内部需要互相通信
    - 集群规模固定
    - 节点有状态，数据会被持久化
    - 如果磁盘损坏，某个节点无法正常运行，集群功能会受到影响

【有状态】   
无状态认为所有的pod都是一样的，没有顺序，也无所谓运行在哪个主机，可以根据需要杀掉或者创建pod。但是在实际生活中会存在主备、主从关系。或者说是数据储存类的pod，重启后就和原来的数据失去了关系。这种实例之间存在关系，对外部数据有依赖关系的就是有状态。
有状态主要体现在两个方面：**拓扑状态**、**储存状态**。因为在实际集群中会存在主从结构、主备关系。核心功能就是：通过某种方式记录这两种状态，pod被新建后能够恢复原有的状态。
- 拓扑状态
    意味着容器的启动和停止之间是存在顺序的，启动时会在容器名称后边编号【容器名称+编号】，必须按照一定的顺序启动容器。此外还给每个pod提供一个固定且唯一的访问入口（DNS记录，即使ip改变，pod仍然能访问到）
- 储存状态
    多个实例分别绑定了不同的数据，podA第一次读到的数据和十分钟后读到的应该是同一份数据。删除一个pod后，远端储存不会删除，当新版本pod升级启动后，和原始pod名称相同，因此还能访问原始的PVC             

statusfulSet控制器直接管理的是pod，会按顺序启动并命名，且维护DNS记录。还会为每个pod创建一个同样编号的PVC，重启后新pod可以接管PVC.

- Headless Service
    和普通service的区别是没有cluster ip，解析其DNS时，返回的是对应全部的pod的endpoint列表。
    为每个pod建立了DNS域名  $(podname).$(headless service name)

#### service
service包含四种类型：   
- ClusterIP    
    默认类型，只支持集群内部访问, 该ip无法ping，因为不存在实体来回包

- NodePort    
    为每个node绑定一个端口，访问内部服务  

- LoadBalancer     
    在NodePort的基础上，借助CloudProvider创建一个外部负载均衡器，并将请求转发到NodePort
    loadbalancer独立于集群之外

- ExternalName     
    把集群外部的服务引入到集群内部来，在集群内部直接使用。没有任何类型代理被创建，这只有 Kubernetes 1.7或更高版本的kube-dns才支持。

#### daemonSet
#### Job

---

### 网络插件
k8s对于网络插件的要求：
- 所有pod在不适用nat的情况下与其他pod通信
- 所有节点都可以在没有nat的情况下与所有pod通信
- pod自己的ip和其他pod看到的ip是一样的

#### 容器到容器的网络
每个pod在最初创立时，会创建一个pause，该pause会创建namespace，后续pod内的全部容器会加入该namespace，可以通过localhost进行通信     
Pod 中的容器还可以访问共享卷，这些卷被定义为 Pod 的一部分，并且可以挂载到每个容器的文件系统中     

#### pod到pod的网络
在k8s集群中，每个pod都拥有自己的ip，因此可以通过ip实现pod和pod之间的通信。

##### 同节点pod通信
##### 不同节点pod通信

#### pod到service的网络
#### 互联网到service的网络

---
### 持久化储存
#### volume
被定义在pod上，属于计算资源的一部分
- enptyDir
    主要用于临时目录，无需永久保存数据

- hostPath
    pod所挂载的宿主机的目录，

- gcePersistentDisk
- awsElasticBlockStore
- NFS
- ......

#### Presistent Volume
属于独立于计算资源的实体资源，理解为集群中的某个网络存储


---
### 服务发现
---
## 故障处理
1. [通过yaml创建pod，状态一直是ContainerCreating](https://www.cnblogs.com/randy-lo/p/13321148.html)
2. 


## 搭建集群
### 环境配置
- windows安装虚拟机   
```
ifconfig    yum install net-tools -y 
export PS1='[\u@\h $PWD]$'
```
- centos配置
```
关闭防火墙
systemctl stop firewalld; systemctl disable firewalld

关闭selinux 重启后生效
    临时：setenforce 0
    永久：sed -i 's/SELINUX=enforcing/SELINUX=disabled/g' /etc/selinux/config
    reboot

关闭交换分区:防止内存交换，提升性能
    临时：swapoff -a
    永久：sed -ri 's/.*swap.*/#&/' /etc/fstab
    
修改内核参数
    #将桥接的IPV4流量传递到iptables的链
    cat > /etc/sysctl.d/k8s.conf << EOF 
    net.bridge.bridge-nf-call-ip6tables = 1 
    net.bridge.bridge-nf-call-iptables = 1 
    EOF
    sysctl --system  #生效
```
### 安装master node节点
[安装kubernetes](https://www.cnblogs.com/spll/p/10033316.html)   
[centos7二进制安装kubernetes](https://www.58jb.com/html/180.html)

#### kubeadm安装
- 系统基本配置
- 所有节点安装docker/kubeadm/kubelet
```
wget https://mirrors.aliyun.com/docker-ce/linux/centos/docker-ce.repo -O /etc/yum.repos.d/docker-ce.repo
yum -y install docker-ce
systemctl enable docker && systemctl start docker


cat > /etc/docker/daemon.json <<EOF
{
  "registry-mirrors": ["https://b9pmyelo.mirror.aliyuncs.com"]
}
EOF
systemctl restart docker
docker info

# 可以指定版本
yum install docker-ce-19.03.8 -y
```
- 安装kubeadm/kubelet/kubectl
```
cat > /etc/yum.repos.d/kubernetes.repo <<EOF
[kubernetes]
name=Kubernetes
baseurl=https://mirrors.aliyun.com/kubernetes/yum/repos/kubernetes-el7-x86_64
enabled=1
gpgcheck=0
repo_gpgcheck=0
gpgkey=https://mirrors.aliyun.com/kubernetes/yum/doc/yum-key.gpg https://mirrors.aliyun.com/kubernetes/yum/doc/rpm-package-key.gpg
EOF


yum install -y kubelet-1.18.0 kubeadm-1.18.0 kubectl-1.18.0
systemctl enable kubelet
```

- master 节点执行
```
kubeadm init \
  --apiserver-advertise-address=192.168.52.100 \
  --image-repository registry.aliyuncs.com/google_containers \
  --kubernetes-version v1.18.0 \
  --service-cidr=10.96.0.0/12 \
  --pod-network-cidr=10.244.0.0/16 \
  --ignore-preflight-errors=all

# 等待执行完约5-10分钟  
# 设置开机启动
systemctl enable kubelet.service

# 只能命名为./kube/config
mkdir -p $HOME/.kube
sudo cp -i /etc/kubernetes/admin.conf $HOME/.kube/config
sudo chown $(id -u):$(id -g) $HOME/.kube/config
$ kubectl get nodes
NAME         STATUS   ROLES    AGE   VERSION
k8s-master   Ready    master   2m   v1.18.0


安装CNI-calico
wget https://docs.projectcalico.org/manifests/calico.yaml --no-check-certificate
kubectl apply -f calico.yaml 
```
- node加入集群
```
master /etc/kubernetes/token.csv查看token

# node节点执行
kubeadm join 192.168.52.100：6443 --token <token> --discovery-token-ca-cert-hash sha256:<hash>

## 生成token
kubeadm token list
kubeadm token create

openssl x509 -pubkey -in /etc/kubernetes/pki/ca.crt | openssl rsa -pubin -outform der 2>/dev/null | \
   openssl dgst -sha256 -hex | sed 's/^.* //'
```