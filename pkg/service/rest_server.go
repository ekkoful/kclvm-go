// Copyright 2021 The KCL Authors. All rights reserved.

package service

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"

	"github.com/golang/protobuf/proto"
	"github.com/julienschmidt/httprouter"

	"kusionstack.io/kclvm-go/pkg/3rdparty/grpc_gateway_util"
	"kusionstack.io/kclvm-go/pkg/spec/gpyrpc"
)

var _ = fmt.Sprint

// Client represents an restful method result.
type RestfulResult struct {
	Error  string        `json:"error"`
	Result proto.Message `json:"result"`
}

type restServer struct {
	address string
	router  *httprouter.Router
	builtin *BuiltinServiceClient
	c       *KclvmServiceClient
}

func RunRestServer(address string) error {
	s := newRestServer(address)
	return s.Run()
}

func newRestServer(address string) *restServer {
	if strings.HasPrefix(address, ":") {
		address = "127.0.0.1" + address
	}
	p := &restServer{
		address: address,
		router:  httprouter.New(),
		builtin: NewBuiltinServiceClient(),
		c:       NewKclvmServiceClient(),
	}
	p.initHttpRrouter()
	return p
}

func (p *restServer) Run() error {
	fmt.Printf("listen on http://%s ...\n", p.address)
	return http.ListenAndServe(p.address, p.router)
}

func (p *restServer) initHttpRrouter() {
	p.router.GET("/api:protorpc/BuiltinService.Ping", p.handle_Ping)
	p.router.GET("/api:protorpc/BuiltinService.ListMethod", p.handle_ListMethod)

	p.router.GET("/api:protorpc/KclvmService.ExecProgram", p.handle_ExecProgram)
	p.router.GET("/api:protorpc/KclvmService.FormatCode", p.handle_FormatCode)
	p.router.GET("/api:protorpc/KclvmService.FormatPath", p.handle_FormatPath)
	p.router.GET("/api:protorpc/KclvmService.LintPath", p.handle_LintPath)
	p.router.GET("/api:protorpc/KclvmService.OverrideFile", p.handle_OverrideFile)
	p.router.GET("/api:protorpc/KclvmService.GetSchemaType", p.handle_GetSchemaType)
	p.router.GET("/api:protorpc/KclvmService.ValidateCode", p.handle_ValidateCode)

	p.router.POST("/api:protorpc/BuiltinService.Ping", p.handle_Ping)
	p.router.POST("/api:protorpc/BuiltinService.ListMethod", p.handle_ListMethod)

	p.router.POST("/api:protorpc/KclvmService.ExecProgram", p.handle_ExecProgram)
	p.router.POST("/api:protorpc/KclvmService.FormatCode", p.handle_FormatCode)
	p.router.POST("/api:protorpc/KclvmService.FormatPath", p.handle_FormatPath)
	p.router.POST("/api:protorpc/KclvmService.LintPath", p.handle_LintPath)
	p.router.POST("/api:protorpc/KclvmService.OverrideFile", p.handle_OverrideFile)
	p.router.POST("/api:protorpc/KclvmService.GetSchemaType", p.handle_GetSchemaType)
	p.router.POST("/api:protorpc/KclvmService.ValidateCode", p.handle_ValidateCode)
}

func (p *restServer) handle(
	w http.ResponseWriter, r *http.Request,
	args proto.Message, fn func() (proto.Message, error),
) {
	switch r.Method {
	case "GET":
		if err := grpc_gateway_util.PopulateQueryParameters(args, r.URL.Query()); err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
	default:
		if err := json.NewDecoder(r.Body).Decode(args); err != nil && err != io.EOF {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
	}

	w.Header().Set("Content-Type", "application/json")

	var result RestfulResult
	if x, err := fn(); err != nil {
		result.Error = err.Error()
	} else {
		result.Result = x // OK
	}

	encoder := json.NewEncoder(w)
	encoder.SetIndent("", "\t")

	if err := encoder.Encode(result); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

func (p *restServer) handle_Ping(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	var args = new(gpyrpc.Ping_Args)
	p.handle(w, r, args, func() (proto.Message, error) {
		return p.builtin.Ping(args)
	})
}

func (p *restServer) handle_ListMethod(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	var args = new(gpyrpc.ListMethod_Args)
	p.handle(w, r, args, func() (proto.Message, error) {
		return p.builtin.ListMethod(args)
	})
}

func (p *restServer) handle_ExecProgram(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	var args = new(gpyrpc.ExecProgram_Args)
	p.handle(w, r, args, func() (proto.Message, error) {
		return p.c.ExecProgram(args)
	})
}

func (p *restServer) handle_FormatCode(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	var args = new(gpyrpc.FormatCode_Args)
	p.handle(w, r, args, func() (proto.Message, error) {
		return p.c.FormatCode(args)
	})
}

func (p *restServer) handle_FormatPath(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	var args = new(gpyrpc.FormatPath_Args)
	p.handle(w, r, args, func() (proto.Message, error) {
		return p.c.FormatPath(args)
	})
}

func (p *restServer) handle_LintPath(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	args := new(gpyrpc.LintPath_Args)
	p.handle(w, r, args, func() (proto.Message, error) {
		return p.c.LintPath(args)
	})
}

func (p *restServer) handle_OverrideFile(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	args := new(gpyrpc.OverrideFile_Args)
	p.handle(w, r, args, func() (proto.Message, error) {
		return p.c.OverrideFile(args)
	})
}

func (p *restServer) handle_GetSchemaType(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	args := new(gpyrpc.GetSchemaType_Args)
	p.handle(w, r, args, func() (proto.Message, error) {
		return p.c.GetSchemaType(args)
	})
}

func (p *restServer) handle_ValidateCode(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	args := new(gpyrpc.ValidateCode_Args)
	p.handle(w, r, args, func() (proto.Message, error) {
		return p.c.ValidateCode(args)
	})
}

func (p *restServer) handle_ListDepFiles(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	args := new(gpyrpc.ListDepFiles_Args)
	p.handle(w, r, args, func() (proto.Message, error) {
		return p.c.ListDepFiles(args)
	})
}

func (p *restServer) handle_LoadSettingsFiles(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	args := new(gpyrpc.LoadSettingsFiles_Args)
	p.handle(w, r, args, func() (proto.Message, error) {
		return p.c.LoadSettingsFiles(args)
	})
}
