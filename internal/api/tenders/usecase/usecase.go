package usecase

import (
	"context"
	"errors"
	"strconv"

	"avito-tenders/internal/api/tenders"
	"avito-tenders/internal/api/tenders/entities"
	"avito-tenders/pkg/apperror"
	"avito-tenders/pkg/query"
)

type usecase struct {
	repo tenders.Repository
}

func (u *usecase) Create(ctx context.Context, request entities.CreateTenderRequest) (entities.ResponseTender, error) {
	return u.repo.Create(ctx, request)
}

func (u *usecase) Edit(ctx context.Context, id string, request entities.EditTenderRequest) (entities.ResponseTender, error) {
	var idInt int
	idInt, err := strconv.Atoi(id)
	if err != nil {
		return entities.ResponseTender{}, apperror.BadRequest(errors.New("tender id is not number"))
	}

	return u.repo.Edit(ctx, idInt, request)
}

func (u *usecase) FindList(ctx context.Context, id string) ([]entities.ResponseTender, error) {
	// TODO implement me
	panic("implement me")
}

func (u *usecase) FindById(ctx context.Context, id string) (entities.ResponseTender, error) {
	var idInt int
	idInt, err := strconv.Atoi(id)
	if err != nil {
		return entities.ResponseTender{}, apperror.BadRequest(errors.New("tender id is not number"))
	}

	return u.repo.FindById(ctx, idInt)
}

func (u *usecase) FindByUsername(ctx context.Context, username string, pagination query.Pagination) ([]entities.ResponseTender, error) {
	return u.repo.FindByUsername(ctx, username, pagination)
}

func (u *usecase) EditStatus(ctx context.Context, id string, request entities.EditTenderStatusRequest) (entities.ResponseTender, error) {
	var idInt int
	idInt, err := strconv.Atoi(id)
	if err != nil {
		return entities.ResponseTender{}, apperror.BadRequest(errors.New("tender id is not number"))
	}

	return u.repo.EditStatus(ctx, idInt, request)
}

func NewUseCase(repo tenders.Repository) tenders.Usecase {
	return &usecase{repo}
}
