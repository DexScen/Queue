package service

import (
	"context"
	"errors"

	"github.com/DexScen/Queue/backend/internal/domain"
	e "github.com/DexScen/Queue/backend/internal/errors"
	"golang.org/x/crypto/bcrypt"
)

type QueuesRepository interface {
	GetAllGames(ctx context.Context, listGames *domain.ListGames) error
	GetGameInfoByID(ctx context.Context, id int) (*domain.Game, error)
	GetGamesByLogin(ctx context.Context, login string, listGames *domain.ListGames) error
	GetIdByLogin(ctx context.Context, login string) (int, error)

	GetPassword(ctx context.Context, login string) (string, error)
	GetRole(ctx context.Context, login string) (string, error)
	UserExists(ctx context.Context, login string) (bool, error)
	Register(ctx context.Context, user *domain.User) error

	AddPlayerToQueue(ctx context.Context, user_id, game_id int) (int, error)
	RemovePlayerFromQueue(ctx context.Context, user_id, game_id int) error

	GetPlayersByGameID(ctx context.Context, game_id int, listUsers *domain.ListUsers) error
}

type Queues struct {
	repo QueuesRepository
}

func NewQueues(repo QueuesRepository) *Queues {
	return &Queues{
		repo: repo,
	}
}

func (q *Queues) GetAllGames(ctx context.Context, listGames *domain.ListGames) error {
	return q.repo.GetAllGames(ctx, listGames)
}

func (q *Queues) GetGameInfoByID(ctx context.Context, id int) (*domain.Game, error) {
	return q.repo.GetGameInfoByID(ctx, id)
}

func (q *Queues) GetGamesByLogin(ctx context.Context, login string, listGames *domain.ListGames) error{
	return q.repo.GetGamesByLogin(ctx, login, listGames)
}

func (q *Queues) LogIn(ctx context.Context, login, password string) (string, error) {
	passwordHash, err := q.repo.GetPassword(ctx, login)

	if err != nil {
		if errors.Is(err, e.ErrUserNotFound) {
			return "", e.ErrUserNotFound
		}
		return "", err
	}

	err = bcrypt.CompareHashAndPassword([]byte(passwordHash), []byte(password))
	if err != nil {
		return "", e.ErrWrongPassword
	}
	return q.repo.GetRole(ctx, login)
}

func (q *Queues) Register(ctx context.Context, user *domain.User) error {
	exists, err := q.repo.UserExists(ctx, user.Login)
	if exists {
		return e.ErrUserExists
	}
	if err != nil {
		return err
	}

	//тут лучше конвертер но увы
	hash, err := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)
	if err != nil {
		return err
	}
	return q.repo.Register(ctx, &domain.User{
		Login:    user.Login,
		Password: string(hash),
	})
}

func (q *Queues) RemovePlayerFromQueue(ctx context.Context, user_id, game_id int) error{
	return q.repo.RemovePlayerFromQueue(ctx, user_id, game_id)
}

func (q *Queues) AddPlayerToQueue(ctx context.Context, user_id, game_id int) (int, error){
	return q.repo.AddPlayerToQueue(ctx , user_id, game_id)
}

func (q *Queues) GetIdByLogin(ctx context.Context, login string) (int,error){
	return q.repo.GetIdByLogin(ctx, login)
}

func (q *Queues) GetPlayersByGameID(ctx context.Context, game_id int, listUsers *domain.ListUsers) error{
	return q.repo.GetPlayersByGameID(ctx, game_id, listUsers)
}