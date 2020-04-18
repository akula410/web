package server

import (
	"context"
	"fmt"
	"google.golang.org/grpc"
	"net"
	"net/http"
)

type Server struct {
	Network string
	Address string
	message []string
	lst net.Listener
	srv *grpc.Server
	isProto bool

	web *http.Server
	webAddress string
	webHandleFunc func(w http.ResponseWriter, r *http.Request)
	webPattern string
	isWeb bool

	done chan struct{}

	end chan bool
}

var start = "start"
var stop = "stop"
var restart = "restart"

func (s *Server) Add(ctx context.Context, req *Request) (*Response, error){
	var command string
	var ok bool
	s.clearMessage()

	if cmd, cmd2 := req.GetCommand(), req.GetCmd(); len(cmd) > 0 {
		command = cmd
	}else if len(cmd2) > 0{
		command = cmd2
	}


	if len(command) > 0 {

		switch command{
		case start:
			s.startWeb()
			ok = true
		case stop:
			//s.end = make(chan bool)
			//s.stopWeb()
			s.Stop()
			ok = true
		case restart:
			s.restartWeb()
			ok = true
		}

		if s.end != nil {
			<-s.end
			s.end = nil
		}
	}
	var message = s.getMessage()
	s.clearMessage()

	return &Response{Result: ok, Message:message}, nil
}

func (s *Server) HandleFunc(address string, pattern string, handle func(w http.ResponseWriter, r *http.Request)) *Server{
	s.webAddress = address
	s.webPattern = pattern
	s.webHandleFunc = handle
	return s
}

func (s *Server) isDone(){
	if s.done == nil{
		s.done = make(chan struct{})
	}
}

func (s *Server) Start(){
	if !s.isProto{
		s.isDone()
		s.isProto = true
		go func()bool{
			var err error
			s.srv = grpc.NewServer()

			RegisterApiServer(s.srv, s)

			s.lst, err = net.Listen(s.Network, s.Address)

			if err != nil {
				s.isProto = false
				s.setMessage(fmt.Sprintf("Fatal error 1: Proto server stop...., %s", err.Error()))
				s.stopWeb()
				return false
			}

			err = s.srv.Serve(s.lst)

			if err != nil {
				s.isProto = false
				s.setMessage(fmt.Sprintf("Fatal error 2: Proto server stop...., %s", err.Error()))
				s.stopWeb()
				return false
			}
			s.done <- struct{}{}
			return true
		}()
		s.setMessage("Proto server start....")
	}

	http.HandleFunc(s.webPattern, s.webHandleFunc)
	s.startWeb()
}

func (s *Server) Stop(){
	if s.srv != nil{
		s.setMessage("Proto server close")
		s.srv.Stop()
	}
}

func (s *Server) startWeb(){
	if !s.isWeb {
		s.web = &http.Server{Addr:s.webAddress}
		go func(){
			s.isWeb = true
			err := s.web.ListenAndServe()
			if err != nil {
				s.setMessage("Web server stop....")
				s.isWeb = false
				if s.end != nil {
					s.end <- true
				}
			}
		}()
		s.setMessage("Web server start....")
	}else{
		s.setMessage("Web server has been start early")
	}

}

func (s *Server) stopWeb(){
	if s.web != nil && s.isWeb {
		_ = s.web.Shutdown(context.Background())
		s.isWeb = false
	}else{
		s.setMessage("Web server has been stop early")
	}
}

func (s *Server) restartWeb(){
	s.stopWeb()
	s.startWeb()
}

func (s *Server) clearMessage(){
	s.message = make([]string, 0)
}

func (s *Server) setMessage(message string){
	if s.message == nil {
		s.message = make([]string, 0)
	}
	s.message = append(s.message, message)
}

func (s *Server) getMessage()[]string{
	return s.message
}

func (s *Server) Block(){
	<-s.done
}




