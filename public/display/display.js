

document.getElementById("confirm-button").addEventListener("click", function() {
    document.getElementById("confirm").remove()

    const ws = new WebSocket("ws://localhost:7380/ws?client=display");

// サーバーから動画全体をロードして、blobを生成する
    const videoSrc = "./files/video.mp4"; // 動画のURL

    let blobUrl;

    async function loadVideoAsBlob() {
        try {
            const response = await fetch(videoSrc);
            if (!response.ok) {
                throw new Error(`HTTP error! status: ${response.status}`);
            }

            const blob = await response.blob(); // blobを取得
            blobUrl = blob; // blobUrlに割り当て
            console.log("Blob generated:", blob); // blobをログに出力
        } catch (error) {
            console.error("動画のロード中にエラーが発生しました:", error);
        }
    }

    function playVideo(current) {
        console.log("playVideo called");
        const videoElement = document.getElementById("video");
        if (blobUrl) {
            videoElement.src = URL.createObjectURL(blobUrl);
            videoElement.currentTime = current / 1000; // 現在の時間を設定
            videoElement.play();
        }else {
            console.error("Blob URLが生成されていません。");
        }

        // canvasで再生
        const canvas = document.getElementById("canvas");
        const ctx = canvas.getContext("2d");
        const video = document.getElementById("video");

        function step() {
            if (!video.paused && !video.ended) {
                ctx.drawImage(video, 0, 0, canvas.width, canvas.height);
            }
        }

        setInterval(() => {
            step();
        }, 1000 / 30); // 30 FPSで描画
    }

    const msg = {
        Command: "ready",
        timestamp: Date.now().toString()
    };

    loadVideoAsBlob().then(r => ws.send(JSON.stringify(msg))).catch(e => console.error(e));

// WebSocketの接続が開かれたときの処理
    ws.addEventListener("open", function() {
        console.log("WebSocket connection opened");
        // ここでサーバーにメッセージを送信することができます
        ws.onmessage = (byte) => {
            console.log("Received message:", JSON.parse(byte.data).command);
            if (JSON.parse(byte.data).command === "start") {
                if(Date.now() < JSON.parse(byte.data).timestamp*1 + 2000){
                    setTimeout(()=>{
                        playVideo(0);
                    }, JSON.parse(byte.data).timestamp*1 + 2000 - Date.now());
                }else{
                    playVideo(Date.now() - JSON.parse(byte.data).timestamp*1 - 2000);
                }
            }
        }
    });

})

