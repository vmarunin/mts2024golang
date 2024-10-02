package pinger

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log/slog"
	"math/rand"
	"net/http"
	"strings"
	"time"

	"github.com/vmarunin/mts2024golang/seminar2/internal/pkg/config"
	"github.com/vmarunin/mts2024golang/seminar2/internal/pkg/models"
)

var ErrTaskFailed = errors.New("task failed")

type Service struct {
	log *slog.Logger
	cfg *config.Config
}

func New(cfg *config.Config, log *slog.Logger) *Service {
	return &Service{
		cfg: cfg,
		log: log,
	}
}

func (s *Service) AddTask(ctx context.Context, url string) (string, error) {
	http.Get(url)
	res, err := http.Get(url + "/api/v1/simple/name")
	if err != nil {
		return "", err
	}
	body, err := io.ReadAll(res.Body)
	if err != nil {
		return "", err
	}
	res.Body.Close()
	if res.StatusCode > 299 {
		return "", fmt.Errorf("%w response failed with status code: %d and\nbody: %s\n", ErrTaskFailed, res.StatusCode, body)
	}
	data := models.ApiNameResponse{}
	err = json.Unmarshal(body, &data)
	if err != nil {
		return "", fmt.Errorf("%w response failed %w", ErrTaskFailed, err)
	}

	s.log.Info("SimpleSolved", slog.String("student", data.Name))

	token := genToken()

	go s.processTask(token, url, data.Name)

	return token, nil
}

func genToken() string {
	var sb strings.Builder
	for i := 0; i < 16; i++ {
		sb.WriteByte(byte('a' + rand.Intn(26)))
	}
	return sb.String()
}

type tHTTP struct {
	fn   func(ctx context.Context, token, url string) error
	name string
}

func (s *Service) processTask(token, url, studentName string) {
	time.Sleep(s.cfg.GetWaitBeforeProcessTask())
	ctx, cancel := context.WithTimeout(context.Background(), s.cfg.GetTaskTimeout())
	defer cancel()

	tests := []tHTTP{
		{
			fn:   s.tAuthSuccess,
			name: "AuthSuccess",
		},
		{
			fn:   s.tAuthFailNoToken,
			name: "AuthFailNoPasswd",
		},
		{
			fn:   s.tAuthFailBadToken,
			name: "AuthFailNoPasswd",
		},
		{
			fn:   s.tDeleteSimple,
			name: "DeleteSi",
		},
		{
			fn:   s.tAuth404,
			name: "Auth404",
		},
	}

	rand.Shuffle(len(tests), func(i, j int) { tests[i], tests[j] = tests[j], tests[i] })

	isSuccess := true
	for _, rec := range tests {
		err := rec.fn(ctx, token, url)
		if err != nil {
			s.log.Info(rec.name+" failed", slog.String("student", studentName), slog.Any("err", err))
			isSuccess = false
		}
	}

	if isSuccess {
		s.log.Info("AdvancedSolved", slog.String("student", studentName))
	}
}

func (s *Service) tDeleteSimple(ctx context.Context, token, url string) error {
	req, err := http.NewRequestWithContext(ctx, "DELETE", url+"/api/v1/simple/name", nil)
	if err != nil {
		return err
	}
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return err
	}
	// body, err := io.ReadAll(resp.Body)
	// if err != nil {
	// 	return err
	// }
	resp.Body.Close()
	if resp.StatusCode != 405 {
		return fmt.Errorf("Expected 405 but got %d", resp.StatusCode)
	}
	return nil
}

func (s *Service) tAuthSuccess(ctx context.Context, token, url string) error {
	req, err := http.NewRequestWithContext(ctx, "GET", url+"/api/v1/auth/name", nil)
	req.Header.Set("Authorization", "Bearer "+token)
	if err != nil {
		return err
	}
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	if resp.StatusCode != 200 {
		return fmt.Errorf("Expected 200 but got %d", resp.StatusCode)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return err
	}

	data := models.ApiNameResponse{}
	err = json.Unmarshal(body, &data)
	if err != nil {
		return fmt.Errorf("%w can't parse %s", err, body)
	}

	return nil
}

func (s *Service) tAuthFailNoToken(ctx context.Context, token, url string) error {
	req, err := http.NewRequestWithContext(ctx, "GET", url+"/api/v1/auth/name", nil)
	if err != nil {
		return err
	}
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	if resp.StatusCode != 401 {
		return fmt.Errorf("Expected 401 but got %d", resp.StatusCode)
	}

	return nil
}

func (s *Service) tAuthFailBadToken(ctx context.Context, token, url string) error {
	req, err := http.NewRequestWithContext(ctx, "GET", url+"/api/v1/auth/name", nil)
	req.Header.Set("Authorization", "Bearer "+token[8:]+token[:8])
	if err != nil {
		return err
	}
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	if resp.StatusCode != 401 {
		return fmt.Errorf("Expected 401 but got %d", resp.StatusCode)
	}

	return nil
}

func (s *Service) tAuth404(ctx context.Context, token, url string) error {
	req, err := http.NewRequestWithContext(ctx, "GET", url+"/api/v1/auth/noname", nil)
	req.Header.Set("Authorization", "Bearer "+token)
	if err != nil {
		return err
	}
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	if resp.StatusCode != 404 {
		return fmt.Errorf("Expected 404 but got %d", resp.StatusCode)
	}

	return nil
}
