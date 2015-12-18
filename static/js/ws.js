var webSocket = new WebSocket("ws://" + location.host + "/ws_chat");

// event handler
webSocket.binaryType = "arraybuffer";
webSocket.onopen = onOpen;
webSocket.onmessage = onMessage;
webSocket.onclose = onClose;
webSocket.onerror = onError;

// conn
function onOpen(event) {
	console.log("onOpen");
	console.log(event);
}

// receive msg
function onMessage(event) {
	console.log("onMessage");
	if(event && event.data ){
		console.log(event);
		document.getElementById("res").innerHTML = msgpack.decode(new Uint8Array(event.data));
	}
}

// err
function onError(event) {
	console.log("onError");
	console.log(event);
}

// close
function onClose(event) {
	console.log("onClose");
	console.log(event);
}

// send
function writeMessage() {
	var val = document.getElementById("message").value;
	console.log(val);
	webSocket.send(msgpack.encode(val));

	return false;
}
