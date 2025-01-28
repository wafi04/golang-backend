package filehandler

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"math/rand"
	"net/http"
	"time"

	"github.com/wafi04/golang-backend/grpc/pb"
	"github.com/wafi04/golang-backend/services/common"
	"github.com/wafi04/golang-backend/services/gateway/server/config"
)

type FileHandler struct {
	FilesClient pb.FileServiceClient
}

func NewFilesGateway(ctx context.Context) (*FileHandler, error) {
	conn, err := config.ConnectWithRetry(config.Load().FilesServiceURL, "files")
	if err != nil {
		return nil, err
	}

	return &FileHandler{
		FilesClient: pb.NewFileServiceClient(conn),
	}, nil
}

func (s *FileHandler) HandleUploadFile(w http.ResponseWriter, r *http.Request) {
	// Pastikan method adalah POST
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	// Parse multipart form
	err := r.ParseMultipartForm(10 << 20) // Maksimal 10 MB
	if err != nil {
		http.Error(w, "Failed to parse form", http.StatusBadRequest)
		return
	}

	file, _, err := r.FormFile("file")
	if err != nil {
		http.Error(w, "Error retrieving file", http.StatusBadRequest)
		return
	}
	defer file.Close()

	// Baca isi file
	fileBytes, err := io.ReadAll(file)
	if err != nil {
		http.Error(w, "Failed to read file", http.StatusInternalServerError)
		return
	}

	// Siapkan context dengan timeout
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	PublicID := fmt.Sprintf("%06d", rand.Intn(1000000))

	uploadRequest := &pb.FileUploadRequest{
		FileData: fileBytes,
		Folder:   "testing",
		PublicId: PublicID,
	}

	response, err := s.FilesClient.UploadFile(ctx, uploadRequest)
	if err != nil {
		http.Error(w, fmt.Sprintf("Upload failed: %v", err), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	resp := common.Success(response, "Files Upload Successfully")
	if err := json.NewEncoder(w).Encode(resp); err != nil {
		log.Printf("Error encoding response: %v", err)
		http.Error(w, "Error encoding response", http.StatusInternalServerError)
		return
	}
}

func (s *FileHandler) HandleUploadFiles(w http.ResponseWriter, r *http.Request) (string, error) {
	if err := r.ParseMultipartForm(10 << 20); err != nil { // 10 MB max
		http.Error(w, "Failed to parse form", http.StatusBadRequest)
	}

	file, _, err := r.FormFile("file")
	if err != nil {
		http.Error(w, "Error retrieving file", http.StatusBadRequest)
	}
	defer file.Close()

	jsonData := r.FormValue("data")
	var req pb.CreateCategoryRequest
	if err := json.Unmarshal([]byte(jsonData), &req); err != nil {
		http.Error(w, "Invalid JSON data", http.StatusBadRequest)
	}

	// Upload file
	fileBytes, err := io.ReadAll(file)
	if err != nil {
		http.Error(w, "Failed to read file", http.StatusInternalServerError)

	}

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	PublicID := fmt.Sprintf("%06d", rand.Intn(1000000))
	uploadRequest := &pb.FileUploadRequest{
		FileData: fileBytes,
		Folder:   "testing",
		PublicId: PublicID,
	}

	response, err := s.FilesClient.UploadFile(ctx, uploadRequest)
	if err != nil {
		return "", err
	}

	return response.Url, nil
}
