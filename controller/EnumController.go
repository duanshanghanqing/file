package controller

import (
	"encoding/json"
	"file/lib/config"
	"file/lib/util"
	"file/service"
	"github.com/julienschmidt/httprouter"
	"net/http"
)

type EnumController struct {
	enumService service.IEnumService
	path        string
}

func (c *EnumController) enumMapList(w http.ResponseWriter, _ *http.Request, _ httprouter.Params) {
	message := util.NewResponseSuccessMessage()
	message.Data = c.enumService.EnumMapList()
	ret, _ := json.Marshal(message)
	_, _ = w.Write(ret)
}

func NewEnumController(router *httprouter.Router) {
	c := &EnumController{
		enumService: service.NewEnumService(),
		path:        "/enum",
	}
	router.GET(config.ApiBaseUrl+c.path+"/enumMapList", c.enumMapList)
}


