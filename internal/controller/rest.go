package controller

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"log/slog"
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	_ "github.com/infinity-ocean/ikakbolit/3-api-grpc-Homework/docs"
	"github.com/infinity-ocean/ikakbolit/internal/dto"
	"github.com/infinity-ocean/ikakbolit/internal/model"
	swagger "github.com/swaggo/http-swagger"
)

type restServer struct {
	service    service
	listenPort string
	log     *slog.Logger
}

type service interface {
	AddSchedule(context.Context, model.Schedule) (int, error)
	GetScheduleIDs(context.Context, int) ([]int, error)
	GetScheduleWithIntake(context.Context, int, int) (model.Schedule, error)
	GetNextTakings(context.Context, int) ([]model.Schedule, error)
}

func NewRestServer(svc service, port string, log *slog.Logger) *restServer {
	return &restServer{service: svc, listenPort: port, log: log}
}

func (c *restServer) Run() error {
	router := chi.NewRouter()

	router.Use(middleware.RequestID)
	router.Use(loggerMiddleware(c.log))
	router.Use(middleware.Recoverer)
	router.Use(middleware.URLFormat)

	router.Mount("/swagger", swagger.Handler(
		swagger.URL("http://localhost:8080/swagger/doc.json"),
		swagger.DeepLinking(true),
		swagger.DocExpansion("none"),
		swagger.DomID("swagger-ui"),
	))

	router.Method("POST", "/schedule", httpWrapper(c.addSchedule))
	router.Method("GET", "/schedules", httpWrapper(c.getScheduleIDs))
	router.Method("GET", "/schedule", httpWrapper(c.getSchedule))
	router.Method("GET", "/next_takings", httpWrapper(c.getNextTakings))

	log.Println("Starting server on", c.listenPort)
	return http.ListenAndServe(c.listenPort, router)
}

// @Summary Add a new schedule [! Предпочтительнее использовать http-клиент для post-запросов, например Postman]
// @Description Create a new schedule for a user
// @Accept  json
// @Produce  json
// @Param   schedule  body  Schedule  true  "Schedule data"
// @Success 201 {object} responseScheduleID
// @Failure 400 {object} APIError "Bad request"
// @Failure 404 {object} APIError "Resource not found"
// @Failure 500 {object} APIError "Internal server error"
// @Router /schedule [post]
func (c *restServer) addSchedule(w http.ResponseWriter, r *http.Request) error {
	schedule := Schedule{}
	if err := json.NewDecoder(r.Body).Decode(&schedule); err != nil {
		return fmt.Errorf("failed to parse schedule into json: %w", err)
	}

	scheduleModel := toModelSchedule(schedule)
	scheduleID, err := c.service.AddSchedule(r.Context(), scheduleModel)
	if err != nil {
		return err
	}

	response := responseScheduleID{Schedule_id: strconv.Itoa(scheduleID)}
	return sendJSONtoHTTP(w, http.StatusCreated, response)
}

// @Summary Get user schedules
// @Description Retrieve schedule IDs for a given user
// @Produce json
// @Param   user_id query int true "User ID"
// @Success 200 {array} int
// @Success 204 "No content"
// @Failure 400 {object} APIError "Bad request"
// @Failure 404 {object} APIError "Resource not found"
// @Failure 500 {object} APIError "Internal server error"
// @Router /schedules [get]
func (c *restServer) getScheduleIDs(w http.ResponseWriter, r *http.Request) error {
	userID, err := strconv.Atoi(r.URL.Query().Get("user_id"))
	if err != nil {
		return fmt.Errorf("invalid user_id: %v; %w", err, dto.ErrBadRequest)
	}
	if userID == 0 {
		return fmt.Errorf("user_id can't be 0: %w", dto.ErrBadRequest)
	}

	response, err := c.service.GetScheduleIDs(r.Context(), userID)
	if err != nil {
		return err
	}

	if len(response) == 0 {
		return sendJSONtoHTTP(w, http.StatusNoContent, response)
	}

	return sendJSONtoHTTP(w, http.StatusOK, response)
}

// @Summary Get a specific schedule
// @Description Retrieve a schedule by user ID and schedule ID
// @Produce json
// @Param   user_id query int true "User ID"
// @Param   schedule_id query int true "Schedule ID"
// @Success 200 {object} Schedule
// @Success 204 "No content"
// @Failure 400 {object} APIError "Bad request"
// @Failure 404 {object} APIError "Resource not found"
// @Failure 500 {object} APIError "Internal server error"
// @Router /schedule [get]
func (c *restServer) getSchedule(w http.ResponseWriter, r *http.Request) error {
	userID, err := strconv.Atoi(r.URL.Query().Get("user_id"))
	if err != nil {
		return fmt.Errorf("invalid user_id: %v; %w", err, dto.ErrBadRequest)
	}

	scheduleID, err := strconv.Atoi(r.URL.Query().Get("schedule_id"))
	if err != nil {
		return fmt.Errorf("invalid schedule_id: %v; %w", err, dto.ErrBadRequest)
	}

	response, err := c.service.GetScheduleWithIntake(r.Context(), userID, scheduleID)
	if err != nil {
		return err
	}

	return sendJSONtoHTTP(w, http.StatusOK, Schedule(response))
}

// @Summary Get next scheduled takings
// @Description Retrieve upcoming medication schedulesResp for a user
// @Produce json
// @Param   user_id query int true "User ID"
// @Success 200 {object} Schedule
// @Success 204 "No content"
// @Failure 400 {object} APIError "Bad request"
// @Failure 404 {object} APIError "Resource not found"
// @Failure 500 {object} APIError "Internal server error"
// @Router /next_takings [get]
func (c *restServer) getNextTakings(w http.ResponseWriter, r *http.Request) error {
	userID, err := strconv.Atoi(r.URL.Query().Get("user_id"))
	if err != nil {
		return fmt.Errorf("invalid user_id: %v; %w", err, dto.ErrBadRequest)
	}

	schedules, err := c.service.GetNextTakings(r.Context(), userID)
	if err != nil {
		return err
	}

	if len(schedules) == 0 {
		return sendJSONtoHTTP(w, http.StatusNoContent, schedules)
	}

	response := SchedulesInWindow{Schedules: fromModelSchedule(schedules)}
	return sendJSONtoHTTP(w, http.StatusOK, response)
}

func fromModelSchedule(schedules []model.Schedule) []Schedule {
	schedulesResp := make([]Schedule, 0, len(schedules))
	for _, s := range schedules {
		schedulesResp = append(schedulesResp, Schedule(s))
	}

	return schedulesResp
}
