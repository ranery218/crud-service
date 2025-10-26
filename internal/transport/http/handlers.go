package http

import (
	"crud/internal/services/user"
	helpers "crud/internal/transport/http/helpers"
	"crud/internal/transport/http/middleware"
	"errors"
	"log"
	"net/http"

	"github.com/samber/mo"
)

type UserHandler struct {
	registerService *user.RegisterService
	loginService    *user.LoginService
	updateService   *user.UpdateService
	deleteService   *user.DeleteService
	logger          *log.Logger
}

func NewUserHandler(
	registerService *user.RegisterService,
	loginService *user.LoginService,
	updateService *user.UpdateService,
	deleteService *user.DeleteService,
	logger *log.Logger) *UserHandler {
	return &UserHandler{
		registerService: registerService,
		loginService:    loginService,
		updateService:   updateService,
		deleteService:   deleteService,
		logger:          logger,
	}
}

func (h *UserHandler) setSessionCookie(w http.ResponseWriter, session user.Session) {
	cookie := &http.Cookie{
		Name:     "session_id",
		Value:    session.ID,
		Path:     "/",
		Expires:  session.ExpiresAt,
		HttpOnly: true,
		Secure:   false,
		SameSite: http.SameSiteLaxMode,
	}
	http.SetCookie(w, cookie)
}

func (h *UserHandler) clearSessionCookie(w http.ResponseWriter) {
	cookie := &http.Cookie{
		Name:     "session_id",
		Value:    "",
		Path:     "/",
		MaxAge:   -1,
		HttpOnly: true,
		Secure:   false,
		SameSite: http.SameSiteLaxMode,
	}
	http.SetCookie(w, cookie)
}

func (h *UserHandler) Register(w http.ResponseWriter, r *http.Request) {
	var registerReq RegisterRequest
	err := helpers.DecodeJSON(r, &registerReq)
	if err != nil {
		h.logger.Printf("update: decode request failed: %v", err)
		helpers.WriteError(w, http.StatusBadRequest, "invalid request")
		return
	}

	ctx := r.Context()

	serviceRequest := user.RegisterRequest{
		Username: registerReq.UserName,
		Email:    registerReq.Email,
		Password: registerReq.Password,
	}

	serviceResponse, err := h.registerService.Register(ctx, serviceRequest)
	if err != nil {
		switch {
		case errors.Is(err, user.ErrEmailRequired) || errors.Is(err, user.ErrPasswordRequired) || errors.Is(err, user.ErrUsernameRequired) ||
			errors.Is(err, user.ErrEmailIncorrect) || errors.Is(err, user.ErrPasswordIncorrect):
			helpers.WriteError(w, http.StatusBadRequest, err.Error())
			return
		case errors.Is(err, user.ErrEmailTaken) || errors.Is(err, user.ErrUsernameTaken):
			helpers.WriteError(w, http.StatusConflict, "conflict")
			return
		default:
			h.logger.Printf("register: internal error: %v", err)
			helpers.WriteError(w, http.StatusInternalServerError, "internal error")
			return
		}
	}

	user := serviceResponse.User
	registerResp := RegisterResponse{
		User: UserDTO{
			ID:       user.ID,
			UserName: user.Username,
			Email:    user.Email,
		},
	}

	err = helpers.WriteJSON(w, http.StatusCreated, registerResp)
	if err != nil {
		h.logger.Printf("register: write response failed: %v", err)
	}
}

func (h *UserHandler) Login(w http.ResponseWriter, r *http.Request) {
	var loginReq LoginRequest
	err := helpers.DecodeJSON(r, &loginReq)
	if err != nil {
		h.logger.Printf("login: decode request failed: %v", err)
		helpers.WriteError(w, http.StatusBadRequest, "invalid request")
		return
	}

	ctx := r.Context()

	serviceRequest := user.LoginRequest{
		Email:    loginReq.Email,
		Password: loginReq.Password,
	}

	serviceResponse, err := h.loginService.Login(ctx, serviceRequest)
	if err != nil {
		switch {
		case errors.Is(err, user.ErrEmailRequired) || errors.Is(err, user.ErrPasswordRequired) ||
			errors.Is(err, user.ErrEmailIncorrect):
			helpers.WriteError(w, http.StatusBadRequest, err.Error())
			return
		case errors.Is(err, user.ErrPasswordIncorrect) || errors.Is(err, user.ErrUserNotFound):
			helpers.WriteError(w, http.StatusUnauthorized, "invalid credentials")
			return
		default:
			h.logger.Printf("login: internal error: %v", err)
			helpers.WriteError(w, http.StatusInternalServerError, "internal error")
			return
		}
	}

	session := serviceResponse.Session
	user := serviceResponse.User

	loginResp := LoginResponse{
		User: UserDTO{
			ID:       user.ID,
			UserName: user.Username,
			Email:    user.Email,
		},
	}

	h.setSessionCookie(w, session)

	err = helpers.WriteJSON(w, http.StatusOK, loginResp)
	if err != nil {
		h.logger.Printf("login: write response failed: %v", err)
	}
}

