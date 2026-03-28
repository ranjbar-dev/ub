package user

const (
	UserLoginHistoryTypeSuccessful = "SUCCESSFUL"
	UserLoginHistoryTypeFailed     = "FAILED"
)

// LoginHistoryService provides operations for recording and querying user login history.
type LoginHistoryService interface {
	// CreateLoginHistory persists a new login history entry (successful or failed).
	CreateLoginHistory(loginHistory *LoginHistory) error
	// GetLastLoginHistoryByUserID retrieves the most recent login record for the given user.
	GetLastLoginHistoryByUserID(userID int, loginHistory *LoginHistory) error
}

type loginHistoryService struct {
	loginHistoryRepository LoginHistoryRepository
}

func (s *loginHistoryService) CreateLoginHistory(loginHistory *LoginHistory) error {
	return s.loginHistoryRepository.Create(loginHistory)
}

func (s *loginHistoryService) GetLastLoginHistoryByUserID(userID int, loginHistory *LoginHistory) error {
	return s.loginHistoryRepository.GetLastLoginHistoryByUserID(userID, loginHistory)
}

func NewLoginHistoryService(loginHistoryRepository LoginHistoryRepository) LoginHistoryService {
	return &loginHistoryService{
		loginHistoryRepository: loginHistoryRepository,
	}
}
