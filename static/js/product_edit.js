document.addEventListener('DOMContentLoaded', () => {
    const productId = window.location.pathname.split('/').pop();
    
    const imageUploadInput = document.getElementById('imageUploadInput');
    const triggerUploadButton = document.getElementById('triggerUploadButton');
    const productImageView = document.getElementById('productImageView');
    const uploadStatus = document.getElementById('upload-status');

    const cameraModal = document.getElementById('cameraModal');
    const triggerCameraButton = document.getElementById('triggerCameraButton');
    const cancelCameraButton = document.getElementById('cancelCameraButton');
    const videoFeed = document.getElementById('videoFeed');
    const captureButton = document.getElementById('captureButton');

    let stream;

    // --- File Upload Logic ---
    if (triggerUploadButton) {
        triggerUploadButton.addEventListener('click', () => {
            imageUploadInput.click();
        });
    }

    if (imageUploadInput) {
        imageUploadInput.addEventListener('change', (event) => {
            const file = event.target.files[0];
            if (file) {
                uploadFile(file);
            }
        });
    }

    async function uploadFile(file) {
        const formData = new FormData();
        formData.append('image', file);

        setStatus('Uploading...', 'text-primary');

        try {
            const response = await fetch(`/api/products/${productId}/image`, {
                method: 'POST',
                body: formData,
            });

            if (!response.ok) {
                const errorData = await response.text();
                throw new Error(`Upload failed: ${errorData}`);
            }

            const result = await response.json();
            
            // Update image view with a cache-busting query parameter
            productImageView.src = result.imageURL + '?t=' + new Date().getTime();
            
            setStatus('Upload successful!', 'text-success');

        } catch (error) {
            console.error('Upload error:', error);
            setStatus(error.message, 'text-danger');
        }
    }

    function setStatus(message, className) {
        uploadStatus.textContent = message;
        uploadStatus.className = `form-text mt-2 ${className}`;
    }


    // --- Camera Logic ---
    if (triggerCameraButton) {
        triggerCameraButton.addEventListener('click', async () => {
            if (!('mediaDevices' in navigator && 'getUserMedia' in navigator.mediaDevices)) {
                alert('Camera API is not available in your browser.');
                return;
            }
            cameraModal.style.display = 'block';
            try {
                stream = await navigator.mediaDevices.getUserMedia({ video: { facingMode: 'environment' } });
                videoFeed.srcObject = stream;
            } catch (err) {
                console.error("Error accessing camera:", err);
                alert('Could not access the camera. Please ensure permissions are granted.');
                cameraModal.style.display = 'none';
            }
        });
    }

    function stopCamera() {
        if (stream) {
            stream.getTracks().forEach(track => track.stop());
        }
        cameraModal.style.display = 'none';
    }

    if (cancelCameraButton) {
        cancelCameraButton.addEventListener('click', stopCamera);
    }

    if (captureButton) {
        captureButton.addEventListener('click', () => {
            const canvas = document.createElement('canvas');
            canvas.width = videoFeed.videoWidth;
            canvas.height = videoFeed.videoHeight;
            const context = canvas.getContext('2d');
            context.drawImage(videoFeed, 0, 0, canvas.width, canvas.height);
            
            stopCamera();

            canvas.toBlob((blob) => {
                const filename = `capture-${new Date().toISOString()}.jpg`;
                const file = new File([blob], filename, { type: 'image/jpeg' });
                uploadFile(file);
            }, 'image/jpeg');
        });
    }

    // Close modal if clicking outside of it
    window.addEventListener('click', (event) => {
        if (event.target == cameraModal) {
            stopCamera();
        }
    });
});