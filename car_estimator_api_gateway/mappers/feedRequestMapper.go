package mappers

import (
	feed "github.com/nikita-itmo-gh-acc/car_estimator_api_contracts/gen/feed_v1"
	"github.com/nikita-itmo-gh-acc/car_estimator_api_gateway/domain"
)

func ToDomain(message *feed.CarListing) *domain.CarListing {
	return &domain.CarListing{
		ListingId:        message.ListingId,
		SellerId:         message.SellerId,
		Description:      message.Description,
		PostedAt:         message.PostedAt.AsTime(),
		Status:           message.Status,
		DealType:         message.DealType,
		Price:            message.Price,
		CarId:            message.CarId,
		Mileage:          message.Mileage,
		OwnersCount:      message.OwnersCount,
		AccidentsCount:   message.AccidentsCount,
		Condition:        message.Condition,
		Color:            message.Color,
		ConfigId:         message.ConfigId,
		EngineType:       message.EngineType,
		EngineVolume:     message.EngineVolume,
		EnginePower:      message.EnginePower,
		Cylinders:        message.Cylinders,
		Transmission:     message.Transmission,
		Drivetrain:       message.Drivetrain,
		ModelId:          message.ModelId,
		ModelName:        message.ModelName,
		Make:             message.Make,
		Year:             message.Year,
		BodyType:         message.BodyType,
		Generation:       message.Generation,
		WeightKg:         message.WeightKg,
		SellerName:       message.SellerName,
		SellerRating:     message.SellerRating,
		SellerSalesCount: message.SellerSalesCount,
		SellerIsBusiness: message.SellerIsBusiness,
	}
}

func ToMessage(listing *domain.CarListing) *feed.CarListing {
	return &feed.CarListing{
		SellerId:    		listing.SellerId,
		Description: 		listing.Description,
		Status:      		listing.Status,
		DealType:    		listing.DealType,
		Price:       		listing.Price,
		Tags:		 		listing.Tags,
		CarId: 		 		listing.CarId,
		Mileage: 	 		listing.Mileage,
		OwnersCount: 		listing.OwnersCount,
		AccidentsCount: 	listing.AccidentsCount,
		Condition: 	 		listing.Condition,
		Color: 		 		listing.Color,
		ConfigId: 	 		listing.ConfigId,
		EngineType:     	listing.EngineType,
		EngineVolume:   	listing.EngineVolume,
		EnginePower:    	listing.EnginePower,
		Cylinders:      	listing.Cylinders,
		Transmission:   	listing.Transmission,
		Drivetrain:      	listing.Drivetrain,
		ModelId:         	listing.ModelId,
		ModelName:       	listing.ModelName,
		Make:            	listing.Make,
		Year:            	listing.Year,
		BodyType:        	listing.BodyType,
		Generation:      	listing.Generation,
		WeightKg:        	listing.WeightKg,
		SellerName:      	listing.SellerName,
		SellerRating:    	listing.SellerRating,
		SellerSalesCount:	listing.SellerSalesCount,
		SellerIsBusiness:	listing.SellerIsBusiness,
	}
}
