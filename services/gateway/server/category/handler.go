package categoryhandler

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"math/rand"
	"net/http"
	"strconv"
	"strings"

	"github.com/gorilla/mux"
	"github.com/wafi04/golang-backend/grpc/pb"
	"github.com/wafi04/golang-backend/services/common"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

type CategoryHandler struct {
	categoryClient pb.CategoryServiceClient
    filesCLient pb.FileServiceClient
}

func NewCategoryGateway(ctx context.Context) (*CategoryHandler, error) {
	log.Println("Attempting to connect to category service...")

	conn, err := grpc.DialContext(ctx,
		common.LoadEnv("CATEGORY_PORT"),
		grpc.WithTransportCredentials(insecure.NewCredentials()),
		grpc.WithBlock(),
	)
	if err != nil {
		log.Printf("Failed to connect to category service: %v", err)
		return nil, fmt.Errorf("failed to connect to category service: %v", err)
	}
	connFile, err := grpc.DialContext(ctx,
		common.LoadEnv("FILES_PORT"),
		grpc.WithTransportCredentials(insecure.NewCredentials()),
		grpc.WithBlock(),
	)
	if err != nil {
		log.Printf("Failed to connect to files service: %v", err)
		return nil, fmt.Errorf("failed to connect to files service: %v", err)
	}

	log.Println("Successfully connected to auth service")
	return &CategoryHandler{
		categoryClient: pb.NewCategoryServiceClient(conn),
        filesCLient: pb.NewFileServiceClient(connFile),
	}, nil
}
func (h *CategoryHandler) HandleCreateCategory(w http.ResponseWriter, r *http.Request) {
    if err := r.ParseMultipartForm(10 << 20); err != nil {
        http.Error(w, "Failed to parse form", http.StatusBadRequest)
        return
    }
    // Get data dari form
    name := r.FormValue("name")
    description := r.FormValue("description")
    parentID := r.FormValue("parent_id")

    file,_, err := r.FormFile("file")
    var imageUrl string
    if err == nil {
        defer file.Close()
        
        // Baca file
        fileBytes, err := io.ReadAll(file)
        if err != nil {
            http.Error(w, "Failed to read file", http.StatusInternalServerError)
            return
        }

        // Upload file
        PublicID := fmt.Sprintf("%06d", rand.Intn(1000000))
        uploadRequest := &pb.FileUploadRequest{
            FileData: fileBytes,
            Folder:   "categories",
            PublicId: PublicID,
        }

        uploadResponse, err := h.filesCLient.UploadFile(r.Context(), uploadRequest)
        if err != nil {
            http.Error(w, fmt.Sprintf("File upload failed: %v", err), http.StatusInternalServerError)
            return
        }

        imageUrl = uploadResponse.Url
    }

    var parentIDPtr *string
    if parentID != "" {
        parentIDPtr = &parentID
    }

    var imageUrlPtr *string
    if imageUrl != "" {
        imageUrlPtr = &imageUrl
    }

    resp, err := h.categoryClient.CreateCategory(r.Context(), &pb.CreateCategoryRequest{
        Name:        name,
        Description: description,
        Image:       imageUrlPtr,
        ParentId:    parentIDPtr,
    })

    if err != nil {
        http.Error(w, fmt.Sprintf("Error creating category: %v", err), http.StatusInternalServerError)
        return
    }
    response :=  common.Success(resp, "created category successfully")

    w.Header().Set("Content-Type", "application/json")
    if err = json.NewEncoder(w).Encode(response); err != nil {
        log.Printf("Error encoding response: %v", err)
        http.Error(w, "Error encoding response", http.StatusInternalServerError)
        return
    }
}
func (h *CategoryHandler) HandleGetCategories(w http.ResponseWriter, r *http.Request) {
    categories, err := h.categoryClient.ListCategories(r.Context(), &pb.ListCategoriesRequest{})
    if err != nil {
        common.Error(http.StatusBadRequest,"Failed to retrieve categories")
        return
    }

    common.Success(categories, "Categories retrieved successfully")
}

