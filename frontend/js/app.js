const API_BASE = '/api';

let currentCategory = 'সব';

// Utility to format date safely
function formatDate(dateString, includeTime = false) {
    if (!dateString) return '';
    const date = new Date(dateString);
    // Check if date is invalid or year 1
    if (isNaN(date.getTime()) || date.getFullYear() <= 1) {
        return '';
    }

    const options = {
        year: 'numeric',
        month: 'long',
        day: 'numeric'
    };
    if (includeTime) {
        options.hour = '2-digit';
        options.minute = '2-digit';
    }

    return date.toLocaleDateString('bn-BD', options);
}

// Utility to truncate text safely
function truncateText(text, length) {
    if (!text) return '';
    if (text.length <= length) return text;
    return text.substring(0, length).trim() + '...';
}

// Load categories
async function loadCategories() {
    try {
        const response = await fetch(`${API_BASE}/categories`);
        const categories = await response.json();

        const categoriesDiv = document.getElementById('categories');
        categoriesDiv.innerHTML = '';

        // Prepend "সব" category
        const allCategories = ['সব', ...categories];

        allCategories.forEach(cat => {
            const btn = document.createElement('button');
            btn.className = 'category-btn' + (cat === currentCategory ? ' active' : '');
            btn.textContent = cat;
            btn.onclick = () => selectCategory(cat);
            categoriesDiv.appendChild(btn);
        });
    } catch (error) {
        console.error('Error loading categories:', error);
    }
}

// Select category
function selectCategory(category) {
    currentCategory = category;
    document.querySelectorAll('.category-btn').forEach(btn => {
        btn.classList.remove('active');
        if (btn.textContent === category) {
            btn.classList.add('active');
        }
    });

    const featuredSection = document.querySelector('.featured-section');
    // Always show the featured section and reload it with the category filter
    if (featuredSection) {
        featuredSection.style.display = 'block';
        loadFeatured();
    }

    loadArticles();
}

// Load featured articles
async function loadFeatured() {
    console.log('Loading featured for category:', currentCategory);
    const featuredDiv = document.getElementById('featured');
    featuredDiv.innerHTML = '<div class="loading">লোড হচ্ছে...</div>';

    try {
        let url = `${API_BASE}/featured?limit=6&_t=${Date.now()}`;
        if (currentCategory !== 'সব') {
            url += `&category=${encodeURIComponent(currentCategory)}`;
        }

        console.log('Fetching featured from:', url);
        const response = await fetch(url);
        const articles = await response.json();

        console.log('Received featured articles:', articles.length);
        featuredDiv.innerHTML = '';

        if (!articles || articles.length === 0) {
            featuredDiv.innerHTML = '<p>কোনো খবর নেই</p>';
            return;
        }

        articles.forEach(article => {
            const div = document.createElement('div');
            div.className = 'featured-article';
            div.onclick = () => openArticle(article.id);

            const date = formatDate(article.created);
            // Use placeholder if image is null or empty to prevent 404 errors
            const imageUrl = article.image && article.image.trim() !== ''
                ? article.image
                : 'https://images.unsplash.com/photo-1585007600263-ad12200a09a5?auto=format&fit=crop&q=80&w=800&h=450';

            // Strictly truncate title if it's too long (headline)
            const displayTitle = truncateText(article.title, 60);

            div.innerHTML = `
                <div class="featured-image">
                    <img src="${imageUrl}" 
                         alt="${article.title}" 
                         onerror="this.onerror=null; this.src='https://images.unsplash.com/photo-1585007600263-ad12200a09a5?auto=format&fit=crop&q=80&w=800&h=450'; this.parentElement.classList.add('broken');">
                </div>
                <h3>${displayTitle}</h3>
                <div class="meta">
                    <span class="author">👤 ${article.author}</span>
                    ${date ? '<span class="date">📅 ' + date + '</span>' : ''}
                </div>
                <div class="preview">${truncateText(article.content, 120)}</div>
                <div class="read-more-btn">আরও পড়ুন ➝</div>
            `;

            featuredDiv.appendChild(div);
        });
    } catch (error) {
        console.error('Error loading featured articles:', error);
    }
}

