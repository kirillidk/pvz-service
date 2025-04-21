package grpc

import (
	"context"

	pvz_v1 "github.com/kirillidk/pvz-service/api/proto/pvz/pvz_v1"
	"github.com/kirillidk/pvz-service/internal/dto"
	"github.com/kirillidk/pvz-service/internal/repository"
	"google.golang.org/protobuf/types/known/timestamppb"
)

type PVZService struct {
	pvzRepository repository.PVZRepositoryInterface
	pvz_v1.UnimplementedPVZServiceServer
}

func NewPVZService(pvzRepo repository.PVZRepositoryInterface) *PVZService {
	return &PVZService{
		pvzRepository: pvzRepo,
	}
}

func (s *PVZService) GetPVZList(ctx context.Context, req *pvz_v1.GetPVZListRequest) (*pvz_v1.GetPVZListResponse, error) {
	filter := dto.PVZFilterQuery{
		Page:  1,
		Limit: 1000,
	}

	pvzList, err := s.pvzRepository.GetPVZList(ctx, filter)
	if err != nil {
		return nil, err
	}

	response := &pvz_v1.GetPVZListResponse{
		Pvzs: make([]*pvz_v1.PVZ, 0, len(pvzList)),
	}

	for _, pvz := range pvzList {
		response.Pvzs = append(response.Pvzs, &pvz_v1.PVZ{
			Id:               pvz.ID,
			RegistrationDate: timestamppb.New(pvz.RegistrationDate),
			City:             pvz.City,
		})
	}

	return response, nil
}
