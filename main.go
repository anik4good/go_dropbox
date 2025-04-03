package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"time"
)

const (
	baseURL     = "https://api.dropboxapi.com/2"
	contentURL  = "https://content.dropboxapi.com/2"
	accessToken = "sl.u.AFp2-RNslxG0LUsUjb5H26ABfT_f7Eg2fDEntlwVjOMD0_-raTN4Hfq6rvILPnnqu5-Lm911SXgVdbFLZqOhSNK8rtwDE_23yU30PPlYQoo2yAPAESde0wFhPfguoY9Pltp6exJjyPBq_Xi9lsFAHJ9ij6vRKpV_R-Arj8lIAt8DAqEU3yAdleT3xQ8NyRhAbK0gHYDRIxgxIJFmyXbvpFXLj2Zl7M350u2b496I9PsSUfa-6jI5LDAceyI_t4Uo-zbuCnYEPfWzPG9QZez-ooYz55SvW9cXz2EGXQR4go5U0-EerYL08Me0vtw0c7TtFWGMPVzZ-7gRtVCkNkm2xlH8Dca4mze-3H1bhh0WTDrk4pqccZpNtwWYfTJmvX7jEWRH7uJG335F-krhIFvosKhemfacy_0DryJnK_5O96_lcgXwKw7AikJCJiD-IKVjUNT7s3r1EHgtprKb06GRXdPsiMYSuYw54i1fd9nW5ItvO_0by0YHGNdmIxIS00uadJcbnzc35ULVHfYlNxC-djBm1H-JIz9tpvmWcUSsKP5sTCocH8pcO1K8NASr7Sm6BNgrlp1ITjTA3fwOvGpjnXkEspnYDH54YChjmrPIRasNeIydvH2qmxEhIESV6UfDrqhGGuJ8g5JUIQRDkQt_hrLhQL49qcqqpB1LoudWqzAVoVSEh_AiTh-pfnHoXxMsoEECK6vXmuFm58sNj91mvhjWduan1XcohtRUBseduWawsSspu7kl9s5U5iPzLUDdIHA4d6AVXIYwMRC4fBl92aB4djsMjTGxbdQLGkebF-GrjeeXpbGAbngZvMlJO9ro-eKLxuKqqCjDevJpuY0yIaKTWGEzyFIAFUGeSsYPQvdNj29lXwb5WCxP5qm1tveuVkNiG2HD0UpHT7FAEjpttoKFZCNLk1FH7DWkl7OccXDbpd_kbBECETeImFa9yCGvxlhPVUvSg8i0__FFKAd75z0U9gQ96PbF4fswB2VRDm2eR01KQnrsdGE-N5d4qcY19X8dN8Fuf5w9_4JPqZuJwewitfTfxeWDH1beMyoB6eiwXDxOvb2F9j3hNkgYfELnPvnXF8ywgeR8BezSCn0WOvb2kmQfrFHboDs334u1wzgy8Fl8jUBKiGOtEDjCQ-rhLB-NEKfh9KPqEawJ03gQNIGHAmvMm3iw0scM7l1k4LrKG-vkq3V9UvNXBWsCAP4O5aumb0OuoLCTBYcOjdJbP92H-iR9AC3-KTjrmJHhVHRDJllZ_PdHS70FHrZhXls-hWkeEAp42YL1W6OBFtLLMU7xIjqAUb_n3fDNscj2BcQo6fMZVL-gnpAKtNc0i83vYZh4PbL9QiwqOuxp4VVNm3nbZGdiAWFu3hrIyqQXPuHlLqMr6wGpbz9kochZeUPJTeXZ0wnMHGWAM2bdDDHNMwTDW6XmDOP-BR_sIVLumuvU1w" // Replace with your actual token
	port        = ":8080"
)

