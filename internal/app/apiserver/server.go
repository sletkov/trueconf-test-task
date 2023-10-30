package apiserver

import (
	"net/http"
	"refactoring/internal/app/models"
	"refactoring/internal/app/store"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/render"
)

type server struct {
	router *chi.Mux
	store  store.Store
}

func newServer(store store.Store) *server {
	s := &server{
		router: chi.NewRouter(),
		store:  store,
	}

	s.configureRouter()

	return s
}

func (s *server) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	s.router.ServeHTTP(w, r)
}

func (s *server) configureRouter() {

	s.router.Use(middleware.RequestID)
	s.router.Use(middleware.RealIP)
	s.router.Use(middleware.Logger)
	// s.router.Use(middleware.Recoverer)
	s.router.Use(middleware.Timeout(60 * time.Second))

	s.router.Get("/", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte(time.Now().String()))
	})

	s.router.Route("/api", func(r chi.Router) {
		s.router.Route("/v1", func(r chi.Router) {
			s.router.Route("/users", func(r chi.Router) {
				s.router.Get("/", s.handleSearchUsers())
				s.router.Post("/", s.handleCreateUser())

				s.router.Route("/{id}", func(r chi.Router) {
					s.router.Get("/", s.handleGetUser())
					s.router.Patch("/", s.handleUpdateUser())
					s.router.Delete("/", s.handleDeleteUser())
				})
			})
		})
	})

}

func (srv *server) handleSearchUsers() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		userList, _ := srv.store.User().SearchUsers()
		render.JSON(w, r, userList)
	}
}

type CreateUserRequest struct {
	DisplayName string `json:"display_name"`
	Email       string `json:"email"`
}

func (c *CreateUserRequest) Bind(r *http.Request) error {
	return nil
}

func (srv *server) handleCreateUser() http.HandlerFunc {

	return func(w http.ResponseWriter, r *http.Request) {

		req := &CreateUserRequest{}

		if err := render.Bind(r, req); err != nil {
			_ = render.Render(w, r, ErrInvalidRequest(err))
			return
		}

		u := &models.User{
			CreatedAt:   time.Now(),
			DisplayName: req.DisplayName,
			Email:       req.DisplayName,
		}

		id, _ := srv.store.User().CreateUser(u)

		render.Status(r, http.StatusCreated)
		render.JSON(w, r, map[string]interface{}{
			"user_id": id,
		})
	}

}

func (srv *server) handleGetUser() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		id := chi.URLParam(r, "id")
		u, _ := srv.store.User().GetUser(id)
		render.JSON(w, r, u)
	}

}

type UpdateUserRequest struct {
	*models.User
}

func (u *UpdateUserRequest) Bind(r *http.Request) error { return nil }

func (srv *server) handleUpdateUser() http.HandlerFunc {

	return func(w http.ResponseWriter, r *http.Request) {

		req := &UpdateUserRequest{}

		if err := render.Bind(r, req); err != nil {
			_ = render.Render(w, r, ErrInvalidRequest(err))
			return
		}

		id := chi.URLParam(r, "id")

		if err := srv.store.User().UpdateUser(id, req.User); err != nil {
			_ = render.Render(w, r, ErrInvalidRequest(store.ErrUserNotFound))
			return
		}

		render.Status(r, http.StatusNoContent)
	}
}

func (srv *server) handleDeleteUser() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		id := chi.URLParam(r, "id")

		if err := srv.store.User().DeleteUser(id); err != nil {
			_ = render.Render(w, r, ErrInvalidRequest(store.ErrUserNotFound))
			return
		}

		render.Status(r, http.StatusNoContent)
	}

}

type ErrResponse struct {
	Err            error `json:"-"`
	HTTPStatusCode int   `json:"-"`

	StatusText string `json:"status"`
	AppCode    int64  `json:"code,omitempty"`
	ErrorText  string `json:"error,omitempty"`
}

func (e *ErrResponse) Render(w http.ResponseWriter, r *http.Request) error {
	render.Status(r, e.HTTPStatusCode)
	return nil
}

func ErrInvalidRequest(err error) render.Renderer {
	return &ErrResponse{
		Err:            err,
		HTTPStatusCode: 400,
		StatusText:     "Invalid request.",
		ErrorText:      err.Error(),
	}
}
