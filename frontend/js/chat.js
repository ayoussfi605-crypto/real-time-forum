import { getCurrentUser } from "./state.js";

let ws;
let chatUsers = []; // Array to hold {id, nickname, online}
let currentChatUserId = null;
let messagesOffset = 0;
let isLoadingMessages = false;
let hasMoreMessages = true;

// Throttle function for scroll
function throttle(func, limit) {
    let lastFunc;
    let lastRan;
    return function(...args) {
        if (!lastRan) {
            func.apply(this, args);
            lastRan = Date.now();
        } else {
            clearTimeout(lastFunc);
            lastFunc = setTimeout(() => {
                if ((Date.now() - lastRan) >= limit) {
                    func.apply(this, args);
                    lastRan = Date.now();
                }
            }, limit - (Date.now() - lastRan));
        }
    }
}
let listenersAttached = false;

export async function initChat() {
    const user = getCurrentUser();
    if (!user) return;

    // Fetch initial sorted user list from API
    try {
        const res = await fetch("/api/users");
        if (res.ok) {
            const users = await res.json();
            if (users) {
                chatUsers = users.map(u => ({ ...u, online: false, unread_count: u.unread_count || 0 }));
            }
        }
    } catch (err) {
        console.error("Failed to load chat users", err);
    }

    renderUserList();
    document.getElementById("chat-sidebar").classList.remove("hidden");
    document.body.classList.add("has-sidebar");

    // Initialize WebSocket (only if not already connected)
    if (!ws) {
        ws = new WebSocket("ws://" + window.location.host + "/ws");

        ws.onopen = () => {
            console.log("WebSocket connected");
        };

        ws.onmessage = (event) => {
            const msg = JSON.parse(event.data);
            handleWSMessage(msg);
        };

        ws.onclose = () => {
            console.log("WebSocket disconnected");
            ws = null; // allow reconnect next initChat()
        };
    }

    // Attach DOM listeners only once for the page's lifetime
    if (listenersAttached) return;
    listenersAttached = true;

    // Chat form submit
    document.getElementById("chat-form").addEventListener("submit", (e) => {
        e.preventDefault();
        const currentUser = getCurrentUser();
        if (!currentUser) return;

        const input = document.getElementById("chat-input");
        const content = input.value.trim();
        if (content && currentChatUserId && ws && ws.readyState === WebSocket.OPEN) {
            const chatMsg = {
                type: "chat_message",
                receiver_id: currentChatUserId,
                content: content
            };
            ws.send(JSON.stringify(chatMsg));

            // Optimistically add to UI
            appendMessage({
                sender_id: currentUser.id,
                content: content,
                created_at: new Date().toISOString()
            }, true);

            input.value = "";
            scrollToBottom();

            // Re-sort users to put this user at top
            moveUserToTop(currentChatUserId);
        }
    });

    document.getElementById("close-chat-btn").addEventListener("click", () => {
        document.getElementById("chat-window").classList.add("hidden");
        currentChatUserId = null;
    });

    // Scroll event for pagination
    const messagesContainer = document.getElementById("chat-messages");
    messagesContainer.addEventListener("scroll", throttle(async () => {
        if (messagesContainer.scrollTop <= 50 && !isLoadingMessages && hasMoreMessages) {
            await loadMoreMessages();
        }
    }, 200));
}

function handleWSMessage(msg) {
    if (msg.type === "initial_status") {
        const onlineIds = msg.online_users || [];
        chatUsers.forEach(u => {
            if (onlineIds.includes(u.id)) u.online = true;
        });
        renderUserList();
    } else if (msg.type === "user_status") {
        const user = chatUsers.find(u => u.id === msg.user_id);
        if (user) {
            user.online = msg.online;
            renderUserList();
        }
    } else if (msg.type === "chat_message") {
        // If message is from current chat, append it and mark as read
        if (currentChatUserId === msg.sender_id) {
            appendMessage(msg, true);
            scrollToBottom();
            fetch(`/api/messages/${msg.sender_id}/read`, { method: "POST" }).catch(console.error);
        } else {
            // Otherwise, increment unread count
            const user = chatUsers.find(u => u.id === msg.sender_id);
            if (user) {
                user.unread_count++;
                renderUserList();
            }
        }
        // Move sender to top of user list
        moveUserToTop(msg.sender_id);
    }
}

