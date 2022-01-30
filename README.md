# GoWolf

---

## 前言

闲来无事，造个轮子，方便日后使用~😀

## 功能

* ICMP协议存活扫描
* 单端口扫描
* 多端口扫描

> 更新预告，优化单端口扫描消耗资源，arp协议存活扫描

## 参数详解

> -a &ensp;&ensp;目标地址
>
> -f &ensp;&ensp;&ensp;目标地址文件
>
> -p &ensp;&ensp;指定扫描的的口，多个端口需要使用`-`隔开，例子`1-100`,默认参数`1-100`端口
>
> -t &ensp;&ensp;设定go程数，默认10
>
> -J &ensp;&ensp;设置工作区缓存数，默认200
>
> -O &ensp;&ensp;设置完成区缓存数，默认200
>
> -i &ensp;&ensp;使用ICMP协议存活扫描，开启0，默认1关闭

## 安装

> sudo apt install golang

> git clone https://github.com/New-arkssac/GoWolf.git

> go build main.go

## 演示

> 单端口扫描

![image.png](https://pwl.stackoverflow.wiki/2022/01/image-bf07f145.png)

> 多端口扫描

![image.png](https://pwl.stackoverflow.wiki/2022/01/image-46a49c9d.png)

> 多地址扫描

![image.png](https://pwl.stackoverflow.wiki/2022/01/image-9d7b00ad.png)

> ICMP协议存活扫描

![image.png](https://pwl.stackoverflow.wiki/2022/01/image-2e595d94.png)

## 注释

因为会使用到原始自定义协议包，所以部分功能无法在`windows`和`windows wsl`上运行，所以建议在linux环境下使用`GoWolf`