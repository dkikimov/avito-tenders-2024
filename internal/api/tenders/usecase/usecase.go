package usecase

import (
	"context"
	"errors"
	"strconv"

	"avito-tenders/internal/api/organization"
	"avito-tenders/internal/api/tenders"
	"avito-tenders/internal/api/tenders/entities"
	"avito-tenders/internal/entity"
	"avito-tenders/pkg/apperror"
	"avito-tenders/pkg/queryparams"
)

type usecase struct {
	repo    tenders.Repository
	orgRepo organization.Repository
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

func (u *usecase) GetAll(ctx context.Context, filter tenders.TenderFilter, pagination queryparams.Pagination) ([]entities.ResponseTender, error) {
	return u.repo.GetAll(ctx, filter, pagination)
}

func (u *usecase) GetTenderStatus(ctx context.Context, id string, request entities.TenderStatus) (entities.ResponseTender, error) {
	var idInt int
	idInt, err := strconv.Atoi(id)
	if err != nil {
		return entities.ResponseTender{}, apperror.BadRequest(errors.New("tender id is not number"))
	}

	tender, err := u.repo.FindById(ctx, idInt)
	if err != nil {
		return entities.ResponseTender{}, err
	}

	// If published return tender.
	if tender.Status == entity.TenderPublished {
		return tender, nil
	}

	// If not published check does user have permissions
	if len(request.Username) == 0 {
		return entities.ResponseTender{}, apperror.Unauthorized(errors.New("username is empty"))
	}

	responsible, err := u.orgRepo.IsOrganizationResponsible(ctx, tender.OrganizationId, request.Username)
	if err != nil {
		return entities.ResponseTender{}, err
	}
	if !responsible {
		return entities.ResponseTender{}, apperror.Unauthorized(errors.New("user is not in organization"))
	}

	return tender, nil
}

func (u *usecase) FindByUsername(ctx context.Context, username string, pagination queryparams.Pagination) ([]entities.ResponseTender, error) {
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

func (u *usecase) Rollback(ctx context.Context, id string, request entities.RollbackTenderRequest) (entities.ResponseTender, error) {
	idInt, err := strconv.Atoi(id)
	if err != nil {
		return entities.ResponseTender{}, apperror.BadRequest(errors.New("tender id is not number"))
	}

	versionInt, err := strconv.Atoi(request.Version)
	if err != nil {
		return entities.ResponseTender{}, apperror.BadRequest(errors.New("version is not number"))
	}

	return u.repo.Rollback(ctx, idInt, entities.RollbackTender{
		Username: request.Username,
		Version:  versionInt,
	})
}

func NewUseCase(repo tenders.Repository, orgRepo organization.Repository) tenders.Usecase {
	return &usecase{repo, orgRepo}
}
