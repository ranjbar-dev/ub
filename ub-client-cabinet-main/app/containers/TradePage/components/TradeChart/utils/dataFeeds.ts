//function webSocketConnect() {
//  //Web socket connection

//  ws.on('reconnect', function () {
//    if (ws.connected) {
//      ws.subscribe(topic, [], (err, res) => {});
//    }
//  });

//  ws.on('packetreceive', (packet) => {
//    let message = packet.payload != null && JSON.parse(packet.payload);
//    if (packet.topic === prepareTopic()) {
//      if (message) {
//        newBarInfo = prepareNewBarInfo(message);
//        if (onRealtimeCallbackFunc !== undefined) {
//          onRealtimeCallbackFunc(newBarInfo);
//        }
//      }
//    }
//  });
//}
