package main

import (
	"github.com/kataras/iris/v12"
	"io"
	"os"
)

const maxSize = 5 << 20 // 5MB
var tpl = `
<html>
<head>
    <title>Upload file</title>
</head>
<body>
    <form enctype="multipart/form-data" action="http://127.0.0.1:8080/upload" method="POST">
        <input type="file" name="uploadfile" />

        <input type="hidden" name="token" value="{{.}}" />

        <input type="submit" value="upload" />
    </form>
</body>
</html>
`

func main() {
	app := iris.New()
	app.Get("/upload", func(ctx iris.Context) {
		ctx.HTML(tpl)
	})
	//处理来自upload_form.html的请求数据处理
	//用iris.LimitRequestBodySize的好处是：在上传过程中，若检测到文件大小超过限制，就会立即切断与客户端的连接。
	app.Post("/upload", iris.LimitRequestBodySize(maxSize+1<<20), func(ctx iris.Context) {
		// Get the file from the request.
		file, info, err := ctx.FormFile("uploadfile")
		if err != nil {
			ctx.StatusCode(iris.StatusInternalServerError)
			ctx.HTML("Error while uploading: <b>" + err.Error() + "</b>")
			return
		}
		defer file.Close()
		fname := info.Filename
		//创建一个具有相同名称的文件
		//假设你有一个名为'uploads'的文件夹
		out, err := os.OpenFile("./uploads/"+fname, os.O_WRONLY|os.O_CREATE, 0666)
		if err != nil {
			ctx.StatusCode(iris.StatusInternalServerError)
			ctx.HTML("Error while uploading: <b>" + err.Error() + "</b>")
			return
		}
		defer out.Close()
		io.Copy(out, file)
			ctx.HTML(tpl)
		ctx.WriteString("上传成功！")
	})
	//在http//localhost:8080启动服务器，上传限制为5MB。
	app.Run(iris.Addr(":8080") /* 0.*/, iris.WithPostMaxMemory(maxSize))
}
