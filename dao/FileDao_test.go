package dao

import (
	"file/dto"
	"file/model"
	"fmt"
	"testing"
)

// go test -v -run TestFileDao_Add FileDao_test.go FileDao.go
func TestFileDao_Add(t *testing.T) {
	fileDao := NewFileDao()
	app := model.File{
		FileId: "7aedfc5d-20d9-4195-be2d-b033c3192dcf",
		Name:   "aaa.png",
		Size:   1000,
		Suffix: ".png",
		State:  "0",
		Bucket: "qinguanjia",
		Prefix: "",
	}
	appId, err := fileDao.Add(app)
	if err != nil {
		t.Logf(`插入数据失败, %s`, err)
		return
	}
	t.Logf(`插入数据成功, %s`, appId)
}

// go test -v -run TestFileDao_SelectFileById FileDao_test.go FileDao.go
func TestFileDao_SelectFileById(t *testing.T) {
	fileDao := NewFileDao()
	file, err := fileDao.SelectFileById("7aedfc5d-20d9-4195-be2d-b033c3192dcf")
	if err != nil {
		t.Logf(`查询数据失败, %s`, err)
		return
	}
	t.Logf(`查询数据成功, %s`, file.FileId)
}

// go test -v -run TestFileDao_SelectFileByIds FileDao_test.go FileDao.go
func TestFileDao_SelectFileByIds(t *testing.T) {
	fileDao := NewFileDao()
	fileIds := []string{"cf895fc5-b12c-41e8-b9cc-5c20f7b6ee78", "1bc0132f-d837-43b6-84ed-de211aa03d7a"}
	list, err := fileDao.SelectFileByIds(fileIds)
	if err != nil {
		t.Logf(`查询数据失败, %s`, err)
		return
	}
	t.Logf(`查询数据成功, %d`, len(list))
}

// go test -v -run TestFileDao_SelectFileAllList FileDao_test.go FileDao.go
func TestFileDao_SelectFileAllList(t *testing.T) {
	fileDao := NewFileDao()
	list, err := fileDao.SelectAllList("qinguanjia")
	if err != nil {
		t.Logf(`查询数据失败, %s`, err)
		return
	}
	t.Logf(`查询数据成功, %s`, fmt.Sprint(len(list)))
}

// go test -v -run TestFileDao_SelectFilePageList FileDao_test.go FileDao.go
func TestFileDao_SelectFilePageList(t *testing.T) {
	fileDao := NewFileDao()
	page := dto.FilePage{}
	page.Data.Bucket = "qinguanjia"
	if page.SortField == "" {
		page.SortField = "create_time"
	}
	if page.Sort == "" {
		page.Sort = "desc"
	}
	if page.Current == 0 {
		page.Current = 1
	}
	if page.PageSize == 0 {
		page.PageSize = 10
	}
	list, err := fileDao.SelectPageList(page)
	if err != nil {
		t.Logf(`查询数据失败, %s`, err)
		return
	}
	allList, err := fileDao.SelectAllList(page.Data.Bucket)
	if err != nil {
		t.Logf(`查询数据失败, %s`, err)
		return
	}
	page.Total = len(allList)
	for _, fileMode := range list {
		page.List = append(page.List, fileMode)
	}
	//fmt.Println(page)
	t.Logf(`查询数据成功,%d`, len(list))
}

// go test -v -run TestFileDao_Delete FileDao_test.go FileDao.go
func TestFileDao_Delete(t *testing.T) {
	fileDao := NewFileDao()
	fileIds := []string{"4d0cb5da-08b2-47a7-b6d1-3d71f4304cbb", "3c2c4d92-bb41-4d10-a142-fea93506c7a0"}
	err := fileDao.Delete(fileIds)
	if err != nil {
		t.Logf(`失败, %s`, err)
		return
	}
	t.Logf(`成功`)
}
