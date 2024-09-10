package query

import (
	"errors"
	"fmt"
	"net/url"
	"strconv"

	"avito-tenders/pkg/apperror"
)

const (
	defaultLimit  = 5
	defaultOffset = 0

	minLimit = 0
	maxLimit = 50

	minOffset = 0
)

type Pagination struct {
	Limit  int
	Offset int
}

func ParsePagination(values url.Values) (Pagination, error) {
	pagination := Pagination{
		Limit:  defaultLimit,
		Offset: defaultOffset,
	}

	limit, offset := values.Get("limit"), values.Get("offset")
	if limit != "" {
		parsedLimit, err := strconv.Atoi(limit)
		if err != nil {
			return Pagination{}, apperror.BadRequest(errors.New("limit is not number"))
		}

		if parsedLimit < minLimit || parsedLimit > maxLimit {
			return Pagination{}, apperror.BadRequest(fmt.Errorf("limit must be between %d and %d", minLimit, maxLimit))
		}

		pagination.Limit = parsedLimit
	}

	if offset != "" {
		parsedOffset, err := strconv.Atoi(offset)
		if err != nil {
			return Pagination{}, apperror.BadRequest(apperror.ErrInvalidInput)
		}

		if parsedOffset < minOffset {
			return Pagination{}, apperror.BadRequest(fmt.Errorf("offset cannot be less than %d", minOffset))
		}

		pagination.Offset = parsedOffset
	}

	return pagination, nil
}
