package server

import (
	commonv1 "github.com/chekrzk/ufanet-home-manager/contracts/gen/go/common/v1"
	newsv1 "github.com/chekrzk/ufanet-home-manager/contracts/gen/go/news/v1"
	"github.com/chekrzk/ufanet-home-manager/news-service/internal/models"
	"google.golang.org/protobuf/types/known/timestamppb"
)

// newsFilterFromProto отделяет фильтры gRPC-контракта от domain-фильтров ленты.
func newsFilterFromProto(req *newsv1.ListNewsRequest) models.NewsFilter {
	return models.NewsFilter{
		Actor: userContextFromProto(req.GetUser()),
		Pagination: models.Pagination{
			Page:  int(req.GetPagination().GetPage()),
			Limit: int(req.GetPagination().GetLimit()),
		},
		DateFrom: req.GetDateFrom(),
		DateTo:   req.GetDateTo(),
	}
}

// createNewsCommandFromProto собирает доменную команду публикации из proto.
func createNewsCommandFromProto(req *newsv1.CreateNewsRequest) models.CreateNewsCommand {
	return models.CreateNewsCommand{
		Author:  userContextFromProto(req.GetAuthor()),
		Title:   req.GetTitle(),
		Body:    req.GetBody(),
		HouseID: req.GetHouseId(),
	}
}

// userContextFromProto переносит trusted auth context между сервисами.
func userContextFromProto(user *commonv1.UserContext) models.UserContext {
	if user == nil {
		return models.UserContext{}
	}
	return models.UserContext{UserID: user.GetUserId(), Role: user.GetRole()}
}

// newsToProto возвращает наружу только contract-представление новости.
func newsToProto(item models.News) *commonv1.News {
	return &commonv1.News{
		Id:        item.ID,
		Title:     item.Title,
		Body:      item.Body,
		HouseId:   item.HouseID,
		CreatedAt: timestamppb.New(item.CreatedAt),
	}
}