func (h *CategoryHandler) HandleListCategories(w http.ResponseWriter, r *http.Request) {

    page, _ := strconv.Atoi(r.URL.Query().Get("page"))
    limit, _ := strconv.Atoi(r.URL.Query().Get("limit"))
    parentID := r.URL.Query().Get("parent_id")
    includeChildren := r.URL.Query().Get("include_children") == "true"

    if page <= 0 {
        page = 1
    }
    if limit <= 0 {
        limit = 10
    }

    req := &pb.ListCategoriesRequest{
        Page:            int32(page),
        Limit:           int32(limit),
        IncludeChildren: includeChildren,
    }
    
    if parentID != "" {
        req.ParentId = &parentID
    }

    resp, err := h.categoryClient.ListCategories(r.Context(), req)
    if err != nil {
        log.Printf("Error calling ListCategories: %v", err)
        http.Error(w, "Error retrieving categories", http.StatusInternalServerError)
        return
    }

    w.Header().Set("Content-Type", "application/json")
    response := common.Success(map[string]interface{}{
        "categories": resp.Categories,
        "total":     resp.Total,
        "page":      page,
        "limit":     limit,
    }, "Categories retrieved successfully")
    
    if err = json.NewEncoder(w).Encode(response); err != nil {
        log.Printf("Error encoding response: %v", err)
        http.Error(w, "Error encoding response", http.StatusInternalServerError)
        return
    }
}


func (h *CategoryHandler) HandleUpdateCategory(w http.ResponseWriter, r *http.Request) {
    w.Header().Set("Content-Type", "application/json")
    vars := mux.Vars(r)
    id, ok := vars["id"]
    if !ok {
        http.Error(w, "Category ID is required", http.StatusBadRequest)
        return
    }

    var request struct {
        Name        *string `json:"name,omitempty"`
        Description *string `json:"description,omitempty"`
        Image       *string `json:"image,omitempty"`
        ParentID    *string `json:"parent_id,omitempty"`
        Depth       *int32  `json:"depth,omitempty"`
    }

    if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
        http.Error(w, "Invalid request body", http.StatusBadRequest)
        return
    }

    // Create gRPC request
    updateReq := &pb.UpdateCategoryRequest{
        Id:          id,
        Name:        request.Name,
        Description: request.Description,
        Image:       request.Image,
        ParentId:    request.ParentID,
        Depth:       request.Depth,
    }

    ctx := r.Context()
    category, err := h.categoryClient.UpdateCategory(ctx, updateReq)
    if err != nil {
        switch {
        case strings.Contains(err.Error(), "not found"):
            http.Error(w, err.Error(), http.StatusNotFound)
        case strings.Contains(err.Error(), "invalid"):
            http.Error(w, err.Error(), http.StatusBadRequest)
        default:
            http.Error(w, "Internal server error", http.StatusInternalServerError)
        }
        return
    }

    response :=  common.Success(category,"Update Category Succcess")
    if err := json.NewEncoder(w).Encode(response); err != nil {
        http.Error(w, "Error encoding response", http.StatusInternalServerError)
        return
    }
}

func (h *CategoryHandler) HandleDeleteCategory(w http.ResponseWriter, r *http.Request) {
    w.Header().Set("Content-Type", "application/json")
    vars := mux.Vars(r)
    id, ok := vars["id"]
    if !ok {
        http.Error(w, "Category ID is required", http.StatusBadRequest)
        return
    }

    updateReq := &pb.DeleteCategoryRequest{
        Id:          id,
        DeleteChildren: false,
    }

    ctx := r.Context()
    category, err := h.categoryClient.DeleteCategory(ctx, updateReq)
    if err != nil {
        switch {
        case strings.Contains(err.Error(), "not found"):
            http.Error(w, err.Error(), http.StatusNotFound)
        case strings.Contains(err.Error(), "invalid"):
            http.Error(w, err.Error(), http.StatusBadRequest)
        default:
            http.Error(w, "Internal server error", http.StatusInternalServerError)
        }
        return
    }

    response :=  common.Success(category,"Delete Category Succcess")
    if err := json.NewEncoder(w).Encode(response); err != nil {
        http.Error(w, "Error encoding response", http.StatusInternalServerError)
        return
    }
}