package supplier

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"github.com/go-playground/validator"
	"github.com/nats-io/nats.go"
)

type CustomResponse struct {
	Message string                 `json:"message"`
	Status  int                    `json:"status"`
	Data    map[string]interface{} `json:"data,omitempty"`
}

type Handler struct {
	service Service
}

type GetSupplierByIDRequest struct {
	Data struct {
		ID uint32 `json:"id" validate:"required"`
	} `json:"data"`
}
type GetAllSuppliersRequest struct {
	Data struct {
		Page  *int `json:"page" validate:"omitempty,min=1"`
		Limit *int `json:"limit" validate:"omitempty,min=1,max=100"`
	} `json:"data"`
}
type CreateSupplierRequest struct {
	Data struct {
		Name    string  `json:"name" validate:"required"`
		Email   *string `json:"email,omitempty" validate:"omitempty,email"`
		Phone   *string `json:"phone,omitempty" validate:"omitempty,max=20"`
		Address *string `json:"address,omitempty" validate:"omitempty,max=255"`
		RFC     *string `json:"rfc,omitempty" validate:"omitempty,len=13"`
	} `json:"data"`
}
type UpdateSupplierRequest struct {
	Data struct {
		ID      uint32  `json:"id" validate:"required"`
		Name    *string `json:"name,omitempty" validate:"omitempty"`
		Email   *string `json:"email,omitempty" validate:"omitempty,email"`
		Phone   *string `json:"phone,omitempty" validate:"omitempty,max=20"`
		Address *string `json:"address,omitempty" validate:"omitempty,max=255"`
		RFC     *string `json:"rfc,omitempty" validate:"omitempty,len=13"`
	} `json:"data"`
}

func NewHandler(service Service) *Handler {
	return &Handler{service}
}

func (h *Handler) Subscribe(nc *nats.Conn) {

	subscribe := func(subject string, handler func([]byte) ([]byte, error)) {
		_, err := nc.Subscribe(subject, func(msg *nats.Msg) {
			resp, err := handler(msg.Data)
			if err != nil {
				log.Printf("Error handling subject %s: %v", subject, err)
			}
			_ = msg.Respond(resp)
		})
		if err != nil {
			log.Printf("Error subscribing to %s: %v", subject, err)
		}
	}
	subscribe("suppliers.getAll", h.getAll)
	subscribe("suppliers.findById", h.getByID)
	subscribe("suppliers.create", h.Create)
	subscribe("suppliers.update", h.Update)

}

func NewSupplierFromRequest(req CreateSupplierRequest) *Supplier {
	return &Supplier{
		Name:    req.Data.Name,
		Email:   req.Data.Email,
		Phone:   req.Data.Phone,
		Address: req.Data.Address,
		RFC:     req.Data.RFC,
	}
}

func ParseSupplierFromRequest(req UpdateSupplierRequest) *Supplier {
	return &Supplier{
		ID:      uint(req.Data.ID),
		Name:    *req.Data.Name,
		Email:   req.Data.Email,
		Phone:   req.Data.Phone,
		Address: req.Data.Address,
		RFC:     req.Data.RFC,
	}
}

func (h *Handler) Update(data []byte) ([]byte, error) {
	var payload UpdateSupplierRequest
	err := parseAndValidate(data, &payload)
	if err != nil {
		return marshalResponse(err.Error(), http.StatusBadRequest, nil)
	}
	supplier := ParseSupplierFromRequest(payload)
	if err := h.service.UpdateSupplier(supplier); err != nil {
		return marshalResponse("couldn't update supplier", http.StatusInternalServerError, nil)
	}
	return marshalResponse("", http.StatusCreated, map[string]interface{}{"supplier": supplier})
}

func (h *Handler) Create(data []byte) ([]byte, error) {
	var supplierRequest CreateSupplierRequest
	err := parseAndValidate(data, &supplierRequest)
	if err != nil {
		return marshalResponse(err.Error(), http.StatusBadRequest, nil)
	}
	supplier := NewSupplierFromRequest(supplierRequest)
	if err := h.service.CreateSupplier(supplier); err != nil {
		return marshalResponse("couldn't create supplier", http.StatusInternalServerError, nil)
	}

	return marshalResponse("", http.StatusCreated, map[string]interface{}{"supplier": supplier})
}

func (h *Handler) getByID(data []byte) ([]byte, error) {
	var payload GetSupplierByIDRequest
	err := parseAndValidate(data, &payload)
	if err != nil {
		return marshalResponse(err.Error(), http.StatusBadRequest, nil)
	}

	supplier, err := h.service.GetSupplierByID(uint(payload.Data.ID))
	if err != nil {
		return marshalResponse(err.Error(), http.StatusNotFound, nil)
	}

	return marshalResponse("", http.StatusCreated, map[string]interface{}{"supplier": supplier})
}

func (h *Handler) getAll(data []byte) ([]byte, error) {
	var payload GetAllSuppliersRequest

	err := parseAndValidate(data, &payload)
	if err != nil {
		return marshalResponse(err.Error(), http.StatusBadRequest, nil)
	}
	suppliers, err := h.service.GetAllSuppliers(*payload.Data.Page, *payload.Data.Limit)
	if err != nil {
		return marshalResponse("Error retrieving suppliers", http.StatusInternalServerError, nil)
	}
	return marshalResponse("", http.StatusOK, map[string]interface{}{"suppliers": suppliers})
}

func parseAndValidate[T any](data []byte, dst *T) error {
	if err := json.Unmarshal(data, dst); err != nil {
		return fmt.Errorf("invalid payload")
	}
	validate := validator.New()
	if err := validate.Struct(dst); err != nil {
		return fmt.Errorf("invalid payload structure")
	}

	return nil
}

func marshalResponse(message string, status int, data map[string]interface{}) ([]byte, error) {
	resp := CustomResponse{
		Message: message,
		Status:  status,
		Data:    data,
	}

	if data != nil {
		return json.Marshal(resp.Data)
	}

	return json.Marshal(resp)
}
