package services

import (
	"context"
	"time"

	"github.com/google/uuid"
	"github.com/mcorrigan89/openmic/internal/domain/entities"
	"github.com/mcorrigan89/openmic/internal/domain/repositories"
	"github.com/mcorrigan89/openmic/internal/infrastructure/postgres/models"

	"github.com/rs/xid"
)

type UserService interface {
	GetUserByID(ctx context.Context, querier models.Querier, userID uuid.UUID) (*entities.UserEntity, error)
	GetUserByEmail(ctx context.Context, querier models.Querier, email string) (*entities.UserEntity, error)
	GetUserByHandle(ctx context.Context, querier models.Querier, userHandle string) (*entities.UserEntity, error)
	GetUserContextBySessionToken(ctx context.Context, querier models.Querier, sessionToken string) (*entities.UserContextEntity, error)
	CreateUser(ctx context.Context, querier models.Querier, user *entities.UserEntity) (*entities.UserEntity, error)
	UpdateUser(ctx context.Context, querier models.Querier, user *entities.UserEntity) (*entities.UserEntity, error)
	CreateSession(ctx context.Context, querier models.Querier, user *entities.UserEntity) (*entities.UserContextEntity, error)
	CreateLoginLink(ctx context.Context, querier models.Querier, email string) (*entities.ReferenceLinkEntity, error)
	CreateInviteLink(ctx context.Context, querier models.Querier, userEntity *entities.UserEntity) (*entities.ReferenceLinkEntity, error)
	LoginWithLink(ctx context.Context, querier models.Querier, token string) (*entities.UserContextEntity, error)
	AcceptInviteLink(ctx context.Context, querier models.Querier, token string) (*entities.UserContextEntity, error)
	SetAvatarImage(ctx context.Context, querier models.Querier, image *entities.ImageEntity, user *entities.UserEntity) (*entities.UserEntity, error)
}

type userService struct {
	userRepo    repositories.UserRepository
	refLinkRepo repositories.ReferenceLinkRepository
	imageRepo   repositories.ImageRepository
}

func NewUserService(userRepo repositories.UserRepository, refLinkRepo repositories.ReferenceLinkRepository, imageRepo repositories.ImageRepository) *userService {
	return &userService{userRepo: userRepo, refLinkRepo: refLinkRepo, imageRepo: imageRepo}
}

func (s *userService) GetUserByID(ctx context.Context, querier models.Querier, userID uuid.UUID) (*entities.UserEntity, error) {
	user, err := s.userRepo.GetUserByID(ctx, querier, userID)
	if err != nil {
		return nil, err
	}

	return user, nil
}

func (s *userService) GetUserByEmail(ctx context.Context, querier models.Querier, email string) (*entities.UserEntity, error) {
	user, err := s.userRepo.GetUserByEmail(ctx, querier, email)
	if err != nil {
		return nil, err
	}

	return user, nil
}

func (s *userService) GetUserByHandle(ctx context.Context, querier models.Querier, userHandle string) (*entities.UserEntity, error) {
	user, err := s.userRepo.GetUserByHandle(ctx, querier, userHandle)
	if err != nil {
		return nil, err
	}

	return user, nil
}

func (s *userService) GetUserContextBySessionToken(ctx context.Context, querier models.Querier, sessionToken string) (*entities.UserContextEntity, error) {
	user, err := s.userRepo.GetUserContextBySessionToken(ctx, querier, sessionToken)
	if err != nil {
		return nil, err
	}

	return user, nil
}

func (s *userService) CreateUser(ctx context.Context, querier models.Querier, user *entities.UserEntity) (*entities.UserEntity, error) {
	user, err := s.userRepo.CreateUser(ctx, querier, user)
	if err != nil {
		return nil, err
	}

	return user, nil
}

func (s *userService) UpdateUser(ctx context.Context, querier models.Querier, user *entities.UserEntity) (*entities.UserEntity, error) {
	user, err := s.userRepo.UpdateUser(ctx, querier, user)
	if err != nil {
		return nil, err
	}

	return user, nil
}

