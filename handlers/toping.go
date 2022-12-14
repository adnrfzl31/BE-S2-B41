package handlers

import (
	dto "BE-S2-B41/dto/result"
	topingdto "BE-S2-B41/dto/toping"
	"BE-S2-B41/models"
	"BE-S2-B41/repositories"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"strconv"

	"github.com/cloudinary/cloudinary-go/v2"
	"github.com/cloudinary/cloudinary-go/v2/api/uploader"
	"github.com/go-playground/validator/v10"
	"github.com/golang-jwt/jwt/v4"
	"github.com/gorilla/mux"
)

type handlerToping struct {
	TopingRepository repositories.TopingRepository
}

// Create `path_file` Global variable
var path_files = os.Getenv("PATH_FILE")

func HandlerToping(TopingRepository repositories.TopingRepository) *handlerToping {
	return &handlerToping{TopingRepository}
}

func (h *handlerToping) FindTopings(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	topings, err := h.TopingRepository.FindTopings()
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		response := dto.ErrorResult{Status: "Server Error", Message: err.Error()}
		json.NewEncoder(w).Encode(response)
		return
	}

	// Create Embed Path File on Image property
	// for i, p := range topings {
	// 	imagePath := os.Getenv("PATH_FILE") + p.Image
	// 	topings[i].Image = imagePath
	// }

	w.WriteHeader(http.StatusOK)
	response := dto.SuccessResult{Status: "Success", Data: topings}
	json.NewEncoder(w).Encode(response)
}

func (h *handlerToping) GetToping(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	id, _ := strconv.Atoi(mux.Vars(r)["id"])

	var toping models.Toping
	toping, err := h.TopingRepository.GetToping(id)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		response := dto.ErrorResult{Status: "Server Error", Message: err.Error()}
		json.NewEncoder(w).Encode(response)
		return
	}

	// Create Embed Path File on Image property
	toping.Image = path_files + toping.Image

	w.WriteHeader(http.StatusOK)
	response := dto.SuccessResult{Status: "Success", Data: toping}
	json.NewEncoder(w).Encode(response)
}

func (h *handlerToping) CreateToping(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	// ambil data user id dari token yang sudah di decode
	userInfo := r.Context().Value("userInfo").(jwt.MapClaims)
	// userId := int(userInfo["id"].(float64))
	userRole := userInfo["role"]

	if userRole != "admin" {
		w.WriteHeader(http.StatusBadRequest)
		response := dto.ErrorResult{Status: "Failed", Message: "you're not admin"}
		json.NewEncoder(w).Encode(response)
		return
	}

	// Get dataFile from midleware and store to filename variable
	dataContex := r.Context().Value("dataFile")
	filepath := dataContex.(string)

	price, _ := strconv.Atoi(r.FormValue("price"))
	request := topingdto.AddToping{
		Nametoping: r.FormValue("nametoping"),
		Price:      price,
		Image:      filepath,
	}

	validation := validator.New()
	err := validation.Struct(request)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		response := dto.ErrorResult{Status: "Server Error", Message: err.Error()}
		json.NewEncoder(w).Encode(response)
		return
	}

	// Declare Context Background, Cloud Name, API Key, API Secret ...
	var ctx = context.Background()
	var CLOUD_NAME = os.Getenv("CLOUD_NAME")
	var API_KEY = os.Getenv("API_KEY")
	var API_SECRET = os.Getenv("API_SECRET")

	// Add your Cloudinary credentials ...
	cld, _ := cloudinary.NewFromParams(CLOUD_NAME, API_KEY, API_SECRET)

	// Upload file to Cloudinary ...
	resp, err := cld.Upload.Upload(ctx, filepath, uploader.UploadParams{Folder: "waysbukcks"})

	if err != nil {
		fmt.Println(err.Error())
	}

	toping := models.Toping{
		Nametoping: request.Nametoping,
		Price:      request.Price,
		Image:      resp.SecureURL,
	}

	toping, err = h.TopingRepository.CreateToping(toping)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		response := dto.ErrorResult{Status: "Server Error", Message: err.Error()}
		json.NewEncoder(w).Encode(response)
		return
	}

	toping, _ = h.TopingRepository.GetToping(toping.ID)

	w.WriteHeader(http.StatusOK)
	response := dto.SuccessResult{Status: "Success", Data: toping}
	json.NewEncoder(w).Encode(response)
}

