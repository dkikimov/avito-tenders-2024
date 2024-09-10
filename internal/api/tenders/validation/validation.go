package validation

import (
	"github.com/asaskevich/govalidator"

	"avito-tenders/internal/entity"
)

func AddValidations() {
	govalidator.CustomTypeTagMap.Set("tender_status", func(i interface{}, o interface{}) bool {
		value, ok := i.(entity.TenderStatus)
		if !ok {
			return false
		}

		switch value {
		case entity.TenderCreated, entity.TenderClosed, entity.TenderPublished:
			return true
		default:
			return false
		}
	})

	govalidator.CustomTypeTagMap.Set("service_type", func(i interface{}, o interface{}) bool {
		value, ok := i.(entity.ServiceType)
		if !ok {
			return false
		}

		switch value {
		case entity.ServiceConstruction, entity.ServiceDelivery, entity.ServiceManufacture:
			return true
		default:
			return false
		}
	})
}
