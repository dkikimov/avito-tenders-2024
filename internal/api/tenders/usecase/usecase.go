package usecase

import (
	"context"
	"errors"

	trm "github.com/avito-tech/go-transaction-manager/trm/v2/manager"

	"avito-tenders/internal/api/organization"
	"avito-tenders/internal/api/tenders"
	"avito-tenders/internal/api/tenders/dtos"
	"avito-tenders/internal/entity"
	"avito-tenders/pkg/apperror"
	"avito-tenders/pkg/queryparams"
)

type Usecase struct {
	repo      tenders.Repository
	orgRepo   organization.Repository
	trManager *trm.Manager
}

type Opts struct {
	Repo      tenders.Repository
	OrgRepo   organization.Repository
	TrManager *trm.Manager
}

func NewUsecase(opts Opts) *Usecase {
	return &Usecase{repo: opts.Repo, orgRepo: opts.OrgRepo, trManager: opts.TrManager}
}

func (u *Usecase) Create(ctx context.Context, request dtos.CreateTenderRequest) (dtos.TenderResponse, error) {
	tender, err := u.repo.Create(ctx, request.ToEntity())
	if err != nil {
		return dtos.TenderResponse{}, err
	}

	return dtos.NewTenderResponse(tender), nil
}

func (u *Usecase) Edit(ctx context.Context, id string, request dtos.EditTenderRequest) (dtos.TenderResponse, error) {
	var tender entity.Tender
	err := u.trManager.Do(ctx, func(ctx context.Context) error {
		oldTender, err := u.repo.FindById(ctx, id)
		if err != nil {
			return err
		}

		if len(request.Name) != 0 {
			oldTender.Name = request.Name
		}
		if len(request.Description) != 0 {
			oldTender.Description = request.Description
		}
		if len(request.ServiceType) != 0 {
			oldTender.ServiceType = request.ServiceType
		}

		tender, err = u.repo.Update(ctx, oldTender)
		if err != nil {
			return err
		}

		return nil
	})
	if err != nil {
		return dtos.TenderResponse{}, err
	}

	return dtos.NewTenderResponse(tender), nil
}

func (u *Usecase) GetAll(ctx context.Context, filter tenders.TenderFilter, pagination queryparams.Pagination) ([]dtos.TenderResponse, error) {
	tendersList, err := u.repo.GetAll(ctx, filter, pagination)
	if err != nil {
		return nil, err
	}

	return dtos.NewTenderResponseList(tendersList), nil
}

func (u *Usecase) GetTenderStatus(ctx context.Context, id string, request dtos.TenderStatus) (dtos.TenderResponse, error) {
	var tender entity.Tender
	err := u.trManager.Do(ctx, func(ctx context.Context) error {
		var err error
		tender, err = u.repo.FindById(ctx, id)
		if err != nil {
			return err
		}

		// If published return tender.
		if tender.Status == entity.TenderPublished {
			return nil
		}

		// Otherwise check if user is responsible.
		responsible, err := u.orgRepo.IsOrganizationResponsible(ctx, tender.OrganizationId, request.Username)
		if err != nil {
			return err
		}
		if !responsible {
			return apperror.Forbidden(apperror.ErrForbidden)
		}

		return nil
	})
	if err != nil {
		return dtos.TenderResponse{}, err
	}

	return dtos.NewTenderResponse(tender), nil
}

func (u *Usecase) FindByUsername(ctx context.Context, username string, pagination queryparams.Pagination) ([]dtos.TenderResponse, error) {
	tendersList, err := u.repo.FindByUsername(ctx, username, pagination)
	if err != nil {
		return nil, err
	}

	return dtos.NewTenderResponseList(tendersList), nil
}

func (u *Usecase) EditStatus(ctx context.Context, id string, request dtos.EditTenderStatusRequest) (dtos.TenderResponse, error) {
	var tender entity.Tender
	err := u.trManager.Do(ctx, func(ctx context.Context) error {
		oldTender, err := u.repo.FindById(ctx, id)
		if err != nil {
			return err
		}

		oldTender.Status = request.Status

		tender, err = u.repo.Update(ctx, oldTender)
		if err != nil {
			return err
		}

		return nil
	})
	if err != nil {
		return dtos.TenderResponse{}, err
	}

	return dtos.NewTenderResponse(tender), nil
}

func (u *Usecase) Rollback(ctx context.Context, id string, request dtos.RollbackTenderRequest) (dtos.TenderResponse, error) {
	var tender entity.Tender
	err := u.trManager.Do(ctx, func(ctx context.Context) error {
		oldTender, err := u.repo.FindByIDFromHistory(ctx, id, request.Version)
		if err != nil {
			return err
		}

		responsible, err := u.orgRepo.IsOrganizationResponsible(ctx, oldTender.OrganizationId, request.Username)
		if err != nil {
			return err
		}
		if !responsible {
			return apperror.Unauthorized(errors.New("user is not in organization"))
		}

		tender, err = u.repo.Update(ctx, oldTender)
		if err != nil {
			return err
		}

		return nil
	})
	if err != nil {
		return dtos.TenderResponse{}, err
	}

	return dtos.NewTenderResponse(tender), nil
}
