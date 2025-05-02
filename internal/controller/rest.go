package controller

import (
	"encoding/json"
	"fmt"
	"log"
	"log/slog"

	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"
	"github.com/infinity-ocean/ikakbolit/internal/model"

	_ "github.com/infinity-ocean/ikakbolit/3-api-grpc-Homework/docs"
	swagger "github.com/swaggo/http-swagger"
)

type restServer struct {
	service    service
	listenPort string
	logger     *slog.Logger
}

type service interface {
	AddSchedule(model.Schedule) (int, error)
	GetScheduleIDs(int) ([]int, error)
	GetScheduleWithIntake(int, int) (model.Schedule, error)
	GetNextTakings(int) ([]model.Schedule, error)
}

func NewRestServer(svc service, port string) *restServer {
	return &restServer{service: svc, listenPort: port}
}

func (c *restServer) Run() error {
	router := chi.NewRouter()

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
		return fmt.Errorf("incorrect input: %w.\ninternal error: %s", model.ErrBadRequest, err)
	}

	scheduleModel := toModelSchedule(schedule)
	scheduleID, err := c.service.AddSchedule(scheduleModel)
	if err != nil {
		return err
	}

	response := responseScheduleID{Schedule_id: strconv.Itoa(scheduleID)}
	return writeJSONtoHTTP(w, http.StatusCreated, response)
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
		return writeJSONtoHTTP(w, http.StatusBadRequest, fmt.Errorf("incorrect user_id: %w", err))
	}
	if userID == 0 {
		return writeJSONtoHTTP(w, http.StatusBadRequest, fmt.Errorf("user_id can't be 0: %w", err))
	}

	response, err := c.service.GetScheduleIDs(userID)
	if err != nil {
		return err
	}

	if len(response) == 0 {
		return writeJSONtoHTTP(w, http.StatusNoContent, response)
	}

	return writeJSONtoHTTP(w, http.StatusOK, response)
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
		return writeJSONtoHTTP(w, http.StatusBadRequest, fmt.Errorf("incorrect user_id: %w", err))
	}

	scheduleID, err := strconv.Atoi(r.URL.Query().Get("schedule_id"))
	if err != nil {
		return writeJSONtoHTTP(w, http.StatusBadRequest, fmt.Errorf("incorrect schedule_id: %w", err))
	}

	response, err := c.service.GetScheduleWithIntake(userID, scheduleID)
	if err != nil {
		return err
	}

	if response.ID == 0 {
		return writeJSONtoHTTP(w, http.StatusNoContent, Schedule(response))
	}

	return writeJSONtoHTTP(w, http.StatusOK, Schedule(response))
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
		return writeJSONtoHTTP(w, http.StatusBadRequest, fmt.Errorf("incorrect user_id: %w", err))
	}

	schedules, err := c.service.GetNextTakings(userID)
	if err != nil {
		return err
	}

	if len(schedules) == 0 {
		return writeJSONtoHTTP(w, http.StatusNoContent, schedules)
	}

	response := SchedulesInWindow{Schedules: fromModelSchedule(schedules)}
	return writeJSONtoHTTP(w, http.StatusOK, response)
}

func fromModelSchedule(schedules []model.Schedule) []Schedule {
	schedulesResp := make([]Schedule, 0, len(schedules))
	for _, s := range schedules {
		schedulesResp = append(schedulesResp, Schedule(s))
	}

	return schedulesResp
}
