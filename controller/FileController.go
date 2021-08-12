package controller

import (
	"encoding/json"
	"file/dto"
	"file/lib/config"
	"file/lib/util"
	"file/service"
	"fmt"
	"github.com/aliyun/aliyun-oss-go-sdk/oss"
	"github.com/julienschmidt/httprouter"
	"html/template"
	"io"
	"io/ioutil"
	"net/http"
	"net/url"
	"os"
	"path"
	"strings"
	"sync"
)

type FileController struct {
	fileService service.IFileService
	path        string
	client      *oss.Client
}

func (c *FileController) index(w http.ResponseWriter, _ *http.Request, _ httprouter.Params) {
	tel, _ := template.ParseFiles(
		"./resources/view/upload.html",
		"./resources/view/template/head.html",
		"./resources/view/template/footer.html",
	)
	dataMap := make(map[string]interface{})
	_ = tel.Execute(w, dataMap)
}

func (c *FileController) info(w http.ResponseWriter, _ *http.Request, ps httprouter.Params) {
	fileId := ps.ByName("fileId")
	if fileId == "" {
		message := util.NewResponseFailMessage()
		message.Message = "fileId 不可为空"
		ret, _ := json.Marshal(message)
		_, _ = w.Write(ret)
		return
	}

	file, err := c.fileService.SelectFileById(fileId)
	if err != nil || file.FileId == "" {
		message := util.NewResponseFailMessage()
		message.Message = "查询文件失败，请重试"
		ret, _ := json.Marshal(message)
		_, _ = w.Write(ret)
		return
	}

	message := util.NewResponseSuccessMessage()
	message.Data = file
	ret, _ := json.Marshal(message)
	_, _ = w.Write(ret)
}

func (c *FileController) pageList(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	var page dto.FilePage
	_ = json.NewDecoder(r.Body).Decode(&page)
	_page, err := c.fileService.SelectPageList(page)
	message := util.NewResponseMessage()
	if err != nil {
		message = util.NewResponseFailMessage()
	}
	message.Data = _page
	ret, _ := json.Marshal(message)
	_, _ = w.Write(ret)
}

func (c *FileController) list(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	values := r.URL.Query()
	bucketName := values.Get("bucket")
	list, err := c.fileService.SelectAllList(bucketName)
	message := util.NewResponseMessage()
	if err != nil {
		message = util.NewResponseFailMessage()
	}
	message.Data = list
	ret, _ := json.Marshal(message)
	_, _ = w.Write(ret)
}

