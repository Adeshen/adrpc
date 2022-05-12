package adrpc

func NewServer(port string) *Server {
	server := Server{
		Services: make(map[string]*service),
	}
	return &server
}

func (server *Server) Register(newStruct interface{}) {

}
