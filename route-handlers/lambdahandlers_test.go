package api

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"path/filepath"
	"testing"

	"github.com/RazorSh4rk/lambdaathome/db"
	"github.com/RazorSh4rk/lambdaathome/types"
	dockerTypes "github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/image"
	"github.com/gin-gonic/gin"
)

type mockDocker struct {
	killCalls            []string
	removeContainerCalls []string
	removeImageCalls     []string
	runDetachedResult    string
	listRunningResult    []dockerTypes.Container
	listInstalledResult  []image.Summary
}

func (m *mockDocker) Kill(ID string)                                { m.killCalls = append(m.killCalls, ID) }
func (m *mockDocker) RemoveContainer(ID string)                     { m.removeContainerCalls = append(m.removeContainerCalls, ID) }
func (m *mockDocker) RemoveImage(tag string)                        { m.removeImageCalls = append(m.removeImageCalls, tag) }
func (m *mockDocker) RunDetached(lambda types.LambdaFun) string     { return m.runDetachedResult }
func (m *mockDocker) ListRunning() []dockerTypes.Container          { return m.listRunningResult }
func (m *mockDocker) ListInstalledImages() []image.Summary          { return m.listInstalledResult }
func (m *mockDocker) BuildImage(function types.LambdaFun)           {}
func (m *mockDocker) Close()                                        {}

func setupTestDB(t *testing.T) db.KV {
	t.Helper()
	return db.New(filepath.Join(t.TempDir(), "test-db"))
}

func setupRouter() *gin.Engine {
	gin.SetMode(gin.TestMode)
	return gin.New()
}

func seedFunction(t *testing.T, store db.KV, lambda types.LambdaFun) {
	t.Helper()
	record, err := json.Marshal(lambda)
	if err != nil {
		t.Fatal(err)
	}
	store.Set(lambda.ID, string(record))
}

var testLambda = types.LambdaFun{
	Name:    "test-fn",
	Tag:     "test-fn:latest",
	Runtime: "test-runtime",
	Port:    "9002",
	ID:      "abc123",
}

// findFunctionByName

func TestFindFunctionByName_Found(t *testing.T) {
	store := setupTestDB(t)
	defer store.Close()
	seedFunction(t, store, testLambda)

	key, lambda, found := findFunctionByName(store, "test-fn")
	if !found {
		t.Fatal("expected to find function")
	}
	if key != "abc123" {
		t.Errorf("expected key abc123, got %s", key)
	}
	if lambda.Name != "test-fn" {
		t.Errorf("expected name test-fn, got %s", lambda.Name)
	}
}

func TestFindFunctionByName_NotFound(t *testing.T) {
	store := setupTestDB(t)
	defer store.Close()
	seedFunction(t, store, testLambda)

	_, _, found := findFunctionByName(store, "nonexistent")
	if found {
		t.Fatal("expected not to find function")
	}
}

func TestFindFunctionByName_EmptyDB(t *testing.T) {
	store := setupTestDB(t)
	defer store.Close()

	_, _, found := findFunctionByName(store, "anything")
	if found {
		t.Fatal("expected not to find function in empty db")
	}
}

// HandleGetFunction

func TestHandleGetFunction_Found(t *testing.T) {
	store := setupTestDB(t)
	defer store.Close()
	seedFunction(t, store, testLambda)

	router := setupRouter()
	HandleGetFunction(router, store)

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/function/get/test-fn", nil)
	router.ServeHTTP(w, req)

	if w.Code != 200 {
		t.Errorf("expected 200, got %d", w.Code)
	}
	var result types.LambdaFun
	json.Unmarshal(w.Body.Bytes(), &result)
	if result.Name != "test-fn" {
		t.Errorf("expected name test-fn, got %s", result.Name)
	}
	if result.Port != "9002" {
		t.Errorf("expected port 9002, got %s", result.Port)
	}
}

func TestHandleGetFunction_NotFound(t *testing.T) {
	store := setupTestDB(t)
	defer store.Close()

	router := setupRouter()
	HandleGetFunction(router, store)

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/function/get/nonexistent", nil)
	router.ServeHTTP(w, req)

	if w.Code != 404 {
		t.Errorf("expected 404, got %d", w.Code)
	}
}

// HandleListFunctions

func TestHandleListFunctions_Empty(t *testing.T) {
	store := setupTestDB(t)
	defer store.Close()

	router := setupRouter()
	HandleListFunctions(router, store)

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/function/list", nil)
	router.ServeHTTP(w, req)

	if w.Code != 200 {
		t.Errorf("expected 200, got %d", w.Code)
	}
}

func TestHandleListFunctions_WithFunctions(t *testing.T) {
	store := setupTestDB(t)
	defer store.Close()
	seedFunction(t, store, testLambda)

	router := setupRouter()
	HandleListFunctions(router, store)

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/function/list", nil)
	router.ServeHTTP(w, req)

	if w.Code != 200 {
		t.Errorf("expected 200, got %d", w.Code)
	}

	var result struct {
		Keys      []string          `json:"keys"`
		Functions []types.LambdaFun `json:"functions"`
	}
	json.Unmarshal(w.Body.Bytes(), &result)
	if len(result.Functions) != 1 {
		t.Fatalf("expected 1 function, got %d", len(result.Functions))
	}
	if result.Functions[0].Name != "test-fn" {
		t.Errorf("expected name test-fn, got %s", result.Functions[0].Name)
	}
}

// HandleDeleteFunction

func TestHandleDeleteFunction_Found(t *testing.T) {
	store := setupTestDB(t)
	defer store.Close()
	seedFunction(t, store, testLambda)

	mock := &mockDocker{}
	original := newDockerClient
	newDockerClient = func() (dockerClient, error) { return mock, nil }
	defer func() { newDockerClient = original }()

	router := setupRouter()
	HandleDeleteFunction(router, store)

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("DELETE", "/function/delete/test-fn", nil)
	router.ServeHTTP(w, req)

	if w.Code != 200 {
		t.Errorf("expected 200, got %d", w.Code)
	}
	if len(mock.killCalls) != 1 || mock.killCalls[0] != "abc123" {
		t.Errorf("expected Kill(abc123), got %v", mock.killCalls)
	}
	if len(mock.removeContainerCalls) != 1 || mock.removeContainerCalls[0] != "abc123" {
		t.Errorf("expected RemoveContainer(abc123), got %v", mock.removeContainerCalls)
	}
	if len(mock.removeImageCalls) != 1 || mock.removeImageCalls[0] != "test-fn:latest" {
		t.Errorf("expected RemoveImage(test-fn:latest), got %v", mock.removeImageCalls)
	}
	if store.HasKey("abc123") {
		t.Error("expected function to be deleted from db")
	}
}

func TestHandleDeleteFunction_NotFound(t *testing.T) {
	store := setupTestDB(t)
	defer store.Close()

	router := setupRouter()
	HandleDeleteFunction(router, store)

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("DELETE", "/function/delete/nonexistent", nil)
	router.ServeHTTP(w, req)

	if w.Code != 404 {
		t.Errorf("expected 404, got %d", w.Code)
	}
}
