# components
主要存放和配置无相关的组件，随意拷贝至此框架的任意一个项目都能正常运转

### recover.go

```
func Go(params *map[string]interface{}, callback func(*map[string]interface{}))
主要解决goroutines的如下问题：
1.goroutines panic后会终止主进程而退出。(封装了recover, 保证程序不会因为程序内部问题而退出，除非接收到signal。)
2.由于主进程退出时，goroutines可能还未运行完。(封装了注册全局信号量机制，就是接收到信号量，也需要等goroutines全部完成)(同时也需要长时间运行的goroutines,需要自己捕获信号量并主动退出)
```
