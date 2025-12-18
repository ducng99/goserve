package files_test

import (
	"io"
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"testing"

	"github.com/ducng99/goserve/internal/files"
	"github.com/ducng99/goserve/internal/server"
)

func createTestServer(t *testing.T, rootDir string) *httptest.Server {
	// Create a server config with the test root directory
	config := server.ServerConfig{
		RootDir: rootDir,
	}

	// Use the existing NewServeMux function instead of manually creating mux
	mux := config.NewServeMux()

	ts := httptest.NewServer(mux)
	t.Cleanup(func() {
		ts.Close()
	})

	return ts
}

func TestServerFileServingFunctionality(t *testing.T) {
	// Test the core file serving functionality by directly calling the route handler
	testRootDir := filepath.Join("..", "testdata", "root")

	// Test that the SanitisePath function works correctly in context
	sanitisedPath, err := files.SanitisePath(testRootDir, "/file1.txt")
	if err != nil {
		t.Fatalf("SanitisePath failed for valid file: %v", err)
	}

	// Verify the file exists
	if _, err := os.Stat(sanitisedPath); os.IsNotExist(err) {
		t.Fatalf("Sanitised path does not exist: %s", sanitisedPath)
	}

	// Read the actual file content to determine expected content
	fileContent, err := os.ReadFile(sanitisedPath)
	if err != nil {
		t.Fatalf("Failed to read file content: %v", err)
	}

	ts := createTestServer(t, testRootDir)

	// Test file serving
	resp, err := http.Get(ts.URL + "/file1.txt")
	if err != nil {
		t.Fatalf("Failed to get file1.txt: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		t.Errorf("Expected status 200, got %d", resp.StatusCode)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		t.Fatalf("Failed to read response body: %v", err)
	}

	// Check that the content matches what we expect (read from actual file)
	expectedContent := string(fileContent)
	if string(body) != expectedContent {
		t.Errorf("Expected content %q, got %q", expectedContent, string(body))
	}
}

func TestServerRootDirectoryAccess(t *testing.T) {
	testRootDir := filepath.Join("..", "testdata", "root")
	ts := createTestServer(t, testRootDir)

	// Test directory access returns 200 (not rendering directory view in test)
	resp, err := http.Get(ts.URL + "/")
	if err != nil {
		t.Fatalf("Failed to get root directory: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		t.Errorf("Expected status 200 for directory, got %d", resp.StatusCode)
	}
}

func TestServerMissingFileReturns404(t *testing.T) {
	testRootDir := filepath.Join("..", "testdata", "root")
	ts := createTestServer(t, testRootDir)

	// Test missing file returns 404
	resp, err := http.Get(ts.URL + "/nonexistent.txt")
	if err != nil {
		t.Fatalf("Failed to get nonexistent file: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusNotFound {
		t.Errorf("Expected status 404 for missing file, got %d", resp.StatusCode)
	}
}

func TestServerSubdirectoryFileAccess(t *testing.T) {
	testRootDir := filepath.Join("..", "testdata", "root")
	ts := createTestServer(t, testRootDir)

	// Test directory access
	resp, err := http.Get(ts.URL + "/dir1/file1.txt")
	if err != nil {
		t.Fatalf("Failed to get file in subdirectory: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		t.Errorf("Expected status 200 for file in subdirectory, got %d", resp.StatusCode)
	}
}
