package producthandler

import (
	"fmt"
	"io"
	"math/rand"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/wafi04/golang-backend/grpc/pb"
	"github.com/wafi04/golang-backend/services/common"
)

func (p *ProductHandler) HandleAddProductImage(w http.ResponseWriter, r *http.Request) {
	if err := r.ParseMultipartForm(10 << 20); err != nil {
		http.Error(w, "Failed to parse form", http.StatusBadRequest)
		return
	}

	vars := mux.Vars(r)
	VariantId, ok := vars["id"]
	if !ok {
		http.Error(w, "Product Variants Id is required", http.StatusBadRequest)
		return
	}

	file, _, err := r.FormFile("file")
	var imageUrl string
	if err == nil {
		defer file.Close()

		fileBytes, err := io.ReadAll(file)
		if err != nil {
			http.Error(w, "Failed to read file", http.StatusInternalServerError)
			return
		}

		PublicID := fmt.Sprintf("%06d", rand.Intn(1000000))
		uploadRequest := &pb.FileUploadRequest{
			FileData: fileBytes,
			Folder:   "products",
			PublicId: PublicID,
		}

		uploadResponse, err := p.filesClient.UploadFile(r.Context(), uploadRequest)
		if err != nil {
			common.SendErrorResponse(w, http.StatusBadRequest, "Failed to upload image")
			return
		}

		imageUrl = uploadResponse.Url
	}

	var imageUrlPtr *string
	if imageUrl != "" {
		imageUrlPtr = &imageUrl
	}

	productImage, err := p.productclient.AddProductImage(r.Context(), &pb.AddProductImageRequest{
		VariantId: VariantId,
		Url:       *imageUrlPtr,
		IsMain:    true,
	})
	if err != nil {
		common.SendErrorResponse(w, http.StatusBadRequest, "Failed to created product image")
		return
	}
	common.SendSuccessResponse(w, http.StatusCreated, "Created Product Image succesfully", productImage)
}

func (p *ProductHandler) HandleDeleteProductImage(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	ProductImageId, ok := vars["id"]
	if !ok {
		http.Error(w, "Product images Id is required", http.StatusBadRequest)
		return
	}

	productImage, err := p.productclient.DeleteProductImage(r.Context(), &pb.DeleteProductImageRequest{
		Id: ProductImageId,
	})

	if err != nil {
		common.SendErrorResponse(w, http.StatusBadRequest, "Failed to Delete product image")
		return
	}
	common.SendSuccessResponse(w, http.StatusOK, "Delete Product Image succesfully", productImage)
}
