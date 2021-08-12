package service

import (
	"file/lib/util"
	"file/model"
	"file/vo"
)

type IEnumService interface {
	EnumMapList() (enumMapList vo.EnumMapListType)
}

type EnumService struct {}

func (s *EnumService) EnumMapList() (enumMapList vo.EnumMapListType) {
	_enumMapList := make(vo.EnumMapListType)
	util.MapMapping(model.EnumMapList, &_enumMapList)
	return _enumMapList
}

func NewEnumService() IEnumService {


	return &EnumService{}
}

