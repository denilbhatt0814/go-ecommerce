package service

import (
	"errors"
	"fmt"
	"go-ecommerce-app/internal/domain"
	"go-ecommerce-app/internal/dto"
	"go-ecommerce-app/internal/helper"
	"go-ecommerce-app/internal/repository"
	"log"
	"time"
)

type UserService struct {
	Repo repository.UserRepository
	Auth helper.Auth
}

func (s UserService) findUserByEmail(email string) (*domain.User, error) {
	// perform some db opperations
	// some buisness logic
	user, err := s.Repo.FindUser(email)

	return &user, err
}

func (s UserService) SignUp(input dto.UserSignUp) (string, error) {

	// hashed password
	hashedPassword, err := s.Auth.CreateHashedPasword(input.Password)
	if err != nil {
		return "", err
	}

	user, err := s.Repo.CreateUser(domain.User{
		Email:    input.Email,
		Password: hashedPassword,
		Phone:    input.Phone,
	})

	// generate token
	userInfo := fmt.Sprintf("%v, %v, %v", user.ID, user.Email, user.UserType)
	log.Println(userInfo)

	return s.Auth.GenerateToken(user.ID, user.Email, user.UserType)
}

func (s UserService) Login(email string, password string) (string, error) {

	user, err := s.findUserByEmail(email)
	if err != nil {
		return "", errors.New("user doesn't exist with given email id")
	}

	err = s.Auth.VerifyPassword(password, user.Password)
	if err != nil {
		return "", err
	}

	return s.Auth.GenerateToken(user.ID, user.Email, user.UserType)
}

func (s UserService) isVerifiedUser(id uint) bool {

	currentUser, err := s.Repo.FindUserById(id)

	return err == nil && currentUser.Verified
}

func (s UserService) GetVerificationCode(e domain.User) (int, error) {

	// if user already verified
	if s.isVerifiedUser(e.ID) {
		return 0, nil
	}

	// generate verification code
	code, err := s.Auth.GenerateCode()
	if err != nil {
		return 0, err
	}

	// update user
	user := domain.User{
		Expiry: time.Now().Add(time.Minute * 30),
		Code:   code,
	}
	_, err = s.Repo.UpdateUser(e.ID, user)
	if err != nil {
		return 0, errors.New("unable to update verification code")
	}

	// send SMS

	// return verification code
	return code, nil
}

func (s UserService) VerifyCode(id uint, code int) error {
	// if user already verified
	if s.isVerifiedUser(id) {
		return errors.New("user already verified")
	}

	user, err := s.Repo.FindUserById(id)
	if err != nil {
		return err
	}

	if !time.Now().Before(user.Expiry) {
		return errors.New("verification code expired")
	}

	if user.Code != code {
		return errors.New("verification code does not match")
	}

	updateUser := domain.User{
		Verified: true,
	}
	_, err = s.Repo.UpdateUser(id, updateUser)
	if err != nil {
		return errors.New("unable to verify user")
	}

	return nil
}

func (s UserService) CreateProfile(id uint, input any) (*domain.User, error) {
	return nil, nil
}

func (s UserService) GetProfile(id uint) *domain.User {
	return nil
}

func (s UserService) UpdateProfile(id uint, input any) error {
	return nil
}

func (s UserService) BecomeSeller(id uint, input any) error {
	return nil
}

func (s UserService) FindCart(id uint) ([]interface{}, error) {
	return nil, nil
}

func (s UserService) CreateCart(input any, u domain.User) ([]interface{}, error) {
	return nil, nil
}

func (s UserService) CreateOrder(u domain.User) (int, error) {
	return 0, nil
}

func (s UserService) GetOrders(u domain.User) ([]interface{}, error) {
	return nil, nil
}

func (s UserService) GetOrderById(id uint, uId uint) (interface{}, error) {
	return nil, nil
}
