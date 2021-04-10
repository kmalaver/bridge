package bridge

import (
	"fmt"
	"net/http"
)

type Server struct {
	port               string
	router             *Router
	generalMiddlewares Stack
}

func NewServer(port string) *Server {
	return &Server{
		port:               port,
		router:             NewRouter(),
		generalMiddlewares: Stack{},
	}
}

func (s *Server) Listen() error {
	http.Handle("/", s.router)
	fmt.Println("server listening at http://localhost" + s.port)
	err := http.ListenAndServe(s.port, nil)
	if err != nil {
		return err
	}
	return nil
}

func (s *Server) Handle(path string, method string, handler http.HandlerFunc) {
	_, exist := s.router.rules[path]
	if !exist {
		s.router.rules[path] = make(map[string]http.HandlerFunc)
	}
	h := s.applyMiddlewares(handler)
	s.router.rules[path][method] = h
}

func (s *Server) Use(middleware MiddlewareFunc) {
	s.generalMiddlewares.push(middleware)
}

func (s *Server) applyMiddlewares(handler http.HandlerFunc) http.HandlerFunc {
	s.generalMiddlewares.forEach(func(m interface{}) {
		fmt.Println("mid")
		handler = SetMiddleware(m.(MiddlewareFunc))(handler)
	})
	return handler
}

func (s *Server) AddMiddleware(f http.HandlerFunc, middlewares ...MiddlewareFunc) http.HandlerFunc {

	for i := len(middlewares) - 1; i >= 0; i-- {
		m := middlewares[i]
		// pass handler to each middleware
		f = SetMiddleware(m)(f)
	}
	return f
}

func (s *Server) Get(path string, handler http.HandlerFunc) {
	s.Handle(path, "GET", handler)
}

func (s *Server) Post(path string, handler http.HandlerFunc) {
	s.Handle(path, "POST", handler)
}

func (s *Server) Put(path string, handler http.HandlerFunc) {
	s.Handle(path, "PUT", handler)
}

func (s *Server) Delete(path string, handler http.HandlerFunc) {
	s.Handle(path, "DELETE", handler)
}
