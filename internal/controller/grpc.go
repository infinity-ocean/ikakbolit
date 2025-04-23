
package controller

import (
    "google.golang.org/protobuf/types/known/timestamppb"
    "context"
    pb "github.com/infinity-ocean/ikakbolit/3-api-grpc-Homework/grpc/ikakbolit"
    "github.com/infinity-ocean/ikakbolit/internal/model"
)

type IkakbolitService interface {
    AddSchedule(model.Schedule) (int, error)
    GetScheduleIDs(int) ([]int, error)
    GetSchedule(int, int) (model.Schedule, error)
    GetNextTakings(int) ([]model.Schedule, error)
}

type server struct {
    pb.UnimplementedIkakbolitServiceServer
    svc IkakbolitService
}

func NewGRPCServer(svc IkakbolitService) *server {
    return &server{svc: svc}
}

func (s *server) AddSchedule(ctx context.Context, req *pb.RequestSchedule) (*pb.ResponseScheduleID, error) {
    schedule := model.Schedule{
        UserID:       int(req.UserId),
        CureName:     req.CureName,
        DosesPerDay:  int(req.DosesPerDay),
        DurationDays: int(req.DurationDays),
    }
    id, err := s.svc.AddSchedule(schedule)
    if err != nil {
        return nil, err
    }
    return &pb.ResponseScheduleID{ScheduleId: int64(id)}, nil
}

func (s *server) GetScheduleIDs(ctx context.Context, req *pb.RequestUserID) (*pb.ResponseScheduleIDs, error) {
    ids, err := s.svc.GetScheduleIDs(int(req.UserId))
    if err != nil {
        return nil, err
    }
    out := &pb.ResponseScheduleIDs{}
    for _, id := range ids {
        out.SchdeduleIds = append(out.SchdeduleIds, int64(id))
    }
    return out, nil
}

func (s *server) GetSchedule(ctx context.Context, req *pb.RequestUserIDScheduleID) (*pb.ResponseSchedule, error) {
    sched, err := s.svc.GetSchedule(int(req.UserId), int(req.ScheduleId))
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

func (s *server) GetNextTakings(ctx context.Context, req *pb.RequestNextTakings) (*pb.ResponseNextTakings, error) {
    schedules, err := s.svc.GetNextTakings(int(req.UserId))
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
