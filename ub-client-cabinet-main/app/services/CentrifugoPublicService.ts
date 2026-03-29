import { Centrifuge } from 'centrifuge';
import { centrifugoUrl } from './constants';
import {
  SideMessageService,
  MarketTradeMessageService,
  OrderBookMessageService,
  MarketWatchMessageService,
  TradeChartMessageService,
} from './message_service';
import { CentrifugoChannels } from 'containers/App/constants';

export class CentrifugoPublicService {
  private static instance: CentrifugoPublicService;
  private centrifuge: Centrifuge;
  private subscriptions: Map<string, any> = new Map();

  private constructor() {
    this.centrifuge = new Centrifuge(centrifugoUrl);

    this.centrifuge.on('connecting', (ctx) => {
      if (process.env.NODE_ENV !== 'production') {
        console.log('Centrifugo public connecting:', ctx);
      }
    });
    this.centrifuge.on('connected', (ctx) => {
      if (process.env.NODE_ENV !== 'production') {
        console.log('Centrifugo public connected:', ctx);
      }
    });
    this.centrifuge.on('disconnected', (ctx) => {
      if (process.env.NODE_ENV !== 'production') {
        console.log('Centrifugo public disconnected:', ctx);
      }
    });

    this.centrifuge.connect();
  }

  public static getInstance(): CentrifugoPublicService {
    if (!CentrifugoPublicService.instance) {
      CentrifugoPublicService.instance = new CentrifugoPublicService();
    }
    return CentrifugoPublicService.instance;
  }

  private routeMessage(channel: string, data: any) {
    if (channel.startsWith(CentrifugoChannels.TickerChannel)) {
      MarketWatchMessageService.send({ name: channel, payload: data });
    } else if (channel.startsWith(CentrifugoChannels.OrderBookPrefix)) {
      OrderBookMessageService.send({ name: channel, payload: data });
    } else if (channel.startsWith(CentrifugoChannels.MarketTradePrefix)) {
      MarketTradeMessageService.send({ name: channel, payload: data });
    } else if (channel.startsWith(CentrifugoChannels.TradeChartPrefix)) {
      TradeChartMessageService.send({ name: channel, payload: data });
    } else {
      SideMessageService.send({ name: channel, payload: data });
    }
  }

  public ConnectToSubject(data: { subject: string }) {
    if (this.subscriptions.has(data.subject)) return;

    const sub = this.centrifuge.newSubscription(data.subject);
    sub.on('publication', (ctx) => {
      this.routeMessage(data.subject, ctx.data);
    });
    sub.on('subscribing', (ctx) => {
      if (process.env.NODE_ENV !== 'production') {
        console.log(`Subscribing to ${data.subject}:`, ctx);
      }
    });
    sub.subscribe();
    this.subscriptions.set(data.subject, sub);
  }

  public DisconnectFromSubject(data: { subject: string }) {
    const sub = this.subscriptions.get(data.subject);
    if (sub) {
      sub.unsubscribe();
      this.centrifuge.removeSubscription(sub);
      this.subscriptions.delete(data.subject);
    }
  }

  public ConnectToNewSubject(data: {
    oldsubject: string;
    newSubject: string;
  }) {
    this.DisconnectFromSubject({ subject: data.oldsubject });
    this.ConnectToSubject({ subject: data.newSubject });
  }
}

export const centrifugoPublicService = CentrifugoPublicService.getInstance();