func (h *UserHandler) Logout(w http.ResponseWriter, r *http.Request) {
	cookie, err := r.Cookie("session_id")
	if err != nil {
		if errors.Is(err, http.ErrNoCookie) {
			h.clearSessionCookie(w)
			w.WriteHeader(http.StatusNoContent)
			return
		}
		h.logger.Printf("logout: read cookie failed: %v", err)
		helpers.WriteError(w, http.StatusBadRequest, "invalid request")
		return
	}

	if err := h.loginService.SessionStore.Delete(r.Context(), cookie.Value); err != nil {
		if !errors.Is(err, user.ErrSessionNotFound) && !errors.Is(err, user.ErrSessionExpired) {
			h.logger.Printf("logout: delete session failed: %v", err)
			helpers.WriteError(w, http.StatusInternalServerError, "internal error")
			return
		}
	}
	h.clearSessionCookie(w)

	w.WriteHeader(http.StatusNoContent)
}

func (h *UserHandler) Update(w http.ResponseWriter, r *http.Request) {
	var updateReq UpdateRequest
	err := helpers.DecodeJSON(r, &updateReq)
	if err != nil {
		h.logger.Printf("update: decode request failed: %v", err)
		helpers.WriteError(w, http.StatusBadRequest, "invalid request")
		return
	}
	ctx := r.Context()
	userID, ok := middleware.UserIDFromContext(ctx)
	if !ok {
		h.logger.Printf("update: userID missing in context")
		helpers.WriteError(w, http.StatusUnauthorized, "missing session")
		return
	}
	email := updateReq.Email
	username := updateReq.UserName
	password := updateReq.Password
	serviceRequest := user.UpdateRequest{
		ID: userID,
	}
	switch {
	case email != nil:
		serviceRequest.Email = mo.Some(*email)
		fallthrough
	case username != nil:
		serviceRequest.Username = mo.Some(*username)
		fallthrough
	case password != nil:
		serviceRequest.Password = mo.Some(*password)
	}
	serviceResp, err := h.updateService.Update(ctx, serviceRequest)
	if err != nil {
		if errors.Is(err, user.ErrUserNotFound) {
			helpers.WriteError(w, http.StatusUnauthorized, "invalid credentials")
			return
		}
		helpers.WriteError(w, http.StatusInternalServerError, "internal error")
		return
	}
	user := serviceResp.User
	updateReps := UpdateResponse{
		User: UserDTO{
			ID:       user.ID,
			UserName: user.Username,
			Email:    user.Email,
		},
	}
	err = helpers.WriteJSON(w, http.StatusOK, updateReps)
	if err != nil {
		h.logger.Printf("update: write response failed: %v", err)
	}
}

func (h *UserHandler) Delete(w http.ResponseWriter, r *http.Request) {
	userID, ok := middleware.UserIDFromContext(r.Context())
	if !ok {
		h.logger.Printf("delete: userID missing in context")
		helpers.WriteError(w, http.StatusUnauthorized, "missing session")
		return
	}

	serviceRequest := user.DeleteRequest{
		ID: userID,
	}

	serviceResp, err := h.deleteService.Delete(r.Context(), serviceRequest)
	if err != nil {
		h.logger.Printf("delete: internal error: %v", err)
		helpers.WriteError(w, http.StatusInternalServerError, "internal error")
		return
	}
	if !serviceResp.Success {
		h.logger.Printf("delete: deletion unsuccessful")
		helpers.WriteError(w, http.StatusInternalServerError, "deletion unsuccessful")
		return
	}

	h.clearSessionCookie(w)

	w.WriteHeader(http.StatusNoContent)
}
