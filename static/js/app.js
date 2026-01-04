const API_BASE = '/api';

let currentCategory = 'সব';

// Load categories
async function loadCategories() {
    try {
        const response = await fetch(`${API_BASE}/categories`);
        const categories = await response.json();
        
        const categoriesDiv = document.getElementById('categories');
        categoriesDiv.innerHTML = '';
        
        categories.forEach(cat => {
            const btn = document.createElement('button');
            btn.className = 'category-btn' + (cat === 'সব' ? ' active' : '');
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
    loadArticles();
}

// Load featured articles
async function loadFeatured() {
    try {
        const response = await fetch(`${API_BASE}/featured`);
        const articles = await response.json();
        
        const featuredDiv = document.getElementById('featured');
        featuredDiv.innerHTML = '';
        
        if (!articles || articles.length === 0) {
            featuredDiv.innerHTML = '<p>কোনো খবর নেই</p>';
            return;
        }
        
        articles.forEach(article => {
            const div = document.createElement('div');
            div.className = 'featured-article';
            div.onclick = () => openArticle(article.id);
            
            const date = new Date(article.created).toLocaleDateString('bn-BD', {
                year: 'numeric',
                month: 'long',
                day: 'numeric'
            });
            
            div.innerHTML = `
                <h3>${article.title}</h3>
                <div class="meta">${article.author} | ${date}</div>
                <div class="preview">${article.content.substring(0, 150)}...</div>
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
            
            const date = new Date(article.created).toLocaleDateString('bn-BD', {
                year: 'numeric',
                month: 'long',
                day: 'numeric'
            });
            
            div.innerHTML = `
                <div class="article-image">📰</div>
                <div class="article-content">
                    <span class="category">${article.category}</span>
                    <h3>${article.title}</h3>
                    <p class="summary">${article.content.substring(0, 100)}...</p>
                    <div class="meta">
                        <span>${article.author}</span>
                        <span>${date}</span>
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
        
        const date = new Date(article.created).toLocaleDateString('bn-BD', {
            year: 'numeric',
            month: 'long',
            day: 'numeric',
            hour: '2-digit',
            minute: '2-digit'
        });
        
        const detailDiv = document.getElementById('articleDetail');
        detailDiv.innerHTML = `
            <div class="article-detail">
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
    let searchTimeout;
    searchInput.addEventListener('input', () => {
        clearTimeout(searchTimeout);
        searchTimeout = setTimeout(() => {
            loadArticles();
        }, 500);
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
