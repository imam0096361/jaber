document.addEventListener('DOMContentLoaded', function () {
    const API_BASE = '/api';
    console.log('Admin script loaded');

    const form = document.getElementById('articleForm');
    if (!form) {
        console.error('Form not found!');
        return;
    }

    // Handle form submission
    form.addEventListener('submit', async function (e) {
        e.preventDefault();
        console.log('Form submitted');

        const submitBtn = form.querySelector('.btn-submit');
        const originalBtnText = submitBtn.textContent;
        submitBtn.disabled = true;
        submitBtn.textContent = 'অপেক্ষা করুন...';

        try {
            const imageFile = document.getElementById('image').files[0];
            let imageUrl = null;

            // Upload image if selected
            if (imageFile) {
                console.log('Uploading image...');
                const imageFormData = new FormData();
                imageFormData.append('image', imageFile);

                const uploadResponse = await fetch(`${API_BASE}/upload-image`, {
                    method: 'POST',
                    body: imageFormData
                });

                if (!uploadResponse.ok) {
                    throw new Error('ছবি আপলোড ব্যর্থ হয়েছে');
                }

                const uploadResult = await uploadResponse.json();
                imageUrl = uploadResult.url;
                console.log('Image uploaded:', imageUrl);
            }

            // Prepare article data
            const articleData = {
                title: document.getElementById('title').value,
                content: document.getElementById('content').value,
                category: document.getElementById('category').value,
                author: document.getElementById('author').value,
                featured: document.getElementById('featured').checked,
                image: imageUrl
            };

            console.log('Sending article data:', articleData);

            // Add article
            const response = await fetch(`${API_BASE}/add-article`, {
                method: 'POST',
                headers: {
                    'Content-Type': 'application/json'
                },
                body: JSON.stringify(articleData)
            });

            const result = await response.json();

            if (!response.ok) {
                throw new Error(result.message || 'খবর যোগ করতে ব্যর্থ হয়েছে');
            }

            console.log('Success:', result);
            alert('খবর সফলভাবে যোগ করা হয়েছে!');

            // Reset form
            form.reset();
            document.getElementById('imagePreview').innerHTML = '';

        } catch (error) {
            console.error('Error:', error);
            alert('ত্রুটি: ' + error.message);
        } finally {
            submitBtn.disabled = false;
            submitBtn.textContent = originalBtnText;
        }
    });

    // Image preview
    const imageInput = document.getElementById('image');
    if (imageInput) {
        imageInput.addEventListener('change', function (e) {
            const file = e.target.files[0];
            const preview = document.getElementById('imagePreview');

            if (file) {
                const reader = new FileReader();
                reader.onload = function (e) {
                    preview.innerHTML = `<img src="${e.target.result}" alt="Preview" style="max-width: 200px; max-height: 200px;">`;
                };
                reader.readAsDataURL(file);
            } else {
                preview.innerHTML = '';
            }
        });
    }
});