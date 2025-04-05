package models

type ApiResponse struct {
  Success bool `json:"success,omitempty"`
  Message string `json:"message,omitempty"`
  Data any `json:"data,omitempty"` 
}