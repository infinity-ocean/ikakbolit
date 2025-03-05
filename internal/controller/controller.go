package controller

import (
	"encoding/json"
	"errors"
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
	AddSchedule(model.ScheduleRequest) (int, error) 
}

func New(svc service, port string) *controller {
	return &controller{service: svc, listenPort: port}
}

func (c *controller) Run() error {
	router := mux.NewRouter()
	router.HandleFunc("/schedule", httpWrapper(c.addSchedule))
	fmt.Println("Starting server on ", c.listenPort)
	if err := http.ListenAndServe(c.listenPort, router); err != nil {
		return err
	}
	return nil
}

func (c *controller) addSchedule(w http.ResponseWriter, r *http.Request) error {
	if r.Method != "POST" {
		return errors.New("method for schedule creation isn't POST")
	}
	schedule := &model.ScheduleRequest{}
	if err := json.NewDecoder(r.Body).Decode(schedule); err != nil {
		return err
	}

	scheduleID, err := c.service.AddSchedule(*schedule)
	if err != nil {
		return err
	}
	response := make(map[string]string)
	response["schedule_id"] = strconv.Itoa(scheduleID)
	return writeJSONtoHTTP(w, http.StatusOK, response)
}

