package dao

import (
	"file/dto"
	"file/lib/db"
	"file/lib/util"
	"file/model"
	"path"
)

type IFileDao interface {
	Add(file model.File) (fileId string, err error)
	SelectFileById(fileId string) (file model.File, err error)
	SelectFileByIds(fileIds []string) (list []model.File, err error)
	SelectAllList(bucket string) (list []model.File, err error)
	SelectPageList(page dto.FilePage) (list []model.File, err error)
	Delete(fileIds []string) (err error)
}

type FileDao struct{}

func (d *FileDao) Add(file model.File) (fileId string, err error) {
	_db := db.GetDB()
	file.OriginName = file.Name
	file.State = "0"
	file.CreateTime = util.CurrentTime()
	file.UpdateTime = file.CreateTime
	if file.Suffix == "" {
		file.Suffix = path.Ext(file.Name)
	}
	result := _db.Create(&file)
	if result.Error != nil {
		return "", result.Error
	}
	return file.FileId, nil
}

func (d *FileDao) SelectFileById(fileId string) (file model.File, err error) {
	_db := db.GetDB()
	result := _db.Where("file_id", fileId).Find(&file)
	if result.Error != nil {
		return file, result.Error
	}
	return file, nil
}

func (d *FileDao) SelectFileByIds(fileIds []string) (list []model.File, err error) {
	_db := db.GetDB()
	query := "file_id IN ("
	for index, fileId := range fileIds {
		if index == len(fileIds)-1 { // 最后一个
			query += "'" + fileId + "'"
		} else {
			query += "'" + fileId + "',"
		}
	}
	query += ");"
	result := _db.Where(query).Find(&list)
	if result.Error != nil {
		return nil, result.Error
	}
	return list, nil
}

func (d *FileDao) SelectAllList(bucket string) (list []model.File, err error) {
	_db := db.GetDB()
	result := _db.Where("bucket=?", bucket).Find(&list)
	if result.Error != nil {
		return nil, result.Error
	}
	if result.RowsAffected == 0 {
		return list, nil
	}
	return list, nil
}

func (d *FileDao) SelectPageList(page dto.FilePage) (list []model.File, err error) {
	_db := db.GetDB()
	// (页码 - 1)*每页条数
	offset := (page.Current - 1) * page.PageSize
	result := _db.Where(
		`
				bucket = ?
			and
				prefix = ?
			and
				name like ?
			and
				suffix like ?
			and
				state like ?
		`,
		page.Data.Bucket,
		page.Data.Prefix,
		"%"+page.Data.Name+"%",
		"%"+page.Data.Suffix+"%",
		"%"+page.Data.State+"%",
	).Offset(offset).Limit(page.PageSize).Order(page.SortField + " " + page.Sort).Find(&list)
	if result.Error != nil {
		return nil, result.Error
	}
	return list, nil
}

func (d *FileDao) Delete(fileIds []string) (err error) {
	_db := db.GetDB()
	query := "file_id IN ("
	for index, fileId := range fileIds {
		if index == len(fileIds)-1 { // 最后一个
			query += "'" + fileId + "'"
		} else {
			query += "'" + fileId + "',"
		}
	}
	query += ");"
	// DELETE FROM file WHERE file_id IN ('a', 'b');
	result := _db.Where(query).Delete(&model.File{})
	if result.Error != nil {
		return result.Error
	}
	return nil
}

func NewFileDao() IFileDao {
	return &FileDao{}
}