func (h *handlerToping) UpdateToping(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	dataContex := r.Context().Value("dataFile")
	filepath := dataContex.(string)

	price, _ := strconv.Atoi(r.FormValue("price"))

	request := topingdto.UpdateToping{
		Nametoping: r.FormValue("nametoping"),
		Price:      price,
		Image:      filepath,
	}

	id, _ := strconv.Atoi(mux.Vars(r)["id"])

	userInfo := r.Context().Value("userInfo").(jwt.MapClaims)
	userRole := userInfo["role"]
	userId := int(userInfo["id"].(float64))

	if userId != id && userRole != "admin" {
		w.WriteHeader(http.StatusBadRequest)
		response := dto.ErrorResult{Status: "Failed", Message: "you're not admin"}
		json.NewEncoder(w).Encode(response)
		return
	}

	// Declare Context Background, Cloud Name, API Key, API Secret ...
	var ctx = context.Background()
	var CLOUD_NAME = os.Getenv("CLOUD_NAME")
	var API_KEY = os.Getenv("API_KEY")
	var API_SECRET = os.Getenv("API_SECRET")

	// Add your Cloudinary credentials ...
	cld, _ := cloudinary.NewFromParams(CLOUD_NAME, API_KEY, API_SECRET)

	// Upload file to Cloudinary ...
	resp, err := cld.Upload.Upload(ctx, filepath, uploader.UploadParams{Folder: "waysbukcks"})

	if err != nil {
		fmt.Println(err.Error())
	}

	toping, err := h.TopingRepository.GetToping(int(id))
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		response := dto.ErrorResult{Status: "Failed", Message: err.Error()}
		json.NewEncoder(w).Encode(response)
		return
	}

	if request.Nametoping != "" {
		toping.Nametoping = request.Nametoping
	}

	if request.Price != 0 {
		toping.Price = request.Price
	}

	if request.Image != "" {
		toping.Image = resp.SecureURL
	}

	data, err := h.TopingRepository.UpdateToping(toping)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		response := dto.ErrorResult{Status: "Server Error", Message: err.Error()}
		json.NewEncoder(w).Encode(response)
		return
	}
	w.WriteHeader(http.StatusOK)
	response := dto.SuccessResult{Status: "Success", Data: data}
	json.NewEncoder(w).Encode(response)
}

func (h *handlerToping) DeleteToping(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	id, _ := strconv.Atoi(mux.Vars(r)["id"])

	userInfo := r.Context().Value("userInfo").(jwt.MapClaims)
	userRole := userInfo["role"]
	userId := int(userInfo["id"].(float64))

	if userId != id && userRole != "admin" {
		w.WriteHeader(http.StatusBadRequest)
		response := dto.ErrorResult{Status: "Failed", Message: "you're not admin"}
		json.NewEncoder(w).Encode(response)
		return
	}

	topings, err := h.TopingRepository.GetToping(id)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		response := dto.ErrorResult{Status: "Failed", Message: err.Error()}
		json.NewEncoder(w).Encode(response)
		return
	}

	deleteToping, err := h.TopingRepository.DeleteToping(topings)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		response := dto.ErrorResult{Status: "Server Error", Message: err.Error()}
		json.NewEncoder(w).Encode(response)
		return
	}

	w.WriteHeader(http.StatusOK)
	response := dto.SuccessResult{Status: "Success", Data: convertResponseToping(deleteToping)}
	json.NewEncoder(w).Encode(response)
}

func convertResponseToping(u models.Toping) topingdto.TopingResponseDelete {
	return topingdto.TopingResponseDelete{
		ID: u.ID,
	}
}
