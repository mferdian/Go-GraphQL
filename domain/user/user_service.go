package user

import (
	"context"

	"github.com/google/uuid"
	"github.com/mferdian/Go-GraphQL/config/jwt"
	"github.com/mferdian/Go-GraphQL/constants"
	"github.com/mferdian/Go-GraphQL/helpers"
	"github.com/mferdian/Go-GraphQL/logging"
)

type (
	IUserService interface {
		Register(ctx context.Context, req RegisterUserRequest) (RegisterUserResponse, error)
		Login(ctx context.Context, req LoginUserRequest) (LoginResponse, error)

		CreateUser(ctx context.Context, req CreateUserRequest) (UserResponse, error)
		GetuserByID(ctx context.Context, userID string) (UserResponse, error)
		GetAllUser(ctx context.Context, search string) ([]UserResponse, error)
		GetAllUserWithPagination(ctx context.Context, req UserPaginationRequest) (UserPaginationResponse, error)
		UpdateUser(ctx context.Context, req UpdateUserRequest) (UserResponse, error)
		DeleteUser(ctx context.Context, req DeleteUserRequest) (UserResponse, error)
	}

	UserService struct {
		userRepo   IUserRepository
		jwtService jwt.InterfaceJWTService
	}
)

func NewUserService(userRepo IUserRepository, jwtService jwt.InterfaceJWTService) *UserService {
	return &UserService{
		userRepo:   userRepo,
		jwtService: jwtService,
	}
}

func (us *UserService) Register(ctx context.Context, req RegisterUserRequest) (RegisterUserResponse, error) {
	if len(req.Name) < 5 {
		logging.Log.Warn(constants.MESSAGE_FAILED_REGISTER + ": name too short")
		return RegisterUserResponse{}, constants.ErrInvalidName
	}

	if !helpers.IsValidEmail(req.Email) {
		logging.Log.Warn(constants.MESSAGE_FAILED_REGISTER + ": invalid email format")
		return RegisterUserResponse{}, constants.ErrInvalidEmail
	}

	_, found, err := us.userRepo.GetUserByEmail(ctx, nil, req.Email)
	if err == nil && found {
		logging.Log.Warn(constants.MESSAGE_FAILED_REGISTER + ": email already exists")
		return RegisterUserResponse{}, constants.ErrEmailAlreadyExists
	}

	if len(req.Password) < 8 {
		logging.Log.Warn(constants.MESSAGE_FAILED_REGISTER + ": password too short")
		return RegisterUserResponse{}, constants.ErrInvalidPassword
	}

	user := User{
		ID:       uuid.New(),
		Name:     req.Name,
		Email:    req.Email,
		Password: req.Password,
		Role:     constants.ENUM_ROLE_USER,
	}

	err = us.userRepo.Register(ctx, nil, user)
	if err != nil {
		logging.Log.WithError(err).Error(constants.MESSAGE_FAILED_REGISTER)
		return RegisterUserResponse{}, constants.ErrRegisterUser
	}

	logging.Log.Infof(constants.MESSAGE_SUCCESS_REGISTER+": %s", user.Email)

	return RegisterUserResponse{
		ID:    user.ID,
		Name:  user.Name,
		Email: user.Email,
	}, nil
}

func (us *UserService) Login(ctx context.Context, req LoginUserRequest) (LoginResponse, error) {
	user, found, err := us.userRepo.GetUserByEmail(ctx, nil, req.Email)
	if err != nil || !found {
		logging.Log.Warn(constants.MESSAGE_FAILED_LOGIN_USER + ": email not found")
		return LoginResponse{}, constants.ErrInvalidLoginCredential
	}

	if ok, err := helpers.CheckPassword(user.Password, []byte(req.Password)); !ok || err != nil {
		logging.Log.Warn(constants.MESSAGE_FAILED_LOGIN_USER + ": password mismatch")
		return LoginResponse{}, constants.ErrInvalidLoginCredential
	}

	accessToken, refreshToken, err := us.jwtService.GenerateToken(user.ID.String(), user.Role)
	if err != nil {
		logging.Log.WithError(err).Error(constants.MESSAGE_FAILED_LOGIN_USER + ": failed generate token")
		return LoginResponse{}, constants.ErrGenerateAccessToken
	}

	logging.Log.Infof(constants.MESSAGE_SUCCESS_LOGIN_USER+": %s", user.Email)

	return LoginResponse{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
	}, nil
}

