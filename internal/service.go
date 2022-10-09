package internal

type AggregatorService struct {
	fs *FileService
}

func NewService(fs *FileService) *AggregatorService {
	return &AggregatorService{fs: fs}
}

func (s *AggregatorService) Init() error {
	return s.fs.Init()
}

func (s *AggregatorService) Get() []Category {
	return s.fs.Get()
}
