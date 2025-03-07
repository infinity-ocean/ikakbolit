package controller

import (
	"encoding/json"

	"fmt"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
	"github.com/infinity-ocean/ikakbolit/internal/model"
)

type controller struct {
	service service
	listenPort string
}

/*
POST /schedule [JSON]
GET /schedules?user_id=
GET /schedule?user_id=&schedule_id=
GET /next_takings?user_id=       [Размер периода за который получаем данные определяем в конфигурации] 
*/

type service interface {
	AddSchedule(model.Schedule) (int, error) 
	GetSchedules(int) ([]int, error)
	GetSchedule(int, int) (model.Schedule, []string, error)
}

func New(svc service, port string) *controller {
	return &controller{service: svc, listenPort: port}
}

func (c *controller) Run() error {
	router := mux.NewRouter()
	router.HandleFunc("/schedule", httpWrapper(c.addSchedule)).Methods("POST")
	router.HandleFunc("/schedules", httpWrapper(c.getSchedules)).Methods("GET")
	router.HandleFunc("/schedule", httpWrapper(c.getSchedule)).Methods("GET")
	fmt.Println("Starting server on ", c.listenPort)
	if err := http.ListenAndServe(c.listenPort, router); err != nil {
		return err
	}
	return nil
}

func (c *controller) addSchedule(w http.ResponseWriter, r *http.Request) error {
	schedule := &ScheduleRequest{}
	if err := json.NewDecoder(r.Body).Decode(schedule); err != nil {
		return err
	}
	scheduleModel := toModelSchedule(*schedule)
	scheduleID, err := c.service.AddSchedule(scheduleModel)
	if err != nil {
		return err
	}

	response := Response{}
	response.Schedule_id = strconv.Itoa(scheduleID)
	return writeJSONtoHTTP(w, http.StatusOK, response)
}

func (c *controller) getSchedules(w http.ResponseWriter, r *http.Request) error {
	userID, err := strconv.Atoi(r.URL.Query().Get("user_id"))
	if err != nil {
		return err
	}
	response, err := c.service.GetSchedules(userID)
	if err != nil {
		return err
	}
	return writeJSONtoHTTP(w, http.StatusOK, response)
}

func (c *controller) getSchedule(w http.ResponseWriter, r *http.Request) error {
	userID, err := strconv.Atoi(r.URL.Query().Get("user_id"))
	if err != nil {
		return err
	}
	scheduleID, err := strconv.Atoi(r.URL.Query().Get("schedule_id"))
	if err != nil {
		return err
	}
	schedule, intakes, err := c.service.GetSchedule(userID, scheduleID)
	if err != nil {
		return err
	}
	response := toScheduleWithIntakes(schedule, intakes)
	return writeJSONtoHTTP(w, http.StatusOK, response)
}