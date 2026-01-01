package org

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/StefanShivarov/gollab-backend/internal/common"
	"github.com/go-chi/chi/v5"
)

type UserHandler struct {
	Service *UserService
}

func NewUserHandler(service *UserService) *UserHandler {
	return &UserHandler{Service: service}
}

func (h *UserHandler) List(w http.ResponseWriter, r *http.Request) {
	page, _ := strconv.Atoi(r.URL.Query().Get("page"))
	size, _ := strconv.Atoi(r.URL.Query().Get("size"))

	if page <= 0 {
		page = 1
	}

	if size <= 0 {
		size = 10
	}

	resp, err := h.Service.List(page, size)
	if err != nil {
		common.WriteError(w, err)
		return
	}

	common.WriteJSON(w, http.StatusOK, resp)
}

func (h *UserHandler) GetByID(w http.ResponseWriter, r *http.Request) {
	id, err := common.ParseUUID(chi.URLParam(r, "userId"))
	if err != nil {
		common.WriteError(w, err)
		return
	}

	resp, err := h.Service.GetByID(id)
	if err != nil {
		common.WriteError(w, err)
		return
	}

	common.WriteJSON(w, http.StatusOK, resp)
}

func (h *UserHandler) Create(w http.ResponseWriter, r *http.Request) {
	var req CreateUserRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		common.WriteError(w, common.BadRequest(err.Error()))
		return
	}

	resp, err := h.Service.Create(req)
	if err != nil {
		common.WriteError(w, err)
		return
	}

	common.WriteJSON(w, http.StatusCreated, resp)
}

func (h *UserHandler) UpdateByID(w http.ResponseWriter, r *http.Request) {
	id, err := common.ParseUUID(chi.URLParam(r, "userId"))
	if err != nil {
		common.WriteError(w, err)
		return
	}

	var req UpdateUserRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		common.WriteError(w, common.BadRequest(err.Error()))
		return
	}

	resp, err := h.Service.UpdateByID(id, req)
	if err != nil {
		common.WriteError(w, err)
		return
	}

	common.WriteJSON(w, http.StatusOK, resp)
}

func (h *UserHandler) DeleteByID(w http.ResponseWriter, r *http.Request) {
	id, err := common.ParseUUID(chi.URLParam(r, "userId"))
	if err != nil {
		common.WriteError(w, err)
		return
	}

	if err := h.Service.DeleteByID(id); err != nil {
		common.WriteError(w, err)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

type TeamHandler struct {
	Service *TeamService
}

func NewTeamHandler(service *TeamService) *TeamHandler {
	return &TeamHandler{Service: service}
}

func (h *TeamHandler) Create(w http.ResponseWriter, r *http.Request) {
	var req CreateTeamRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		common.WriteError(w, common.BadRequest(err.Error()))
		return
	}

	creatorId, err := common.ParseUUID(r.URL.Query().Get("creatorId"))
	if err != nil {
		common.WriteError(w, err)
		return
	}
	resp, err := h.Service.Create(creatorId, req)
	if err != nil {
		common.WriteError(w, err)
		return
	}
	common.WriteJSON(w, http.StatusCreated, resp)
}

func (h *TeamHandler) List(w http.ResponseWriter, r *http.Request) {
	page, _ := strconv.Atoi(r.URL.Query().Get("page"))
	size, _ := strconv.Atoi(r.URL.Query().Get("size"))

	if page <= 0 {
		page = 1
	}

	if size <= 0 {
		size = 10
	}

	resp, err := h.Service.List(page, size)
	if err != nil {
		common.WriteError(w, err)
		return
	}

	common.WriteJSON(w, http.StatusOK, resp)
}

func (h *TeamHandler) GetByID(w http.ResponseWriter, r *http.Request) {
	id, err := common.ParseUUID(chi.URLParam(r, "teamId"))
	if err != nil {
		common.WriteError(w, err)
		return
	}
	resp, err := h.Service.GetByID(id)
	if err != nil {
		common.WriteError(w, err)
		return
	}
	common.WriteJSON(w, http.StatusOK, resp)
}

func (h *TeamHandler) UpdateByID(w http.ResponseWriter, r *http.Request) {
	id, err := common.ParseUUID(chi.URLParam(r, "teamId"))
	if err != nil {
		common.WriteError(w, err)
		return
	}

	var req UpdateTeamRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		common.WriteError(w, common.BadRequest(err.Error()))
		return
	}

	resp, err := h.Service.UpdateByID(id, req)
	if err != nil {
		common.WriteError(w, err)
		return
	}

	common.WriteJSON(w, http.StatusOK, resp)
}

func (h *TeamHandler) DeleteByID(w http.ResponseWriter, r *http.Request) {
	id, err := common.ParseUUID(chi.URLParam(r, "teamId"))
	if err != nil {
		common.WriteError(w, err)
		return
	}

	if err := h.Service.DeleteByID(id); err != nil {
		common.WriteError(w, err)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

func (h *TeamHandler) ListTeamMembers(w http.ResponseWriter, r *http.Request) {
	teamID, err := common.ParseUUID(chi.URLParam(r, "teamId"))
	if err != nil {
		common.WriteError(w, err)
		return
	}
	members, err := h.Service.ListMembers(teamID)
	if err != nil {
		common.WriteError(w, err)
		return
	}

	common.WriteJSON(w, http.StatusOK, members)
}

func (h *TeamHandler) AddMember(w http.ResponseWriter, r *http.Request) {
	var req CreateMembershipRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		common.WriteError(w, common.BadRequest(err.Error()))
		return
	}

	if err := h.Service.AddMembership(req); err != nil {
		common.WriteError(w, err)
		return
	}

	w.WriteHeader(http.StatusCreated)
}

func (h *TeamHandler) RemoveMember(w http.ResponseWriter, r *http.Request) {
	var req DeleteMembershipRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		common.WriteError(w, common.BadRequest(err.Error()))
		return
	}

	if err := h.Service.RemoveMembership(req.TeamID, req.UserID); err != nil {
		common.WriteError(w, err)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}
