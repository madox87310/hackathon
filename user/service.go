package user

import (
	"context"
	"fmt"
	"hackathon/pkg"

	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
)

type Service struct {
	userRepo   *Repository
	jwtService *pkg.JWTService
}

func NewUserService(userRepo *Repository, jwtService *pkg.JWTService) *Service {
	return &Service{
		userRepo:   userRepo,
		jwtService: jwtService,
	}
}

func (s *Service) SignUp(ctx context.Context, req SignUpRequest) (SignUpResponse, error) {
	if user, err := s.findByPhoneNumber(ctx, req.PhoneNumber); user != nil {
		return SignUpResponse{}, fmt.Errorf("s.findByPhoneNumber: %w", err)
	}

	user, err := s.createUser(req)
	if err != nil {
		return SignUpResponse{}, fmt.Errorf("s.createUser: %w", err)
	}

	accessToken, refreshToken, err := s.jwtService.GenerateTokens(user.ID.String())
	if err != nil {
		return SignUpResponse{}, fmt.Errorf("s.jwtService.GenerateTokens: %w", err)
	}

	user.SetRefreshToken(refreshToken)

	if err := s.userRepo.Save(ctx, user); err != nil {
		return SignUpResponse{}, fmt.Errorf("s.userRepo.Save: %w", err)
	}

	return s.createSignUpResponse(user, accessToken), nil
}

func (s *Service) SignIn(ctx context.Context, req SignInRequest) (SignInResponse, error) {
	user, err := s.Authenticate(ctx, req)
	if err != nil {
		return SignInResponse{}, fmt.Errorf("s.Authenticate: %w", err)
	}

	accessToken, _, err := s.jwtService.GenerateTokens(user.ID.String())
	if err != nil {
		return SignInResponse{}, fmt.Errorf("s.jwtService.GenerateTokens: %w", err)
	}

	return s.createSignInResponse(user, accessToken), nil
}

func (s *Service) Refresh(ctx context.Context, req RefreshRequest) (RefreshResponse, error) {
	userID, err := s.jwtService.ValidateRefreshToken(req.RefreshToken)
	if err != nil {
		return RefreshResponse{}, fmt.Errorf("s.jwtService.ValidateRefreshToken: %w", err)
	}

	uuidID, err := uuid.Parse(userID)
	if err != nil {
		return RefreshResponse{}, fmt.Errorf("uuid.Parse: %w", err)
	}

	user, err := s.userRepo.FindByID(ctx, uuidID)
	if err != nil {
		return RefreshResponse{}, fmt.Errorf("s.userRepo.FindByID: %w", err)
	}

	accessToken, refreshToken, err := s.jwtService.GenerateTokens(userID)
	if err != nil {
		return RefreshResponse{}, fmt.Errorf("s.jwtService.GenerateTokens: %w", err)
	}

	user.SetRefreshToken(refreshToken)

	if err := s.userRepo.Update(ctx, user); err != nil {
		return RefreshResponse{}, fmt.Errorf("s.userRepo.Update: %w", err)
	}

	return s.createRefreshResponse(accessToken, refreshToken), nil
}

func (s *Service) Logout(ctx context.Context, req LogoutRequest) error {
	userID, err := s.jwtService.ValidateAccessToken(req.AccessToken)
	if err != nil {
		return fmt.Errorf("s.jwtService.ValidateAccessToken: %w", err)
	}

	uuidID, err := uuid.Parse(userID)
	if err != nil {
		return fmt.Errorf("uuid.Parse: %w", err)
	}

	user, err := s.userRepo.FindByID(ctx, uuidID)
	if err != nil {
		return fmt.Errorf("s.userRepo.FindByID: %w", err)
	}

	user.SetRefreshToken("")

	if err := s.userRepo.Update(ctx, user); err != nil {
		return fmt.Errorf("s.userRepo.Update: %w", err)
	}

	return nil
}

func (s *Service) Authenticate(ctx context.Context, req SignInRequest) (*Model, error) {
	user, err := s.findByPhoneNumber(ctx, req.PhoneNumber)
	if err != nil {
		return nil, fmt.Errorf("s.findByPhoneNumber: %w", err)
	}

	if err := s.verifyPassword(user.Password, req.Password); err != nil {
		return nil, fmt.Errorf("s.verifyPassword: %w", err)
	}

	return user, nil
}

func (s *Service) findByPhoneNumber(ctx context.Context, phoneNumber string) (*Model, error) {
	return s.userRepo.FindByPhoneNumber(ctx, phoneNumber)
}

func (s *Service) verifyPassword(userPassword, password string) error {
	if err := bcrypt.CompareHashAndPassword([]byte(userPassword), []byte(password)); err != nil {
		return fmt.Errorf("bcrypt.CompareHashAndPassword: %w", err)
	}
	return nil
}

func (s *Service) hashPassword(password string) (string, error) {
	hash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return "", fmt.Errorf("bcrypt.GenerateFromPassword: %w", err)
	}
	return string(hash), nil
}

func (s *Service) createUser(req SignUpRequest) (*Model, error) {
	hashedPassword, err := s.hashPassword(req.Password)
	if err != nil {
		return nil, fmt.Errorf("s.hashPassword: %w", err)
	}
	return NewModel(req.DisplayName, req.PhoneNumber, hashedPassword), nil
}

func (s *Service) createSignUpResponse(user *Model, accessToken string) SignUpResponse {
	return SignUpResponse{
		ID:           user.ID.String(),
		AccessToken:  accessToken,
		RefreshToken: user.RefreshToken,
	}
}

func (s *Service) createSignInResponse(user *Model, accessToken string) SignInResponse {
	return SignInResponse{
		ID:           user.ID.String(),
		AccessToken:  accessToken,
		RefreshToken: user.RefreshToken,
	}
}

func (s *Service) createRefreshResponse(accessToken, refreshToken string) RefreshResponse {
	return RefreshResponse{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
	}
}
