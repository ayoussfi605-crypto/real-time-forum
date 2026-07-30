import { navigate } from "./navigate.js";
import { renderPost } from "./post.js";

export async function renderfeed() {
    const app = document.getElementById("app");
    
    // 1. Render the modern layout: header with a "Create Post" button, hidden modal for the form, and a posts grid
    app.innerHTML = `
        <div class="feed-container">
            <div class="feed-header">
                <div class="feed-title-section">
                    <h1>Explore Discussions</h1>
                    <p>Join the conversation with the community.</p>
                </div>
                <button id="open-create-modal" class="btn-primary create-btn">
                    <span style="font-size: 20px; font-weight: bold; margin-right: 6px; line-height: 1;">+</span>
                    Create Post
                </button>
            </div>
            
            <!-- Posts Grid -->
            <div id="posts-container" class="posts-grid">
                <div class="loading-spinner"></div>
            </div>
        </div>

        <!-- Create Post Modal -->
        <div id="create-modal-overlay" class="modal-overlay hidden">
            <div class="modal-content">
                <div class="modal-header">
                    <h2>Create a New Post</h2>
                    <button id="close-create-modal" class="close-btn">&times;</button>
                </div>
                <form id="create-post-form" class="modern-form">
                    <div class="form-group">
                        <label for="post-title">Title</label>
                        <input type="text" id="post-title" placeholder="Give your post a catchy title" required>
                    </div>
                    
                    <div class="form-group">
                        <label for="post-content">Content</label>
                        <textarea id="post-content" placeholder="What do you want to talk about?" required rows="5"></textarea>
                    </div>
                    
                    <div class="form-group">
                        <label>Select Categories</label>
                        <div id="categories-container" class="categories-selector">
                            <span class="loading-text">Loading categories...</span>
                        </div>
                    </div>
                    
                    <p id="post-error" class="error-msg"></p>
                    
                    <div class="modal-actions">
                        <button type="button" id="cancel-post" class="btn-secondary">Cancel</button>
                        <button type="submit" class="btn-primary">Publish Post</button>
                    </div>
                </form>
            </div>
        </div>
    `;

    // Setup Modal Logic
    const overlay = document.getElementById("create-modal-overlay");
    const openBtn = document.getElementById("open-create-modal");
    const closeBtn = document.getElementById("close-create-modal");
    const cancelBtn = document.getElementById("cancel-post");

    const closeModal = () => {
        overlay.classList.add("hidden");
        document.getElementById("post-error").textContent = ""; // clear errors
    };

    const openModal = () => {
        overlay.classList.remove("hidden");
    };

    openBtn.addEventListener("click", openModal);
    closeBtn.addEventListener("click", closeModal);
    cancelBtn.addEventListener("click", closeModal);
    
    // Close on outside click
    overlay.addEventListener("click", (e) => {
        if (e.target === overlay) closeModal();
    });

    // 2. Fetch and render categories for the checkboxes
    try {
        const catRes = await fetch("/api/categories", { credentials: "include" });
        if (catRes.ok) {
            const categories = await catRes.json();
            const catContainer = document.getElementById("categories-container");
            catContainer.innerHTML = categories.map(cat => `
                <label class="category-checkbox">
                    <input type="checkbox" name="category" value="${cat.id}">
                    <span class="checkmark"></span>
                    ${escapeHtml(cat.name)}
                </label>
            `).join("");
        }
    } catch (err) {
        console.error("Failed to load categories", err);
    }

    // 3. Fetch and render all posts
    loadPosts();

    // 4. Handle Create Post form submission
    document.getElementById("create-post-form").addEventListener("submit", async (e) => {
        e.preventDefault();
        const errorBox = document.getElementById("post-error");
        const submitBtn = e.target.querySelector('button[type="submit"]');
        
        errorBox.textContent = "";
        submitBtn.disabled = true;
        submitBtn.textContent = "Publishing...";

        const title = document.getElementById("post-title").value.trim();
        const content = document.getElementById("post-content").value.trim();
        
        const checkedCategories = [];
        document.querySelectorAll('input[name="category"]:checked').forEach(checkbox => {
            checkedCategories.push(parseInt(checkbox.value));
        });

        if (checkedCategories.length === 0) {
            errorBox.textContent = "Please select at least one category.";
            submitBtn.disabled = false;
            submitBtn.textContent = "Publish Post";
            return;
        }

        try {
            const res = await fetch("/api/posts", {
                method: "POST",
                headers: { "Content-Type": "application/json" },
                credentials: "include",
                body: JSON.stringify({ title, content, categories: checkedCategories })
            });

            if (!res.ok) {
                const data = await res.json();
                errorBox.textContent = data.message || "Failed to create post.";
                submitBtn.disabled = false;
                submitBtn.textContent = "Publish Post";
                return;
            }

            // Success! Close modal and reload the feed
            closeModal();
            renderfeed();
        } catch (err) {
            console.error(err);
            errorBox.textContent = "Something went wrong.";
            submitBtn.disabled = false;
            submitBtn.textContent = "Publish Post";
        }
    });
}