func (c *FileController) upload(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	_ = r.ParseForm()

	// 1.从请求头读出文件流
	_ = r.ParseMultipartForm(32 << 20)
	bucketName := r.FormValue("bucket")
	prefix := r.FormValue("prefix")
	if bucketName == "" {
		message := util.NewResponseFailMessage()
		message.Message = "bucket 不能为空"
		ret, _ := json.Marshal(message)
		_, _ = w.Write(ret)
		return
	}

	// 判断bucket是否存在
	lsRes, err := c.client.ListBuckets()
	if err != nil {
		message := util.NewResponseFailMessage()
		message.Message = err
		ret, _ := json.Marshal(message)
		_, _ = w.Write(ret)
		return
	}
	isBucket := false
	for _, bucketItem := range lsRes.Buckets {
		if bucketItem.Name == bucketName {
			isBucket = true
		}
	}
	if isBucket == false {
		message := util.NewResponseFailMessage()
		message.Message = "你使用的bucket不存在"
		ret, _ := json.Marshal(message)
		_, _ = w.Write(ret)
		return
	}

	// 读取文件流
	file, handler, err := r.FormFile("file")
	if err != nil {
		message := util.NewResponseFailMessage()
		message.Message = err.Error()
		ret, _ := json.Marshal(message)
		_, _ = w.Write(ret)
		return
	}
	_ = file.Close()
	// 写入文件
	//fmt.Fprintf(w, "%v", handler.Header)
	filename := handler.Filename
	//filename := path.Base(filename)
	fileSuffix := path.Ext(filename)
	//filePrefix := filename[0:len(filename) - len(fileSuffix)]
	fileId := util.NewUUID()

	// 2.讲请求流写入到本地
	directory := "./files/" // 存放文件目录
	_, err = os.Stat(directory)
	if err != nil {
		_ = os.MkdirAll(directory, 0766)
		message := util.NewResponseFailMessage()
		message.Message = "打开本地目录错误"
		ret, _ := json.Marshal(message)
		_, _ = w.Write(ret)
		return
	}
	// xxx-xxx-xxx.text
	_filename := fileId + fileSuffix
	filePath := directory + _filename // 文件路径
	f, err := os.OpenFile(filePath, os.O_WRONLY|os.O_CREATE, 0666)
	if err != nil {
		message := util.NewResponseFailMessage()
		message.Message = err.Error()
		ret, _ := json.Marshal(message)
		_, _ = w.Write(ret)
		return
	}
	_, _ = io.Copy(f, file)
	_ = f.Close()

	//message := util.NewResponseSuccessMessage()
	//message.Data = "上传成功"
	//ret, _ := json.Marshal(message)
	//_, _ = w.Write(ret)
	//return

	// 3.本地文件上传到OSS
	// 创建OSSClient实例。
	//client, err := oss.New("https://oss-cn-shanghai.aliyuncs.com", "LTAI5tD38nDccyJapbeafV91", "SiykXcKWLRVVrerNLqrOM19E7wtXFV")
	//if err != nil {
	//	fmt.Println("Error oss.New:", err)
	//	message := util.NewResponseFailMessage()
	//	message.Message = err.Error()
	//	ret, _ := json.Marshal(message)
	//	_, _ = w.Write(ret)
	//	return
	//}

	// 获取存储空间。
	//bucket, err := c.client.Bucket("qinguanjia")
	//if err != nil {
	//	fmt.Println("Error bucket:", err)
	//	message := util.NewResponseFailMessage()
	//	message.Message = err.Error()
	//	ret, _ := json.Marshal(message)
	//	_, _ = w.Write(ret)
	//	return
	//}

	// 3.本地文件上传到OSS
	// 读取本地文件。
	fd, err := os.Open(filePath)
	if err != nil {
		fmt.Println("Error:", err)
		message := util.NewResponseFailMessage()
		message.Message = "读取本地文件错误"
		ret, _ := json.Marshal(message)
		_, _ = w.Write(ret)
		return
	}
	defer fd.Close()

	// 使用哪个bucket
	bucket, err := c.client.Bucket(bucketName)
	if err != nil {
		fmt.Println("Error bucket:", err)
		message := util.NewResponseFailMessage()
		message.Message = "bucket 使用错误"
		ret, _ := json.Marshal(message)
		_, _ = w.Write(ret)
		return
	}
	// test/xxx-xxx.png
	err = bucket.PutObject(prefix+_filename, fd)
	if err != nil {
		fmt.Println("Error PutObject:", err)
		message := util.NewResponseFailMessage()
		message.Message = err.Error()
		ret, _ := json.Marshal(message)
		_, _ = w.Write(ret)
		return
	}

	// 4.文件信息保存到数据库
	fileInfo, _ := os.Stat(filePath)
	fileDto := dto.File{}
	fileDto.Bucket = bucketName
	fileDto.Prefix = prefix
	fileDto.FileId = fileId
	fileDto.Name = filename
	fileDto.Suffix = fileSuffix
	fileDto.Size = fileInfo.Size()
	_fileId, err := c.fileService.Add(fileDto)
	if err != nil {
		message := util.NewResponseFailMessage()
		message.Message = err.Error()
		ret, _ := json.Marshal(message)
		_, _ = w.Write(ret)
		return
	}

	// 5.删除临时文件
	_ = os.Remove(filePath)

	// 6.查询文件信息返回
	fileVo, err := c.fileService.SelectFileById(_fileId)
	if err != nil {
		message := util.NewResponseFailMessage()
		message.Message = err.Error()
		ret, _ := json.Marshal(message)
		_, _ = w.Write(ret)
		return
	}

	message := util.NewResponseMessage()
	message.Data = fileVo
	ret, _ := json.Marshal(message)
	_, _ = w.Write(ret)
}
/*
func (c *FileController) preview(w http.ResponseWriter, _ *http.Request, ps httprouter.Params) {
	fileId := ps.ByName("fileId")
	if fileId == "" {
		message := util.NewResponseFailMessage()
		message.Message = "fileId 不可为空"
		ret, _ := json.Marshal(message)
		_, _ = w.Write(ret)
		return
	}

	file, err := c.fileService.SelectFileById(fileId)
	if err != nil {
		fmt.Println("Error:", err)
		message := util.NewResponseFailMessage()
		message.Message = "查询文件信息错误"
		ret, _ := json.Marshal(message)
		_, _ = w.Write(ret)
		return
	}

	// 使用哪个bucket
	bucket, err := c.client.Bucket(file.Bucket)
	if err != nil {
		fmt.Println("Error bucket:", err)
		message := util.NewResponseFailMessage()
		message.Message = "bucket 使用错误"
		ret, _ := json.Marshal(message)
		_, _ = w.Write(ret)
		return
	}

	// 下载文件到流。
	body, err := bucket.GetObject(file.Prefix + file.FileId + file.Suffix)
	if err != nil {
		fmt.Println("Error:", err)
		message := util.NewResponseFailMessage()
		message.Message = "下载文件失败"
		ret, _ := json.Marshal(message)
		_, _ = w.Write(ret)
		return
	}
	// 数据读取完成后，获取的流必须关闭，否则会造成连接泄漏，导致请求无连接可用，程序无法正常工作。
	defer body.Close()
	_, _ = io.Copy(w, body)
}
*/
func (c *FileController) preview(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	fileId := ps.ByName("fileId")
	if fileId == "" {
		message := util.NewResponseFailMessage()
		message.Message = "fileId 不可为空"
		ret, _ := json.Marshal(message)
		_, _ = w.Write(ret)
		return
	}

	file, err := c.fileService.SelectFileById(fileId)
	if err != nil {
		fmt.Println("Error:", err)
		message := util.NewResponseFailMessage()
		message.Message = "查询文件信息错误"
		ret, _ := json.Marshal(message)
		_, _ = w.Write(ret)
		return
	}
	// https://qinguanjia.oss-cn-shanghai.aliyuncs.com/test/a/494bb893-8d9e-4b5d-be51-1695c9b0a242.png
	redirectUrl := fmt.Sprintf(
		"%s%s.%s/%s%s%s",
		config.OSS_protocol,
		file.Bucket,
		config.OSS_domain,
		file.Prefix,
		file.FileId,
		file.Suffix,
	)
	http.Redirect(w, r, redirectUrl, http.StatusTemporaryRedirect)
}