func (us *UserService) CreateUser(ctx context.Context, req CreateUserRequest) (UserResponse, error) {
	if len(req.Name) < 5 {
		logging.Log.Warn(constants.MESSAGE_FAILED_CREATE_USER + ": name too short")
		return UserResponse{}, constants.ErrInvalidName
	}

	if !helpers.IsValidEmail(req.Email) {
		logging.Log.Warn(constants.MESSAGE_FAILED_CREATE_USER + ": invalid email")
		return UserResponse{}, constants.ErrInvalidEmail
	}

	_, found, err := us.userRepo.GetUserByEmail(ctx, nil, req.Email)
	if err == nil && found {
		logging.Log.Warn(constants.MESSAGE_FAILED_CREATE_USER + ": email already exists")
		return UserResponse{}, constants.ErrEmailAlreadyExists
	}

	if len(req.Password) < 8 {
		logging.Log.Warn(constants.MESSAGE_FAILED_CREATE_USER + ": password too short")
		return UserResponse{}, constants.ErrInvalidPassword
	}

	user := User{
		ID:          uuid.New(),
		Name:        req.Name,
		Email:       req.Email,
		Password:    req.Password,
		PhoneNumber: req.PhoneNumber,
		Address:     req.Address,
		Role:        constants.ENUM_ROLE_ADMIN,
	}

	err = us.userRepo.CreateUser(ctx, nil, user)
	if err != nil {
		logging.Log.WithError(err).Error(constants.MESSAGE_FAILED_CREATE_USER)
		return UserResponse{}, constants.ErrCreateUser
	}

	logging.Log.Infof(constants.MESSAGE_SUCCESS_CREATE_USER+": %s", user.Email)

	return UserResponse{
		ID:          user.ID,
		Name:        user.Name,
		Email:       user.Email,
		PhoneNumber: user.PhoneNumber,
		Address:     user.Address,
	}, nil
}

func (us *UserService) GetAllUser(ctx context.Context, search string) ([]UserResponse, error) {
	users, err := us.userRepo.GetAllUser(ctx, nil, search)

	if err != nil {
		return nil, constants.ErrGetAllUser
	}

	var datas []UserResponse
	for _, user := range users {
		data := UserResponse{
			ID:          user.ID,
			Name:        user.Name,
			Email:       user.Email,
			PhoneNumber: user.PhoneNumber,
			Address:     user.Address,
		}

		datas = append(datas, data)
	}
	return datas, nil
}

func (us *UserService) GetAllUserWithPagination(ctx context.Context, req UserPaginationRequest) (UserPaginationResponse, error) {
	dataWithPaginate, err := us.userRepo.GetAllUserWithPagination(ctx, nil, req)
	if err != nil {
		logging.Log.WithError(err).Error(constants.MESSAGE_FAILED_GET_LIST_USER)
		return UserPaginationResponse{}, constants.ErrGetAllUserWithPagination
	}

	logging.Log.Infof(constants.MESSAGE_SUCCESS_GET_LIST_USER+": page %d", req.Page)

	var datas []UserResponse
	for _, user := range dataWithPaginate.Users {
		datas = append(datas, UserResponse{
			ID:          user.ID,
			Name:        user.Name,
			Email:       user.Email,
			PhoneNumber: user.PhoneNumber,
			Address:     user.Address,
		})
	}

	return UserPaginationResponse{
		Data: datas,
		PaginationResponse: PaginationResponse{
			Page:    dataWithPaginate.Page,
			PerPage: dataWithPaginate.PerPage,
			MaxPage: dataWithPaginate.MaxPage,
			Count:   dataWithPaginate.Count,
		},
	}, nil
}

