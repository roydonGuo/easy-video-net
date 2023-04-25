# easy-video-net

#### **项目介绍** : easy-video-net 是基于 golang + vue 开发的前后端分离项目,server端采用 golang + gin + gorm 进行开发 web 端采用vue3 + typescript + element-plus 进行开发

**前言** 

项目为自己在学习golang和vue3时进行编写的学习项目可能依然存在许多问题，欢迎提lssuse，后续可能依然进行迭代，如果对你有帮助的话希望可以得到您的star

体验地址 [easy-video](http://124.220.20.83:9999/)

**主要功能模块**

- **视频上传分享支持转码及弹幕功能**
- **稿件投稿使用富文本编辑器进行简单发布**
- **一个简单直播功能 需要使用livego搭建直播服务**
- **简单的消息通知 及其im功能**
- **个人相关及其相关发布信息的CRUD**


#### 项目环境 

**server**

-  golang  1.18
-  mysql  8.0
-  reids  3.0
-  ffmpeg  4.2

**web**

-  npm 8.18
-  node v16.16

**项目特点**

- 完成上传分接口类型使用不用存储，不同质量，实现本地极其阿里云oss存储 支持分片上传，断点续传 oss 直传
- 视频本地存储使用ffnpeg进行视频转码 , oss使用阿里云智能媒体转码
- 简单实现直播功能并且采用protobuf进行通信
- 支持docker部署

**项目目录**

```
easy-vide-net
│  .gitignore 
│  README.md  
├─service 服务端代码
│  │  .gitignore
│  │  go.mod
│  │  go.sum
│  │  main.go   
│  ├─assets 静态资源
│  ├─config 配置文件
│  │  │  config.ini
│  ├─consts 常量定义
│  ├─controllers 控制器
│  ├─global 全局使用
│  ├─interaction 结构体定义
│  │  ├─receive   请求  
│  │  └─response   响应
│  ├─logic  逻辑处理
│  ├─middlewares 中间件
│  ├─models  模型定义
│  ├─proto   proto文件
│  ├─router   路由定义
│  ├─runtime 运行时文件如日志
│  └─utils 工具类
│              
└─web
    │  .env  配置文件
    │  .gitignore
    │  .hintrc
    │  auto-imports.d.ts
    │  components.d.ts
    │  index.html
    │  package-lock.json
    │  package.json
    │  README.md
    │  tsconfig.json
    │  tsconfig.node.json
    │  vite.config.ts
    │  
    ├─node_modules
    ├─public
    │      vite.svg  
    └─src
        │  App.vue 
        │  main.ts
        │  shime-vue.d.ts
        │  style.css
        │  style.scss
        │  vite-env.d.ts
        │   
        ├─apis 接口定义
        ├─assets 静态资源
        ├─components 组件
        ├─hooks 
        ├─logic 逻辑代码
        ├─proto proto文件
        ├─router 路由定义
        ├─store 状态管理
        ├─types 类型定义
        ├─utils 工具类
        └─views 视图文件
```

**server端启动**

```
//填写完成项目config文件夹内config.ini
go mod tidy 安装所需依赖
go build  打包项目
./easy-video-net.exe 运行项目
```

**web端启动**

```
//填写完成项目env文件配置请求地址
npm i 安装所需依赖
npm run dev 运行项目
```

#### **项目展示**

![image-20230419151645612](https://user-images.githubusercontent.com/64412088/233002215-359b2337-6224-4318-811c-b2195f3cef4a.png)
![image-20230419151937238](https://user-images.githubusercontent.com/64412088/233002263-ff599b43-00c7-4d9a-8caf-2797500b1787.png)
![image-20230419151819626](https://user-images.githubusercontent.com/64412088/233002291-0ff90253-5e13-4240-9d89-43fff9e455b5.png)
![image-20230419151958183](https://user-images.githubusercontent.com/64412088/233002317-6bb54307-b696-48a7-9f73-4a24bdd65261.png)
![image-20230419152335308](https://user-images.githubusercontent.com/64412088/233002344-96b837f1-8174-4d21-9bb7-5d1ea4fad625.png)
![image-20230419151844222](https://user-images.githubusercontent.com/64412088/233002384-374e5375-dad6-4516-9a45-2466ad63d1bb.png)
