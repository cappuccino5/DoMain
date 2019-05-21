package controller

import (
	"bytes"
	"encoding/json"
	"fmt"
	"image"
	"image/jpeg"
	"image/png"
	"io/ioutil"
	"mime/multipart"
	"net/http"
	"net/url"
	"path"
	"time"
)
//*web或者app上传文件到主服务器，主服务器将文件转存到fastdfs文件存储服务器,本done是本地测试环境*//
func ImagesTest(f []*multipart.FileHeader) {
	nowTime := time.Now().Unix()
	var lostTime int64
	var size int64
	for i := range f {
		file, err := f[i].Open()
		if err != nil {
			fmt.Println("file open fail:", err)
		}
		r, err := ImageFileTest(file, f[i])
		if err != nil {
			fmt.Println("image deal:", err)
		}
		lostTime = time.Now().Unix()
		size += f[i].Size
		fmt.Println(r)
	}
	dealTime := lostTime - nowTime
	fmt.Printf("处理时长%v 秒,总大小%v kb\n", dealTime, size)

}

// 图片存储至fastdfs文件服务器,返回url和缩略图url
func ImageFileTest(f multipart.File, header *multipart.FileHeader) (map[string]interface{}, error) {
	f1, err := header.Open()
	if err != nil {
		return nil, err
	}
	defer f1.Close()
	defer f.Close()
	width, height := ReloadThumbnailSize(header.Filename, f1)
	imageUrl, err := FileUpdateRequest(f, header)
	if err != nil {
		return nil, err
	}
	resMap := make(map[string]interface{}, 0)
	resMap["url"] = imageUrl + "?download=0"
	resMap["thumb_url"] = fmt.Sprintf("%s?download=0&width=%v&height=%v", imageUrl, width, height)
	return resMap, nil
}

// 文件存储至fastdfs文件服务器,返回url
func FileUpdateRequest(f multipart.File, header *multipart.FileHeader) (string, error) {
	var resUrl string
	data := make([]byte, header.Size)
	n, _ := f.Read(data)
	data = data[:n]
	buf := new(bytes.Buffer)
	w := multipart.NewWriter(buf)
	fw, err := w.CreateFormFile("file", header.Filename)
	if err != nil {
		return resUrl, err
	}
	fw.Write(data)
	mod := w.FormDataContentType()
	w.Close()
	res, err := HttpPostMux("http://23.94.160.30:8080/upload", buf, mod)
	if err != nil {
		return resUrl, err
	}
	resUrl = fmt.Sprintf("%s", res)
	return resUrl, nil
}

// 重新加载缩略图大小
func ReloadThumbnailSize(imageName string, f multipart.File) (int, int) {
	var width, height int
	suffix := path.Ext(imageName)
	switch suffix {
	case ".png":
		image, _ := png.Decode(f)
		width, height = resizeImage(image)
	case ".jpg", ".jpeg":
		image, _ := jpeg.Decode(f)
		width, height = resizeImage(image)
	default:
		width, height = 200, 200
	}
	return width, height
}

// 重新编辑图片大小,获取缩略图大小
func resizeImage(img image.Image) (int, int) {
	width := img.Bounds().Max.X - img.Bounds().Min.X
	height := img.Bounds().Max.Y - img.Bounds().Min.Y
	min := width
	if width > height {
		min = height
	}
	if width > 200 && height > 200 {
		width = int(float64(width) / float64(min) * 200)
		height = int(float64(height) / float64(min) * 200)
	}
	return width, height
}

// http post 表单请求
func HttpPostMux(url string, buf *bytes.Buffer, mod string) ([]byte, error) {
	req, err := http.NewRequest("POST", url, buf)
	if err != nil {
		return nil, err
	}
	req.Header.Set("Content-Type", mod)
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	return body, nil
}

//Post 请求方法
func HttpPostProxy(reqUrl string, params map[string]interface{}) ([]byte, error) {
	proxy := func(_ *http.Request) (*url.URL, error) {
		return url.Parse("http://127.0.0.1:1080")
	}
	transport := &http.Transport{Proxy: proxy}
	c := &http.Client{Transport: transport}
	query := url.Values{}
	for key, value := range params {
		query.Add(key, fmt.Sprintf("%v", value))
	}
	resp, err := c.PostForm(reqUrl, query)
	if err != nil {
		fmt.Println(err)
		return nil, err
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		fmt.Println(err)
		return nil, err
	}
	return body, nil
}
func HttpGet(getUrl string) ([]byte, error) {
	fmt.Println(getUrl)
	resp, err := http.Get(getUrl)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	return body, nil
}

func HttpPostJson(url string, data interface{}) ([]byte, error) {
	buf, err := json.Marshal(data)
	if err != nil {
		return nil, err
	}

	req, err := http.NewRequest("POST", url, bytes.NewBuffer(buf))
	if err != nil {
		return nil, err
	}
	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	return body, nil
}