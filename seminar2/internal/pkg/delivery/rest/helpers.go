package rest

import (
	"encoding/json"
	"log/slog"
	"net/http"
)

func (s *Server) responseJSON(
	writer http.ResponseWriter,
	data any,
	status int, //nolint:unparam // status may be not 200 in future
) {
	respData, err := json.Marshal(data)
	if err != nil {
		s.responseError(writer, err, http.StatusInternalServerError)
	}
	writer.Header().Set("Content-Type", "application/json")
	writer.WriteHeader(status)
	_, err = writer.Write(respData)
	if err != nil {
		s.log.Error("response write error", slog.Int("status", status), slog.String("error", err.Error()))
	}
}

func (s *Server) responseError(
	writer http.ResponseWriter,
	err error,
	status int,
) {
	s.log.Error("request error", slog.Int("status", status), slog.String("error", err.Error()))
	writer.WriteHeader(status)
}

func (s *Server) responseErrorNoLog(
	writer http.ResponseWriter,
	_ error,
	status int,
) {
	// s.log.Error("request error", slog.Int("status", status), slog.String("error", err.Error()))
	writer.WriteHeader(status)
}

func (s *Server) parseData(r *http.Request, data any) error {
	return json.NewDecoder(r.Body).Decode(data)
}
