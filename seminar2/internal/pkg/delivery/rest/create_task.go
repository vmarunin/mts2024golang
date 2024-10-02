package rest

import (
	"errors"
	"io"
	"net/http"
	"os"

	"github.com/vmarunin/mts2024golang/seminar2/internal/pkg/models"
	"github.com/vmarunin/mts2024golang/seminar2/internal/pkg/service/pinger"
)

func (s *Server) CreateTask(writer http.ResponseWriter, request *http.Request) {
	// POST /api/v1/task
	var err error

	data := &models.ApiCreateTask{}
	err = s.parseData(request, data)
	if err != nil {
		s.responseError(writer, err, http.StatusBadRequest)
		return
	}

	token, err := s.pinger.AddTask(request.Context(), data.URL)

	if errors.Is(err, pinger.ErrTaskFailed) {
		response := models.ApiCreateTaskResponse{
			Token:  token,
			Status: "fail",
			Info:   err.Error(),
		}

		s.responseJSON(writer, response, http.StatusOK)
		return
	}
	if err != nil {
		s.responseError(writer, err, http.StatusInternalServerError)
		return
	}

	response := models.ApiCreateTaskResponse{
		Token:  token,
		Status: "ok",
		Info:   "Task created",
	}

	s.responseJSON(writer, response, http.StatusOK)
}

func (s *Server) GetLog(writer http.ResponseWriter, request *http.Request) {
	logFile, err := os.OpenFile(s.cfg.GetLogPath(), os.O_RDONLY, 0644)
	if err != nil {
		s.responseError(writer, err, http.StatusInternalServerError)
		return
	}
	defer logFile.Close()

	logFile.Seek(int64(io.SeekEnd), 10_000)
	// _, err = logFile.Seek(int64(io.SeekEnd), 10_000)
	// if err != nil {
	// 	s.responseError(writer, err, http.StatusInternalServerError)
	// 	return
	// }
	data, err := io.ReadAll(logFile)
	if err != nil {
		s.responseError(writer, err, http.StatusInternalServerError)
		return
	}
	writer.Header().Set("Content-Type", "text/plain")
	writer.WriteHeader(http.StatusOK)
	writer.Write(data)
}
