package rest

import (
	"context"
	"encoding/json"
	"errors"
	"log"
	"net/http"
	"strconv"

	"github.com/DexScen/Queue/backend/internal/domain"
	e "github.com/DexScen/Queue/backend/internal/errors"
	"github.com/gorilla/mux"
)

type Queues interface {
	GetAllGames(ctx context.Context, listGames *domain.ListGames) error
	GetGameInfoByID(ctx context.Context, id int) (*domain.Game, error)
	GetGamesByLogin(ctx context.Context, login string, listGames *domain.ListGames) error
	GetIdByLogin(ctx context.Context, login string) (int, error)

	Register(ctx context.Context, user *domain.User) error
	LogIn(ctx context.Context, login, password string) (string, error)

	RemovePlayerFromQueue(ctx context.Context, user_id, game_id int) error
	AddPlayerToQueue(ctx context.Context, user_id, game_id int) (int, error)
}

type Handler struct {
	queuesService Queues
}

func NewQueues(queues Queues) *Handler {
	return &Handler{
		queuesService: queues,
	}
}

func setHeaders(w http.ResponseWriter) {
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
	w.Header().Set("Access-Control-Allow-Headers", "Content-Type")
}

func (h *Handler) OptionsHandler(w http.ResponseWriter, r *http.Request) {
	setHeaders(w)
	w.WriteHeader(http.StatusOK)
}

func (h *Handler) InitRouter() *mux.Router {
	r := mux.NewRouter().StrictSlash(true)
	r.Use(loggingMiddleware)
	r.Use(corsMiddleware)

	links := r.PathPrefix("").Subrouter()
	{
		links.HandleFunc("/games", h.GetAllGames).Methods(http.MethodGet)
		links.HandleFunc("/games/{id}", h.GetGameInfoByID).Methods(http.MethodGet)
		links.HandleFunc("/queue/{login}", h.GetGamesByLogin).Methods(http.MethodGet)
		links.HandleFunc("/auth/{login}", h.GetIdByLogin).Methods(http.MethodGet)
		links.HandleFunc("/auth/register", h.Register).Methods(http.MethodPost)
		links.HandleFunc("/auth/login", h.LogIn).Methods(http.MethodPost)

		links.HandleFunc("/remove", h.RemovePlayerFromQueue).Methods(http.MethodDelete)
		links.HandleFunc("/add", h.AddPlayerToQueue).Methods(http.MethodPost)

		links.HandleFunc("", h.OptionsHandler).Methods(http.MethodOptions)
		links.PathPrefix("/").HandlerFunc(h.OptionsHandler).Methods(http.MethodOptions)

	}
	return r
}

func (h *Handler) AddPlayerToQueue(w http.ResponseWriter, r *http.Request) {
	setHeaders(w)
	var addInfo domain.ChangeInfo
	var pos domain.PosInfo
	if err := json.NewDecoder(r.Body).Decode(&addInfo); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		log.Println("addPlayerToQueue error:", err)
		return
	}

	position, err := h.queuesService.AddPlayerToQueue(r.Context(), addInfo.UserID, addInfo.GameID)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		log.Println("addPlayerToQueue error:", err)
		return
	}
	pos.Pos = position
	if jsonResp, err := json.Marshal(pos); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		log.Println("addPlayerToQueue error:", err)
		return
	} else {
		w.Header().Set("Content-Type", "application/json")
		w.Write(jsonResp)
	}
}

func (h *Handler) RemovePlayerFromQueue(w http.ResponseWriter, r *http.Request) {
	setHeaders(w)
	var removeInfo domain.ChangeInfo

	if err := json.NewDecoder(r.Body).Decode(&removeInfo); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		log.Println("RemovePlayerFromQueue error:", err)
		return
	}

	err := h.queuesService.RemovePlayerFromQueue(r.Context(), removeInfo.UserID, removeInfo.GameID)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		log.Println("RemovePlayerFromQueue error:", err)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

func (h *Handler) GetAllGames(w http.ResponseWriter, r *http.Request) {
	setHeaders(w)
	var list domain.ListGames
	if err := h.queuesService.GetAllGames(context.TODO(), &list); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		log.Println("getAllGames error:", err)
		return
	}

	if jsonResp, err := json.Marshal(list); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		log.Println("getAllGames error:", err)
		return
	} else {
		w.Header().Set("Content-Type", "application/json")
		w.Write(jsonResp)
	}
}

