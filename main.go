package main

import (
	"bytes"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"os"
	"strings"
)

//跳转到首页
func sayHello(w http.ResponseWriter, r *http.Request) {
	f, err := ioutil.ReadFile("./view/index.html")
	if err != nil {
		w.Write([]byte("打开首页失败"))
		return
	}
	w.Write(f)
}

//跳转到上传页
func upload(w http.ResponseWriter, r *http.Request) {
	//GET 请求  给表单准备上传
	if r.Method == "GET" {
		f, err := ioutil.ReadFile("./view/upload.html")
		if err != nil {
			w.Write([]byte("打开上传页面失败"))
			return
		}
		w.Write(f)
	}
	//POST请求 填写好表单点击提交后处理上传过来的图片
	if r.Method == "POST" {

		f, h, err := r.FormFile("file")
		if err != nil {
			w.Write([]byte("上传失败"))
			return
		}
		//限制上传图片
		t := h.Header.Get("Content-Type")
		// fmt.Println(t)
		if !strings.Contains(t, "image") {
			w.Write([]byte("请上传图片"))
			return
		}

		os.Mkdir("./static", 0666)
		out, err := os.Create("./static/" + h.Filename)
		if err != nil {
			w.Write([]byte("创建文件失败"))
			return
		}
		defer out.Close()

		_, err = io.Copy(out, f)
		if err != nil {
			w.Write([]byte("文件保存失败"))
			return
		}
		//服务器保存成功，跳转到相册列表页
		http.Redirect(w, r, "/list", 302)
	}
}

// 相册详情页
func iamgeView(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()
	name := r.Form.Get("name")
	// fmt.Println(name)
	f, err := os.Open("./static/" + name)
	if err != nil {
		w.WriteHeader(404)
		return
	}
	defer f.Close()
	w.Header().Set("Content-Type", "image")
	io.Copy(w, f)

}

//相册列表页
func list(w http.ResponseWriter, r *http.Request) {
	names, err := ioutil.ReadDir("./static")
	if err != nil || len(names) == 0 {
		w.Write([]byte("展示出错"))
		return
	}
	var html string
	for i := 0; i < len(names); i++ {
		fmt.Println(names[i].Name())
		html += `
		<li><a href="/detail?name=` + names[i].Name() + `">
		<img src="/image?name=` + names[i].Name() + `" alt="暂无图片" ></a> </li>
		`
	}
	f, err := ioutil.ReadFile("view/list.html")
	if err != nil {
		w.Write([]byte("打开列表页出错"))
		return
	}
	f = bytes.Replace(f, []byte("@html"), []byte(html), 1)
	w.Write(f)
}

//图片详情
func detail(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()
	name := r.Form.Get("name")
	f, err := ioutil.ReadFile("./view/detail.html")
	if err != nil {
		w.Write([]byte("图片打开失败"))
	}
	f = bytes.Replace(f, []byte("@src"), []byte("/image?name="+name), 1)
	w.Write(f)
}

func main() {
	//路由
	http.HandleFunc("/", sayHello)
	http.HandleFunc("/index", sayHello)
	http.HandleFunc("/upload", upload)
	http.HandleFunc("/image", iamgeView)
	http.HandleFunc("/list", list)
	http.HandleFunc("/detail", detail)

	err := http.ListenAndServe("127.0.0.1:8080", nil)

	if err != nil {
		fmt.Println("http server failed:", err)
		return
	}

}
