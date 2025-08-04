// API Configuration
const API_BASE_URL = 'http://localhost:8080/api/v1';
const WORKER_BASE_URL = 'http://localhost:8081';

// DOM Elements
const uploadArea = document.getElementById('uploadArea');
const fileInput = document.getElementById('fileInput');
const fileInfo = document.getElementById('fileInfo');
const fileName = document.getElementById('fileName');
const fileSize = document.getElementById('fileSize');
const uploadBtn = document.getElementById('uploadBtn');
const jobsList = document.getElementById('jobsList');
const refreshBtn = document.getElementById('refreshBtn');
const loadingOverlay = document.getElementById('loadingOverlay');
const toastContainer = document.getElementById('toastContainer');
const apiStatus = document.getElementById('apiStatus');
const workerStatus = document.getElementById('workerStatus');

// State
let selectedFile = null;
let jobs = [];

// Initialize
document.addEventListener('DOMContentLoaded', () => {
    initializeEventListeners();
    checkSystemStatus();
    loadJobs();
    
    // Auto-refresh jobs every 5 seconds
    setInterval(loadJobs, 5000);
    
    // Check system status every 10 seconds
    setInterval(checkSystemStatus, 10000);
});

// Event Listeners
function initializeEventListeners() {
    // Upload area events
    uploadArea.addEventListener('click', () => fileInput.click());
    uploadArea.addEventListener('dragover', handleDragOver);
    uploadArea.addEventListener('dragleave', handleDragLeave);
    uploadArea.addEventListener('drop', handleDrop);
    
    // File input change
    fileInput.addEventListener('change', handleFileSelect);
    
    // Upload button
    uploadBtn.addEventListener('click', handleUpload);
    
    // Refresh button
    refreshBtn.addEventListener('click', loadJobs);
}

// File Upload Handlers
function handleDragOver(e) {
    e.preventDefault();
    uploadArea.classList.add('dragover');
}

function handleDragLeave(e) {
    e.preventDefault();
    uploadArea.classList.remove('dragover');
}

function handleDrop(e) {
    e.preventDefault();
    uploadArea.classList.remove('dragover');
    
    const files = e.dataTransfer.files;
    if (files.length > 0) {
        handleFileSelect({ target: { files } });
    }
}

function handleFileSelect(e) {
    const file = e.target.files[0];
    if (!file) return;
    
    // Validate file type
    if (!file.name.toLowerCase().endsWith('.csv')) {
        showToast('Please select a CSV file', 'error');
        return;
    }
    
    selectedFile = file;
    fileName.textContent = file.name;
    fileSize.textContent = formatFileSize(file.size);
    fileInfo.style.display = 'flex';
}

function formatFileSize(bytes) {
    if (bytes === 0) return '0 Bytes';
    const k = 1024;
    const sizes = ['Bytes', 'KB', 'MB', 'GB'];
    const i = Math.floor(Math.log(bytes) / Math.log(k));
    return parseFloat((bytes / Math.pow(k, i)).toFixed(2)) + ' ' + sizes[i];
}

// Upload Handler
async function handleUpload() {
    if (!selectedFile) {
        showToast('Please select a file first', 'error');
        return;
    }
    
    showLoading(true);
    
    try {
        const formData = new FormData();
        formData.append('file', selectedFile);
        
        // Determine harmonization type from filename
        const harmonizationType = getHarmonizationType(selectedFile.name);
        formData.append('harmonization_type', harmonizationType);
        
        const response = await fetch(`${API_BASE_URL}/jobs/`, {
            method: 'POST',
            body: formData
        });
        
        if (!response.ok) {
            throw new Error(`Upload failed: ${response.statusText}`);
        }
        
        const result = await response.json();
        showToast('File uploaded successfully!', 'success');
        
        // Reset form
        selectedFile = null;
        fileInfo.style.display = 'none';
        fileInput.value = '';
        
        // Refresh jobs
        loadJobs();
        
    } catch (error) {
        console.error('Upload error:', error);
        showToast(`Upload failed: ${error.message}`, 'error');
    } finally {
        showLoading(false);
    }
}

function getHarmonizationType(filename) {
    const name = filename.toLowerCase();
    if (name.includes('patient')) return 'patients';
    if (name.includes('vital')) return 'vitals';
    if (name.includes('medication')) return 'medications';
    if (name.includes('lab')) return 'lab_results';
    return 'generic';
}

// Jobs Management
async function loadJobs() {
    try {
        const response = await fetch(`${API_BASE_URL}/jobs/`);
        if (!response.ok) {
            throw new Error('Failed to load jobs');
        }
        
        jobs = await response.json();
        renderJobs();
        
    } catch (error) {
        console.error('Failed to load jobs:', error);
        // Don't show toast for this error as it might be too frequent
    }
}

function renderJobs() {
    if (jobs.length === 0) {
        jobsList.innerHTML = `
            <div class="empty-state">
                <div class="empty-icon">📋</div>
                <p>No jobs yet. Upload a file to get started.</p>
            </div>
        `;
        return;
    }
    
    jobsList.innerHTML = jobs.map(job => `
        <div class="job-item">
            <div class="job-info">
                <div class="job-id">${job.id}</div>
                <div class="job-details">
                    ${job.harmonization_type} • ${job.file_size} • ${formatDate(job.created_at)}
                </div>
            </div>
            <div class="job-status">
                <span class="status-badge status-${job.status}">${job.status}</span>
                ${job.status === 'completed' && job.output_file ? 
                    `<button class="download-btn" onclick="downloadResult('${job.id}')">Download</button>` : 
                    ''
                }
            </div>
        </div>
    `).join('');
}

function formatDate(dateString) {
    const date = new Date(dateString);
    return date.toLocaleString();
}

async function downloadResult(jobId) {
    try {
        const response = await fetch(`${API_BASE_URL}/jobs/${jobId}/result`);
        if (!response.ok) {
            throw new Error('Failed to download result');
        }
        
        const blob = await response.blob();
        const url = window.URL.createObjectURL(blob);
        const a = document.createElement('a');
        a.href = url;
        a.download = `harmonized_${jobId}.csv`;
        document.body.appendChild(a);
        a.click();
        window.URL.revokeObjectURL(url);
        document.body.removeChild(a);
        
        showToast('Download started!', 'success');
        
    } catch (error) {
        console.error('Download error:', error);
        showToast(`Download failed: ${error.message}`, 'error');
    }
}

// System Status
async function checkSystemStatus() {
    // Check Python API
    try {
        const response = await fetch(`${API_BASE_URL.replace('/api/v1', '')}/health`);
        if (response.ok) {
            apiStatus.className = 'status-indicator online';
        } else {
            apiStatus.className = 'status-indicator offline';
        }
    } catch (error) {
        apiStatus.className = 'status-indicator offline';
    }
    
    // Check Go Worker
    try {
        const response = await fetch(`${WORKER_BASE_URL}/health`);
        if (response.ok) {
            workerStatus.className = 'status-indicator online';
        } else {
            workerStatus.className = 'status-indicator offline';
        }
    } catch (error) {
        workerStatus.className = 'status-indicator offline';
    }
}

// Utility Functions
function showLoading(show) {
    loadingOverlay.style.display = show ? 'flex' : 'none';
}

function showToast(message, type = 'success') {
    const toast = document.createElement('div');
    toast.className = `toast ${type}`;
    toast.textContent = message;
    
    toastContainer.appendChild(toast);
    
    // Auto-remove after 3 seconds
    setTimeout(() => {
        toast.remove();
    }, 3000);
}

// Make downloadResult available globally
window.downloadResult = downloadResult;
