### 安装连码

> 注：根据kubectl版本选择使用v1还是v2, 测试集群以及正式环境使用v1, 自建minikube使用v2


```bash
./setupcc.sh v2 init
```

### 升级链码

升级前请务必修改`setupcc.sh`文件中`CHAINCODE_VERSION`变量值
```bash
./setupcc.sh v2 upgrade
```

### 链码设计


### 流程图

