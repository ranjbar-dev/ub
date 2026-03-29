import { Centrifuge } from 'centrifuge';
import { centrifugoUrl, BaseUrl } from './constants';
import { RegisteredUserMessageService } from './message_service';
import { cookies, CookieKeys } from './cookie';

export class CentrifugoAuthService {
  private static instance: CentrifugoAuthService;
  private centrifuge: Centrifuge;
  private subscriptions: Map<string, any> = new Map();

  private constructor() {
    this.centrifuge = new Centrifuge(centrifugoUrl, {
      getToken: async () => {
        const token = cookies.get(CookieKeys.Token);
        const response = await fetch(BaseUrl + 'auth/centrifugo-token', {
          headers: {
            'Content-Type': 'application/json',
            ...(token && { Authorization: `Bearer ${token}` }),
          },
        });
        const data = await response.json();
        return data.token;
      },
    });

    this.centrifuge.on('connecting', (ctx) => {
      if (process.env.NODE_ENV !== 'production') {
        console.log('Centrifugo auth connecting:', ctx);
      }
    });
    this.centrifuge.on('connected', (ctx) => {
      if (process.env.NODE_ENV !== 'production') {
        console.log('Centrifugo auth connected:', ctx);
      }
    });
    this.centrifuge.on('disconnected', (ctx) => {
      if (process.env.NODE_ENV !== 'production') {
        console.log(
          `%cCentrifugo auth disconnected ${new Date().toISOString()}`,
          'color:red;font-size:12px;',
          ctx,
        );
      }
    });

    this.centrifuge.connect();
  }

  public static getInstance(updatedToken?: string): CentrifugoAuthService {
    if (!CentrifugoAuthService.instance) {
      CentrifugoAuthService.instance = new CentrifugoAuthService();
    } else if (updatedToken) {
      // Tear down old connection and create fresh instance
      CentrifugoAuthService.instance.subscriptions.forEach((sub) => {
        sub.unsubscribe();
        CentrifugoAuthService.instance.centrifuge.removeSubscription(sub);
      });
      CentrifugoAuthService.instance.subscriptions.clear();
      CentrifugoAuthService.instance.centrifuge.disconnect();
      // @ts-ignore
      delete CentrifugoAuthService.instance;

      CentrifugoAuthService.instance = new CentrifugoAuthService();
    }
    return CentrifugoAuthService.instance;
  }

  public ConnectToSubject(data: { subject: string }) {
    if (data.subject === undefined) return;
    if (this.subscriptions.has(data.subject)) return;

    if (process.env.NODE_ENV !== 'production') {
      console.log('connecting', data);
    }

    const channel = data.subject;
    const sub = this.centrifuge.newSubscription(channel, {
      getToken: async () => {
        const token = cookies.get(CookieKeys.Token);
        const response = await fetch(
          BaseUrl +
            `auth/centrifugo-subscribe-token?channel=${encodeURIComponent(channel)}`,
          {
            headers: {
              'Content-Type': 'application/json',
              ...(token && { Authorization: `Bearer ${token}` }),
            },
          },
        );
        const data = await response.json();
        return data.token;
      },
    });

    sub.on('publication', (ctx) => {
      RegisteredUserMessageService.send({
        name: channel,
        payload: ctx.data,
      });
    });
    sub.on('subscribing', (ctx) => {
      if (process.env.NODE_ENV !== 'production') {
        console.log(`Auth subscribing to ${channel}:`, ctx);
      }
    });

    sub.subscribe();
    this.subscriptions.set(channel, sub);
  }

  public DisconnectFromSubject(data: { subject: string }) {
    if (process.env.NODE_ENV !== 'production') {
      console.log('disconnecting', data);
    }
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

  public disconnect() {
    this.subscriptions.forEach((sub) => {
      sub.unsubscribe();
      this.centrifuge.removeSubscription(sub);
    });
    this.subscriptions.clear();
    this.centrifuge.disconnect();
  }
}

export const centrifugoAuthService = CentrifugoAuthService.getInstance();
