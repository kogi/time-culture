const ws = new WebSocket("ws://localhost:7380/ws?client=operation");

ws.onopen = () => {
    console.log("WebSocket connection opened");
    document.getElementById("start").addEventListener("click", () => {
        console.log("send start command");
        const msg = {
            Command: "start",
            timestamp: Date.now().toString()
        };
        ws.send(JSON.stringify(msg));
    });
}