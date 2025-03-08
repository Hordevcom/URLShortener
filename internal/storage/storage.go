package storage

type Storage interface {
	Set(key, value string)
	Get(key string) (string, bool)
}

type MapStorage struct {
	data map[string]string
}

func NewStorage() *MapStorage {
	return &MapStorage{data: make(map[string]string)}
}

func (s *MapStorage) Set(key, value string) {
	s.data[key] = value
}

func (s *MapStorage) Get(key string) (string, bool) {
	value, exist := s.data[key]
	return value, exist
}

type JSONStorage struct {
	URL string `json:"url"`
}

func NewJSONStorage() *JSONStorage {
	return &JSONStorage{URL: ""}
}

func (s *JSONStorage) Set(key, value string) {
	s.URL = value
}

func (s *JSONStorage) Get() string {
	return s.URL
}
