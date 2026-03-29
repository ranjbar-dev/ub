import {MqttAdditionalConfig,mqttServer} from './constants';
import {connect as mqttConnect} from 'mqtt';
import {
	SideMessageService,
	MarketTradeMessageService,
	OrderBookMessageService,
	MarketWatchMessageService,
	TradeChartMessageService,
} from './message_service';
import {MqttTopicsPrefixes} from 'containers/App/constants';
import {cookies,CookieKeys} from './cookie';

export const mqttCipher=(salt: string) => {
	const textToChars=(text: string) => text.split('').map(c => c.charCodeAt(0));
	const byteHex=(n: string) => ('0'+Number(n).toString(16)).substr(-2);
	const applySaltToChar=code => textToChars(salt).reduce((a,b) => a^b,code);

	return (text: string) =>
		text
			.split('')
			.map(textToChars)
			.map(applySaltToChar)
			.map(byteHex)
			.join('');
};

export class MqttService {
	private static instance: MqttService;
	private static mqttCl;
	private static cipher;
	private static clId;
	private static activeSubscriptions: Set<string> = new Set();
	public static Encoder;
	private constructor() {
		let TEncoder;
		if(!window['TextDecoder']) {
			TEncoder=require('text-encoding').TextDecoder;
			window['TextDecoder']=TEncoder;
			MqttService.Encoder=new TextDecoder('utf-8');
		} else {
			MqttService.Encoder=new TextDecoder('utf-8');
		}
	}
	public static getInstance(): MqttService {
		if(!MqttService.instance) {
			MqttService.instance=new MqttService();
			MqttService.cipher=mqttCipher('ubSalt');
			MqttService.clId=MqttService.cipher(
				'mqttjs_'+
				Math.random()
					.toString(16)
					.substr(2,8),
			);
			MqttService.mqttCl=mqttConnect(mqttServer,{
				password: MqttService.clId,
				username: cookies.get(CookieKeys.Token)??MqttService.clId,
				clientId: MqttService.clId,
				...MqttAdditionalConfig,
			});
			MqttService.mqttCl.on('connect',function(params) {
				if(process.env.NODE_ENV!=='production') {
					console.log(params);
				}
			});
			MqttService.mqttCl.on('reconnect', function() {
				MqttService.activeSubscriptions.forEach(topic => {
					MqttService.mqttCl.subscribe(topic);
				});
			});

			MqttService.mqttCl.on('message',function(topic: string,message) {
				// show received message
				const string=MqttService.Encoder.decode(message);
				if(topic.includes(MqttTopicsPrefixes.MarketWatchAddress)) {
					MarketWatchMessageService.send({
						name: topic,
						payload: JSON.parse(string),
					});
				} else if(topic.includes(MqttTopicsPrefixes.OrderBookAddress)) {
					OrderBookMessageService.send({
						name: topic,
						payload: JSON.parse(string),
					});
				} else if(topic.includes(MqttTopicsPrefixes.MarketTradeAddress)) {
					MarketTradeMessageService.send({
						name: topic,
						payload: JSON.parse(string),
					});
				} else if(topic.includes(MqttTopicsPrefixes.TradeChartAddress)) {
					TradeChartMessageService.send({
						name: topic,
						payload: JSON.parse(string),
					});
				} else {
					SideMessageService.send({
						name: topic,
						payload: JSON.parse(string),
					});
				}
			});
		}

		return MqttService.instance;
	}

	public ConnectToSubject(data: {subject: string}) {
		MqttService.activeSubscriptions.add(data.subject);
		MqttService.mqttCl.subscribe(data.subject);
	}
	public DisconnectFromSubject(data: {subject: string}) {
		MqttService.activeSubscriptions.delete(data.subject);
		MqttService.mqttCl.unsubscribe(data.subject);
	}
	public ConnectToNewSubject(data: {
		oldsubject: string;
		newSubject: string;
	}) {
		MqttService.activeSubscriptions.delete(data.oldsubject);
		MqttService.activeSubscriptions.add(data.newSubject);
		MqttService.mqttCl.unsubscribe(data.oldsubject).subscribe(data.newSubject);
	}
}

export const mqttService2=MqttService.getInstance();
