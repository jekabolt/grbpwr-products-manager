package app

import (
	"encoding/json"
	"net/http"

	"github.com/go-chi/chi"
	"github.com/jekabolt/grbpwr-manager/store"
	"github.com/rs/zerolog/log"
)

// admin panel
func (s *Server) addProduct(w http.ResponseWriter, r *http.Request) {

	w.Header().Set("Content-type", "application/json")

	product := &store.Product{}

	defer r.Body.Close()

	if err := json.NewDecoder(r.Body).Decode(product); err != nil {
		log.Error().Err(err).Msgf("addProduct:json.NewDecoder [%v]", err.Error())
		err := map[string]interface{}{"addProduct:Decode": err}
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(err)
		return
	}

	if errors := product.Validate(); len(errors) > 0 {
		log.Error().Msgf("addProduct:product.Validate [%v]", errors)
		err := map[string]interface{}{"validationError": errors}
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(err)
	}

	err := s.DB.AddProduct(product)
	if err != nil {
		log.Error().Msgf("addProduct:AddProduct [%v]", err)
		err := map[string]interface{}{"addProduct:AddProduct": err}
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(err)
	}

	resp := map[string]interface{}{"status": http.StatusText(http.StatusCreated)}
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(resp)

}

func (s *Server) deleteProductById(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	w.Header().Set("Content-type", "application/json")

	err := s.DB.DeleteProductById(id)
	if err != nil {
		log.Error().Msgf("addProduct:AddProduct [%v]", err)
		err := map[string]interface{}{"addProduct:AddProduct": err}
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(err)
	}

	resp := map[string]interface{}{"status": http.StatusText(http.StatusCreated)}
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(resp)
}

func (s *Server) getProductsById(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	w.Header().Set("Content-type", "application/json")

	product, err := s.DB.GetProductsById(id)
	if err != nil {
		log.Error().Msgf("addProduct:AddProduct [%v]", err)
		err := map[string]interface{}{"addProduct:AddProduct": err}
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(err)
	}

	json.NewEncoder(w).Encode(product)
}

func (s *Server) modifyProductsById(w http.ResponseWriter, r *http.Request) {

	id := chi.URLParam(r, "id")

	w.Header().Set("Content-type", "application/json")

	product := &store.Product{}

	defer r.Body.Close()

	if err := json.NewDecoder(r.Body).Decode(product); err != nil {
		log.Error().Err(err).Msgf("addProduct:json.NewDecoder [%v]", err.Error())
		err := map[string]interface{}{"addProduct:Decode": err}
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(err)
		return
	}

	if errors := product.Validate(); len(errors) > 0 {
		log.Error().Msgf("addProduct:product.Validate [%v]", errors)
		err := map[string]interface{}{"validationError": errors}
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(err)
	}

	err := s.DB.ModifyProductById(id, product)
	if err != nil {
		log.Error().Msgf("addProduct:AddProduct [%v]", err)
		err := map[string]interface{}{"addProduct:AddProduct": err}
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(err)
	}

	resp := map[string]interface{}{"status": http.StatusText(http.StatusCreated)}
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(resp)
}

// site
func (s *Server) getProductsByCategory(w http.ResponseWriter, r *http.Request) {
	category := chi.URLParam(r, "category")
	w.Header().Set("Content-type", "application/json")

	product, err := s.DB.GetAllProductsInCategory(category)
	if err != nil {
		log.Error().Msgf("addProduct:AddProduct [%v]", err)
		err := map[string]interface{}{"addProduct:AddProduct": err}
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(err)
	}

	json.NewEncoder(w).Encode(product)
}

func (s *Server) getAllProductsList(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-type", "application/json")

	product, err := s.DB.GetAllProducts()
	if err != nil {
		log.Error().Msgf("addProduct:AddProduct [%v]", err)
		err := map[string]interface{}{"addProduct:AddProduct": err}
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(err)
	}

	json.NewEncoder(w).Encode(product)
}
