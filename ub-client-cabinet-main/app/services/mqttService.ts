import { connect as mqttConnect } from 'mqtt';
import { mqttServer } from './constants';
import { CookieKeys, cookies } from './cookie';
import { SideMessageService } from './message_service';

//connect to our broker
export const mqttCipher = (salt: string) => {
  const textToChars = (text: string) => text.split('').map(c => c.charCodeAt(0));
  const byteHex = (n: string) => ('0' + Number(n).toString(16)).substr(-2);
  const applySaltToChar = code => textToChars(salt).reduce((a, b) => a ^ b, code);

  return (text: string) =>
    text
      .split('')
      .map(textToChars)
      .map(applySaltToChar)
      .map(byteHex)
      .join('');
};
let TEncoder;
//export MQTT
if (!window['TextDecoder']) {
  TEncoder = require('text-encoding').TextDecoder;
  window['TextDecoder'] = TEncoder;
}

export function useStartMQTTMessages (data: { subject: string }) {
  const ubCipher = mqttCipher('ubSalt');
  const clientId = ubCipher(
    'mqttjs_' +
      Math.random()
        .toString(16)
        .substr(2, 8),
  );
  const mqttClient = mqttConnect(mqttServer, {
    password: clientId,
    username: cookies.get(CookieKeys.Token) ?? clientId,
    clientId: clientId,
    connectTimeout: 30 * 1000,
    reconnectPeriod: 2000,
    keepalive: 360,
  });

  mqttClient.on('connect', function () {
    if (process.env.NODE_ENV !== 'production') {
      console.log(
        `%cconnected to %c${data.subject}`,
        'color: #009c27;font-size:15px;font-weight:600;',
        'color: red',
      );
    }
    // subscribe on a topic after connected
    mqttClient.subscribe(data.subject);
  });

  mqttClient.on('message', function (topic, message) {
    // show received message

    const string = new TextDecoder('utf-8').decode(message);

    SideMessageService.send({
      name: data.subject,
      payload: JSON.parse(string),
    });
  });
  mqttClient.on('error', () => {
    alert('error');
  });
  //  mqttClient.on('packetreceive', (packet: any) => {
  //    let message = packet.payload != null && JSON.parse(packet.payload);
  //    console.log(message);
  //  });

  mqttClient.on('reconnect', function () {
    if (mqttClient.connected) {
      console.log(`reconnect`);
      //mqttClient.subscribe(topic, [], (err,res) => {})
    }
  });

  return mqttClient;
}
