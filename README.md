# IM即时通信系统
- 1.九个小版本迭代开发 ，有读写分离模型特点的一个项目。
- 2.基于 TCP 协议的聊天室 ，并实现了聊天室中用户的上下线、消息的广播、私聊公聊等功能 。
- 3.server文件两个最主要的两个数据结构 ，一个是OnlineMap是用来记录当前在线用户 ，一个是channel ，用来进行广播 ，接受传过来的消息 ；然后发给对应的user类所绑定的channel来进行通信 ，user文件里有两个GO程一直跑 ，一个是永
久阻塞等待的channel ，一直在读server端的消息 ，有消息就回写给对应的client ，一个是handlerGO ，永远阻塞等待client
发来的消息。
