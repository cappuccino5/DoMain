package fastdfs

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"go-hisens/commons"
	"image"
	"image/jpeg"
	"image/png"
	"io/ioutil"
	"mime/multipart"
	"net/http"
	"path"
)

type FileServer struct {
	Url   string
	Scene string
}

var con = commons.AppConfig

func initServer() *FileServer {
	return &FileServer{
		Url:   con.FileServer + con.GroupName + "/upload",
		Scene: con.Scene,
	}
}

var fastdfs *FileServer

func init() {
	fastdfs = initServer()
}

//*web或者app上传文件到主服务器，主服务器将文件转存到fastdfs文件存储服务器,本done是本地测试环境*//
// rootPath = 保存的路径名字 fromName = 表单key
func UploadFile(FileDir string, f *multipart.FileHeader) (map[string]interface{}, error) {
	if FileDir == "" {
		return nil, errors.New("param is empty")
	}
	imageUrl, err := fastdfs.fastdfsRequest(f, FileDir)
	if err != nil {
		return nil, err
	}
	resMap := make(map[string]interface{}, 0)
	resMap["url"] = imageUrl
	if path.Ext(f.Filename) != ".amr" {
		width, height := tailorImage(f)
		resMap["url"] = imageUrl + "?download=0"
		resMap["thumb_url"] = fmt.Sprintf("%s?download=0&width=%v&height=%v", imageUrl, width, height)
	}
	return resMap, nil
}

// 自适应图片大小，缩略图
func tailorImage(header *multipart.FileHeader) (int, int) {
	var width, height int
	f, _ := header.Open()
	defer f.Close()
	suffix := path.Ext(header.Filename)
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

// 文件存储至fastdfs文件服务器,返回url
func (ser *FileServer) fastdfsRequest(header *multipart.FileHeader, showDir string) (string, error) {
	var resUrl string
	file, err := header.Open()
	if err != nil {
		return resUrl, err
	}
	defer file.Close()
	tempbuf := make([]byte, header.Size+1)
	n, _ := file.Read(tempbuf)
	tempbuf = tempbuf[:n]
	filebuf := new(bytes.Buffer)
	w := multipart.NewWriter(filebuf)
	w.WriteField("scene", ser.Scene)
	w.WriteField("file_dir", showDir)
	w.WriteField("output", "json")
	fw, err := w.CreateFormFile("file", header.Filename)
	fw.Write(tempbuf)
	methods := w.FormDataContentType()
	w.Close()
	res, err := HttpPostMux(ser.Url, filebuf, methods)
	if err != nil {
		return resUrl, err
	}
	fileInfo := make(map[string]interface{}, 0)
	err = json.Unmarshal(res, &fileInfo)
	if err != nil {
		return resUrl, err
	}
	if size, ok := fileInfo["size"].(float64); !ok || size != float64(header.Size) {
		return resUrl, errors.New("file size error")
	}
	resUrl, _ = fileInfo["url"].(string)
	return resUrl, nil
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
func HttpPostMux(posturl string, buf *bytes.Buffer, mod string) ([]byte, error) {
	fmt.Println("fastdfs req:", posturl)
	req, err := http.NewRequest("POST", posturl, buf)
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