// Helper to fetch and display posts
async function loadPosts() {
    const container = document.getElementById("posts-container");
    try {
        const res = await fetch("/api/posts", { credentials: "include" });
        if (!res.ok) {
            container.innerHTML = `<div class="empty-state"><p>Failed to load posts.</p></div>`;
            return;
        }
        
        const posts = await res.json();
        
        if (posts.length === 0) {
            container.innerHTML = `
                <div class="empty-state">
                    <div style="font-size: 48px; margin-bottom: 16px;">💬</div>
                    <h3>No posts yet</h3>
                    <p>Be the first to start a discussion!</p>
                </div>`;
            return;
        }

        container.innerHTML = posts.map(post => `
            <div class="modern-post-card" data-id="${post.id}">
                <div class="card-header">
                    <div class="author-avatar">${escapeHtml(post.author.charAt(0).toUpperCase())}</div>
                    <div class="author-info">
                        <span class="author-name">@${escapeHtml(post.author)}</span>
                        <span class="post-time">${timeAgo(new Date(post.created_at))}</span>
                    </div>
                </div>
                <h3 class="post-title">${escapeHtml(post.title)}</h3>
                <p class="post-preview">${escapeHtml(post.content.length > 150 ? post.content.substring(0, 150) + '...' : post.content)}</p>
                <div class="post-tags">
                    ${post.categories.map(c => `<span class="tag">${escapeHtml(c)}</span>`).join("")}
                </div>
            </div>
        `).join("");

        // Add click events to go to the single post page
        document.querySelectorAll(".modern-post-card").forEach(card => {
            card.addEventListener("click", () => {
                const postId = card.getAttribute("data-id");
                renderPost(postId);
            });
        });

    } catch (err) {
        console.error(err);
        container.innerHTML = `<div class="empty-state"><p>Error loading posts.</p></div>`;
    }
}

// Simple security helper
function escapeHtml(str) {
    if (!str) return "";
    return str.replace(/&/g, "&amp;").replace(/</g, "&lt;").replace(/>/g, "&gt;");
}

// Time formatter
function timeAgo(date) {
    const seconds = Math.floor((new Date() - date) / 1000);
    let interval = seconds / 31536000;
    if (interval > 1) return Math.floor(interval) + " years ago";
    interval = seconds / 2592000;
    if (interval > 1) return Math.floor(interval) + " months ago";
    interval = seconds / 86400;
    if (interval > 1) return Math.floor(interval) + " days ago";
    interval = seconds / 3600;
    if (interval > 1) return Math.floor(interval) + " hours ago";
    interval = seconds / 60;
    if (interval > 1) return Math.floor(interval) + " minutes ago";
    return Math.floor(seconds) + " seconds ago";
}