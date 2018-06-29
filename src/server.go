package jumphelper

import (
	"fmt"
	"log"
	"net/http"
	"strings"

	"golang.org/x/time/rate"
	"time"
)

// Server is a TCP service that responds to addressbook requests
type Server struct {
	host    string
	port    string
	samHost string
	samPort string

	addressBookPath  string
	jumpHelper       *JumpHelper
	localService     *http.ServeMux
	ext              bool
	verbose          bool
	subscriptionURLs []string
	listing          bool

	rate  int
	burst int

	limiter *rate.Limiter
	err     error
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
	p := strings.TrimPrefix(strings.Replace(r.URL.Path, "request/", "", 1), "/")
	if s.jumpHelper.SearchAddressBook(p) != nil {
		line := "http://" + s.jumpHelper.SearchAddressBook(p)[1] + ".b32.i2p"
		w.Header().Set("Location", line)
		w.WriteHeader(301)
		fmt.Fprintln(w, line)
		return
	}
	fmt.Fprintln(w, "FALSE")
	return
}

// HandleListing lists all synced remote jumphelper urls.
func (s *Server) HandleListing(w http.ResponseWriter, r *http.Request) {
	if s.listing {
		for _, s := range s.jumpHelper.Subs() {
			fmt.Fprintln(w, s)
		}
		return
	}
	fmt.Fprintln(w, "Listings disabled for this server")
	return
}

// NewMux sets up a new ServeMux with handlers
func (s *Server) NewMux() (*http.ServeMux, error) {
	s.localService = http.NewServeMux()
	s.localService.HandleFunc("/check/", s.HandleExists)
	s.localService.HandleFunc("/request/", s.HandleJump)
	s.localService.HandleFunc("/sub/", s.HandleListing)
	s.localService.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintln(w, "Dave's not here man.")
	})
	if s.err != nil {
		return nil, fmt.Errorf("Local mux configuration error: %s", s.err)
	}
	return s.localService, nil
}

// Rate returns the rate
func (s *Server) Rate() rate.Limit {
	r := time.Duration(s.rate) * time.Second
	return rate.Every(r)
}

// NewServer creates a new Server that answers jump-related queries
func NewServer(host, port, book, samhost, samport string, subs []string, useh, verbose, share bool) (*Server, error) {
	return NewServerFromOptions(
		SetServerHost(host),
		SetServerPort(port),
		SetServerAddressBookPath(book),
		SetServerJumpHelperHost(samhost),
		SetServerJumpHelperPort(samport),
		SetServerUseHelper(useh),
		SetServerSubscription(subs),
		SetServerJumpHelperVerbosity(verbose),
		SetServerEnableListing(share),
	)
}

// NewServerFromOptions creates a new Server that answers jump-related queries
func NewServerFromOptions(opts ...func(*Server) error) (*Server, error) {
	var s Server
	s.host = "127.0.0.1"
	s.port = "7054"
	s.samHost = "127.0.0.1"
	s.samPort = "7056"
	s.addressBookPath = "/var/lib/i2pd/addressbook/addresses.csv"
	s.rate = 1
	s.burst = 1
	s.ext = true
	s.verbose = false
	s.listing = false
	s.subscriptionURLs = []string{"http://joajgazyztfssty4w2on5oaqksz6tqoxbduy553y34mf4byv6gpq.b32.i2p/export/alive-hosts.txt"}
	for _, o := range opts {
		if err := o(&s); err != nil {
			return nil, fmt.Errorf("Service configuration error: %s", err)
		}
	}
	s.limiter = rate.NewLimiter(s.Rate(), s.burst)
	log.Println("Configured Rate Limiter")
	s.jumpHelper, s.err = NewJumpHelperFromOptions(
		SetJumpHelperAddressBookPath(s.addressBookPath),
		SetJumpHelperHost(s.samHost),
		SetJumpHelperPort(s.samPort),
		SetJumpHelperUseHelper(s.ext),
		SetJumpHelperVerbosity(s.verbose),
	)
	if len(s.subscriptionURLs) < 1 {
		s.ext = false
	}
	if s.err != nil {
		return nil, fmt.Errorf("Jump helper load error: %s", s.err)
	}
	return &s, s.err
}

// Service quickly generates a service with the defaults.
func Service() {
	s, err := NewServerFromOptions()
	if err != nil {
		log.Fatal(err, "Error starting server")
	}
	go s.Serve()
}

// NewService quickly generates a service with host, port, book strings and fires off a goroutine
func NewService(host, port, book, samhost, samport string, subs []string, useh bool) {
	s, err := NewServerFromOptions(
		SetServerHost(host),
		SetServerPort(port),
		SetServerAddressBookPath(book),
		SetServerUseHelper(useh),
		SetServerJumpHelperHost(samhost),
		SetServerJumpHelperPort(samport),
		SetServerSubscription(subs),
	)
	if err != nil {
		log.Fatal(err, "Error starting server")
	}
	go s.Serve()
}

func service() {
	s, err := NewServerFromOptions(
		SetServerHost("0.0.0.0"),
		SetServerPort("7054"),
		SetServerAddressBookPath("../addresses.csv"),
		SetServerRate(1000),
		SetServerBurst(1000),
		SetServerUseHelper(false),
		SetServerJumpHelperHost("127.0.0.1"),
		SetServerJumpHelperPort("7656"),
		SetServerSubscription([]string{"http://joajgazyztfssty4w2on5oaqksz6tqoxbduy553y34mf4byv6gpq.b32.i2p/export/alive-hosts.txt"}),
	)
	if err != nil {
		log.Fatal(err, "Error starting server")
	}
	go s.Serve()
}
