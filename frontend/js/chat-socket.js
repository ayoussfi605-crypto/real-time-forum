let socket = null;
let onMessageCallback = null;   

export function connectSocket(onMessageReceived) {   
  socket = new WebSocket("ws://localhost:8080/ws");
  onMessageCallback = onMessageReceived;

  socket.onopen = () => {
    console.log("WebSocket connected");
  };

  socket.onmessage = (event) => {
    const msg = JSON.parse(event.data);
    if (onMessageCallback) {
      onMessageCallback(msg);  
    }
  };

  socket.onclose = () => {
    console.log("WebSocket disconnected");
  };

  socket.onerror = (err) => {
    console.error("WebSocket error:", err);
  };
}

export function sendMessage(receiverID, content) {
  const messageData = { receiver_id: receiverID, content: content };
  socket.send(JSON.stringify(messageData));
}