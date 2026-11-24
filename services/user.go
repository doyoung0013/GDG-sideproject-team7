package services

import (
	"context"

	"github.com/GDG-KHU-Side/backend-side-project/apperr"
	"github.com/GDG-KHU-Side/backend-side-project/auth"
	"github.com/GDG-KHU-Side/backend-side-project/models"
	"github.com/GDG-KHU-Side/backend-side-project/repository"
)

// TokenIssuer는 로그인 성공 시 액세스 토큰을 발급합니다.
type TokenIssuer interface {
	Issue(userID int64) (string, error)
}

type UserService struct {
	users  repository.UserRepository
	tokens TokenIssuer
}

func NewUserService(users repository.UserRepository, tokens TokenIssuer) *UserService {
	return &UserService{users: users, tokens: tokens}
}

var errInvalidCredentials = apperr.Unauthorized("이메일 또는 비밀번호가 올바르지 않습니다.")

// 존재하지 않는 이메일로 로그인할 때도 bcrypt 비교 시간을 비슷하게 맞추기 위한 더미 해시
var dummyHash, _ = auth.HashPassword("dummy-password-for-timing")

func (s *UserService) Register(ctx context.Context, req models.RegisterRequest) (*models.User, error) {
	email, err := normalizeEmail(req.Email)
	if err != nil {
		return nil, err
	}
	name, err := validateName(req.Name)
	if err != nil {
		return nil, err
	}
	if err := validatePassword(req.Password); err != nil {
		return nil, err
	}
	phone, err := normalizeOptionalPhone(req.PhoneNum)
	if err != nil {
		return nil, err
	}

	hash, err := auth.HashPassword(req.Password)
	if err != nil {
		return nil, err
	}
	u := &models.User{Email: email, PasswordHash: hash, PhoneNum: phone, Name: name}
	if err := s.users.Create(ctx, u); err != nil {
		return nil, err
	}
	return s.Get(ctx, u.ID)
}

func (s *UserService) Login(ctx context.Context, req models.LoginRequest) (*models.LoginResponse, error) {
	if req.Email == "" || req.Password == "" {
		return nil, apperr.Invalid("이메일과 비밀번호를 입력해주세요.")
	}
	email, err := normalizeEmail(req.Email)
	if err != nil {
		return nil, errInvalidCredentials
	}

	u, err := s.users.GetByEmail(ctx, email)
	if apperr.Is(err, apperr.KindNotFound) {
		auth.CheckPassword(req.Password, dummyHash)
		return nil, errInvalidCredentials
	}
	if err != nil {
		return nil, err
	}
	if !auth.CheckPassword(req.Password, u.PasswordHash) {
		return nil, errInvalidCredentials
	}

	if u.RestIDs, err = s.users.ListRestaurantIDs(ctx, u.ID); err != nil {
		return nil, err
	}
	token, err := s.tokens.Issue(u.ID)
	if err != nil {
		return nil, err
	}
	return &models.LoginResponse{Token: token, User: u}, nil
}

func (s *UserService) Get(ctx context.Context, id int64) (*models.User, error) {
	u, err := s.users.GetByID(ctx, id)
	if err != nil {
		return nil, err
	}
	if u.RestIDs, err = s.users.ListRestaurantIDs(ctx, id); err != nil {
		return nil, err
	}
	return u, nil
}

func (s *UserService) Update(ctx context.Context, id int64, req models.UpdateUserRequest) (*models.User, error) {
	email, err := normalizeEmail(req.Email)
	if err != nil {
		return nil, err
	}
	name, err := validateName(req.Name)
	if err != nil {
		return nil, err
	}
	phone, err := normalizeOptionalPhone(req.PhoneNum)
	if err != nil {
		return nil, err
	}

	if err := s.users.Update(ctx, &models.User{ID: id, Email: email, Name: name, PhoneNum: phone}); err != nil {
		return nil, err
	}
	return s.Get(ctx, id)
}

func (s *UserService) Delete(ctx context.Context, id int64) error {
	return s.users.Delete(ctx, id)
}

// LinkRestaurant는 로그인한 사용자를 식당의 사장님으로 등록합니다.
// 식당 하나에는 사장님 한 명만 등록할 수 있습니다(먼저 등록한 사람). 이미 주인이 있으면 409.
func (s *UserService) LinkRestaurant(ctx context.Context, userID, restID int64) (*models.UserRest, error) {
	if restID <= 0 {
		return nil, apperr.Invalid("rest_id 가 필요합니다.")
	}
	return s.users.LinkRestaurant(ctx, userID, restID)
}