func (c *FileController) download(w http.ResponseWriter, _ *http.Request, ps httprouter.Params) {
	fileId := ps.ByName("fileId")
	if fileId == "" {
		message := util.NewResponseFailMessage()
		message.Message = "fileId 不可为空"
		ret, _ := json.Marshal(message)
		_, _ = w.Write(ret)
		return
	}

	file, err := c.fileService.SelectFileById(fileId)
	if err != nil {
		fmt.Println("Error:", err)
		message := util.NewResponseFailMessage()
		message.Message = "查询文件信息错误"
		ret, _ := json.Marshal(message)
		_, _ = w.Write(ret)
		return
	}

	// 使用哪个bucket
	bucket, err := c.client.Bucket(file.Bucket)
	if err != nil {
		fmt.Println("Error bucket:", err)
		message := util.NewResponseFailMessage()
		message.Message = "bucket 使用错误"
		ret, _ := json.Marshal(message)
		_, _ = w.Write(ret)
		return
	}

	// 下载文件到流。
	body, err := bucket.GetObject(file.Prefix + file.FileId + file.Suffix)
	if err != nil {
		fmt.Println("Error:", err)
		message := util.NewResponseFailMessage()
		message.Message = "下载文件失败, OSS文件不存在"
		ret, _ := json.Marshal(message)
		_, _ = w.Write(ret)
		return
	}
	// 数据读取完成后，获取的流必须关闭，否则会造成连接泄漏，导致请求无连接可用，程序无法正常工作。
	defer body.Close()

	content, err := ioutil.ReadAll(body)

	if err != nil {
		fmt.Println("Error:", err)
		message := util.NewResponseFailMessage()
		message.Message = "读取文件失败"
		ret, _ := json.Marshal(message)
		_, _ = w.Write(ret)
		return
	}

	filename := url.QueryEscape(file.Name) // 重命名前的文件名，防止中文乱码
	w.Header().Add("Content-Type", "application/octet-stream")
	w.Header().Add("Content-Disposition", "attachment; filename=\""+filename+"\"")
	_, _ = w.Write(content)
}

