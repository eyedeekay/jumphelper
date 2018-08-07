package jumphelper

import (
	"fmt"
	"log"
	"net/http"
	"strings"
	"time"
)

import (
	"github.com/LarryBattle/nonce-golang"
	"github.com/bwesterb/go-pow"
	"github.com/eyedeekay/gosam"
	"golang.org/x/time/rate"
)

// Server is a TCP service that responds to addressbook requests
type Server struct {
	host    string
	port    string
	samHost string
	samPort string

	pusher    *goSam.Client
	transport *http.Transport
	client    *http.Client

	addressBookPath  string
	jumpHelper       *JumpHelper
	localService     *http.ServeMux
	ext              bool
	verbose          bool
	subscriptionURLs []string
	listing          bool
	base32           string
	difficulty       int

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

// ServeLocal sets up a listening server on the specified port
func (s *Server) ServeLocal() {
	s.localService, s.err = s.NewMux()
	if s.err != nil {
		log.Fatal(s.err)
	}
	s.err = http.ListenAndServe(s.address(), s.limit(s.localService))
	if s.err != nil {
		log.Fatal(s.err)
	}
}

// Serve sets up a listening server on the specified port
func (s *Server) Serve() {
	s.ServeLocal()
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

// HandleLookup redirects to a b32.i2p URL instead of behaving like a traditional jump service.
func (s *Server) HandleLookup(w http.ResponseWriter, r *http.Request) {
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

// HandleJump redirects to a base64 URL like a traditional jump service.
func (s *Server) HandleJump(w http.ResponseWriter, r *http.Request) {
	p := strings.TrimPrefix(strings.Replace(r.URL.Path, "jump/", "", 1), "/")
	if s.jumpHelper.SearchAddressBook(p) != nil {
		array := s.jumpHelper.SearchAddressBook(p)
		if len(array) == 3 {
			line := "http://" + s.host + "/?i2paddresshelper=" + array[2]
			w.Header().Set("Location", line)
			w.WriteHeader(301)
			fmt.Fprintln(w, line)
			return
		}
		fmt.Fprintln(w, "no, it's me, dave, man. let me up")
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

// HandleBase32 lists all synced remote jumphelper urls.
func (s *Server) HandleBase32(w http.ResponseWriter, r *http.Request) {
	if s.listing {
		fmt.Fprintln(w, s.base32)
		return
	}
	fmt.Fprintln(w, "Listings disabled for this server")
	return
}

// HandlePush creates a signed list of addresses and pushes it to a requested URL
func (s *Server) HandlePush(w http.ResponseWriter, r *http.Request) {
	if s.listing {
		p := strings.TrimPrefix(strings.Replace(r.URL.Path, "push/", "", 1), "/")
		if p != "" {
			send, err := http.NewRequest("POST", p, strings.NewReader(strings.Join(s.jumpHelper.Subs(), ",")))
			if err != nil {
				return
			}
			s.client.Do(send)
			fmt.Fprintln(w, "Your push was sent to", p)
			return
		}
		fmt.Fprintln(w, "FALSE")
		return
	}
	fmt.Fprintln(w, "Listings disabled for this server")
	return
}

// HandleRecv recieves a signed list of URL's from another server's HandlePush
func (s *Server) HandleRecv(w http.ResponseWriter, r *http.Request) {
	if s.listing {
		p := strings.TrimPrefix(strings.Replace(r.URL.Path, "recv/", "", 1), "/")
		if p != "" {
			fmt.Fprintln(w, "I recieved a push from:",
				r.Header.Get("X-I2p-Destb32"),
				"And for now, I did nothing with it because I am dumb)")
			return
		}
		fmt.Fprintln(w, "FALSE")
		return
	}
	fmt.Fprintln(w, "Listings disabled for this server")
	return
}

// HandleProof emits a problem for proof-of-work on the client
func (s *Server) HandleProof(w http.ResponseWriter, r *http.Request) {
	if s.listing {
		fmt.Fprintln(w, pow.NewRequest(uint32(s.difficulty), []byte(nonce.NewToken())))
		return
	}
	fmt.Fprintln(w, "Listings disabled for this server")
	return
}

// HandleAccount emits a problem for proof-of-work on the client
func (s *Server) HandleAccount(w http.ResponseWriter, r *http.Request) {
	if s.listing {
		p := strings.TrimPrefix(strings.Replace(r.URL.Path, "acct/", "", 1), "/")
		if p != "" {
			reqproof := strings.SplitN(p, ",", 4)
			if len(reqproof) == 4 {
				ok, err := pow.Check(reqproof[0], reqproof[1], []byte(reqproof[2]))
				if err != nil {
					fmt.Fprintln(w, err.Error())
					return
				}
				if ok {
					s.jumpHelper.TrustedAddressBook.AddAddress(reqproof[2], reqproof[3])
					fmt.Fprintln(w, "proof-of-work valid")
					return
				}
				fmt.Fprintln(w, "proof-of-work invalid")
				return
			}
			fmt.Fprintln(w, "Invalid length for proof-of-work check=", len(reqproof))
			return
		}
		fmt.Fprintln(w, pow.NewRequest(uint32(s.difficulty*2), []byte(nonce.NewToken())))
		return
	}
	fmt.Fprintln(w, "Listings disabled for this server")
	return
}

// NewMux sets up a new ServeMux with handlers
func (s *Server) NewMux() (*http.ServeMux, error) {
	s.localService = http.NewServeMux()
	s.localService.HandleFunc("/check/", s.HandleExists)
	s.localService.HandleFunc("/request/", s.HandleLookup)
	s.localService.HandleFunc("/jump/", s.HandleJump)
	s.localService.HandleFunc("/sub/", s.HandleListing)
	s.localService.HandleFunc("/addr/", s.HandleBase32)
	s.localService.HandleFunc("/push/", s.HandlePush)
	s.localService.HandleFunc("/recv/", s.HandleRecv)
	s.localService.HandleFunc("/acct/", s.HandleAccount)
	//s.localService.HandleFunc("/update/", s.HandleUpdate)
	s.localService.HandleFunc("/pow", s.HandleProof)
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
func NewServer(host, port, book, samhost, samport string, subs []string, useh, verbose, share bool, base32 string) (*Server, error) {
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
		SetServerBase32(base32),
	)
}

// NewServerFromOptions creates a new Server that answers jump-related queries
func NewServerFromOptions(opts ...func(*Server) error) (*Server, error) {
	var s Server
	s.host = "127.0.0.1"
	s.port = "7854"
	s.samHost = "127.0.0.1"
	s.samPort = "7656"
	s.addressBookPath = "/var/lib/i2pd/addressbook/addresses.csv"
	s.rate = 1
	s.burst = 1
	s.ext = true
	s.verbose = false
	s.listing = false
	s.base32 = ""
	s.difficulty = 1
	s.subscriptionURLs = []string{"http://joajgazyztfssty4w2on5oaqksz6tqoxbduy553y34mf4byv6gpq.b32.i2p/export/alive-hosts.txt"}
	for _, o := range opts {
		if err := o(&s); err != nil {
			return nil, fmt.Errorf("Service configuration error: %s", err)
		}
	}
	s.limiter = rate.NewLimiter(s.Rate(), s.burst)
	log.Println("Configured Rate Limiter")
	if s.listing {
		s.pusher, s.err = goSam.NewClientFromOptions(
			goSam.SetHost(s.samHost),
			goSam.SetPort(s.samPort),
			goSam.SetUnpublished(true),
			goSam.SetInLength(uint(3)),
			goSam.SetOutLength(uint(3)),
			goSam.SetInQuantity(uint(6)),
			goSam.SetOutQuantity(uint(6)),
			goSam.SetInBackups(uint(2)),
			goSam.SetOutBackups(uint(2)),
			goSam.SetCloseIdle(true),
			goSam.SetCloseIdleTime(uint(300000)),
		)
		if s.err != nil {
			return nil, s.err
		}
		s.transport.Dial = s.pusher.Dial
		s.client.Transport = s.transport
	}
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
		SetServerPort("7854"),
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
