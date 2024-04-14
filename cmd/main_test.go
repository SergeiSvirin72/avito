package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"os"
	"strings"
	"testing"

	"avito/internal/dto"
	"avito/internal/model"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/suite"
)

type MainTestSuite struct {
	suite.Suite
	server      *httptest.Server
	db          *sql.DB
	banner      *model.Banner
	accessToken string
	content     string
	tagID       int64
	featureID   int64
}

func (s *MainTestSuite) SetupSuite() {
	r := gin.Default()
	setup(r)
	server := httptest.NewServer(r)
	s.server = server

	db, err := sql.Open("postgres", os.Getenv("APP_DB_DSN"))
	if err != nil {
		s.T().Fatalf("can't connect to database: %v", err)
	}

	s.db = db
	s.content = `{"test":"test"}`
	s.tagID = s.createTag()
	s.featureID = s.createFeature()
	s.createBanner(s.featureID, s.tagID, s.content)
	s.banner = s.getBanner(s.featureID)

	s.accessToken = s.auth(`{"email":"admin@test.com","password":"1234"}`)
}

func (s *MainTestSuite) TestCreateBanner() {
	featureID := s.createFeature()

	data := []struct {
		input, expected string
		featureID       int64
	}{
		{
			input: fmt.Sprintf(
				`{"tag_ids":[%d],"feature_id":%d,"content":{},"is_active":true}`,
				s.tagID,
				featureID,
			),
			featureID: featureID,
			expected:  `{"banner_id":%d}`,
		},
	}

	for i := range data {
		resp, body := s.request(http.MethodPost, "/banner", data[i].input)

		b := s.getBanner(data[i].featureID)

		s.Equal(http.StatusCreated, resp.StatusCode)
		s.Equal(fmt.Sprintf(data[i].expected, b.ID), string(body))
	}
}

func (s *MainTestSuite) TestGetUserBanner() {
	data := []struct {
		expected         string
		featureID, tagID int64
	}{
		{
			featureID: s.featureID,
			tagID:     s.tagID,
			expected:  s.content,
		},
	}

	for i := range data {
		url := fmt.Sprintf(
			"/user_banner?feature_id=%d&tag_id=%d&use_last_revision=true",
			data[i].featureID,
			data[i].tagID,
		)
		resp, body := s.request(http.MethodGet, url, "")

		s.Equal(http.StatusOK, resp.StatusCode)
		s.Equal(data[i].expected, string(body))
	}
}

func (s *MainTestSuite) TestGetBanners() {
	createdAt, err := s.banner.CreatedAt.MarshalJSON()
	if err != nil {
		s.T().Fatal(err)
	}

	updatedAt, err := s.banner.UpdatedAt.MarshalJSON()
	if err != nil {
		s.T().Fatal(err)
	}

	data := []struct {
		expected         string
		featureID, tagID int64
	}{
		{
			featureID: s.featureID,
			tagID:     s.tagID,
			expected: fmt.Sprintf(
				`[{"created_at":%s,"updated_at":%s,"content":%s,"tag_ids":[%d],"banner_id":%d,"feature_id":%d,"is_active":true}]`,
				createdAt,
				updatedAt,
				s.banner.Content,
				s.tagID,
				s.banner.ID,
				s.featureID,
			),
		},
	}

	for i := range data {
		url := fmt.Sprintf("/banner?feature_id=%d&tag_id=%d", data[i].featureID, data[i].tagID)
		resp, body := s.request(http.MethodGet, url, "")

		s.Equal(http.StatusOK, resp.StatusCode)
		s.Equal(data[i].expected, string(body))
	}
}

func (s *MainTestSuite) TestUpdateBanner() {
	featureID := s.createFeature()
	bannerID := s.createBanner(featureID, s.tagID, s.content)

	data := []struct {
		input    string
		bannerID int64
	}{
		{
			input: fmt.Sprintf(
				`{"tag_ids":[%d],"feature_id":%d,"content":%s,"is_active":false}`,
				s.tagID,
				featureID,
				s.content,
			),
			bannerID: bannerID,
		},
	}

	for i := range data {
		resp, _ := s.request(http.MethodPatch, fmt.Sprintf("/banner/%d", data[i].bannerID), data[i].input)

		s.Equal(http.StatusOK, resp.StatusCode)
	}
}

func (s *MainTestSuite) TestDeleteBanner() {
	featureID := s.createFeature()
	bannerID := s.createBanner(featureID, s.tagID, s.content)

	data := []struct {
		bannerID int64
	}{
		{
			bannerID: bannerID,
		},
	}

	for i := range data {
		resp, _ := s.request(http.MethodDelete, fmt.Sprintf("/banner/%d", data[i].bannerID), "")

		s.Equal(http.StatusNoContent, resp.StatusCode)
	}
}

func TestMainTestSuite(t *testing.T) {
	suite.Run(t, new(MainTestSuite))
}

func (s *MainTestSuite) auth(input string) string {
	req, err := http.NewRequest(http.MethodPost, s.server.URL+"/auth", strings.NewReader(input))
	if err != nil {
		s.T().Fatal(err)
	}

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		s.T().Fatal(err)
	}

	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		s.T().Fatal(err)
	}

	if ok := s.Equal(http.StatusOK, resp.StatusCode); !ok {
		s.T().Fatal(string(body))
	}

	var output dto.AuthOutput
	if err = json.Unmarshal(body, &output); err != nil {
		s.T().Fatal(err)
	}

	return output.AccessToken
}

func (s *MainTestSuite) getBanner(featureID int64) *model.Banner {
	var b model.Banner

	row := s.db.QueryRow("SELECT b.* FROM banners b WHERE b.feature_id = $1", featureID)
	if err := row.Scan(&b.ID, &b.FeatureID, &b.Content, &b.IsActive, &b.CreatedAt, &b.UpdatedAt); err != nil {
		s.T().Fatal(err)
	}

	return &b
}

func (s *MainTestSuite) createBanner(featureID, tagID int64, content string) int64 {
	var bannerID int64

	row := s.db.QueryRow(
		"INSERT INTO banners (feature_id, content, is_active) VALUES ($1, $2, true) RETURNING id",
		featureID,
		content,
	)
	if err := row.Scan(&bannerID); err != nil {
		s.T().Fatal(err)
	}

	_, err := s.db.Exec("INSERT INTO banner_tag (banner_id, tag_id) VALUES ($1, $2)", bannerID, tagID)
	if err != nil {
		s.T().Fatal(err)
	}

	return bannerID
}

func (s *MainTestSuite) createTag() int64 {
	var tagID int64

	row := s.db.QueryRow("INSERT INTO tags (id) VALUES (DEFAULT) RETURNING id")
	if err := row.Scan(&tagID); err != nil {
		s.T().Fatalf("create tag error: %v", err)
	}

	return tagID
}

func (s *MainTestSuite) createFeature() int64 {
	var featureID int64

	row := s.db.QueryRow("INSERT INTO features (id) VALUES (DEFAULT) RETURNING id")
	if err := row.Scan(&featureID); err != nil {
		s.T().Fatalf("create feature error: %v", err)
	}

	return featureID
}

func (s *MainTestSuite) request(method, url string, input string) (*http.Response, []byte) {
	var reqBody io.Reader
	if input != "" {
		reqBody = strings.NewReader(input)
	}

	req, err := http.NewRequest(method, s.server.URL+url, reqBody)
	if err != nil {
		s.T().Fatal(err)
	}

	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", s.accessToken))

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		s.T().Fatal(err)
	}

	defer resp.Body.Close()

	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		s.T().Fatal(err)
	}

	return resp, respBody
}
