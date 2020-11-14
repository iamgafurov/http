package app

import (
	"net/http"
	"github.com/iamgafurov/http/pkg/banners"
	"log"
	"encoding/json"
	"strconv"
	"io/ioutil"
	"fmt"
	"os"
	"strings"
	"path/filepath"
)

type Server struct{
	mux *http.ServeMux
	bannersSvc *banners.Service
}

func NewServer(mux *http.ServeMux, bannersSvc *banners.Service) *Server{
	return &Server{mux:mux, bannersSvc: bannersSvc}
}

func (s *Server) ServeHTTP(writer http.ResponseWriter, request *http.Request){
	s.mux.ServeHTTP(writer, request)
}

func (s *Server) Init() {
	s.mux.HandleFunc("/banners.getAll",s.handleGetAllBanners)
	s.mux.HandleFunc("/banners.getById",s.handleGetBannerByID)
	s.mux.HandleFunc("/banners.save", s.handleSaveBanner)
	s.mux.HandleFunc("/banners.removeById",s.handleRemoveByID)
}

func (s *Server) handleGetBannerByID(writer http.ResponseWriter, request *http.Request){
	idParam := request.URL.Query().Get("id")
	id, err := strconv.ParseInt(idParam,10,64)
	if err != nil {
		log.Print(err)
		http.Error(writer, http.StatusText(http.StatusInternalServerError),http.StatusInternalServerError)
		return
	}

	item,err := s.bannersSvc.ByID(request.Context(),id)
	if err != nil {
		log.Print(err)
		http.Error(writer, http.StatusText(http.StatusInternalServerError),http.StatusInternalServerError)
		return
	}

	data,err := json.Marshal(item)
	if err != nil {
		log.Print(err)
		http.Error(writer, http.StatusText(http.StatusInternalServerError),http.StatusInternalServerError)
		return
	}

	writer.Header().Set("Content-Type", "application/json")
	_, err = writer.Write(data)
	if err != nil {
		log.Print(err)
		http.Error(writer, http.StatusText(http.StatusInternalServerError),http.StatusInternalServerError)
		return
	}
}


func (s *Server) handleSaveBanner(writer http.ResponseWriter, request *http.Request){
	idParam := request.FormValue("id")
	log.Print(idParam)
	titleParam := request.FormValue("title")
	log.Print(titleParam)
	contentParam := request.FormValue("content")
	log.Print(contentParam)
	buttonParam := request.FormValue("button")
	log.Print(buttonParam)
	linkParam := request.FormValue("link")
	log.Print(linkParam)
	err := request.ParseMultipartForm(10*1024*1024)
	if err != nil {
		log.Print(err)
	}
	bannerID, err := strconv.ParseInt(idParam,10,64)
		if err != nil {
			log.Print(err)
			http.Error(writer, http.StatusText(http.StatusInternalServerError),http.StatusInternalServerError)
			return
		}
	file,header,err := request.FormFile("image")
	banner := &banners.Banner{}
	if err == nil {
		log.Print(err)
	
		log.Print(header.Filename)
		buf, err := ioutil.ReadAll(file)
		if err != nil {
		log.Print(err)
		}
		
		fmt.Println(buf)
		
		imgExt := strings.Split(header.Filename,".")[1]
		banner = &banners.Banner{
			ID : bannerID,
			Title: titleParam,
			Content: contentParam,
			Button: buttonParam,
			Link: linkParam,
			Image: imgExt,
		}
		banner,err = s.bannersSvc.Save(request.Context(),banner)
		if err != nil {
			log.Print(err)
			http.Error(writer, http.StatusText(http.StatusInternalServerError),http.StatusInternalServerError)
			return
		}
		idParam = strconv.FormatInt(banner.ID,10) 
		
		f,err := os.Create( "./web/banners/" + idParam +"."+ imgExt)
		pth ,_ := filepath.Abs("./")
		log.Print("paththththth:  ",pth)
		if err != nil {
			log.Print(err)
		}
		defer f.Close()
		_, err = f.Write(buf)
		if err != nil {
			log.Print(err)
		}
	}else {
		banner = &banners.Banner{
			ID : bannerID,
			Title: titleParam,
			Content: contentParam,
			Button: buttonParam,
			Link: linkParam,
			Image: "",
		}
		banner,err = s.bannersSvc.Save(request.Context(),banner)
		if err != nil {
			log.Print(err)
			http.Error(writer, http.StatusText(http.StatusInternalServerError),http.StatusInternalServerError)
			return
		}
	}

	if idParam == "" || titleParam =="" || contentParam == "" || buttonParam == "" || linkParam == "" {
		log.Print("incorrect params")
		http.Error(writer, http.StatusText(http.StatusInternalServerError),http.StatusInternalServerError)
		return
	}

	data,err := json.Marshal(banner)
	if err != nil {
		log.Print(err)
		http.Error(writer, http.StatusText(http.StatusInternalServerError),http.StatusInternalServerError)
		return
	}

	writer.Header().Set("Content-Type", "application/json")
	_, err = writer.Write(data)
	if err != nil {
		log.Print(err)
		http.Error(writer, http.StatusText(http.StatusInternalServerError),http.StatusInternalServerError)
		return
	}
}

func (s *Server) handleGetAllBanners(writer http.ResponseWriter, request *http.Request){
	banners := s.bannersSvc.AllBanners()
	data,err := json.Marshal(banners)
	if err != nil {
		log.Print(err)
		http.Error(writer, http.StatusText(http.StatusInternalServerError),http.StatusInternalServerError)
		return
	}

	writer.Header().Set("Content-Type", "application/json")
	_, err = writer.Write(data)
	if err != nil {
		log.Print(err)
		http.Error(writer, http.StatusText(http.StatusInternalServerError),http.StatusInternalServerError)
		return
	}
}

func (s *Server) handleRemoveByID(writer http.ResponseWriter, request *http.Request){
	idParam := request.URL.Query().Get("id")
	id, err := strconv.ParseInt(idParam,10,64)
	if err != nil {
		log.Print(err)
		http.Error(writer, http.StatusText(http.StatusInternalServerError),http.StatusInternalServerError)
		return
	}
	banner, err := s.bannersSvc.RemoveById(request.Context(),id)
	if err != nil {
		log.Print(err)
		http.Error(writer, http.StatusText(http.StatusInternalServerError),http.StatusInternalServerError)
		return
	}
	data,err := json.Marshal(banner)
	if err != nil {
		log.Print(err)
		http.Error(writer, http.StatusText(http.StatusInternalServerError),http.StatusInternalServerError)
		return
	}

	writer.Header().Set("Content-Type", "application/json")
	_, err = writer.Write(data)
	if err != nil {
		log.Print(err)
		http.Error(writer, http.StatusText(http.StatusInternalServerError),http.StatusInternalServerError)
		return
	}
}