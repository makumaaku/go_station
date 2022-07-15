package handler

import (
	"context"
	"encoding/json"
	"log"
	"net/http"
	"strconv"

	"github.com/TechBowl-japan/go-stations/model"
	"github.com/TechBowl-japan/go-stations/service"
)

// A TODOHandler implements handling REST endpoints.
type TODOHandler struct {
	svc *service.TODOService
}

// NewTODOHandler returns TODOHandler based http.Handler.
func NewTODOHandler(svc *service.TODOService) *TODOHandler {
	return &TODOHandler{
		svc: svc,
	}
}

// Create handles the endpoint that creates the TODO.
func (h *TODOHandler) Create(ctx context.Context, req *model.CreateTODORequest) (*model.CreateTODOResponse, error) {
	todo, err := h.svc.CreateTODO(ctx, req.Subject, req.Description)
	return &model.CreateTODOResponse{
		TODO: *todo,
	}, err
}

// Read handles the endpoint that reads the TODOs.
func (h *TODOHandler) Read(ctx context.Context, req *model.ReadTODORequest) (*model.ReadTODOResponse, error) {
	todoPointers, _ := h.svc.ReadTODO(ctx,req.PrevID , req.Size)

	todos := make([]model.TODO, 0)

	for _, todo := range todoPointers {
		todos = append(todos, *todo)
}
	return &model.ReadTODOResponse{
		TODOs: todos,
	}, nil
}

// Update handles the endpoint that updates the TODO.
func (h *TODOHandler) Update(ctx context.Context, req *model.UpdateTODORequest) (*model.UpdateTODOResponse, error) {
	todo, err := h.svc.UpdateTODO(ctx, req.ID, req.Subject, req.Description)
	return &model.UpdateTODOResponse{
		TODO: *todo,
	}, err
}

// Delete handles the endpoint that deletes the TODOs.
func (h *TODOHandler) Delete(ctx context.Context, req *model.DeleteTODORequest) (*model.DeleteTODOResponse, error) {
	_ = h.svc.DeleteTODO(ctx, nil)
	return &model.DeleteTODOResponse{}, nil
}

func (h *TODOHandler) ServeHTTP(w http.ResponseWriter, r *http.Request){
	//HTTP メソッドが Post の場合を判定
	switch r.Method{
	case http.MethodPost:
		serveHTTPPost(w,r,h)
	case http.MethodPut:
		serveHTTPPut(w,r,h)
	case http.MethodGet:
		serveHTTPGet(w,r,h)
	}	
}

func serveHTTPPost(w http.ResponseWriter, r *http.Request,h *TODOHandler){
//CreateTODORequest に JSON Decode を行おう
var request = model.CreateTODORequest{}
if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
	log.Println(err)
	return
}

if request.Subject == "" {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusBadRequest)
	return
}

//引数の r から r.Context() を呼び出し
context := r.Context()

//CreateTODO メソッドに呼び出し DB に TODO を保存
resp,err:= h.Create(context,&request);
if err != nil {
	log.Println(err)
	return
}

//保存した TODO を CreateTODOResponse に代入し
w.Header().Set("Content-Type", "application/json")
w.WriteHeader(http.StatusOK)

//JSON Encode を行い HTTP Response を返そう
err = json.NewEncoder(w).Encode(resp)
if(err != nil){
	println(err)
	return
}
}

func serveHTTPPut(w http.ResponseWriter, r *http.Request,h *TODOHandler){
	var request = model.UpdateTODORequest{}
	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
		log.Println(err)
		return
	}

	if request.Subject == "" || request.ID == 0 {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	context := r.Context()
	resp,err := h.Update(context,&request)
	if err != nil {
		log.Println(err)
		return
	}
	
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	
	//JSON Encode を行い HTTP Response を返そう
	err = json.NewEncoder(w).Encode(resp)
	if(err != nil){
		println(err)
		return
	}
}

func serveHTTPGet(w http.ResponseWriter, r *http.Request,h *TODOHandler){
	prevIDStr := r.URL.Query().Get("prev_id")
	prevID, err := strconv.Atoi(prevIDStr)
	if(err != nil){
		log.Println(err,"prev")
		prevID = 0
	}

	sizeStr := r.URL.Query().Get("size")
	size, err := strconv.Atoi(sizeStr)
	if(err != nil){
		log.Println(err,"size")
		size = 5
	}

	request := model.ReadTODORequest{
		PrevID: int64(prevID),
		Size: int64(size),
	}

	context := r.Context()
	resp,err := h.Read(context,&request)

	if(err != nil){
		log.Println(err)
	  return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	err = json.NewEncoder(w).Encode(resp)
	if(err != nil){
		log.Println(err)
		return
	}
}
