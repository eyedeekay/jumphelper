package jumphelper

import (
	"fmt"
	"log"
	"net/http"
	"strings"

	"golang.org/x/time/rate"
)

// Server is a TCP service that responds to addressbook requests
type Server struct {
	host string
	port string

	addressBookPath string
	jumpHelper      *JumpHelper
	localService    *http.ServeMux

	limiter *rate.Limiter

	err error
}

func (s *Server) address() string {
	return s.host + ":" + s.port
}

func (s *Server) limit(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if s.limiter.Allow() == false {
			http.Error(w, http.StatusText(429), http.StatusTooManyRequests)
			return
		}

		next.ServeHTTP(w, r)
	})
}

// Serve sets up a listening server on the specified port
func (s *Server) Serve() {
	s.localService, s.err = s.NewMux()
	if s.err != nil {
		log.Fatal(s.err)
	}
	s.err = http.ListenAndServe(s.address(), s.limit(s.localService))
	if s.err != nil {
		log.Fatal(s.err)
	}
}

// HandleExists prints true:address if an antecedent URL exists in the addressbook, false if not
func (s *Server) HandleExists(w http.ResponseWriter, r *http.Request) {
	p := strings.TrimPrefix(strings.Replace(r.URL.Path, "check/", "", 1), "/")
	if s.jumpHelper.CheckAddressBook(p) {
		fmt.Fprintln(w, "TRUE", p)
		return
	}
	fmt.Fprintln(w, "FALSE", p)
	return
}

// HandleJump redirects to a b32.i2p URL instead of behaving like a traditional jump service.
func (s *Server) HandleJump(w http.ResponseWriter, r *http.Request) {
	p := strings.TrimPrefix(r.URL.Path, "/")
	if s.jumpHelper.CheckAddressBook(p) {
		line := "http://" + s.jumpHelper.SearchAddressBook(p)[1] + ".b32.i2p"
		w.Header().Set("Location", line)
		w.WriteHeader(301)
		fmt.Fprintln(w, line)
		return
	}
	fmt.Fprintln(w, "FALSE")
	return
}

// NewMux sets up a new ServeMux with handlers
func (s *Server) NewMux() (*http.ServeMux, error) {
	s.localService = http.NewServeMux()
	s.localService.Handle("/check/", http.HandlerFunc(s.HandleExists))
	s.localService.Handle("/", http.HandlerFunc(s.HandleJump))
	if s.err != nil {
		return nil, fmt.Errorf("Local mux configuration error: %s", s.err)
	}
	return s.localService, nil
}

// NewServer creates a new Server that answers jump-related queries
func NewServer(Host, Port, addressBookPath string) (*Server, error) {
	return NewServerFromOptions(SetServerHost(Host), SetServerPort(Port), SetServerAddressBookPath(addressBookPath))
}

// NewServerFromOptions creates a new Server that answers jump-related queries
func NewServerFromOptions(opts ...func(*Server) error) (*Server, error) {
	var s Server
	s.host = "127.0.0.1"
	s.port = "7054"
	s.addressBookPath = "/var/lib/i2pd/addressbook/addresses.csv"
	for _, o := range opts {
		if err := o(&s); err != nil {
			return nil, fmt.Errorf("Service configuration error: %s", err)
		}
	}
	s.limiter = rate.NewLimiter(1, 1)
	s.jumpHelper, s.err = NewJumpHelper(s.addressBookPath)
	if s.err != nil {
		return nil, fmt.Errorf("Jump helper load error: %s", s.err)
	}
	return &s, s.err
}