func (s *userService) CreateSession(ctx context.Context, querier models.Querier, user *entities.UserEntity) (*entities.UserContextEntity, error) {

	token := xid.New().String()
	expiresAt := time.Now().Add(time.Hour * 24 * 30)

	userSession, err := s.userRepo.CreateSession(ctx, querier, user, token, expiresAt)
	if err != nil {
		return nil, err
	}

	return userSession, nil
}

func (s *userService) CreateLoginLink(ctx context.Context, querier models.Querier, email string) (*entities.ReferenceLinkEntity, error) {
	userEntity, err := s.userRepo.GetUserByEmail(ctx, querier, email)
	if err != nil {
		return nil, err
	}

	newLinkEntity := entities.ReferenceLinkEntity{
		ID:        uuid.New(),
		LinkID:    userEntity.ID,
		Token:     xid.New().String(),
		Type:      entities.RefLinkTypeLogin,
		ExpiresAt: time.Now().Add(time.Minute * 30),
	}

	loginLinkEntity, err := s.refLinkRepo.CreateReferenceLink(ctx, querier, &newLinkEntity)
	if err != nil {
		return nil, err
	}

	return loginLinkEntity, nil
}

func (s *userService) CreateInviteLink(ctx context.Context, querier models.Querier, userEntity *entities.UserEntity) (*entities.ReferenceLinkEntity, error) {
	newLinkEntity := entities.ReferenceLinkEntity{
		ID:        uuid.New(),
		LinkID:    userEntity.ID,
		Token:     xid.New().String(),
		Type:      entities.RefLinkTypeInvite,
		ExpiresAt: time.Now().Add(time.Hour * 24 * 30),
	}

	loginLinkEntity, err := s.refLinkRepo.CreateReferenceLink(ctx, querier, &newLinkEntity)
	if err != nil {
		return nil, err
	}

	return loginLinkEntity, nil
}

func (s *userService) LoginWithLink(ctx context.Context, querier models.Querier, token string) (*entities.UserContextEntity, error) {
	refLinkEntity, err := s.refLinkRepo.GetReferenceLinkByToken(ctx, querier, token)
	if err != nil {
		return nil, err
	}

	if refLinkEntity.IsExpired() {
		return nil, entities.ErrLinkExpired
	}

	if refLinkEntity.Type != entities.RefLinkTypeLogin {
		return nil, entities.ErrLinkInvalid
	}

	userEntity, err := s.userRepo.GetUserByID(ctx, querier, refLinkEntity.LinkID)
	if err != nil {
		return nil, err
	}

	userSession, err := s.CreateSession(ctx, querier, userEntity)
	if err != nil {
		return nil, err
	}

	err = s.refLinkRepo.DeleteReferenceLink(ctx, querier, refLinkEntity)
	if err != nil {
		return nil, err
	}

	return userSession, nil
}

func (s *userService) AcceptInviteLink(ctx context.Context, querier models.Querier, token string) (*entities.UserContextEntity, error) {
	refLinkEntity, err := s.refLinkRepo.GetReferenceLinkByToken(ctx, querier, token)
	if err != nil {
		return nil, err
	}

	if refLinkEntity.IsExpired() {
		return nil, entities.ErrLinkExpired
	}

	if refLinkEntity.Type != entities.RefLinkTypeInvite {
		return nil, entities.ErrLinkInvalid
	}

	userEntity, err := s.userRepo.GetUserByID(ctx, querier, refLinkEntity.LinkID)
	if err != nil {
		return nil, err
	}

	userEntity.Claimed = true
	userEntity.EmailVerified = true

	userEntity, err = s.userRepo.UpdateUser(ctx, querier, userEntity)
	if err != nil {
		return nil, err
	}

	userSession, err := s.CreateSession(ctx, querier, userEntity)
	if err != nil {
		return nil, err
	}

	err = s.refLinkRepo.DeleteReferenceLink(ctx, querier, refLinkEntity)
	if err != nil {
		return nil, err
	}

	return userSession, nil
}

func (s *userService) SetAvatarImage(ctx context.Context, querier models.Querier, image *entities.ImageEntity, user *entities.UserEntity) (*entities.UserEntity, error) {
	userEntity, err := s.userRepo.SetAvatarImage(ctx, querier, image, user)
	if err != nil {
		return nil, err
	}

	return userEntity, nil
}
