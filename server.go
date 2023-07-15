package main

import (
	"net/http"
	"fmt"
	"io/ioutil"
	"context"
	"net/http/pprof"
	"time"
	"sync"
)

// server describes http server
type Server struct {
	http http.Server
	mux *http.ServeMux
}

const (
	httpApiTimeout  = time.Second * 3
	shutdownTimeout = time.Second * 10
)

func New(port int) *Server {
	s   := &Server{}
	mux := http.NewServeMux()

	mux.Handle("/debug/pprof", http.HandlerFunc(pprof.Index))
	mux.Handle("/debug/pprof/cmdline", http.HandlerFunc(pprof.Cmdline))
	mux.Handle("/debug/pprof/profile", http.HandlerFunc(pprof.Profile))
	mux.Handle("/debug/pprof/symbol", http.HandlerFunc(pprof.Symbol))
	mux.Handle("/debug/pprof/trace", http.HandlerFunc(pprof.Trace))

	s.mux = mux

	s.http = http.Server{
		Addr: fmt.Sprintf(":%d", port),
		Handler: http.TimeoutHandler(s.mux, httpApiTimeout, ""),
	}
	return s
}

// Run starts the HTTP server
func (s *Server) Run(stopCh <-chan struct{}, wg *sync.WaitGroup) {
	defer wg.Done()
	var err error
	// start 
	go func() {
		if err2 := s.http.ListenAndServe(); err2 != http.ErrServerClosed {
			panic(fmt.Errorf("could not start http server: %v", err2))
		}
	}()
	fmt.Printf("Listening on %s \n", s.http.Addr)
	// stop
	<-stopCh // wait for signal to be written in the channel
	// turning the server off
	ctx, cancel := context.WithTimeout(context.Background(), shutdownTimeout)
	defer cancel()
	if err = s.http.Shutdown(ctx); err == nil {
		fmt.Println("Http server shut down")
		return
	}
	if err == context.DeadlineExceeded {
		fmt.Println("Shutdown timeout exceeded. closing http server")
		if err = s.http.Close(); err != nil {
			fmt.Printf("could not close http connection: %v \n", err)
		}
		return
	}
}

func (s *Server) HandleFunc(pattern string, handler func(http.ResponseWriter, *http.Request) ) {
	s.mux.HandleFunc(pattern, handler)
}

func main() {
	wg := new(sync.WaitGroup)
	wg.Add(1)
	
	server := New(8000)
	server.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		// mimicking that the request takes 4s
		// <-time.After(time.Second * 4)
		// fmt.Println(r.Header)
		b,_ := ioutil.ReadFile("resources/go.png")
		w.Header().Add("content-type", "image/png")
		w.Write(b)
	})

	stop := make(chan struct{})
	go server.Run(stop, wg)

	go func ()  {
		<-time.After(time.Second * 30)
		stop <- struct{}{}
	}()

	wg.Wait()
	fmt.Println("Party is over")
}