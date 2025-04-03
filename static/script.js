document.addEventListener('DOMContentLoaded', function () {
    const uploadForm = document.getElementById('uploadForm');
    const fileInput = document.getElementById('fileInput');
    const fileList = document.getElementById('fileList');
    const noFilesMessage = document.getElementById('noFilesMessage');

    uploadForm.addEventListener('submit', function(e) {
        e.preventDefault();
        const file = fileInput.files[0];
        if (!file) return;

        const formData = new FormData();
        formData.append('file', file);

        // Create progress bar container if it doesn't exist
        let progressContainer = document.getElementById('progressContainer');
        if (!progressContainer) {
            progressContainer = document.createElement('div');
            progressContainer.id = 'progressContainer';
            progressContainer.className = 'progress-container';
            progressContainer.style.marginTop = '15px';
            progressContainer.style.display = 'none';
            
            // Insert after the file input's parent div
            const fileInputDiv = fileInput.parentElement;
            fileInputDiv.insertAdjacentElement('afterend', progressContainer);
        }
        
        // Show the progress container
        progressContainer.style.display = 'block';

        // Create progress bar elements
        const progressBar = document.createElement('div');
        progressBar.className = 'progress-bar';
        const progressText = document.createElement('div');
        progressText.className = 'progress-text';
        progressContainer.innerHTML = '';
        progressContainer.appendChild(progressBar);
        progressContainer.appendChild(progressText);

        const xhr = new XMLHttpRequest();
        xhr.open('POST', '/upload', true);

        // Track upload progress
        xhr.upload.onprogress = function(e) {
            if (e.lengthComputable) {
                const percent = Math.round((e.loaded / e.total) * 100);
                progressBar.style.width = percent + '%';
                progressText.textContent = percent + '%';
            }
        };

        xhr.onload = function() {
            if (xhr.status === 200) {
                const data = JSON.parse(xhr.responseText);
                progressText.textContent = 'Upload complete!';
                setTimeout(() => {
                    progressContainer.style.display = 'none';
                    fileInput.value = '';
                    loadFiles();
                }, 1000);
            } else {
                progressText.textContent = 'Upload failed';
                console.error('Error:', xhr.statusText);
            }
        };

        xhr.onerror = function() {
            progressText.textContent = 'Upload failed';
            console.error('Error:', xhr.statusText);
        };

        xhr.send(formData);
});

// Powered by Go

    loadFiles();

    function loadFiles() {
        fetch('/files')
        .then(response => response.json())
        .then(files => {
            if (files.length === 0) {
                noFilesMessage.style.display = 'block';
                fileList.innerHTML = '';
                return;
            }

            noFilesMessage.style.display = 'none';
            let html = '';
            files.forEach(file => {
                const icon = getFileIcon(file.name);
                html += `
                    <div class="file-item">
                        <div class="file-info">
                            <i class="${icon} file-icon"></i>
                            <div>
                                <div class="fw-bold">${file.name}</div>
                                <small class="text-muted">${file.size} • ${file.date}</small>
                            </div>
                        </div>
                        <div class="file-actions">
                            <a href="/download?path=${encodeURIComponent(file.path)}" class="btn btn-sm btn-outline-primary">
                                <i class="fas fa-download"></i> Download
                            </a>
                            <button class="btn btn-sm btn-outline-danger" onclick="deleteFile('${encodeURIComponent(file.path)}')">
                                <i class="fas fa-trash"></i>
                            </button>
                        </div>
                    </div>
                `;
            });
            fileList.innerHTML = html;
        })
        .catch(error => console.error('Error loading files:', error));
    }

    function getFileIcon(filename) {
        const ext = filename.split('.').pop().toLowerCase();
        const icons = {
            pdf: 'fas fa-file-pdf',
            doc: 'fas fa-file-word',
            docx: 'fas fa-file-word',
            xls: 'fas fa-file-excel',
            xlsx: 'fas fa-file-excel',
            ppt: 'fas fa-file-powerpoint',
            pptx: 'fas fa-file-powerpoint',
            jpg: 'fas fa-file-image',
            jpeg: 'fas fa-file-image',
            png: 'fas fa-file-image',
            gif: 'fas fa-file-image',
            txt: 'fas fa-file-alt',
            zip: 'fas fa-file-archive',
            rar: 'fas fa-file-archive',
            mp3: 'fas fa-file-audio',
            mp4: 'fas fa-file-video'
        };
        return icons[ext] || 'fas fa-file';
    }

    window.deleteFile = function(path) {
        if (confirm('Are you sure you want to delete this file?')) {
            fetch('/delete', {
                method: 'POST',
                headers: {'Content-Type': 'application/x-www-form-urlencoded'},
                body: `path=${encodeURIComponent(path)}`
            })
            .then(response => response.json())
            .then(data => data.status === 'success' && loadFiles())
            .catch(error => console.error('Error deleting file:', error));
        }
    };
});
