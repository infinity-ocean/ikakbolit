package controller

import (
	"context"
	"log"
	"net"

	"github.com/infinity-ocean/ikakbolit/3-api-grpc-Homework/grpc/ikakbolit"
	"github.com/infinity-ocean/ikakbolit/internal/model"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/timestamppb"
)

type GrpcServer struct {
	ikakbolit.UnimplementedIkakbolitServiceServer
	service service
}

func NewGrpcServer(service service) *GrpcServer {
	return &GrpcServer{service: service}
}

func StartGrpcServer(srv *GrpcServer, port string) error {
	listener, err := net.Listen("tcp", ":"+port)
	if err != nil {
		return err
	}

	grpcServer := grpc.NewServer()
	ikakbolit.RegisterIkakbolitServiceServer(grpcServer, srv)
	reflection.Register(grpcServer)

	log.Printf("gRPC server listening on port %s", port)
	return grpcServer.Serve(listener)
}

func protoToModelSchedule(s *ikakbolit.Schedule) model.Schedule {
	return model.Schedule{
		ID:           int(s.Id),
		UserID:       int(s.UserId),
		CureName:     s.CureName,
		DosesPerDay:  int(s.DosesPerDay),
		DurationDays: int(s.DurationDays),
		CreatedAt:    s.CreatedAt.AsTime(),
		DayStart:     s.DayStart.AsTime(),
		DayFinish:    s.DayFinish.AsTime(),
		Intakes:      s.Intakes,
	}
}

func (s *GrpcServer) AddSchedule(ctx context.Context, req *ikakbolit.Schedule) (*ikakbolit.ResponseScheduleID, error) {
	schedule := protoToModelSchedule(req)

	id, err := s.service.AddSchedule(schedule)
	if err != nil {
		return nil, err
	}

	return &ikakbolit.ResponseScheduleID{
		ScheduleId: int64(id),
	}, nil
}

func (s *GrpcServer) GetScheduleIDs(ctx context.Context, req *ikakbolit.UserIDRequest) (*ikakbolit.ScheduleIDsResponse, error) {
	ids, err := s.service.GetScheduleIDs(int(req.UserId))
	if err != nil {
		return nil, status.Errorf(codes.Internal, "could not get schedule IDs: %v", err)
	}

	idList := make([]int64, len(ids))
	for i, id := range ids {
		idList[i] = int64(id)
	}

	return &ikakbolit.ScheduleIDsResponse{
		Ids: idList,
	}, nil
}

func (s *GrpcServer) GetSchedule(ctx context.Context, req *ikakbolit.ScheduleRequest) (*ikakbolit.Schedule, error) {
	schedule, err := s.service.GetScheduleWithIntake(int(req.UserId), int(req.ScheduleId))
	if err != nil {
		return nil, status.Errorf(codes.Internal, "could not get schedule: %v", err)
	}

	return &ikakbolit.Schedule{
		Id:           int64(schedule.ID),
		UserId:       int64(schedule.UserID),
		CureName:     schedule.CureName,
		DosesPerDay:  int64(schedule.DosesPerDay),
		DurationDays: int64(schedule.DurationDays),
		CreatedAt:    timestamppb.New(schedule.CreatedAt),
		DayStart:     timestamppb.New(schedule.DayStart),
		DayFinish:    timestamppb.New(schedule.DayFinish),
		Intakes:      schedule.Intakes,
	}, nil
}

func (s *GrpcServer) GetNextTakings(ctx context.Context, req *ikakbolit.UserIDRequest) (*ikakbolit.SchedulesInWindow, error) {
	schedules, err := s.service.GetNextTakings(int(req.UserId))
	if err != nil {
		return nil, status.Errorf(codes.Internal, "could not get next takings: %v", err)
	}

	result := make([]*ikakbolit.Schedule, 0, len(schedules))
	for _, sch := range schedules {
		result = append(result, &ikakbolit.Schedule{
			Id:           int64(sch.ID),
			UserId:       int64(sch.UserID),
			CureName:     sch.CureName,
			DosesPerDay:  int64(sch.DosesPerDay),
			DurationDays: int64(sch.DurationDays),
			CreatedAt:    timestamppb.New(sch.CreatedAt),
			DayStart:     timestamppb.New(sch.DayStart),
			DayFinish:    timestamppb.New(sch.DayFinish),
			Intakes:      sch.Intakes,
		})
	}

	return &ikakbolit.SchedulesInWindow{
		Schedules: result,
	}, nil
}