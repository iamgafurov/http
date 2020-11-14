package banners
import (
	"context"
	"sync"
	"errors"
	"strconv"
	//"log"
	//"net/http"
)

type Service struct {
	mu sync.RWMutex
	items []*Banner
	nextID int64
}

func NewService() *Service{
	return &Service{items: make([]*Banner, 0)}
}

type Banner struct {
	ID 	int64
	Title	string
	Content	string
	Button	string
	Link	string
	Image 	string
}

//All Return all existing banners
func (s *Service) All(ctx context.Context)([]*Banner, error){
	return s.items,nil
}

// ByID return banner by id
func (s *Service) ByID(ctx context.Context, id int64)(*Banner, error){
	s.mu.RLock()
	defer s.mu.RUnlock()
	for _,banner := range s.items {
		if banner.ID == id {
			return banner, nil
		}
	}

	return nil, errors.New("item not found")
}

func (s *Service) Save(ctx context.Context, banner *Banner)(*Banner,error){
	
	s.mu.RLock()
	defer s.mu.RUnlock()
	if banner.ID ==0 {
		s.nextID++
		banner.ID = s.nextID
		banner.Image = strconv.FormatInt(banner.ID,10) + "." +  banner.Image 
		s.items = append(s.items,banner) 
		return banner,nil
	}
	for i,item := range s.items {
		if banner.ID == item.ID {
			banner.Image = item.Image
			s.items[i] = banner
			return banner,nil
		}
	}

	return nil,errors.New("item not found") 
}
func (s *Service) RemoveById(ctx context.Context, id int64)(*Banner,error){
	s.mu.RLock()
	defer s.mu.RUnlock()
	for i,banner := range s.items {
		if banner.ID == id {
			if i == len(s.items){
				s.items = s.items[:i]
			}else{
				s.items = append(s.items[:i],s.items[i+1:]...)
			}
			return banner, nil
		}
	}
	return nil, errors.New("item not found")
}

func (s *Service) AllBanners()[]*Banner{
	return s.items
}
