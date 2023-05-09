package handler

import (
	"net/http"

	"github.com/aidarbek436/forum/internal/repository"
)

type handler struct {
	storage repository.Storage
}

func NewHandler(storage repository.Storage) *handler {
	return &handler{
		storage: storage,
	}
}

func (h *handler) Router() http.Handler { // DONE
	mux := http.NewServeMux()
	mux.HandleFunc("/", middleware(h.RootHandle))
	mux.HandleFunc("/login", middleware(h.LoginHandle))
	mux.HandleFunc("/register", middleware(h.RegistrationHandle))
	mux.HandleFunc("/logout", middleware(h.LogoutHandle))
	mux.HandleFunc("/createpost", middleware(h.CreatepostHandle))
	mux.HandleFunc("/showpost/", middleware(h.ShowpostHandle))
	mux.HandleFunc("/createdposts", middleware(h.CreatedPostsHandle))
	mux.HandleFunc("/likedposts", middleware(h.LikedPostsHandle))
	return mux
}
