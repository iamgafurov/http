package main
import (
	"net/http"
	"os"
	"log"
	"net"
	"github.com/iamgafurov/http/pkg/banners"
	"github.com/iamgafurov/http/cmd/app"
)

func main(){
	host := "localhost"
	port := "9999"

	if err := execute(host, port); err != nil {
		os.Exit(1)
	}
}

func execute(host string, port string)(err error){
	mux := http.NewServeMux()
	bannersSvc := banners.NewService()
	server := app.NewServer(mux,bannersSvc)
	server.Init()

	srv := &http.Server {
		Addr: net.JoinHostPort(host,port),
		Handler: server,
	}
	log.Print("server run on " +srv.Addr)
	return srv.ListenAndServe()
}