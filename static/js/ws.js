var webSocket = new WebSocket("ws://" + location.host + "/ws_chat");

// イベントハンドラの設定
webSocket.onopen = onOpen;
webSocket.onmessage = onMessage;
webSocket.onclose = onClose;
webSocket.onerror = onError;

// 接続イベント
function onOpen(event) {
	console.log("onOpen");
	console.log(event);
}

// メッセージ受信イベント
function onMessage(event) {
	console.log("onMessage");
	if(event && event.data ){
		console.log(event);
		document.getElementById("res").innerHTML = event.data;
	}
}

// エラーイベント
function onError(event) {
	console.log("onError");
	console.log(event);
}

// 切断イベント
function onClose(event) {
	console.log("onClose");
	console.log(event);
}

function writeMessage() {
	var val = document.getElementById("message").value;
	console.log(val);
	webSocket.send(val);

	return false;
}
