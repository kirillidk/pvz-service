package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/kirillidk/pvz-service/internal/app"
	"github.com/kirillidk/pvz-service/internal/config"
	"github.com/kirillidk/pvz-service/internal/dto"
	"github.com/kirillidk/pvz-service/internal/model"
)

func init() {
	gin.SetMode(gin.ReleaseMode)
}

func setupTestEnv(t *testing.T) {
	t.Setenv("SERVER_PORT", "8080")
	t.Setenv("DB_HOST", "localhost")
	t.Setenv("DB_PORT", "5432")
	t.Setenv("DB_USER", "postgres")
	t.Setenv("DB_PASSWORD", "postgres")
	t.Setenv("DB_NAME", "pvz_service")
	t.Setenv("DB_SSLMODE", "disable")
	t.Setenv("JWT_SECRET", "secret-key")
}

func TestPVZFlow(t *testing.T) {
	setupTestEnv(t)

	cfg := config.NewConfig()

	app, err := app.NewApp(cfg)
	if err != nil {
		t.Fatalf("Failed to initialize app: %v", err)
	}

	moderatorToken := getModeratorToken(t, app)

	pvz := createPVZ(t, app, moderatorToken)
	log.Printf("Created PVZ with ID: %s", pvz.ID)

	employeeToken := getEmployeeToken(t, app)

	reception := createReception(t, app, employeeToken, pvz.ID)
	log.Printf("Created reception with ID: %s", reception.ID)

	productTypes := []string{"электроника", "одежда", "обувь"}
	for i := 0; i < 50; i++ {
		productType := productTypes[i%len(productTypes)]
		product := createProduct(t, app, employeeToken, productType, pvz.ID)
		log.Printf("Added product %d with ID: %s, Type: %s", i+1, product.ID, product.Type)
	}

	closedReception := closeReception(t, app, employeeToken, pvz.ID)
	log.Printf("Closed reception with ID: %s, Status: %s", closedReception.ID, closedReception.Status)

	log.Println("Integration test completed successfully!")
}

func getModeratorToken(t *testing.T, app *app.App) string {
	return getDummyToken(t, app, model.ModeratorRole)
}

func getEmployeeToken(t *testing.T, app *app.App) string {
	return getDummyToken(t, app, model.EmployeeRole)
}

func getDummyToken(t *testing.T, app *app.App, role model.UserRole) string {
	reqBody := dto.DummyLoginRequest{
		Role: role,
	}
	jsonBody, _ := json.Marshal(reqBody)

	req := httptest.NewRequest("POST", "/dummyLogin", bytes.NewBuffer(jsonBody))
	req.Header.Set("Content-Type", "application/json")

	resp := httptest.NewRecorder()
	app.Router.ServeHTTP(resp, req)

	if resp.Code != http.StatusOK {
		t.Fatalf("Failed to get token for role %s: %d - %s", role, resp.Code, resp.Body.String())
	}

	var tokenResp model.Token
	err := json.Unmarshal(resp.Body.Bytes(), &tokenResp)
	if err != nil {
		t.Fatalf("Failed to unmarshal token response: %v", err)
	}

	return tokenResp.Value
}

func createPVZ(t *testing.T, app *app.App, token string) *model.PVZ {
	reqBody := dto.PVZCreateRequest{
		RegistrationDate: time.Now(),
		City:             "Москва",
	}
	jsonBody, _ := json.Marshal(reqBody)

	req := httptest.NewRequest("POST", "/pvz", bytes.NewBuffer(jsonBody))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", token))

	resp := httptest.NewRecorder()
	app.Router.ServeHTTP(resp, req)

	if resp.Code != http.StatusCreated {
		t.Fatalf("Failed to create PVZ: %d - %s", resp.Code, resp.Body.String())
	}

	var pvz model.PVZ
	err := json.Unmarshal(resp.Body.Bytes(), &pvz)
	if err != nil {
		t.Fatalf("Failed to unmarshal PVZ response: %v", err)
	}

	return &pvz
}

func createReception(t *testing.T, app *app.App, token string, pvzID string) *model.Reception {
	reqBody := dto.ReceptionCreateRequest{
		PVZID: pvzID,
	}
	jsonBody, _ := json.Marshal(reqBody)

	req := httptest.NewRequest("POST", "/receptions", bytes.NewBuffer(jsonBody))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", token))

	resp := httptest.NewRecorder()
	app.Router.ServeHTTP(resp, req)

	if resp.Code != http.StatusCreated {
		t.Fatalf("Failed to create reception: %d - %s", resp.Code, resp.Body.String())
	}

	var reception model.Reception
	err := json.Unmarshal(resp.Body.Bytes(), &reception)
	if err != nil {
		t.Fatalf("Failed to unmarshal reception response: %v", err)
	}

	return &reception
}

func createProduct(t *testing.T, app *app.App, token string, productType string, pvzID string) *model.Product {
	reqBody := dto.ProductCreateRequest{
		Type:  productType,
		PVZID: pvzID,
	}
	jsonBody, _ := json.Marshal(reqBody)

	req := httptest.NewRequest("POST", "/products", bytes.NewBuffer(jsonBody))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", token))

	resp := httptest.NewRecorder()
	app.Router.ServeHTTP(resp, req)

	if resp.Code != http.StatusCreated {
		t.Fatalf("Failed to create product: %d - %s", resp.Code, resp.Body.String())
	}

	var product model.Product
	err := json.Unmarshal(resp.Body.Bytes(), &product)
	if err != nil {
		t.Fatalf("Failed to unmarshal product response: %v", err)
	}

	return &product
}

func closeReception(t *testing.T, app *app.App, token string, pvzID string) *model.Reception {
	req := httptest.NewRequest("POST", fmt.Sprintf("/pvz/%s/close_last_reception", pvzID), nil)
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", token))

	resp := httptest.NewRecorder()
	app.Router.ServeHTTP(resp, req)

	if resp.Code != http.StatusOK {
		t.Fatalf("Failed to close reception: %d - %s", resp.Code, resp.Body.String())
	}

	var reception model.Reception
	err := json.Unmarshal(resp.Body.Bytes(), &reception)
	if err != nil {
		t.Fatalf("Failed to unmarshal reception response: %v", err)
	}

	return &reception
}
