package user

// ConfigService provides access to user configuration settings.
type ConfigService interface {
	// GetUserConfig retrieves the configuration preferences for the given user.
	GetUserConfig(userID int) (Config, error)
}

type configService struct {
	userConfigRepo ConfigRepository
}

func (s *configService) GetUserConfig(userID int) (Config, error) {
	uc := Config{}
	err := s.userConfigRepo.GetUserConfigByUserID(userID, &uc)
	return uc, err
}

func NewUserConfigService(userConfigRepo ConfigRepository) ConfigService {
	return &configService{
		userConfigRepo: userConfigRepo,
	}
}
