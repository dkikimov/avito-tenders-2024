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

type Usecase struct {
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

func NewUsecase(createOpts Opts) *Usecase {
	return &Usecase{
		repo:      createOpts.Repo,
		trManager: createOpts.TrManager,
		orgRepo:   createOpts.OrgRepo,
		empRepo:   createOpts.EmpRepo,
		tendRepo:  createOpts.TenderRepo,
	}
}

func (u Usecase) Create(ctx context.Context, req dtos.CreateBidRequest) (dtos.BidResponse, error) {
	var result entity.Bid
	err := u.trManager.Do(ctx, func(ctx context.Context) error {
		// Check does user exist.
		_, err := u.empRepo.FindByID(ctx, req.AuthorID)
		if err != nil {
			return err
		}

		// Check does author exist.
		switch req.AuthorType {
		case entity.AuthorOrganization:
			_, err := u.orgRepo.GetUserOrganization(ctx, req.AuthorID)
			if err != nil {
				if errors.Is(err, apperror.ErrNotFound) {
					return apperror.Forbidden(apperror.ErrForbidden)
				}

				return err
			}

		case entity.AuthorUser:
			break
		default:
			return apperror.BadRequest(errors.New("author type is invalid"))
		}

		// Check tender status.
		tender, err := u.tendRepo.FindByID(ctx, req.TenderID)
		if err != nil {
			return err
		}
		if tender.Status != entity.TenderPublished {
			return apperror.Forbidden(apperror.ErrForbidden)
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

func (u Usecase) FindByUsername(ctx context.Context, username string, pagination queryparams.Pagination) ([]dtos.BidResponse, error) {
	bidsList, err := u.repo.FindByUsername(ctx, models.FindByUsername{
		Username:   username,
		Pagination: pagination,
	})
	if err != nil {
		return nil, err
	}

	return dtos.NewBidResponseList(bidsList), nil
}

func (u Usecase) FindByTenderID(ctx context.Context, req dtos.FindByTenderIDRequest) ([]dtos.BidResponse, error) {
	var filteredBidsList []dtos.BidResponse

	err := u.trManager.Do(ctx, func(ctx context.Context) error {
		bidsList, err := u.repo.FindByTenderID(ctx, models.FindByTenderID{
			TenderID:   req.TenderID,
			Pagination: req.Pagination,
		})
		if err != nil {
			return err
		}

		tender, err := u.tendRepo.FindByID(ctx, req.TenderID)
		if err != nil {
			return err
		}

		isResponsible, err := u.orgRepo.IsOrganizationResponsible(ctx, tender.OrganizationID, req.Username)
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

func (u Usecase) GetStatusByID(ctx context.Context, bidID, username string) (entity.BidStatus, error) {
	bid, err := u.repo.FindByID(ctx, bidID)
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

func (u Usecase) UpdateStatusByID(ctx context.Context, req dtos.UpdateStatusRequest) (dtos.BidResponse, error) {
	bid, err := u.repo.FindByID(ctx, req.BidID)
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

func (u Usecase) SubmitDecision(ctx context.Context, req dtos.SubmitDecisionRequest) (dtos.BidResponse, error) {
	var resultBid dtos.BidResponse
	err := u.trManager.Do(ctx, func(ctx context.Context) error {
		bid, err := u.repo.FindByID(ctx, req.BidID)
		if err != nil {
			return err
		}

		tender, err := u.tendRepo.FindByID(ctx, bid.TenderID)
		if err != nil {
			return err
		}

		isResponsible, err := u.orgRepo.IsOrganizationResponsible(ctx, tender.OrganizationID, req.Username)
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
		}

		if req.Decision == entity.DecisionApproved {
			user, err := u.empRepo.FindByUsername(ctx, req.Username)
			if err != nil {
				return err
			}

			err = u.repo.SubmitApproveDecision(ctx, bid.ID, user.ID)
			if err != nil {
				return err
			}

			approveBidCount, err := u.repo.GetBidApproveAmount(ctx, bid.ID)
			if err != nil {
				return err
			}

			responsibleList, err := u.orgRepo.GetOrganizationResponsible(ctx, tender.OrganizationID)
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
			} else {
				bid, err = u.repo.FindByID(ctx, req.BidID)
				if err != nil {
					return err
				}

				resultBid = dtos.NewBidResponse(bid)
			}
		}

		return nil
	})
	if err != nil {
		return dtos.BidResponse{}, err
	}

	return resultBid, nil
}

func (u Usecase) SendFeedback(ctx context.Context, req dtos.SendFeedbackRequest) (dtos.BidResponse, error) {
	var resultBid entity.Bid
	err := u.trManager.Do(ctx, func(ctx context.Context) error {
		bid, err := u.repo.FindByID(ctx, req.BidID)
		if err != nil {
			return err
		}
		resultBid = bid

		tender, err := u.tendRepo.FindByID(ctx, bid.TenderID)
		if err != nil {
			slog.Error("couldn't find tender by bid id")
			return err
		}

		isResponsible, err := u.orgRepo.IsOrganizationResponsible(ctx, tender.OrganizationID, req.Username)
		if err != nil {
			return err
		}
		if !isResponsible {
			return apperror.Forbidden(apperror.ErrForbidden)
		}

		err = u.repo.SendFeedback(ctx, models.SendFeedback{
			BidID:    req.BidID,
			Feedback: req.Feedback,
		})
		if err != nil {
			return err
		}

		return nil
	})
	if err != nil {
		return dtos.BidResponse{}, err
	}

	return dtos.NewBidResponse(resultBid), nil
}

func (u Usecase) Rollback(ctx context.Context, req dtos.RollbackRequest) (dtos.BidResponse, error) {
	oldBid, err := u.repo.FindByIDFromHistory(ctx, req.BidID, req.Version)
	if err != nil {
		return dtos.BidResponse{}, err
	}

	responsible, err := u.AuthorHasPermissions(ctx, oldBid, req.Username)
	if err != nil {
		return dtos.BidResponse{}, err
	}
	if !responsible {
		return dtos.BidResponse{}, apperror.Forbidden(apperror.ErrUnauthorized)
	}

	updatedBid, err := u.repo.Update(ctx, oldBid)
	if err != nil {
		return dtos.BidResponse{}, err
	}

	return dtos.NewBidResponse(updatedBid), nil
}

func (u Usecase) FindReviewsByTenderID(ctx context.Context, req dtos.FindReviewsRequest) ([]dtos.ReviewResponse, error) {
	var resultReviews []entity.Review
	err := u.trManager.Do(ctx, func(ctx context.Context) error {
		emp, err := u.empRepo.FindByUsername(ctx, req.RequesterUsername)
		if err != nil {
			return err
		}

		author, err := u.empRepo.FindByUsername(ctx, req.AuthorUsername)
		if err != nil {
			return err
		}

		org, err := u.orgRepo.GetUserOrganization(ctx, emp.ID)
		if err != nil {
			return err
		}

		tender, err := u.tendRepo.FindByID(ctx, req.TenderID)
		if err != nil {
			return err
		}
		if tender.OrganizationID != org.ID {
			return apperror.Forbidden(apperror.ErrForbidden)
		}

		bidsList, err := u.repo.FindBidsByOrganization(ctx, org.ID)
		if err != nil {
			return err
		}

		filteredBidsList := make([]entity.Bid, 0)
		for _, bid := range bidsList {
			if bid.AuthorID == author.ID {
				filteredBidsList = append(filteredBidsList, bid)
			}
		}

		reviews, err := u.repo.FindReviews(ctx, models.FindReview{
			Bids:       filteredBidsList,
			Pagination: req.Pagination,
		})
		if err != nil {
			return err
		}

		resultReviews = reviews

		return nil
	})
	if err != nil {
		return nil, err
	}

	return dtos.NewReviewResponseList(resultReviews), nil
}

func (u Usecase) AuthorHasPermissions(ctx context.Context, bid entity.Bid, username string) (bool, error) {
	switch bid.AuthorType {
	case entity.AuthorOrganization:
		org, err := u.orgRepo.GetUserOrganization(ctx, bid.AuthorID)
		if err != nil {
			return false, err
		}

		isResponsible, err := u.orgRepo.IsOrganizationResponsible(ctx, org.ID, username)
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

		if emp.ID != bid.AuthorID {
			return false, apperror.Forbidden(apperror.ErrForbidden)
		}
	default:
		slog.Error("Unknown author type", "bid", bid)
		return false, apperror.InternalServerError(apperror.ErrInternal)
	}

	return true, nil
}

func (u Usecase) Edit(ctx context.Context, req dtos.EditBidRequest) (dtos.BidResponse, error) {
	bid, err := u.repo.FindByID(ctx, req.BidID)
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
