package usecase

import "friendsManagement/internal/repository"

type FriendUsecase struct {
	userRepo repository.UserRepository
}

func NewFriendUsecase(repo repository.UserRepository) *FriendUsecase {
	return &FriendUsecase{userRepo: repo}
}
