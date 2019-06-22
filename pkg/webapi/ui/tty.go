package ui

import (
	"github.com/emicklei/go-restful"
	"io"
)

func Tty(request *restful.Request, response *restful.Response) {
	getclient(request.Request.Body,response)

}


func getclient(in io.Reader,out io.Writer) {


}
