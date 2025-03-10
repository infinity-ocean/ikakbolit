package controller

import (
	"encoding/json"
	"errors"
	"log"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
	"github.com/infinity-ocean/ikakbolit/internal/model"

	_ "github.com/infinity-ocean/ikakbolit/internal/docs"
	swagger "github.com/swaggo/http-swagger"
)

type controller struct {
	service    service
	listenPort string
}

type service interface {
	AddSchedule(model.Schedule) (int, error)
	GetScheduleIDs(int) ([]int, error)
	GetScheduleWithIntake(int, int) (model.Schedule, error)
	GetNextTakings(int) ([]model.Schedule, error)
}

var (
    ErrInvalidUserID      = errors.New("user_id cannot be zero")
    ErrScheduleNotFound   = errors.New("schedule not found")
)

func New(svc service, port string) *controller {
	return &controller{service: svc, listenPort: port}
}

func (c *controller) Run() error {
	
	router := mux.NewRouter()
	router.PathPrefix("/swagger/").Handler(swagger.Handler(
		swagger.URL("http://localhost:8080/swagger/doc.json"),
		swagger.DeepLinking(true),
		swagger.DocExpansion("none"),
		swagger.DomID("swagger-ui"),
	)).Methods(http.MethodGet)
	router.Handle("/swagger/", http.StripPrefix("/swagger/", swagger.WrapHandler))
	router.HandleFunc("/schedule", httpWrapper(c.addSchedule)).Methods("POST")
	router.HandleFunc("/schedules", httpWrapper(c.getSchedules)).Methods("GET")
	router.HandleFunc("/schedule", httpWrapper(c.getSchedule)).Methods("GET")
	router.HandleFunc("/next_takings", httpWrapper(c.getNextTakings)).Methods("GET")

	log.Println("Starting server on ", c.listenPort)
	return http.ListenAndServe(c.listenPort, router)
}

// @Summary Add a new schedule
// @Description Create a new schedule for a user
// @Accept  json
// @Produce  json
// @Param   schedule  body  Schedule  true  "Schedule data"
// @Success 200 {object} Response
// @Failure 400 {object} APIError
// @Router /schedule [post]
func (c *controller) addSchedule(w http.ResponseWriter, r *http.Request) error {
	schedule := &Schedule{}
	if err := json.NewDecoder(r.Body).Decode(schedule); err != nil {
		return err
	}

	scheduleModel := toModelSchedule(*schedule)
	scheduleID, err := c.service.AddSchedule(scheduleModel)
	if err != nil {
		return err
	}

	response := Response{Schedule_id: strconv.Itoa(scheduleID)}
	return writeJSONtoHTTP(w, http.StatusOK, response)
}

// @Summary Get user schedules
// @Description Retrieve schedule IDs for a given user
// @Produce json
// @Param   user_id query int true "User ID"
// @Success 200 {array} int
// @Failure 400 {object} APIError
// @Router /schedules [get]
func (c *controller) getSchedules(w http.ResponseWriter, r *http.Request) error {
	userID, err := strconv.Atoi(r.URL.Query().Get("user_id"))
	if err != nil {
		return err
	}
	if userID == 0 {
		return ErrInvalidUserID
	}

	response, err := c.service.GetScheduleIDs(userID)
	if err != nil {
		return err
	}
	return writeJSONtoHTTP(w, http.StatusOK, response)
}

// @Summary Get a specific schedule
// @Description Retrieve a schedule by user ID and schedule ID
// @Produce json
// @Param   user_id query int true "User ID"
// @Param   schedule_id query int true "Schedule ID"
// @Success 200 {object} model.Schedule
// @Failure 404 {object} APIError
// @Router /schedule [get]
func (c *controller) getSchedule(w http.ResponseWriter, r *http.Request) error {
	userID, err := strconv.Atoi(r.URL.Query().Get("user_id"))
	if err != nil {
		return err
	}
	scheduleID, err := strconv.Atoi(r.URL.Query().Get("schedule_id"))
	if err != nil {
		return err
	}

	response, err := c.service.GetScheduleWithIntake(userID, scheduleID)
	if err != nil {
		if errors.Is(err, ErrScheduleNotFound) {
			return writeJSONtoHTTP(w, http.StatusNotFound, APIError{Err: err.Error()})
		}
		return err
	}
	return writeJSONtoHTTP(w, http.StatusOK, response)
}

// @Summary Get next scheduled takings
// @Description Retrieve upcoming medication schedules for a user
// @Produce json
// @Param   user_id query int true "User ID"
// @Success 200 {object} SchedulesInWindow
// @Failure 400 {object} APIError
// @Router /next_takings [get]
func (c *controller) getNextTakings(w http.ResponseWriter, r *http.Request) error {
	userID, err := strconv.Atoi(r.URL.Query().Get("user_id"))
	if err != nil {
		return err
	}

	schedulesRaw, err := c.service.GetNextTakings(userID)
	if err != nil {
		return err
	}

	schedules := make([]Schedule, 0, len(schedulesRaw))
	for _, s := range schedulesRaw {
		schedules = append(schedules, Schedule(s))
	}

	response := SchedulesInWindow{Schedules: schedules}
	return writeJSONtoHTTP(w, http.StatusOK, response)
}
