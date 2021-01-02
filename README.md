# img-cdn-go

服务器中越来越多的项目需要使用图片上传和下载功能，所以开发一个专门负责图片上传下载的API

## 项目架构
<del>采用iris框架，分为两个部分：上传图片和下载图片<del>

不用本地了，OSS+CDN走起！

由于最近比较忙，所以API暂不公开使用

***

## API

**上传图片**
`https://upload.ri-co.cn/img/upload`

**返回值** 图片的完整URL，例如：

`https://cdn.ri-co.cn/img/754e0dc8-13a9-4aac-b0cf-0b8015085fc5rpionfish-Overwatch-Hentai.jpg`

![实例图片](https://cdn.ri-co.cn/img/754e0dc8-13a9-4aac-b0cf-0b8015085fc5rpionfish-Overwatch-Hentai.jpg?x-oss-process=style/middle_cut)

