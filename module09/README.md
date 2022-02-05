####  总结几种 socket 粘包的解包方式：fix length/delimiter based/length field based frame decoder。尝试举例其应用。

##### 黏包原因
```
TCP是面向连接的传输协议，TCP传输的数据是以流的形式，而流数据是没有明确的开始结尾边界，所以TCP也没办法判断哪一段流属于一条消息。
粘包产生的主要原因：
    发送方每次发送的数据 < socket缓冲区大小
    接收方读取socket缓冲区数据不够及时

```
##### fix length
 - 发送端和接收端规定固定大小的缓冲区，当字符长度不够时使用空字符弥补
 - fix_length.go

##### delimiter based
 - 使用某几个特殊字符组合作为分隔符，每次读取的时候读到分隔符就停止，然后解析包。比较适合传递长度不固定、文本格式的数据， 分隔符可以用文本中不会出现的字符。
 - delimiter_based.go

##### length field based frame decoder
 - 给数据包增加一个消息头，里面包含了包的长度信息。读取的时候先读取消息头（可以是固定长度，比如4个字节）， 然后再按这个长度读取消息体。这种方式可以传递几乎所有信息
 - delimiter_based.go