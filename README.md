#read_comp
read_comp is simple uncompress file command in Go .

read_comp uncompress .tar.gz files.

#Installation

```
> git clone https://github.com/akms/read_comp
> cd read_comp
> go install
```

#Example

```
> read_comp /home/user/file.tar.gz target
file/target/hoge
file/target/fuga
file/hogehoge/target

> read_comp ../../file.tar.gz target
file/target/hoge
file/target/fuga
file/hogehoge/target
```