/*
func (c *FileController) download(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	fileId := ps.ByName("fileId")
	if fileId == "" {
		message := util.NewResponseFailMessage()
		message.Message = "fileId 不可为空"
		ret, _ := json.Marshal(message)
		_, _ = w.Write(ret)
		return
	}

	file, err := c.fileService.SelectFileById(fileId)
	if err != nil {
		fmt.Println("Error:", err)
		message := util.NewResponseFailMessage()
		message.Message = "查询文件信息错误"
		ret, _ := json.Marshal(message)
		_, _ = w.Write(ret)
		return
	}

	userFile := file.Prefix + file.FileId + file.Suffix
	fut, _ := os.Create(userFile)
	defer fut.Close()
	res, ossErr := http.Get(config.OSS_endpoint + "/" + file.Prefix + file.FileId + file.Suffix)
	if ossErr != nil {
		fmt.Println("Error:", ossErr)
		message := util.NewResponseFailMessage()
		message.Message = "oss获取文件失败"
		ret, _ := json.Marshal(message)
		_, _ = w.Write(ret)
		return
	}
	w.Header().Add("Content-Type", "application/octet-stream")
	w.Header().Add("Content-Disposition", "attachment; filename=\""+file.Name+"\"")
	content, _ := ioutil.ReadAll(res.Body)
	w.Write(content)
}
*/

func (c *FileController) delete(w http.ResponseWriter, _ *http.Request, ps httprouter.Params) {
	fileIds := ps.ByName("fileIds")
	if fileIds == "" {
		message := util.NewResponseFailMessage()
		message.Message = "fileId 不可为空"
		ret, _ := json.Marshal(message)
		_, _ = w.Write(ret)
		return
	}

	// 数据库查询文件
	countSplit := strings.Split(fileIds, ",")              // 字符串转数组
	fileIdArr := util.RemoveDuplicatesAndEmpty(countSplit) // 去空，去重复
	voFileList, err := c.fileService.SelectFileByIds(fileIdArr)
	if err != nil {
		message := util.NewResponseFailMessage()
		message.Message = "查询文件失败"
		ret, _ := json.Marshal(message)
		_, _ = w.Write(ret)
		return
	}

	if len(voFileList) == 0 {
		message := util.NewResponseFailMessage()
		message.Message = "文件不存在"
		ret, _ := json.Marshal(message)
		_, _ = w.Write(ret)
		return
	}

	// 统计删除成功的文件
	var deleteSuccessFileIds []string
	var mutex sync.Mutex
	add := func(fileId string) {
		mutex.Lock()
		deleteSuccessFileIds = append(deleteSuccessFileIds, fileId)
		defer mutex.Unlock()
	}

	// 并发删除, 支持删除其他bucket目录下的文件
	var wg sync.WaitGroup
	wg.Add(len(voFileList))
	for _, voFile := range voFileList {
		go func(_bucketName string, _prefix string, _fileId string, _suffix string, wg *sync.WaitGroup) {
			bucket, err := c.client.Bucket(_bucketName)
			if err != nil {
				fmt.Println("c.client.Bucket Error:", err)
				wg.Done()
				return
			}
			err = bucket.DeleteObject(_prefix + _fileId + _suffix)
			if err != nil {
				fmt.Println("bucket.DeleteObject Error:", err)
			} else {
				// 删除成功的
				add(_fileId)
				fmt.Println(_fileId + " 删除成功")
			}
			wg.Done()
		}(voFile.Bucket, voFile.Prefix, voFile.FileId, voFile.Suffix, &wg)
	}
	wg.Wait()

	// 使用哪个bucket
	//bucket, err := c.client.Bucket(voFileList[0].Bucket)
	//if err != nil {
	//	fmt.Println("Error bucket:", err)
	//	message := util.NewResponseFailMessage()
	//	message.Message = "bucket 使用错误"
	//	ret, _ := json.Marshal(message)
	//	_, _ = w.Write(ret)
	//	return
	//}

	// 删除文件的数组
	//var ossFileNameArr []string
	//for _, voFile := range voFileList {
	//	ossFileName := voFile.Prefix + voFile.FileId + voFile.Suffix
	//	ossFileNameArr = append(ossFileNameArr, ossFileName)
	//}

	// 返回删除成功的文件。
	//delRes, err := bucket.DeleteObjects(ossFileNameArr)
	//if err != nil {
	//	fmt.Println("Error:", err)
	//	message := util.NewResponseFailMessage()
	//	message.Message = "OSS删除失败请重试"
	//	ret, _ := json.Marshal(message)
	//	_, _ = w.Write(ret)
	//	return
	//}

	// 删除数据库
	err = c.fileService.Delete(deleteSuccessFileIds)
	if err != nil {
		message := util.NewResponseFailMessage()
		message.Message = "删除记录失败请重试"
		ret, _ := json.Marshal(message)
		_, _ = w.Write(ret)
		return
	}

	message := util.NewResponseSuccessMessage()
	message.Data = deleteSuccessFileIds
	ret, _ := json.Marshal(message)
	_, _ = w.Write(ret)
}