type FileMetadata struct {
	Name           string `json:"name"`
	Path           string `json:"path_display"`
	Size           uint64 `json:"size"`
	ClientModified string `json:"client_modified"`
	ContentHash    string `json:"content_hash"`
}

func main() {
	// Serve static files
	fs := http.FileServer(http.Dir("./static"))
	http.Handle("/static/", http.StripPrefix("/static/", fs))

	// Routes
	http.HandleFunc("/", handleIndex)
	http.HandleFunc("/upload", handleUpload)
	http.HandleFunc("/files", handleFiles)
	http.HandleFunc("/download", handleDownload)
	http.HandleFunc("/delete", handleDelete)

	log.Printf("Server running on http://localhost%s\n", port)
	log.Fatal(http.ListenAndServe(port, nil))
}

func handleIndex(w http.ResponseWriter, r *http.Request) {
    html := `<!DOCTYPE html>
<html lang="en">
<head>
    <meta charset="UTF-8">
    <meta name="viewport" content="width=device-width, initial-scale=1.0">
    <title>Dropbox File Manager</title>
    <link href="https://cdn.jsdelivr.net/npm/bootstrap@5.3.0/dist/css/bootstrap.min.css" rel="stylesheet">
    <link rel="stylesheet" href="https://cdnjs.cloudflare.com/ajax/libs/font-awesome/6.4.0/css/all.min.css">
    <link rel="stylesheet" href="https://cdnjs.cloudflare.com/ajax/libs/animate.css/4.1.1/animate.min.css">
    <link rel="stylesheet" href="/static/style.css">
</head>
<body>
    <div class="container py-5">
        <div class="row justify-content-center">
            <div class="col-lg-10">
                <div class="card shadow-lg">
                    <div class="card-header py-3" style="background: linear-gradient(135deg, #667eea 0%, #764ba2 100%);">
                        <div class="d-flex justify-content-between align-items-center">
                            <h4 class="m-0 text-white">
                                <i class="fas fa-cloud-upload-alt me-2"></i>
                                Dropbox File Manager
                            </h4>
                            <span class="badge bg-white text-primary" id="storageInfo">0 MB used</span>
                        </div>
                    </div>
                    
                    <div class="card-body">
                        <!-- Upload Section -->
                        <div class="upload-section mb-5">
                            <h5 class="mb-4 text-primary">
                                <i class="fas fa-upload me-2"></i>
                                Upload Files
                            </h5>
                            <form id="uploadForm" enctype="multipart/form-data">
                                <div class="mb-3">
                                    <input type="file" id="fileInput" name="file" class="form-control form-control-lg">
                                </div>
                                <button type="submit" class="btn btn-primary btn-lg w-100">
                                    <i class="fas fa-cloud-upload-alt me-2"></i>
                                    Upload to Cloud
                                </button>
                            </form>
                        </div>

                        <!-- Files List -->
                        <div class="files-section">
                            <h5 class="mb-4 text-primary">
                                <i class="fas fa-folder-open me-2"></i>
                                Your Files
                            </h5>
                            <div class="card shadow-sm">
                                <div class="card-body p-0">
                                    <div class="file-list" id="fileList">
                                        <!-- Files will be loaded here -->
                                        <div class="text-center py-5 text-muted" id="noFilesMessage">
                                            <i class="fas fa-cloud fa-3x mb-3 text-secondary"></i>
                                            <h5 class="mb-2">Your cloud storage is empty</h5>
                                            <p>Upload your first file to get started</p>
                                        </div>
                                    </div>
                                </div>
                            </div>
                        </div>
                    </div>
                    
                    <div class="card-footer text-center py-3 bg-light">
                        <small class="text-muted">
                            <i class="fas fa-code me-1"></i>
                            Powered by Go
                        </small>
                    </div>
                </div>
            </div>
        </div>
    </div>

    <!-- Loading Spinner (hidden by default) -->
    <div id="loadingSpinner" class="position-fixed top-0 start-0 w-100 h-100 d-flex justify-content-center align-items-center" style="background: rgba(0,0,0,0.5); z-index: 9999; display: none !important;">
        <div class="spinner-border text-primary" style="width: 3rem; height: 3rem;" role="status">
            <span class="visually-hidden">Loading...</span>
        </div>
    </div>

    <script src="https://cdn.jsdelivr.net/npm/bootstrap@5.3.0/dist/js/bootstrap.bundle.min.js"></script>
    <script src="/static/script.js"></script>
</body>
</html>`

	w.Header().Set("Content-Type", "text/html")
	w.Write([]byte(html))
}

