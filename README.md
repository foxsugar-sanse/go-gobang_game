# go-gobang-game

#### 介绍

>  使用go编写的基于控制台的五子棋程序，支持局域网对战

#### TODO

> 1.将配置文件中的ip和port修改

#### build

```bash
bash make.sh
```

#### run

```bash
cd ./_build
./gobang
```

#### docker build

```bash
bash make.sh
cp Dockerfile.dockerfile ./_build
cd ./_build
docker build -t gobang:0.2.2 .
```


