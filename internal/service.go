package internal

import "simplydash/internal/models"

type Service struct {
	fw         *FileWatcher
	fwChannels []chan string
}

func NewService(fw *FileWatcher) *Service {
	return &Service{fw: fw, fwChannels: make([]chan string, 0)}
}

func (s *Service) Init() {
	s.fw.Load()
	s.fw.Watch(&s.fwChannels)
}

func (s *Service) AddUpdateChannel(c chan string) {
	s.fwChannels = append(s.fwChannels, c)
}

func (s *Service) RemoveUpdateChannel(ch chan string) {
	for i, c := range s.fwChannels {
		if ch == c {
			s.fwChannels = append(s.fwChannels[:i], s.fwChannels[i+1:]...)
			return
		}
	}
}

func (s *Service) Get() []models.Category {
	return s.fw.Get()
}