func (us *UserService) GetuserByID(ctx context.Context, userID string) (UserResponse, error) {
	if _, err := uuid.Parse(userID); err != nil {
		logging.Log.Warn(constants.MESSAGE_FAILED_GET_DETAIL_USER + ": invalid UUID")
		return UserResponse{}, constants.ErrInvalidUUID
	}

	user, _, err := us.userRepo.GetUserByID(ctx, nil, userID)
	if err != nil {
		logging.Log.WithError(err).WithField("id", userID).Error(constants.MESSAGE_FAILED_GET_DETAIL_USER)
		return UserResponse{}, constants.ErrGetUserByID
	}

	logging.Log.Infof(constants.MESSAGE_SUCCESS_GET_DETAIL_USER+": %s", userID)

	return UserResponse{
		ID:          user.ID,
		Name:        user.Name,
		Email:       user.Email,
		Address:     user.Address,
		PhoneNumber: user.PhoneNumber,
	}, nil
}

func (us *UserService) UpdateUser(ctx context.Context, req UpdateUserRequest) (UserResponse, error) {
	user, _, err := us.userRepo.GetUserByID(ctx, nil, req.ID)
	if err != nil {
		logging.Log.WithError(err).WithField("id", req.ID).Error(constants.MESSAGE_FAILED_UPDATE_USER)
		return UserResponse{}, constants.ErrGetUserByID
	}

	if req.Name != nil && len(*req.Name) < 5 {
		logging.Log.Warn(constants.MESSAGE_FAILED_UPDATE_USER + ": invalid name")
		return UserResponse{}, constants.ErrInvalidName
	} else if req.Name != nil {
		user.Name = *req.Name
	}

	if req.Email != nil {
		if !helpers.IsValidEmail(*req.Email) {
			logging.Log.Warn(constants.MESSAGE_FAILED_UPDATE_USER + ": invalid email format")
			return UserResponse{}, constants.ErrInvalidEmail
		}

		existingUser, found, err := us.userRepo.GetUserByEmail(ctx, nil, *req.Email)
		if err == nil && found && existingUser.ID != user.ID {
			logging.Log.Warn(constants.MESSAGE_FAILED_UPDATE_USER + ": email already used by other user")
			return UserResponse{}, constants.ErrEmailAlreadyExists
		}

		user.Email = *req.Email
	}

	if req.Password != nil {
		if ok, _ := helpers.CheckPassword(user.Password, []byte(*req.Password)); ok {
			logging.Log.Warn(constants.MESSAGE_FAILED_UPDATE_USER + ": new password same as old")
			return UserResponse{}, constants.ErrPasswordSame
		}

		hashed, err := helpers.HashPassword(*req.Password)
		if err != nil {
			logging.Log.WithError(err).Error(constants.MESSAGE_FAILED_UPDATE_USER + ": hash password error")
			return UserResponse{}, constants.ErrHashPassword
		}
		user.Password = hashed
	}

	if req.PhoneNumber != nil {
		user.PhoneNumber = *req.PhoneNumber
	}

	if req.Address != nil {
		user.Address = *req.Address
	}

	err = us.userRepo.UpdateUser(ctx, nil, user)
	if err != nil {
		logging.Log.WithError(err).Error(constants.MESSAGE_FAILED_UPDATE_USER)
		return UserResponse{}, constants.ErrUpdateUser
	}

	logging.Log.Infof(constants.MESSAGE_SUCCESS_UPDATE_USER+": %s", user.ID)

	return UserResponse{
		ID:          user.ID,
		Name:        user.Name,
		Email:       user.Email,
		PhoneNumber: user.PhoneNumber,
		Address:     user.Address,
	}, nil
}

func (us *UserService) DeleteUser(ctx context.Context, req DeleteUserRequest) (UserResponse, error) {
	user, _, err := us.userRepo.GetUserByID(ctx, nil, req.UserID)
	if err != nil {
		logging.Log.WithError(err).Error(constants.MESSAGE_FAILED_DELETE_USER)
		return UserResponse{}, constants.ErrGetUserByID
	}

	err = us.userRepo.DeleteUserByID(ctx, nil, req.UserID)
	if err != nil {
		logging.Log.WithError(err).Error(constants.MESSAGE_FAILED_DELETE_USER)
		return UserResponse{}, constants.ErrDeleteUserByID
	}

	logging.Log.Infof(constants.MESSAGE_SUCCESS_DELETE_USER+": %s", req.UserID)

	return UserResponse{
		ID:          user.ID,
		Name:        user.Name,
		Email:       user.Email,
		PhoneNumber: user.PhoneNumber,
		Address:     user.Address,
	}, nil
}
