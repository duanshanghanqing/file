package service

import (
	"errors"
	"file/dao"
	"file/dto"
	"file/lib/util"
	"file/model"
	"file/vo"
)

type IFileService interface {
	Add(file dto.File) (fileId string, err error)
	SelectFileById(fileId string) (file vo.File, err error)
	SelectFileByIds(fileIds []string) (list []vo.File, err error)
	SelectAllList(bucket string) (list []vo.File, err error)
	SelectPageList(page dto.FilePage) (_page dto.FilePage, err error)
	Delete(fileIs []string) (err error)
}

type FileService struct {
	fileDao dao.IFileDao
}

func (s *FileService) Add(file dto.File) (fileId string, err error) {
	fileModel := model.File{}
	util.EntityMapping(file, &fileModel)
	fileId, err = s.fileDao.Add(fileModel)
	return fileId, err
}

func (s *FileService) SelectFileById(fileId string) (file vo.File, err error) {
	fileModel, _err := s.fileDao.SelectFileById(fileId)
	util.EntityMapping(fileModel, &file)
	return file, _err
}

func (s *FileService) SelectFileByIds(fileIds []string) (list []vo.File, err error) {
	fileListModel, err := s.fileDao.SelectFileByIds(fileIds)
	for _, fileModel := range fileListModel {
		var voFile vo.File
		util.EntityMapping(fileModel, &voFile)
		list = append(list, voFile)
	}
	return list, err
}

func (s *FileService) SelectAllList(bucket string) (list []vo.File, err error) {
	if bucket == "" {
		return nil, errors.New("bucket不可为空")
	}
	fileModelList, err := s.fileDao.SelectAllList(bucket)
	if err != nil {
		return nil, err
	}
	// Model => Dto
	for _, fileModel := range fileModelList {
		var file vo.File
		util.EntityMapping(fileModel, &file)
		list = append(list, file)
	}
	return list, err
}

func (s *FileService) SelectPageList(page dto.FilePage) (_page dto.FilePage, err error) {
	if page.Data.Bucket == "" {
		return _page, errors.New("bucket 不能为空")
	}
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
	list, err := s.fileDao.SelectPageList(page)
	if err != nil {
		return _page, err
	}
	allList, err := s.fileDao.SelectAllList(page.Data.Bucket)
	if err != nil {
		return _page, err
	}
	page.Total = len(allList)
	// 总页数
	if page.Total%page.PageSize == 0 {
		page.TotalPage = page.Total / page.PageSize
	} else {
		page.TotalPage = page.Total/page.PageSize + 1
	}

	// Model => vo
	for _, fileModel := range list {
		var fileVo vo.File
		util.EntityMapping(fileModel, &fileVo)
		page.List = append(page.List, fileVo)
	}
	return page, nil
}

func (s *FileService) Delete(fileIs []string) (err error)  {
	err = s.fileDao.Delete(fileIs)
	return err
}

func NewFileService() IFileService {
	return &FileService{
		fileDao: dao.NewFileDao(),
	}
}
