package grpc

import (
	"context"
	"log/slog"
	"net"

	pb "github.com/infinity-ocean/ikakbolit/3-api-grpc-Homework/grpc/ikakbolit"
	"github.com/infinity-ocean/ikakbolit/internal/domain/entity"
	"github.com/infinity-ocean/ikakbolit/pkg/middlewarex/interceptor"
	"google.golang.org/grpc"
	"google.golang.org/protobuf/types/known/timestamppb"
)

type gRPCServer struct {
	pb.UnimplementedIkakbolitServiceServer
	svc  service
	port string
	log  *slog.Logger
}

type service interface {
	AddSchedule(context.Context, entity.Schedule) (int, error)
	GetScheduleIDs(context.Context, int) ([]int, error)
	GetScheduleWithIntake(context.Context, int, int) (entity.Schedule, error)
	GetNextTakings(context.Context, int) ([]entity.Schedule, error)
}

func NewGRPCServer(svc service, port string, logger *slog.Logger) *gRPCServer {
	return &gRPCServer{
		svc:  svc,
		port: port,
		log:  logger,
	}
}

func (s *gRPCServer) Run() error {
	lis, err := net.Listen("tcp", s.port)
	if err != nil {
		return err
	}

	grpcServer := grpc.NewServer(
		grpc.UnaryInterceptor(interceptor.UnaryLoggingInterceptor(s.log)),
	)

	pb.RegisterIkakbolitServiceServer(grpcServer, s)

	s.log.Info("starting gRPC server", slog.String("port", s.port))
	return grpcServer.Serve(lis)
}

func (s *gRPCServer) AddSchedule(ctx context.Context, req *pb.RequestSchedule) (*pb.ResponseScheduleID, error) {
    schedule := entity.Schedule{
        UserID:       int(req.UserId),
        CureName:     req.CureName,
        DosesPerDay:  int(req.DosesPerDay),
        DurationDays: int(req.DurationDays),
    }
    id, err := s.svc.AddSchedule(ctx, schedule)
    if err != nil {
        return nil, err
    }
    return &pb.ResponseScheduleID{ScheduleId: int64(id)}, nil
}

func (s *gRPCServer) GetScheduleIDs(ctx context.Context, req *pb.RequestUserID) (*pb.ResponseScheduleIDs, error) {
    ids, err := s.svc.GetScheduleIDs(ctx, int(req.UserId))
    if err != nil {
        return nil, err
    }
    out := &pb.ResponseScheduleIDs{}
    for _, id := range ids {
        out.SchdeduleIds = append(out.SchdeduleIds, int64(id))
    }
    return out, nil
}

func (s *gRPCServer) GetSchedule(ctx context.Context, req *pb.RequestUserIDScheduleID) (*pb.ResponseSchedule, error) {
    sched, err := s.svc.GetScheduleWithIntake(ctx, int(req.UserId), int(req.ScheduleId))
    if err != nil {
        return nil, err
    }
    return &pb.ResponseSchedule{
        Id:           int64(sched.ID),
        UserId:       int64(sched.UserID),
        CureName:     sched.CureName,
        DosesPerDay:  int64(sched.DosesPerDay),
        DurationDays: int64(sched.DurationDays),
        CreatedAt:    timestamppb.New(sched.CreatedAt),
        Intakes:      sched.Intakes,
    }, nil
}

func (s *gRPCServer) GetNextTakings(ctx context.Context, req *pb.RequestNextTakings) (*pb.ResponseNextTakings, error) {
    schedules, err := s.svc.GetNextTakings(ctx, int(req.UserId))
    if err != nil {
        return nil, err
    }
    var grpcSchedules []*pb.ResponseSchedule
    for _, sched := range schedules {
        grpcSchedules = append(grpcSchedules, &pb.ResponseSchedule{
            Id:           int64(sched.ID),
            UserId:       int64(sched.UserID),
            CureName:     sched.CureName,
            DosesPerDay:  int64(sched.DosesPerDay),
            DurationDays: int64(sched.DurationDays),
            CreatedAt:    timestamppb.New(sched.CreatedAt),
            Intakes:      sched.Intakes,
        })
    }
    return &pb.ResponseNextTakings{Schedules: grpcSchedules}, nil
}