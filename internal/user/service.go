package user

type UserService struct {
 repo UserRepository
}

func NewUserService(repo *UserRepository) *UserService {
	return &UserService{repo: *repo}
}

func (svc *UserService) GetAllUsers() ([]User, error) {
	return svc.repo.GetAll();
}