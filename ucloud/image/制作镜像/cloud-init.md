[toc]

# 镜像安装cloud-init踩坑合集

## ubuntu1604

使用2004打的包安装1604:

https://uhost-package.cn-bj.ufileos.com/cloud-init/Ubuntu20.04/cloud-init_24.1.3-1ubuntu1~20.04.5ucloud_all.deb

```
# 卸载历史的cloudinit
cloud-init clean
apt-get remove --purge -y cloud-init
apt-get autoremove -y
rm -rf /etc/cloud /var/lib/cloud /var/log/cloud-init*


# 安装
apt-get install -y cloud-guest-utils ethtool 
wget https://uhost-package.cn-bj.ufileos.com/cloud-init/Ubuntu20.04/cloud-init_24.1.3-1ubuntu1~20.04.5ucloud_all.deb
apt-get install -y ./cloud-init_24.1.3-0ubuntu3.3ucloud_all.deb
rm -f cloud-init_24.1.3-0ubuntu3.3ucloud_all.deb
apt-mark hold cloud-init
apt-get install -y tzdata-legacy
```