// Load articles
async function loadArticles() {
    try {
        const searchQuery = document.getElementById('searchInput').value;
        let url = `${API_BASE}/articles`;

        const params = new URLSearchParams();
        params.append('limit', '15');
        if (currentCategory !== 'সব') {
            params.append('category', currentCategory);
        }
        if (searchQuery) {
            params.append('search', searchQuery);
        }

        if (params.toString()) {
            url += '?' + params.toString();
        }

        const response = await fetch(url);
        const articles = await response.json();

        const articlesDiv = document.getElementById('articles');
        articlesDiv.innerHTML = '';

        if (!articles || articles.length === 0) {
            articlesDiv.innerHTML = '<p>কোনো খবর নেই</p>';
            return;
        }

        articles.forEach(article => {
            const div = document.createElement('div');
            div.className = 'article-card';
            div.onclick = () => openArticle(article.id);

            const date = formatDate(article.created);
            // Use placeholder if image is null or empty to prevent 404 errors
            const imageUrl = article.image && article.image.trim() !== ''
                ? article.image
                : 'https://images.unsplash.com/photo-1504711331083-9c895941bf81?auto=format&fit=crop&q=80&w=400&h=225';
            const displayTitle = truncateText(article.title, 50);

            div.innerHTML = `
                <div class="article-image">
                    <img src="${imageUrl}" 
                         alt="${article.title}" 
                         loading="lazy"
                         onerror="this.onerror=null; this.src='https://images.unsplash.com/photo-1504711331083-9c895941bf81?auto=format&fit=crop&q=80&w=400&h=225'; this.parentElement.classList.add('broken');">
                </div>
                <div class="article-content">
                    <span class="category">${article.category}</span>
                    <h3>${displayTitle}</h3>
                    <p class="summary">${truncateText(article.content, 80)}</p>
                    <div class="read-more-btn">আরও পড়ুন ➝</div>
                    <div class="meta">
                        <span>👤 ${article.author}</span>
                        ${date ? `<span>📅 ${date}</span>` : ''}
                    </div>
                </div>
            `;

            articlesDiv.appendChild(div);
        });
    } catch (error) {
        console.error('Error loading articles:', error);
    }
}

// Open article detail
async function openArticle(id) {
    try {
        const response = await fetch(`${API_BASE}/article?id=${id}`);
        const article = await response.json();

        const date = formatDate(article.created, true);

        const detailDiv = document.getElementById('articleDetail');
        detailDiv.innerHTML = `
            <div class="article-detail">
                ${article.image ? `<img src="${article.image}" alt="${article.title}" style="width: 100%; max-height: 400px; object-fit: cover; border-radius: 10px; margin-bottom: 20px;">` : ''}
                <h1>${article.title}</h1>
                <div class="detail-meta">
                    <strong>${article.category}</strong> | লেখক: ${article.author} | ${date}
                </div>
                <div class="content">${article.content}</div>
            </div>
        `;

        const modal = document.getElementById('articleModal');
        modal.classList.add('show');
    } catch (error) {
        console.error('Error loading article:', error);
    }
}

// Close modal
function closeModal() {
    document.getElementById('articleModal').classList.remove('show');
}

// Event listeners
document.addEventListener('DOMContentLoaded', () => {
    loadCategories();
    loadFeatured();
    loadArticles();

    // Search functionality
    const searchInput = document.getElementById('searchInput');
    searchInput.addEventListener('keypress', (e) => {
        if (e.key === 'Enter') {
            const query = searchInput.value.trim();
            const featuredSection = document.querySelector('.featured-section');

            // If searching, hide featured section to show only search results
            if (query) {
                featuredSection.style.display = 'none';
            } else if (currentCategory === 'সব') {
                featuredSection.style.display = 'block';
                loadFeatured();
            }

            loadArticles();
        }
    });

    // Modal close
    const modal = document.getElementById('articleModal');
    const closeBtn = document.querySelector('.close');

    closeBtn.onclick = closeModal;
    window.onclick = (event) => {
        if (event.target === modal) {
            closeModal();
        }
    };
});
