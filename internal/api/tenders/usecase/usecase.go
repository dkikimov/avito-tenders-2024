package usecase

import (
	"context"
	"errors"
	"strconv"

	"avito-tenders/internal/api/organization"
	"avito-tenders/internal/api/tenders"
	"avito-tenders/internal/api/tenders/dtos"
	"avito-tenders/internal/api/tenders/models"
	"avito-tenders/internal/entity"
	"avito-tenders/pkg/apperror"
	"avito-tenders/pkg/queryparams"
)

type usecase struct {
	repo    tenders.Repository
	orgRepo organization.Repository
}

func (u *usecase) Create(ctx context.Context, request dtos.CreateTenderRequest) (dtos.TenderResponse, error) {
	tender, err := u.repo.Create(ctx, request.ToEntity())
	if err != nil {
		return dtos.TenderResponse{}, err
	}

	return dtos.NewTenderResponse(tender), nil
}

func (u *usecase) Edit(ctx context.Context, id string, request dtos.EditTenderRequest) (dtos.TenderResponse, error) {
	tender, err := u.repo.Edit(ctx, models.EditTender{
		TenderID:    id,
		Name:        request.Name,
		Description: request.Description,
		ServiceType: request.ServiceType,
		Username:    request.Username,
	})
	if err != nil {
		return dtos.TenderResponse{}, err
	}

	return dtos.NewTenderResponse(tender), nil
}

func (u *usecase) GetAll(ctx context.Context, filter tenders.TenderFilter, pagination queryparams.Pagination) ([]dtos.TenderResponse, error) {
	tendersList, err := u.repo.GetAll(ctx, filter, pagination)
	if err != nil {
		return nil, err
	}

	return dtos.NewTenderResponseList(tendersList), nil
}

func (u *usecase) GetTenderStatus(ctx context.Context, id string, request dtos.TenderStatus) (dtos.TenderResponse, error) {
	tender, err := u.repo.FindById(ctx, id)
	if err != nil {
		return dtos.TenderResponse{}, err
	}

	// If published return tender.
	if tender.Status == entity.TenderPublished {
		return dtos.NewTenderResponse(tender), nil
	}

	// If not published check does user have permissions
	if len(request.Username) == 0 {
		return dtos.TenderResponse{}, apperror.Unauthorized(errors.New("username is empty"))
	}

	responsible, err := u.orgRepo.IsOrganizationResponsible(ctx, tender.OrganizationId, request.Username)
	if err != nil {
		return dtos.TenderResponse{}, err
	}
	if !responsible {
		return dtos.TenderResponse{}, apperror.Unauthorized(errors.New("user is not in organization"))
	}

	return dtos.NewTenderResponse(tender), nil
}

func (u *usecase) FindByUsername(ctx context.Context, username string, pagination queryparams.Pagination) ([]dtos.TenderResponse, error) {
	tendersList, err := u.repo.FindByUsername(ctx, username, pagination)
	if err != nil {
		return nil, err
	}

	return dtos.NewTenderResponseList(tendersList), nil
}

func (u *usecase) EditStatus(ctx context.Context, id string, request dtos.EditTenderStatusRequest) (dtos.TenderResponse, error) {
	tender, err := u.repo.EditStatus(ctx, models.EditTenderStatus{
		TenderId: id,
		Status:   request.Status,
		Username: request.Username,
	})
	if err != nil {
		return dtos.TenderResponse{}, err
	}

	return dtos.NewTenderResponse(tender), nil
}

func (u *usecase) Rollback(ctx context.Context, id string, request dtos.RollbackTenderRequest) (dtos.TenderResponse, error) {
	versionInt, err := strconv.Atoi(request.Version)
	if err != nil {
		return dtos.TenderResponse{}, apperror.BadRequest(errors.New("version is not number"))
	}

	oldTender, err := u.repo.FindByIDFromHistory(ctx, id, versionInt)
	if err != nil {
		return dtos.TenderResponse{}, err
	}

	responsible, err := u.orgRepo.IsOrganizationResponsible(ctx, oldTender.OrganizationId, request.Username)
	if err != nil {
		return dtos.TenderResponse{}, err
	}
	if !responsible {
		return dtos.TenderResponse{}, apperror.Unauthorized(errors.New("user is not in organization"))
	}

	updatedTender, err := u.repo.Update(ctx, oldTender)
	if err != nil {
		return dtos.TenderResponse{}, err
	}

	return dtos.NewTenderResponse(updatedTender), nil
}

func NewUseCase(repo tenders.Repository, orgRepo organization.Repository) tenders.Usecase {
	return &usecase{repo, orgRepo}
}
