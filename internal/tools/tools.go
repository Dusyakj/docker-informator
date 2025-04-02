package tools

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"project8/internal/models"
)

type TagsResponse struct {
	Name string   `json:"name"`
	Tags []string `json:"tags"`
}

type OCIImageIndex struct {
	Manifests []ManifestEntry `json:"manifests"`
}

type ManifestEntry struct {
	Digest   string   `json:"digest"`
	Platform Platform `json:"platform"`
}

type Platform struct {
	Architecture string `json:"architecture"`
	OS           string `json:"os"`
}

type OCIManifest struct {
	SchemaVersion int               `json:"schemaVersion"`
	MediaType     string            `json:"mediaType"`
	Config        Layer             `json:"config"`
	Layers        []Layer           `json:"layers"`
	Annotations   map[string]string `json:"annotations,omitempty"`
}

type Layer struct {
	MediaType string `json:"mediaType"`
	Digest    string `json:"digest"`
	Size      int64  `json:"size"`
}

func analyzeManifest(manifest OCIManifest) (layersCount int, totalSize int64) {
	layersCount = len(manifest.Layers)

	for _, layer := range manifest.Layers {
		totalSize += layer.Size
	}
	return layersCount, totalSize
}

func contains(slice []string, item string) bool {
	for _, v := range slice {
		if v == item {
			return true
		}
	}
	return false
}

func GetProperty(info *models.Information, information models.Data, digest string) error {
	url := fmt.Sprintf("https://%s/v2/%s/blobs/%s", information.Repository, information.Name, digest)

	resp, err := http.Get(url)
	if err != nil {
		return fmt.Errorf("error durong request: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("ошибка HTTP: %d", resp.StatusCode)
	}

	var oci OCIManifest
	if err := json.NewDecoder(resp.Body).Decode(&oci); err != nil {
		return fmt.Errorf("failed to parse JSON: %v", err)
	}

	layersCount, totalSize := analyzeManifest(oci)
	info.LayersCount = layersCount
	info.TotalSize = totalSize

	return nil
}

func GetManifest(information models.Data) (string, error) {
	url := fmt.Sprintf("https://%s/v2/%s/manifests/%s", information.Repository, information.Name, information.Tag)

	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return "", fmt.Errorf("error creating request: %v", err)
	}

	req.Header.Add("Accept", "application/vnd.oci.image.index.v1+json")
	req.Header.Add("Accept", "application/vnd.oci.image.manifest.v1+json")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return "", fmt.Errorf("error making request: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return "", fmt.Errorf("request failed with status %d: %s", resp.StatusCode, body)
	}

	var oci OCIImageIndex
	if err := json.NewDecoder(resp.Body).Decode(&oci); err != nil {
		return "", fmt.Errorf("failed to parse JSON: %v", err)
	}

	for _, manifest := range oci.Manifests {
		if manifest.Digest != "" && manifest.Platform.Architecture == "amd64" && manifest.Platform.OS == "linux" {
			return manifest.Digest, nil
		}
	}

	return "", fmt.Errorf("no suitable manifest found")
}

func ValidateTag(information models.Data) error {
	url := fmt.Sprintf("https://%s/v2/%s/tags/list", information.Repository, information.Name)

	resp, err := http.Get(url)
	if err != nil {
		return fmt.Errorf("error durong request: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("ошибка HTTP: %d, %s", resp.StatusCode, information.Repository)
	}

	var tags TagsResponse
	if err := json.NewDecoder(resp.Body).Decode(&tags); err != nil {
		return fmt.Errorf("failed to parse JSON: %v", err)
	}

	if contains(tags.Tags, information.Tag) {
		return nil
	} else {
		return fmt.Errorf("incorrect tag")
	}
}
