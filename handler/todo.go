package handler

import (
	"context"
	"encoding/json"
	"log"
	"net/http"

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
	_, _ = h.svc.ReadTODO(ctx, 0, 0)
	return &model.ReadTODOResponse{}, nil
}

// Update handles the endpoint that updates the TODO.
func (h *TODOHandler) Update(ctx context.Context, req *model.UpdateTODORequest) (*model.UpdateTODOResponse, error) {
	_, _ = h.svc.UpdateTODO(ctx, 0, "", "")
	return &model.UpdateTODOResponse{}, nil
}

// Delete handles the endpoint that deletes the TODOs.
func (h *TODOHandler) Delete(ctx context.Context, req *model.DeleteTODORequest) (*model.DeleteTODOResponse, error) {
	_ = h.svc.DeleteTODO(ctx, nil)
	return &model.DeleteTODOResponse{}, nil
}

func (h *TODOHandler) ServeHTTP(w http.ResponseWriter, r *http.Request){
	//HTTP メソッドが Post の場合を判定
	if(r.Method != http.MethodPost){
		return
	}

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