func (c *FileController) bucketList(w http.ResponseWriter, _ *http.Request, _ httprouter.Params) {
	lsRes, err := c.client.ListBuckets()
	if err != nil {
		message := util.NewResponseFailMessage()
		message.Message = err
		ret, _ := json.Marshal(message)
		_, _ = w.Write(ret)
		return
	}

	message := util.NewResponseSuccessMessage()
	message.Data = lsRes
	ret, _ := json.Marshal(message)
	_, _ = w.Write(ret)
}

func (c *FileController) bucketPrefixList(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	values := r.URL.Query()
	bucketName := values.Get("bucket")
	prefix := values.Get("prefix")
	if bucketName == "" {
		message := util.NewResponseFailMessage()
		message.Message = "bucket 不能为空"
		ret, _ := json.Marshal(message)
		_, _ = w.Write(ret)
		return
	}

	// 判断bucket是否存在
	lsRes, err := c.client.ListBuckets()
	if err != nil {
		message := util.NewResponseFailMessage()
		message.Message = err
		ret, _ := json.Marshal(message)
		_, _ = w.Write(ret)
		return
	}
	isBucket := false
	for _, bucketItem := range lsRes.Buckets {
		if bucketItem.Name == bucketName {
			isBucket = true
		}
	}
	if isBucket == false {
		message := util.NewResponseFailMessage()
		message.Message = "你使用的bucket不存在"
		ret, _ := json.Marshal(message)
		_, _ = w.Write(ret)
		return
	}

	// 使用哪个bucket
	bucket, err := c.client.Bucket(bucketName)
	if err != nil {
		fmt.Println("Error bucket:", err)
		message := util.NewResponseFailMessage()
		message.Message = "bucket 使用错误"
		ret, _ := json.Marshal(message)
		_, _ = w.Write(ret)
		return
	}

	// 列举所有文件。
	_lsRes, err := bucket.ListObjectsV2(
		oss.Prefix(prefix), // 获取指定目录下的目录 test/ 目录下
		oss.ContinuationToken(""),
		oss.Delimiter("/"), // 目录一 / 分割
	)
	if err != nil {
		fmt.Println("Error bucket:", err)
		message := util.NewResponseFailMessage()
		message.Message = "列举所有文件目录错误"
		ret, _ := json.Marshal(message)
		_, _ = w.Write(ret)
		return
	}
	message := util.NewResponseSuccessMessage()
	message.Data = _lsRes.CommonPrefixes
	ret, _ := json.Marshal(message)
	_, _ = w.Write(ret)
}

func NewFileController(router *httprouter.Router) {
	client, err := oss.New(
		config.OSS_endpoint,
		config.OSS_accessKeyID,
		config.OSS_accessKeySecret,
	)
	if err != nil {
		fmt.Println("Error oss.New:", err)
		return
	}

	c := &FileController{
		fileService: service.NewFileService(),
		path:        "/file",
		client:      client,
	}
	router.GET(config.BaseUrl+c.path, c.index)
	router.GET(config.ApiBaseUrl+c.path+"/info/:fileId", c.info)
	router.POST(config.ApiBaseUrl+c.path+"/pageList", c.pageList)
	router.GET(config.ApiBaseUrl+c.path+"/list", c.list)
	router.POST(config.ApiBaseUrl+c.path+"/upload", c.upload)
	router.GET(config.ApiBaseUrl+c.path+"/preview/:fileId", c.preview)
	router.GET(config.ApiBaseUrl+c.path+"/download/:fileId", c.download)
	router.GET(config.ApiBaseUrl+c.path+"/delete/:fileIds", c.delete)
	router.GET(config.ApiBaseUrl+c.path+"/bucketList", c.bucketList)
	router.GET(config.ApiBaseUrl+c.path+"/bucketPrefixList", c.bucketPrefixList)
}
