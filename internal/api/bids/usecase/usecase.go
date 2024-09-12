package usecase

import (
	"context"
	"errors"
	"log/slog"

	trm "github.com/avito-tech/go-transaction-manager/trm/v2/manager"

	"avito-tenders/internal/api/bids"
	"avito-tenders/internal/api/bids/dtos"
	"avito-tenders/internal/api/bids/models"
	"avito-tenders/internal/api/employee"
	"avito-tenders/internal/api/organization"
	"avito-tenders/internal/api/tenders"
	"avito-tenders/internal/entity"
	"avito-tenders/pkg/apperror"
	"avito-tenders/pkg/queryparams"
)

type usecase struct {
	repo      bids.Repository
	orgRepo   organization.Repository
	empRepo   employee.Repository
	tendRepo  tenders.Repository
	trManager *trm.Manager
}

type Opts struct {
	Repo       bids.Repository
	OrgRepo    organization.Repository
	EmpRepo    employee.Repository
	TenderRepo tenders.Repository
	TrManager  *trm.Manager
}

func NewUsecase(createOpts Opts) bids.Usecase {
	return &usecase{
		repo:      createOpts.Repo,
		trManager: createOpts.TrManager,
		orgRepo:   createOpts.OrgRepo,
		empRepo:   createOpts.EmpRepo,
		tendRepo:  createOpts.TenderRepo,
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

func (u usecase) FindByUsername(ctx context.Context, username string, pagination queryparams.Pagination) ([]dtos.BidResponse, error) {
	bidsList, err := u.repo.FindByUsername(ctx, models.FindByUsername{
		Username:   username,
		Pagination: pagination,
	})
	if err != nil {
		return nil, err
	}

	return dtos.NewBidResponseList(bidsList), nil
}

func (u usecase) FindByTenderId(ctx context.Context, req dtos.FindByTenderIdRequest) ([]dtos.BidResponse, error) {
	var filteredBidsList []dtos.BidResponse

	err := u.trManager.Do(ctx, func(ctx context.Context) error {
		bidsList, err := u.repo.FindByTenderId(ctx, models.FindByTenderId{
			TenderId:   req.TenderId,
			Pagination: req.Pagination,
		})
		if err != nil {
			return err
		}

		tender, err := u.tendRepo.FindById(ctx, req.TenderId)
		if err != nil {
			return err
		}

		isResponsible, err := u.orgRepo.IsOrganizationResponsible(ctx, tender.OrganizationId, req.Username)
		if err != nil {
			return err
		}

		filteredBidsList = make([]dtos.BidResponse, 0, len(bidsList))
		for _, bid := range bidsList {
			// If user is responsible for tender we need to add `Published` bids
			if isResponsible && bid.Status == entity.BidPublished {
				filteredBidsList = append(filteredBidsList, dtos.NewBidResponse(bid))
				continue
			}

			// If user created bid or bid was created by user from his company, then we need to add this bid.
			// Status doesn't matter
			has, _ := u.AuthorHasPermissions(ctx, bid, req.Username)
			if has {
				filteredBidsList = append(filteredBidsList, dtos.NewBidResponse(bid))
			}
		}

		return nil
	})
	if err != nil {
		return nil, err
	}

	return filteredBidsList, nil
}

func (u usecase) GetStatusById(ctx context.Context, bidId string, username string) (entity.BidStatus, error) {
	bid, err := u.repo.FindByID(ctx, bidId)
	if err != nil {
		return "", err
	}

	has, err := u.AuthorHasPermissions(ctx, bid, username)
	if err != nil {
		return "", err
	}
	if !has {
		return "", apperror.Forbidden(apperror.ErrForbidden)
	}

	return bid.Status, nil
}

func (u usecase) UpdateStatusById(ctx context.Context, req dtos.UpdateStatusRequest) (dtos.BidResponse, error) {
	bid, err := u.repo.FindByID(ctx, req.BidId)
	if err != nil {
		return dtos.BidResponse{}, err
	}

	has, err := u.AuthorHasPermissions(ctx, bid, req.Username)
	if err != nil {
		return dtos.BidResponse{}, err
	}
	if !has {
		return dtos.BidResponse{}, apperror.Forbidden(apperror.ErrForbidden)
	}

	newBid := bid
	newBid.Status = req.Status

	updatedBid, err := u.repo.Update(ctx, newBid)
	if err != nil {
		return dtos.BidResponse{}, err
	}

	return dtos.NewBidResponse(updatedBid), nil
}

func (u usecase) SubmitDecision(ctx context.Context, req dtos.SubmitDecisionRequest) (dtos.BidResponse, error) {
	var resultBid dtos.BidResponse
	err := u.trManager.Do(ctx, func(ctx context.Context) error {
		bid, err := u.repo.FindByID(ctx, req.BidId)
		if err != nil {
			return err
		}

		tender, err := u.tendRepo.FindById(ctx, bid.TenderId)
		if err != nil {
			return err
		}

		isResponsible, err := u.orgRepo.IsOrganizationResponsible(ctx, tender.OrganizationId, req.Username)
		if err != nil {
			return err
		}
		if !isResponsible {
			return apperror.Forbidden(apperror.ErrForbidden)
		}

		if req.Decision == entity.DecisionRejected {
			newBid := bid
			newBid.Status = entity.BidRejected

			updatedBid, err := u.repo.Update(ctx, newBid)
			if err != nil {
				return err
			}

			resultBid = dtos.NewBidResponse(updatedBid)

			return nil
		} else if req.Decision == entity.DecisionApproved {
			user, err := u.empRepo.FindByUsername(ctx, req.Username)
			if err != nil {
				return err
			}

			err = u.repo.SubmitApproveDecision(ctx, bid.Id, user.Id)
			if err != nil {
				return err
			}

			approveBidCount, err := u.repo.GetBidApproveAmount(ctx, bid.Id)
			if err != nil {
				return err
			}

			responsibleList, err := u.orgRepo.GetOrganizationResponsible(ctx, tender.OrganizationId)
			if err != nil {
				return err
			}

			if approveBidCount >= min(3, len(responsibleList)) {
				// Update bid status
				newBid := bid
				newBid.Status = entity.BidApproved

				updatedBid, err := u.repo.Update(ctx, newBid)
				if err != nil {
					return err
				}
				resultBid = dtos.NewBidResponse(updatedBid)

				// Update tender status
				newTender := tender
				newTender.Status = entity.TenderClosed
				_, err = u.tendRepo.Update(ctx, newTender)
				if err != nil {
					return err
				}
			}
		}

		return nil
	})
	if err != nil {
		return dtos.BidResponse{}, err
	}

	return resultBid, nil
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

func (u usecase) AuthorHasPermissions(ctx context.Context, bid entity.Bid, username string) (bool, error) {
	switch bid.AuthorType {
	case entity.AuthorOrganization:
		org, err := u.orgRepo.GetUserOrganization(ctx, bid.AuthorId)
		if err != nil {
			return false, err
		}

		isResponsible, err := u.orgRepo.IsOrganizationResponsible(ctx, org.Id, username)
		if err != nil {
			return false, err
		}
		if !isResponsible {
			return false, apperror.Forbidden(apperror.ErrForbidden)
		}
	case entity.AuthorUser:
		emp, err := u.empRepo.FindByUsername(ctx, username)
		if err != nil {
			return false, err
		}

		if emp.Id != bid.AuthorId {
			return false, apperror.Forbidden(apperror.ErrForbidden)
		}
	default:
		slog.Error("Unknown author type", "bid", bid)
		return false, apperror.InternalServerError(apperror.ErrInternal)
	}

	return true, nil
}

func (u usecase) Edit(ctx context.Context, req dtos.EditBidRequest) (dtos.BidResponse, error) {
	bid, err := u.repo.FindByID(ctx, req.BidId)
	if err != nil {
		return dtos.BidResponse{}, err
	}

	has, err := u.AuthorHasPermissions(ctx, bid, req.Username)
	if err != nil {
		return dtos.BidResponse{}, err
	}
	if !has {
		return dtos.BidResponse{}, apperror.Forbidden(apperror.ErrForbidden)
	}

	newBid := bid
	if len(req.Name) != 0 {
		newBid.Name = req.Name
	}
	if len(req.Description) != 0 {
		newBid.Description = req.Description
	}

	updatedBid, err := u.repo.Update(ctx, newBid)
	if err != nil {
		return dtos.BidResponse{}, err
	}

	return dtos.NewBidResponse(updatedBid), nil
}