func (h *Handler) GetGameInfoByID(w http.ResponseWriter, r *http.Request) {
	setHeaders(w)
	vars := mux.Vars(r)
	idStr := vars["id"]
	id, err := strconv.Atoi(idStr)

	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		log.Println("getGameInfoByID error:", err)
		return
	}

	game, err := h.queuesService.GetGameInfoByID(context.TODO(), id)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		log.Println("getGameInfoByID error:", err)
		return
	}

	if jsonResp, err := json.Marshal(game); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		log.Println("getGameInfoByID error:", err)
		return
	} else {
		w.Header().Set("Content-Type", "application/json")
		w.Write(jsonResp)
	}
}

func (h *Handler) GetGamesByLogin(w http.ResponseWriter, r *http.Request) {
	setHeaders(w)
	vars := mux.Vars(r)
	loginStr := vars["login"]

	var list domain.ListGames
	if err := h.queuesService.GetGamesByLogin(context.TODO(), loginStr, &list); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		log.Println("GetGamesByLogin error:", err)
		return
	}

	if jsonResp, err := json.Marshal(list); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		log.Println("GetGamesByLogin error:", err)
		return
	} else {
		w.Header().Set("Content-Type", "application/json")
		w.Write(jsonResp)
	}
}

func (h *Handler) LogIn(w http.ResponseWriter, r *http.Request) {
	setHeaders(w)
	var info domain.LoginInfo
	var roleInfo domain.RoleInfo

	if err := json.NewDecoder(r.Body).Decode(&info); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		log.Println("Login error:", err)
		return
	}

	role, err := h.queuesService.LogIn(context.TODO(), info.Login, info.Password)
	if err != nil {
		if errors.Is(err, e.ErrUserNotFound) {
			role = "user not found"
			log.Println("Login error:", err)
		} else if errors.Is(err, e.ErrWrongPassword) {
			role = "wrong password"
			log.Println("Login error:", err)
		} else {
			w.WriteHeader(http.StatusInternalServerError)
			log.Println("Login error:", err, e.ErrUserNotFound)
			return
		}
	}
	roleInfo.Role = role
	if jsonResp, err := json.Marshal(roleInfo); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		log.Println("Login error:", err)
		return
	} else {
		w.Header().Set("Content-Type", "application/json")
		w.Write(jsonResp)
	}
}

func (h *Handler) Register(w http.ResponseWriter, r *http.Request) {
	setHeaders(w)
	var user domain.User
	var roleInfo domain.RoleInfo

	if err := json.NewDecoder(r.Body).Decode(&user); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		log.Println("Register error:", err)
		return
	}

	if err := h.queuesService.Register(context.TODO(), &user); err != nil {
		if errors.Is(err, e.ErrUserExists) {
			roleInfo.Role = "user exists"
			if jsonResp, err := json.Marshal(roleInfo); err != nil {
				w.WriteHeader(http.StatusInternalServerError)
				log.Println("Register error:", err)
				return
			} else {
				w.Header().Set("Content-Type", "application/json")
				w.Write(jsonResp)
			}
		} else {
			w.WriteHeader(http.StatusInternalServerError)
			log.Println("Register error:", err)
			return
		}
	} else { // user registered success
		role, err := h.queuesService.LogIn(context.TODO(), user.Login, user.Password)
		roleInfo.Role = role
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			log.Println("Register error:", err)
			return
		}
		if jsonResp, err := json.Marshal(roleInfo); err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			log.Println("Register error:", err)
			return
		} else {
			w.Header().Set("Content-Type", "application/json")
			w.Write(jsonResp)
		}
	}
}

func (h *Handler) GetIdByLogin(w http.ResponseWriter, r *http.Request){
	vars := mux.Vars(r)
	loginStr := vars["login"]

	var idStr domain.IdInfo
	id, err := h.queuesService.GetIdByLogin(context.TODO(), loginStr)

	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		log.Println("GetIdByLogin error:", err)
		return
	}
	idStr.Id = id
	if jsonResp, err := json.Marshal(idStr); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		log.Println("GetIdByLogin error:", err)
		return
	} else {
		w.Header().Set("Content-Type", "application/json")
		w.Write(jsonResp)
	}
}