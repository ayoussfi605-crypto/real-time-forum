import { navigate } from "./navigate.js";

// ─────────────────────────────────────────────
// Renders the detail view of a single post
// postId — the ID from the URL (e.g. 5)
// ─────────────────────────────────────────────
export async function renderPost(postId) {
  const app = document.getElementById("app");

  // Show loading state
  app.innerHTML = `<p class="loading-msg">Loading post...</p>`;

  try {
    const response = await fetch(`/api/posts/${postId}`, { credentials: "include" });

    if (!response.ok) {
      app.innerHTML = `<p class="error-msg">Post not found.</p>`;
      return;
    }

    const post = await response.json();

    // Render the full post + comments
    app.innerHTML = `
      <div class="feed-container">
        
        <!-- Back button -->
        <button id="back-btn" class="btn-secondary" style="margin-bottom: 20px; display: inline-flex; align-items: center; gap: 8px;">
          <span style="font-size: 18px; line-height: 1;">&larr;</span>
          Back to Feed
        </button>

        <!-- The post itself -->
        <div class="modern-post-card post-detail-main">
          <div class="card-header">
            <div class="author-avatar">${escapeHtml(post.author.charAt(0).toUpperCase())}</div>
            <div class="author-info">
              <span class="author-name">@${escapeHtml(post.author)}</span>
              <span class="post-time">${formatDate(post.created_at)}</span>
            </div>
          </div>
          <h2 class="post-detail-title">${escapeHtml(post.title)}</h2>
          <p class="post-detail-content" style="white-space: pre-wrap; font-size: 16px; line-height: 1.7; color: #334155; margin-bottom: 24px;">${escapeHtml(post.content)}</p>
          <div class="post-tags" style="border-top: 1px solid #e2e8f0; padding-top: 16px;">
            ${post.categories.map(cat => `<span class="tag">${escapeHtml(cat)}</span>`).join("")}
          </div>
        </div>

        <!-- Reactions -->
      <div class="post-reactions">

      <button
      id="like-btn"
      class="reaction-btn ${post.reactions.user_reaction === "like" ? "active" : ""}"
    >
    👍
    <span id="likes-count">
      ${post.reactions.likes}
    </span>
    </button>

    <button
    id="dislike-btn"
    class="reaction-btn ${post.reactions.user_reaction === "dislike" ? "active" : ""}"
    >
    👎
    <span id="dislikes-count">
      ${post.reactions.dislikes}
    </span>
    </button>

    </div>

        <!-- Comments section -->
        <div class="comments-section">
          <h3 class="comments-title">Comments <span class="comment-count">${post.comments.length}</span></h3>

          <!-- Add comment form -->
          <div class="modern-post-card comment-input-card">
            <form id="comment-form" class="comment-form modern-form">
              <div class="form-group" style="margin-bottom: 12px;">
                <textarea id="comment-content" placeholder="Share your thoughts..." rows="3" required style="width: 100%; border: none; resize: none; background: #f8fafc; outline: none; box-shadow: none;" onfocus="this.style.background='#fff'; this.parentElement.style.borderColor='#3a6ea5';" onblur="this.style.background='#f8fafc'; this.parentElement.style.borderColor='#e2e8f0';"></textarea>
              </div>
              <p id="comment-error" class="error-msg" style="margin-top:0; margin-bottom:10px;"></p>
              <div style="display:flex; justify-content: flex-end;">
                <button type="submit" class="btn-primary" style="padding: 8px 20px; border-radius: 20px;">Reply</button>
              </div>
            </form>
          </div>

          <!-- List of existing comments -->
          <div id="comments-list" class="comments-list">
            ${post.comments.length === 0
        ? `<div class="empty-state" style="padding: 30px;"><p>No comments yet. Be the first to reply!</p></div>`
        : post.comments.map(c => `
                  <div class="comment-card">
                    <div class="card-header" style="margin-bottom: 8px;">
                      <div class="author-avatar" style="width: 32px; height: 32px; font-size: 14px;">${escapeHtml(c.author.charAt(0).toUpperCase())}</div>
                      <div class="author-info">
                        <span class="author-name" style="font-size: 14px;">@${escapeHtml(c.author)}</span>
                        <span class="post-time" style="font-size: 12px;">${formatDate(c.created_at)}</span>
                      </div>
                    </div>
                    <p class="comment-content">${escapeHtml(c.content)}</p>
                  </div>
                `).join("")
      }
          </div>

        </div>

      </div>
    `;

    // Back button — go back to the feed
    document.getElementById("back-btn").addEventListener("click", () => {
      navigate("feed");
    });

    document.getElementById("like-btn").addEventListener("click", () => {
	  handleReaction(postId, "like");
  });

  document.getElementById("dislike-btn").addEventListener("click", () => {
	handleReaction(postId, "dislike");
  });
    // Comment form submission
    document.getElementById("comment-form").addEventListener("submit", (e) => {
      handleAddComment(e, postId);
    });

  } catch (err) {
    console.error(err);
    app.innerHTML = `<p class="error-msg">Something went wrong.</p>`;
  }
}

async function handleReaction(postId, reaction) {
	try {
		const response = await fetch(
			`/api/posts/${postId}/reaction`,
			{
				method: "POST",
				headers: {
					"Content-Type": "application/json"
				},
				credentials: "include",
				body: JSON.stringify({
					reaction: reaction
				})
			}
		);

		const result = await response.json();

		console.log(result);

		if (!response.ok) {
			console.error(result.message);
			return;
		}

		const reactionData = result.data[0];

		// Update numbers
		document.getElementById("likes-count").textContent =
			reactionData.likes;

		document.getElementById("dislikes-count").textContent =
			reactionData.dislikes;

		// Get buttons
		const likeBtn = document.getElementById("like-btn");
		const dislikeBtn = document.getElementById("dislike-btn");

		// Remove active from both
		likeBtn.classList.remove("active");
		dislikeBtn.classList.remove("active");

		// Add active to current reaction
		if (reactionData.userReaction === "like") {
			likeBtn.classList.add("active");
		}

		if (reactionData.userReaction === "dislike") {
			dislikeBtn.classList.add("active");
		}

	} catch (err) {
		console.error("Reaction error:", err);
	}
}

// ─────────────────────────────────────────────
// Handles submitting a new comment
// ─────────────────────────────────────────────
async function handleAddComment(e, postId) {
  e.preventDefault();

  const errorBox = document.getElementById("comment-error");
  errorBox.textContent = "";

  const content = document.getElementById("comment-content").value.trim();
  if (!content) {
    errorBox.textContent = "Comment cannot be empty.";
    return;
  }

  try {
    const response = await fetch(`/api/posts/${postId}/comments`, {
      method: "POST",
      headers: { "Content-Type": "application/json" },
      credentials: "include",
      body: JSON.stringify({ content }),
    });

    const result = await response.json();

    if (!response.ok) {
      errorBox.textContent = result.message;
      return;
    }

    // Re-render the whole post page to show the new comment
    renderPost(postId);

  } catch (err) {
    console.error(err);
    errorBox.textContent = "Something went wrong. Try again.";
  }
}

// ─────────────────────────────────────────────
// Helpers
// ─────────────────────────────────────────────
function formatDate(dateStr) {
  return new Date(dateStr).toLocaleDateString("en-US", {
    year: "numeric", month: "short", day: "numeric", hour: "2-digit", minute: "2-digit"
  });
}

function escapeHtml(str) {
  return str
    .replace(/&/g, "&amp;")
    .replace(/</g, "&lt;")
    .replace(/>/g, "&gt;");
}
