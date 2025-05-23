package rest

import (
	"context"
	"encoding/json"
	"fmt"
	"log/slog"
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	_ "github.com/infinity-ocean/ikakbolit/3-api-grpc-Homework/docs"
	"github.com/infinity-ocean/ikakbolit/internal/domain/entity"
	"github.com/infinity-ocean/ikakbolit/pkg/errcodes"
	"github.com/infinity-ocean/ikakbolit/pkg/httpx/reply"
	mdw "github.com/infinity-ocean/ikakbolit/pkg/middlewarex/middleware"
	"github.com/infinity-ocean/ikakbolit/pkg/rest"
	swagger "github.com/swaggo/http-swagger"
)

type HTTPRouter struct {
	service    service
	ListenPort string
	log        *slog.Logger
}

type service interface {
	AddSchedule(context.Context, entity.Schedule) (int, error)
	GetScheduleIDs(context.Context, int) ([]int, error)
	GetScheduleWithIntake(context.Context, int, int) (entity.Schedule, error)
	GetNextTakings(context.Context, int) ([]entity.Schedule, error)
}

func NewHTTPRouter(svc service, port string, log *slog.Logger) *HTTPRouter {
	portStr := ":" + port
	return &HTTPRouter{service: svc, ListenPort: portStr, log: log}
}

func (c *HTTPRouter) GetRouter() *chi.Mux {
	router := chi.NewRouter()

	router.Use(middleware.RequestID)
	router.Use(mdw.LoggerMiddleware(c.log))
	router.Use(middleware.Recoverer)
	router.Use(middleware.URLFormat)

	router.Mount("/swagger", swagger.Handler(
		swagger.URL("http://localhost:8080/swagger/doc.json"),
		swagger.DeepLinking(true),
		swagger.DocExpansion("none"),
		swagger.DomID("swagger-ui"),
	))

	router.Method("POST", "/schedule", reply.ErrorDecorator(c.addSchedule))
	router.Method("GET", "/schedules", reply.ErrorDecorator(c.getScheduleIDs))
	router.Method("GET", "/schedule", reply.ErrorDecorator(c.getSchedule))
	router.Method("GET", "/next_takings", reply.ErrorDecorator(c.getNextTakings))

	return router
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
func (c *HTTPRouter) addSchedule(w http.ResponseWriter, r *http.Request) error {
	schedule := rest.Schedule{}
	if err := json.NewDecoder(r.Body).Decode(&schedule); err != nil {
		return fmt.Errorf("failed to parse schedule into json: %w", err)
	}

	scheduleModel := rest.ToModelSchedule(schedule)
	scheduleID, err := c.service.AddSchedule(r.Context(), scheduleModel)
	if err != nil {
		return err
	}

	response := rest.ResponseScheduleID{ScheduleID: strconv.Itoa(scheduleID)}
	return reply.JSON(w, http.StatusCreated, response)
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
func (c *HTTPRouter) getScheduleIDs(w http.ResponseWriter, r *http.Request) error {
	userID, err := strconv.Atoi(r.URL.Query().Get("user_id"))
	if err != nil {
		return fmt.Errorf("invalid user_id: %v; %w", err, errcodes.ErrBadRequest)
	}
	if userID == 0 {
		return fmt.Errorf("user_id can't be 0: %w", errcodes.ErrBadRequest)
	}

	response, err := c.service.GetScheduleIDs(r.Context(), userID)
	if err != nil {
		return err
	}

	if len(response) == 0 {
		return reply.JSON(w, http.StatusNoContent, response)
	}

	return reply.JSON(w, http.StatusOK, response)
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
func (c *HTTPRouter) getSchedule(w http.ResponseWriter, r *http.Request) error {
	userID, err := strconv.Atoi(r.URL.Query().Get("user_id"))
	if err != nil {
		return fmt.Errorf("invalid user_id: %v; %w", err, errcodes.ErrBadRequest)
	}

	scheduleID, err := strconv.Atoi(r.URL.Query().Get("schedule_id"))
	if err != nil {
		return fmt.Errorf("invalid schedule_id: %v; %w", err, errcodes.ErrBadRequest)
	}

	response, err := c.service.GetScheduleWithIntake(r.Context(), userID, scheduleID)
	if err != nil {
		return err
	}

	return reply.JSON(w, http.StatusOK, rest.Schedule(response))
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
func (c *HTTPRouter) getNextTakings(w http.ResponseWriter, r *http.Request) error {
	userID, err := strconv.Atoi(r.URL.Query().Get("user_id"))
	if err != nil {
		return fmt.Errorf("invalid user_id: %v; %w", err, errcodes.ErrBadRequest)
	}

	schedules, err := c.service.GetNextTakings(r.Context(), userID)
	if err != nil {
		return err
	}

	if len(schedules) == 0 {
		return reply.JSON(w, http.StatusNoContent, schedules)
	}

	response := rest.SchedulesInWindow{Schedules: fromModelSchedule(schedules)}
	return reply.JSON(w, http.StatusOK, response)
}

func fromModelSchedule(schedules []entity.Schedule) []rest.Schedule {
	schedulesResp := make([]rest.Schedule, 0, len(schedules))
	for _, s := range schedules {
		schedulesResp = append(schedulesResp, rest.Schedule(s))
	}

	return schedulesResp
}
