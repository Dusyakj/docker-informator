package requester

import (
	"encoding/json"
	"fmt"
	"net/http"
	"project8/internal/models"
	"project8/internal/tools"
)

func GetInformation(body []byte) ([]byte, int, error) {
	var data models.Data
	if err := json.Unmarshal(body, &data); err != nil {
		return nil, http.StatusBadRequest, fmt.Errorf("invalid request data: %w", err)
	}

	if data.Repository == "" {
		return nil, http.StatusBadRequest, fmt.Errorf("field 'repository' is required")
	}

	if data.Name == "" {
		return nil, http.StatusBadRequest, fmt.Errorf("field 'name' is required")
	}

	if data.Tag == "" {
		data.Tag = "latest"
	}

	// if err := tools.ValidateTag(data); err != nil {
	// 	return nil, http.StatusNotFound, fmt.Errorf("tag validation failed: %w", err)
	// }

	digest, err := tools.GetManifest(data)
	if err != nil {
		return nil, http.StatusNotFound, fmt.Errorf("manifest error: %w", err)
	}

	var info models.Information
	if err := tools.GetProperty(&info, data, digest); err != nil {
		return nil, http.StatusInternalServerError, fmt.Errorf("property error: %w", err)
	}

	jsonData, err := json.Marshal(info)
	if err != nil {
		return nil, http.StatusInternalServerError, fmt.Errorf("marshal error: %w", err)
	}

	return jsonData, http.StatusOK, nil
}
