package handler

import (
	"context"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/emptypb"
	"srvs/goods_srv/global"
	"srvs/goods_srv/model"

	"srvs/goods_srv/proto/proto"
)

// 轮播图
func (s *GoodsServer) BannerList(context.Context, *emptypb.Empty) (*proto.BannerListResponse, error) {
	banners := []model.Banner{}

	result := global.DB.Find(&banners)

	if result.Error != nil {
		return nil, result.Error
	}
	rsp := &proto.BannerListResponse{}

	rsp.Total = int32(result.RowsAffected)

	for _, banner := range banners {
		rsp.Data = append(rsp.Data, &proto.BannerResponse{
			Id:    banner.ID,
			Index: banner.Index,
			Image: banner.Image,
			Url:   banner.Url,
		})
	}

	return rsp, nil
}
func (s *GoodsServer) CreateBanner(c context.Context, r *proto.BannerRequest) (*proto.BannerResponse, error) {
	var banner model.Banner

	banner.Index = r.Index
	banner.Image = r.Image
	banner.Url = r.Url

	result := global.DB.Save(&banner)

	if result.Error != nil {
		return nil, result.Error
	}

	rsp := &proto.BannerResponse{
		Id:    banner.ID,
		Index: banner.Index,
		Image: banner.Image,
		Url:   banner.Url,
	}

	return rsp, nil
}
func (s *GoodsServer) DeleteBanner(c context.Context, r *proto.BannerRequest) (*emptypb.Empty, error) {
	result := global.DB.Delete(&model.Banner{}, r.Id)

	if result.Error != nil {
		return nil, result.Error
	}

	if result.RowsAffected == 0 {
		return nil, status.Errorf(codes.NotFound, "轮播图不存在")
	}

	return &emptypb.Empty{}, nil
}
func (s *GoodsServer) UpdateBanner(c context.Context, req *proto.BannerRequest) (*emptypb.Empty, error) {
	banner := model.Banner{}

	if result := global.DB.First(&banner, req.Id); result.RowsAffected == 0 {
		return nil, status.Errorf(codes.NotFound, "轮播图不存在")
	}

	if req.Image != "" {
		banner.Image = req.Image
	}

	if req.Url != "" {
		banner.Url = req.Url
	}

	if req.Index > 0 {
		banner.Index = req.Index
	}

	result := global.DB.Save(&banner)

	if result.Error != nil {
		return nil, result.Error
	}

	return &emptypb.Empty{}, nil
}
