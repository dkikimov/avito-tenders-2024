package usecase

import (
	"context"
	"errors"

	trm "github.com/avito-tech/go-transaction-manager/trm/v2/manager"

	"avito-tenders/internal/api/bids"
	"avito-tenders/internal/api/bids/dtos"
	"avito-tenders/internal/api/employee"
	"avito-tenders/internal/api/organization"
	"avito-tenders/internal/entity"
	"avito-tenders/pkg/apperror"
)

type usecase struct {
	repo      bids.Repository
	orgRepo   organization.Repository
	empRepo   employee.Repository
	trManager *trm.Manager
}

type Opts struct {
	Repo      bids.Repository
	OrgRepo   organization.Repository
	EmpRepo   employee.Repository
	TrManager *trm.Manager
}

func NewUsecase(createOpts Opts) bids.Usecase {
	return &usecase{
		repo:      createOpts.Repo,
		trManager: createOpts.TrManager,
		orgRepo:   createOpts.OrgRepo,
		empRepo:   createOpts.EmpRepo,
	}
}

func (u usecase) Create(ctx context.Context, req dtos.CreateBidRequest) (dtos.BidResponse, error) {
	var result entity.Bid
	err := u.trManager.Do(ctx, func(ctx context.Context) error {
		// Check does author exist.
		switch req.AuthorType {
		case entity.AuthorOrganization:
			_, err := u.orgRepo.GetUserOrganization(ctx, req.AuthorId)
			if err != nil {
				return err
			}

		case entity.AuthorUser:
			_, err := u.empRepo.FindById(ctx, req.AuthorId)
			if err != nil {
				return err
			}
		default:
			return apperror.BadRequest(errors.New("author type is invalid"))
		}

		// Create bid.
		createdBid, err := u.repo.Create(ctx, req.ToEntity())
		if err != nil {
			return err
		}

		result = createdBid

		return nil
	})
	if err != nil {
		return dtos.BidResponse{}, err
	}
	return dtos.NewBidResponse(result), nil
}

func (u usecase) FindByUsername(ctx context.Context, username string) ([]dtos.BidResponse, error) {
	// TODO implement me
	panic("implement me")
}

func (u usecase) FindByTenderId(ctx context.Context, req dtos.FindByTenderIdRequest) ([]dtos.BidResponse, error) {
	// TODO implement me
	panic("implement me")
}

func (u usecase) GetStatusById(ctx context.Context, bidId string, username string) (entity.BidStatus, error) {
	// TODO implement me
	panic("implement me")
}

func (u usecase) UpdateStatusById(ctx context.Context, req dtos.UpdateStatusRequest) (dtos.BidResponse, error) {
	// TODO implement me
	panic("implement me")
}

func (u usecase) SubmitDecision(ctx context.Context, req dtos.SubmitDecisionRequest) (dtos.BidResponse, error) {
	// TODO implement me
	panic("implement me")
}

func (u usecase) SendFeedback(ctx context.Context, req dtos.SendFeedbackRequest) (dtos.BidResponse, error) {
	// TODO implement me
	panic("implement me")
}

func (u usecase) Rollback(ctx context.Context, req dtos.RollbackRequest) (dtos.BidResponse, error) {
	// TODO implement me
	panic("implement me")
}

func (u usecase) FindReviewsByTenderId(ctx, req dtos.FindReviewsRequest) ([]entity.Review, error) {
	// TODO implement me
	panic("implement me")
}
