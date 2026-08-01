let socket = null; 

export function connectSocket() {   
  socket = new WebSocket("ws://localhost:8080/ws");

  socket.onopen = () => {
    console.log("WebSocket connected");
    sendMessage(2, "Salam");
  };


  socket.onmessage = (event) => {
    const msg = JSON.parse(event.data);
  };

  socket.onclose = () => {
    console.log("WebSocket disconnected");
  };

  socket.onerror = (err) => {
    console.error("WebSocket error:", err);
  };
}

export function sendMessage(receiverID, content) {
    if (!socket || socket.readyState !== WebSocket.OPEN) {
        console.log("Socket not connected");
        return;
    }

    const message = {
        event_type: "private_message",
        data: {
            receiver_id: receiverID,
            content: content,
        },
    };

    console.log("Sending:", message);

    socket.send(JSON.stringify(message));
}
 
window.sendMessage = sendMessage;