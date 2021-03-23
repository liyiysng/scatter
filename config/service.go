package config

// SubService 子服务
type SubService struct {
	SelectPolicy string
}

// Service 服务
type Service struct {
	SelectPolicy string
	Meta         map[string]string
}
