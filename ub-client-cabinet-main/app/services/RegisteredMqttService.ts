import { MqttAdditionalConfig, mqttServer } from './constants';
import { connect as mqttConnect, MqttClient } from 'mqtt';
import { RegisteredUserMessageService } from './message_service';
import { cookies, CookieKeys } from './cookie';

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

let connectionTestTimer;

export class RegisteredMqttService {
  private static instance: RegisteredMqttService;
  private static mqttCl: MqttClient;
  private static cipher;
  private static clId;
  private static activeSubscriptions: Set<string> = new Set();
  public static Encoder;
  private constructor () {
    let TEncoder;
    if (!window['TextDecoder']) {
      TEncoder = require('text-encoding').TextDecoder;
      window['TextDecoder'] = TEncoder;
      RegisteredMqttService.Encoder = new TextDecoder('utf-8');
    } else {
      RegisteredMqttService.Encoder = new TextDecoder('utf-8');
    }
  }

  public static getInstance (updatedToken?: string): RegisteredMqttService {
    if (!RegisteredMqttService.instance) {
      RegisteredMqttService.instance = new RegisteredMqttService();
      RegisteredMqttService.cipher = mqttCipher('ubSalt');
      RegisteredMqttService.clId = RegisteredMqttService.cipher(
        'mqttjs_' +
          Math.random()
            .toString(16)
            .substr(2, 8),
      );
      RegisteredMqttService.mqttCl = mqttConnect(mqttServer, {
        password: RegisteredMqttService.clId,
        username: cookies.get(CookieKeys.Token) ?? RegisteredMqttService.clId,
        clientId: RegisteredMqttService.clId,
        ...MqttAdditionalConfig,
      });
      RegisteredMqttService.mqttCl.on('connect', function (params) {
        if (process.env.NODE_ENV !== 'production') {
          console.log(params);
        }
      });
      RegisteredMqttService.mqttCl.on('reconnect', function () {
        if (process.env.NODE_ENV !== 'production') {
          console.log(
            `%cRegistered mqtt Reconnected ${new Date().toISOString()}`,
            'color:magenta;font-size:16px;',
          );
        }
        RegisteredMqttService.activeSubscriptions.forEach(topic => {
          RegisteredMqttService.mqttCl.subscribe(topic);
        });
      });
      RegisteredMqttService.mqttCl.on('close', () => {
        if (process.env.NODE_ENV !== 'production') {
          console.log(
            `%cRegistered mqtt disconnected1 ${new Date().toISOString()}`,
            'color:red;font-size:12px;',
          );
        }
      });

      RegisteredMqttService.mqttCl.on('message', function (topic, message) {
        // show received message
        const string = RegisteredMqttService.Encoder.decode(message);
        RegisteredUserMessageService.send({
          name: topic,
          payload: JSON.parse(string),
        });
      });
    } else if (updatedToken) {
      RegisteredMqttService.mqttCl.end(true);
      RegisteredMqttService.activeSubscriptions.clear();
      //@ts-ignore
      delete RegisteredMqttService.mqttCl;
      //@ts-ignore
      delete RegisteredMqttService.instance;
      RegisteredMqttService.instance = new RegisteredMqttService();
      RegisteredMqttService.cipher = mqttCipher('ubSalt');
      RegisteredMqttService.clId = RegisteredMqttService.cipher(
        'mqttjs_' +
          Math.random()
            .toString(16)
            .substr(2, 8),
      );
      RegisteredMqttService.mqttCl = mqttConnect(mqttServer, {
        password: RegisteredMqttService.clId,
        username: updatedToken,
        clientId: RegisteredMqttService.clId,
        ...MqttAdditionalConfig,
      });
      RegisteredMqttService.mqttCl.on('connect', function (params) {
        clearInterval(connectionTestTimer);
        connectionTestTimer = setInterval(() => {
          RegisteredMqttService.mqttCl.publish('/testing', 'connectionCheck');
        }, 5000);

        if (process.env.NODE_ENV !== 'production') {
          console.log(params);
        }
      });
      RegisteredMqttService.mqttCl.on('reconnect', function () {
        if (process.env.NODE_ENV !== 'production') {
          console.log(
            `%cRegistered mqtt Reconnected ${new Date().toISOString()}`,
            'color:magenta;font-size:18px;',
          );
        }
        RegisteredMqttService.activeSubscriptions.forEach(topic => {
          RegisteredMqttService.mqttCl.subscribe(topic);
        });
      });
      RegisteredMqttService.mqttCl.on('close', () => {
        if (process.env.NODE_ENV !== 'production') {
          console.log(
            `%cRegistered mqtt disconnected2 ${new Date().toISOString()}`,
            'color:red;font-size:12px;',
          );
        }
      });

      RegisteredMqttService.mqttCl.on('message', function (topic, message) {
        // show received message
        const string = RegisteredMqttService.Encoder.decode(message);
        RegisteredUserMessageService.send({
          name: topic,
          payload: JSON.parse(string),
        });
      });
    }
    return RegisteredMqttService.instance;
  }

  public ConnectToSubject (data: { subject: string }) {
    if (process.env.NODE_ENV !== 'production') {
      console.log('connecting', data);
    }
    if (data.subject !== undefined) {
      RegisteredMqttService.activeSubscriptions.add(data.subject);
      RegisteredMqttService.mqttCl.subscribe(data.subject);
    }
  }
  public DisconnectFromSubject (data: { subject: string }) {
    if (process.env.NODE_ENV !== 'production') {
      console.log('disConnecting', data);
    }
    RegisteredMqttService.activeSubscriptions.delete(data.subject);
    RegisteredMqttService.mqttCl.unsubscribe(data.subject);
  }
  public ConnectToNewSubject (data: {
    oldsubject: string;
    newSubject: string;
  }) {
    RegisteredMqttService.activeSubscriptions.delete(data.oldsubject);
    RegisteredMqttService.activeSubscriptions.add(data.newSubject);
    RegisteredMqttService.mqttCl
      .unsubscribe(data.oldsubject)
      .subscribe(data.newSubject);
  }
}

export const registeredMqttService = RegisteredMqttService.getInstance();
