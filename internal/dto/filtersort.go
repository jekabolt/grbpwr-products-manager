package dto

import (
	"github.com/jekabolt/grbpwr-manager/internal/entity"
	pb_common "github.com/jekabolt/grbpwr-manager/proto/gen/common"
	"github.com/shopspring/decimal"
)

// ConvertEntityPriceFromToToPBCommon converts PriceFromTo from entity to pb_common
func ConvertEntityPriceFromToToPBCommon(pft entity.PriceFromTo) *pb_common.PriceFromTo {
	return &pb_common.PriceFromTo{
		From: pft.From.String(),
		To:   pft.To.String(),
	}
}

// ConvertEntityFilterConditionsToPBCommon converts FilterConditions from entity to pb_common
func ConvertEntityFilterConditionsToPBCommon(fc entity.FilterConditions) *pb_common.FilterConditions {
	sizes := make([]int32, len(fc.SizesIds))
	for i, v := range fc.SizesIds {
		sizes[i] = int32(v)
	}

	return &pb_common.FilterConditions{
		PriceFromTo: ConvertEntityPriceFromToToPBCommon(fc.PriceFromTo),
		OnSale:      fc.OnSale,
		Color:       fc.Color,
		CategoryId:  int32(fc.CategoryId),
		SizesIds:    sizes,
		Preorder:    fc.Preorder,
		ByTag:       fc.ByTag,
	}
}

// ConvertPBCommonOrderFactorToEntity converts OrderFactor from pb_common to entity
func ConvertPBCommonOrderFactorToEntity(of pb_common.OrderFactor) entity.OrderFactor {
	switch of {
	case pb_common.OrderFactor_ASC:
		return entity.Ascending
	case pb_common.OrderFactor_DESC:
		return entity.Descending
	default:
		return entity.Ascending // default value
	}
}

// ConvertPBCommonSortFactorToEntity converts SortFactor from pb_common to entity
func ConvertPBCommonSortFactorToEntity(sf pb_common.SortFactor) entity.SortFactor {
	switch sf {
	case pb_common.SortFactor_CREATED_AT:
		return entity.CreatedAt
	case pb_common.SortFactor_UPDATED_AT:
		return entity.UpdatedAt
	case pb_common.SortFactor_NAME:
		return entity.Name
	case pb_common.SortFactor_PRICE:
		return entity.Price
	default:
		return entity.CreatedAt // default value
	}
}

// ConvertPBCommonPriceFromToToEntity converts PriceFromTo from pb_common to entity
func ConvertPBCommonPriceFromToToEntity(pft *pb_common.PriceFromTo) entity.PriceFromTo {
	from, _ := decimal.NewFromString(pft.From)
	to, _ := decimal.NewFromString(pft.To)
	return entity.PriceFromTo{
		From: from,
		To:   to,
	}
}

// ConvertPBCommonFilterConditionsToEntity converts FilterConditions from pb_common to entity
func ConvertPBCommonFilterConditionsToEntity(fc *pb_common.FilterConditions) *entity.FilterConditions {
	if fc == nil {
		return nil
	}

	sizes := make([]int, len(fc.SizesIds))
	for i, v := range fc.SizesIds {
		sizes[i] = int(v)
	}

	return &entity.FilterConditions{
		PriceFromTo: ConvertPBCommonPriceFromToToEntity(fc.PriceFromTo),
		OnSale:      fc.OnSale,
		Color:       fc.Color,
		CategoryId:  int(fc.CategoryId),
		SizesIds:    sizes,
		Preorder:    fc.Preorder,
		ByTag:       fc.ByTag,
	}
}