function renderUserList() {
    const list = document.getElementById("user-list");
    list.innerHTML = "";
    
    // Sort logic: Keep original order (which is by last message from DB), 
    // but we update it via moveUserToTop when new messages arrive.
    // So we just render the array as is.
    
    chatUsers.forEach(u => {
        const li = document.createElement("li");
        li.className = "user-item";
        li.innerHTML = `
            <div class="status-indicator ${u.online ? 'online' : ''}"></div>
            <div class="user-item-name">${u.nickname}</div>
            <span class="unread-badge ${u.unread_count > 0 ? '' : 'hidden'}">${u.unread_count}</span>
        `;
        li.addEventListener("click", () => openChat(u.id, u.nickname));
        list.appendChild(li);
    });
}

function moveUserToTop(userId) {
    const index = chatUsers.findIndex(u => u.id === userId);
    if (index > 0) {
        const user = chatUsers.splice(index, 1)[0];
        chatUsers.unshift(user);
        renderUserList();
    }
}

async function openChat(userId, nickname) {
    currentChatUserId = userId;
    messagesOffset = 0;
    hasMoreMessages = true;
    
    // Clear unread count
    const user = chatUsers.find(u => u.id === userId);
    if (user && user.unread_count > 0) {
        user.unread_count = 0;
        renderUserList();
        fetch(`/api/messages/${userId}/read`, { method: "POST" }).catch(console.error);
    }
    
    document.getElementById("chat-window").classList.remove("hidden");
    document.getElementById("chat-user-name").textContent = nickname;
    document.getElementById("chat-messages").innerHTML = "";
    document.getElementById("chat-input").focus();
    
    await loadMoreMessages();
    scrollToBottom();
}

async function loadMoreMessages() {
    if (!currentChatUserId || isLoadingMessages || !hasMoreMessages) return;
    
    isLoadingMessages = true;
    const container = document.getElementById("chat-messages");
    const oldScrollHeight = container.scrollHeight;
    
    try {
        const res = await fetch(`/api/messages/${currentChatUserId}?offset=${messagesOffset}`);
        if (res.ok) {
            const messages = await res.json();
            if (messages && messages.length > 0) {
                // messages are chronologically ordered by backend
                // prepend them to UI
                messages.reverse().forEach(msg => {
                    appendMessage(msg, false); // false = prepend
                });
                
                messagesOffset += messages.length;
                if (messages.length < 10) {
                    hasMoreMessages = false; // reached the end
                }
                
                // Keep scroll position stable
                container.scrollTop = container.scrollHeight - oldScrollHeight;
            } else {
                hasMoreMessages = false;
            }
        }
    } catch (err) {
        console.error("Failed to load messages", err);
    }
    
    isLoadingMessages = false;
}

function appendMessage(msg, append = true) {
    const container = document.getElementById("chat-messages");
    const user = getCurrentUser();
    
    const isSent = msg.sender_id === user.id;
    const div = document.createElement("div");
    div.className = `chat-message ${isSent ? 'sent' : 'received'}`;
    
    const dateStr = new Date(msg.created_at).toLocaleTimeString([], { hour: '2-digit', minute: '2-digit' });

    let senderName = isSent ? "You" : (document.getElementById("chat-user-name").textContent);
    
    div.innerHTML = `
        <div class="chat-message-sender">${senderName}</div>
        <div class="chat-message-content">${escapeHTML(msg.content)}</div>
        <div class="chat-message-meta">${dateStr}</div>
    `;
    
    if (append) {
        container.appendChild(div);
    } else {
        container.insertBefore(div, container.firstChild);
    }
}

function scrollToBottom() {
    const container = document.getElementById("chat-messages");
    container.scrollTop = container.scrollHeight;
}

function escapeHTML(str) {
    const div = document.createElement('div');
    div.textContent = str;
    return div.innerHTML;
}

export function closeChatAndWS() {
    if (ws) {
        ws.close();
        ws = null;
    }
    document.getElementById("chat-sidebar").classList.add("hidden");
    document.getElementById("chat-window").classList.add("hidden");
    document.body.classList.remove("has-sidebar");
    currentChatUserId = null;
    chatUsers = [];
}