func handleUpload(w http.ResponseWriter, r *http.Request) {
	if r.Method != "POST" {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	// Parse the multipart form
	err := r.ParseMultipartForm(10 << 20) // 10 MB max
	if err != nil {
		http.Error(w, "Error parsing form", http.StatusBadRequest)
		return
	}

	// Get the file from the form
	file, handler, err := r.FormFile("file")
	if err != nil {
		http.Error(w, "Error retrieving file", http.StatusBadRequest)
		return
	}
	defer file.Close()

	// Upload to Dropbox
	url := contentURL + "/files/upload"
	client := &http.Client{}

	apiArg := map[string]interface{}{
		"path":            "/" + handler.Filename,
		"mode":            "add",
		"autorename":      true,
		"mute":            false,
		"strict_conflict": false,
	}

	apiArgBytes, err := json.Marshal(apiArg)
	if err != nil {
		http.Error(w, "Error preparing upload", http.StatusInternalServerError)
		return
	}

	req, err := http.NewRequest("POST", url, file)
	if err != nil {
		http.Error(w, "Error creating request", http.StatusInternalServerError)
		return
	}

	req.Header.Set("Authorization", "Bearer "+accessToken)
	req.Header.Set("Dropbox-API-Arg", string(apiArgBytes))
	req.Header.Set("Content-Type", "application/octet-stream")

	resp, err := client.Do(req)
	if err != nil {
		http.Error(w, "Error uploading file", http.StatusInternalServerError)
		return
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		http.Error(w, fmt.Sprintf("Upload failed: %s", string(body)), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Write([]byte(`{"status": "success", "message": "File uploaded successfully"}`))
}

func handleFiles(w http.ResponseWriter, r *http.Request) {
	if r.Method != "GET" {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	log.Println("Attempting to list Dropbox files...")

	// List files from Dropbox
	url := baseURL + "/files/list_folder"
	client := &http.Client{}

	data := map[string]string{
		"path": "",
	}

	jsonData, err := json.Marshal(data)
	if err != nil {
		log.Printf("Error marshaling request data: %v", err)
		http.Error(w, "Error preparing request", http.StatusInternalServerError)
		return
	}

	req, err := http.NewRequest("POST", url, bytes.NewBuffer(jsonData))
	if err != nil {
		log.Printf("Error creating request: %v", err)
		http.Error(w, "Error creating request", http.StatusInternalServerError)
		return
	}

	req.Header.Set("Authorization", "Bearer "+accessToken)
	req.Header.Set("Content-Type", "application/json")

	log.Println("Sending request to Dropbox API...")
	resp, err := client.Do(req)
	if err != nil {
		log.Printf("Error making request to Dropbox API: %v", err)
		http.Error(w, "Error listing files", http.StatusInternalServerError)
		return
	}
	defer resp.Body.Close()

	body, _ := io.ReadAll(resp.Body)
	log.Printf("Dropbox API response status: %d", resp.StatusCode)
	log.Printf("Dropbox API response body: %s", string(body))

	if resp.StatusCode != http.StatusOK {
		log.Printf("Dropbox API error: %s", string(body))
		http.Error(w, fmt.Sprintf("List failed: %s", string(body)), http.StatusInternalServerError)
		return
	}

	var result struct {
		Entries []FileMetadata `json:"entries"`
	}
	err = json.Unmarshal(body, &result)
	if err != nil {
		log.Printf("Error parsing Dropbox response: %v", err)
		http.Error(w, "Error parsing response", http.StatusInternalServerError)
		return
	}

	// Filter out folders and format response
	var files []map[string]interface{}
	for _, entry := range result.Entries {
		if entry.ContentHash != "" { // This indicates it's a file, not a folder
			files = append(files, map[string]interface{}{
				"name": entry.Name,
				"path": entry.Path,
				"size": formatFileSize(entry.Size),
				"date": formatDate(entry.ClientModified),
			})
		}
	}

	json.NewEncoder(w).Encode(files)
}

func handleDownload(w http.ResponseWriter, r *http.Request) {
	if r.Method != "GET" {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	path := r.URL.Query().Get("path")
	if path == "" {
		http.Error(w, "Path parameter is required", http.StatusBadRequest)
		return
	}

	// Get temporary download link
	tempLink, err := getTemporaryLink(path)
	if err != nil {
		http.Error(w, "Error getting download link", http.StatusInternalServerError)
		return
	}

	// Redirect to the temporary link
	http.Redirect(w, r, tempLink, http.StatusFound)
}

func handleDelete(w http.ResponseWriter, r *http.Request) {
	if r.Method != "POST" {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	path := r.FormValue("path")
	if path == "" {
		http.Error(w, "Path parameter is required", http.StatusBadRequest)
		return
	}

	// Delete from Dropbox
	url := baseURL + "/files/delete_v2"
	client := &http.Client{}

	data := map[string]string{
		"path": path,
	}

	jsonData, err := json.Marshal(data)
	if err != nil {
		http.Error(w, "Error preparing request", http.StatusInternalServerError)
		return
	}

	req, err := http.NewRequest("POST", url, bytes.NewBuffer(jsonData))
	if err != nil {
		http.Error(w, "Error creating request", http.StatusInternalServerError)
		return
	}

	req.Header.Set("Authorization", "Bearer "+accessToken)
	req.Header.Set("Content-Type", "application/json")

	resp, err := client.Do(req)
	if err != nil {
		http.Error(w, "Error deleting file", http.StatusInternalServerError)
		return
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		http.Error(w, fmt.Sprintf("Delete failed: %s", string(body)), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Write([]byte(`{"status": "success", "message": "File deleted successfully"}`))
}

func getTemporaryLink(path string) (string, error) {
	url := baseURL + "/files/get_temporary_link"
	client := &http.Client{}

	data := map[string]string{
		"path": path,
	}

	jsonData, err := json.Marshal(data)
	if err != nil {
		return "", err
	}

	req, err := http.NewRequest("POST", url, bytes.NewBuffer(jsonData))
	if err != nil {
		return "", err
	}

	req.Header.Set("Authorization", "Bearer "+accessToken)
	req.Header.Set("Content-Type", "application/json")

	resp, err := client.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}

	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("API error: %s - %s", resp.Status, string(body))
	}

	var result struct {
		Link string `json:"link"`
	}
	err = json.Unmarshal(body, &result)
	if err != nil {
		return "", err
	}

	return result.Link, nil
}

func formatFileSize(size uint64) string {
	if size < 1024 {
		return fmt.Sprintf("%d B", size)
	} else if size < 1024*1024 {
		return fmt.Sprintf("%.1f KB", float64(size)/1024)
	} else if size < 1024*1024*1024 {
		return fmt.Sprintf("%.1f MB", float64(size)/(1024*1024))
	} else {
		return fmt.Sprintf("%.1f GB", float64(size)/(1024*1024*1024))
	}
}

func formatDate(dateStr string) string {
	t, err := time.Parse(time.RFC3339, dateStr)
	if err != nil {
		return dateStr
	}
	return t.Format("Jan 2, 2006 at 15:04")
}
