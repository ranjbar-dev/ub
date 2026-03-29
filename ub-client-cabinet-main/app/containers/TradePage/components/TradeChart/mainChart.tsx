import { ChartApiPrefix, LocalStorageKeys } from 'services/constants';
import {
  MessageNames,
  MessageService,
  SideSubscriber,
  Subscriber,
  TradeChartSubscriber,
} from 'services/message_service';
import { CentrifugoChannels, Themes } from 'containers/App/constants';
import React, { memo, useEffect, useMemo, useRef } from 'react';

import { ChartingLibraryWidgetOptions } from 'charting_library/charting_library.min';
import { CentrifugoPublicService } from 'services/CentrifugoPublicService';
import axios from 'axios';
import { prepareSymbolName } from './utils/methods';
import { savedPairName } from 'utils/sharedData';
import { storage } from 'utils/storage';
import styled from 'styles/styled-components';
import { widget } from './charting_library/charting_library.min';
import { widgetOptions } from './constants';

export interface ChartContainerState {}

const greenCandleColor = '#06BA61';
const redCandleColor = '#E64141';

let configs;

const TVChartContainer = () => {
  configs = storage.read(LocalStorageKeys.TRADE_CONFIGS);

  let barsInfo = [],
    newBarInfo,
    lastBarInfo;

  let historyCallBack;

  const chart: any = useRef();

  const selectedPair = useRef(savedPairName());

  const mqtt2 = useRef(CentrifugoPublicService.getInstance());

  const tInterval = useRef(
    localStorage[LocalStorageKeys.TIME_FRAME]
      ? localStorage[LocalStorageKeys.TIME_FRAME]
      : '60',
  );

  const Theme =
    localStorage[LocalStorageKeys.Theme] &&
    localStorage[LocalStorageKeys.Theme] === Themes.LIGHT
      ? 'Light'
      : 'Dark';

  const Pair = selectedPair.current;

  const onRealtimeCallbackFunc: any = useRef();

  const UBWidget: any = useRef();

  const timeFrame = () => {
    switch (tInterval.current) {
      case '1':
      case '3':
        return '1minute';
      case '5':
      case '15':
      case '30':
      case '45':
        return '5minutes';
      case '60':
      case '120':
      case '180':
      case '240':
        return '1hour';
      case '1D':
      case '1W':
      case '1M':
        return '1day';
      default:
        return '1minute';
    }
  };

  const DataFeed = useMemo(() => {
    return {
      onReady: callback => {
        if (configs) {
          setTimeout(() => {
            callback(configs);
          }, 0);
        }
        mqtt2.current.ConnectToSubject({
          subject: `${CentrifugoChannels.TradeChartPrefix}${timeFrame()}:${
            selectedPair.current
          }`,
        });

        axios({
          method: 'GET',
          url: ChartApiPrefix + 'get-configuration',
        })
          .then(res => {
            const { data } = res.data;
            storage.write(LocalStorageKeys.TRADE_CONFIGS, data);
            if (!configs) {
              callback(data);
            }
          })
          .catch(error => {
            console.log(error);
          });
      },
      searchSymbols: (
        userInput,
        exchange,
        symbolType,
        onResultReadyCallback,
      ) => {
        const url = `${ChartApiPrefix}get-search-result?query=${userInput}&limit=12`;
        axios({
          method: 'GET',
          url: url,
        })
          .then(res => {
            onResultReadyCallback(res.data);
          })
          .catch(error => {
            console.log(error);
          });
      },
      resolveSymbol: (
        symbolName,
        onSymbolResolvedCallback,
        onResolveErrorCallback,
      ) => {
        const url = ChartApiPrefix + `get-symbol-info?symbol=${symbolName}`;
        return axios({
          method: 'GET',
          url: url,
        })
          .then(res => {
            setTimeout(() => onSymbolResolvedCallback(res.data), 0);
          })
          .catch(error => {
            setTimeout(() => onResolveErrorCallback(error), 0);
          });
      },
      getBars: (
        symbolInfo,
        resolution,
        from,
        to,
        onHistoryCallback,
        onErrorCallback,
        firstDataRequest,
      ) => {
        historyCallBack = onHistoryCallback;
        let res = resolution;
        if (resolution === 'D') {
          res = '1D';
        }
        const url =
          ChartApiPrefix +
          `get-bars?symbol=${symbolInfo.name}&resolution=${res}&from=${from}&to=${to}`;

        axios({
          method: 'GET',
          url: url,
        })
          .then((response: any) => {
            const { bars } = response.data;
            barsInfo = bars;
            if (bars?.length > 0) {
              MessageService.send({
                name: MessageNames.MAIN_CHART_LAST_PRICE,
                payload: {
                  price: bars[bars.length - 1]?.close,
                },
              });
            }
            let meta = { noData: false };
            // Check if for there is any more bars or not
            if (symbolInfo.symbol_first_bar_timestamp > from) {
              meta = { noData: true };
            }
            //set lastBarInfo
            if (firstDataRequest) {
              lastBarInfo = barsInfo[barsInfo.length - 1];
            }
            onHistoryCallback(barsInfo, meta);
          })
          .catch(error => {
            console.log(error);
          });
      },
      subscribeBars: (
        symbolInfo,
        resolution,
        onRealtimeCallback,
        subscribeUID,
        onResetCacheNeededCallback,
      ) => {
        //  topics[subscribeUID] = prepareTopic(props.pair);
        onRealtimeCallbackFunc.current = onRealtimeCallback;
      },
      unsubscribeBars: subscriberUID => {
        //  onRealtimeCallbackFunc = null;
      },
    };
  }, []);
  //@ts-ignore
  const options: ChartingLibraryWidgetOptions = useMemo(() => {
    return {
      ...widgetOptions,

      theme: Theme,
      symbol: Pair.replace('-', '/'),
      datafeed: DataFeed,
      interval: tInterval.current,
      timezone: Intl.DateTimeFormat().resolvedOptions().timeZone ?? 'exchange',
      overrides: {
        'paneProperties.background':
          Theme === 'Dark' ? 'rgb(28, 28, 33)' : 'rgb(255, 255, 255)',
        'mainSeriesProperties.candleStyle.upColor': greenCandleColor,
        'mainSeriesProperties.candleStyle.downColor': redCandleColor,
        'mainSeriesProperties.candleStyle.wickUpColor': greenCandleColor,
        'mainSeriesProperties.candleStyle.wickDownColor': redCandleColor,
        'mainSeriesProperties.candleStyle.borderColor': greenCandleColor,
        'mainSeriesProperties.candleStyle.borderUpColor': greenCandleColor,
        'mainSeriesProperties.candleStyle.borderDownColor': redCandleColor,
        'scalesProperties.fontSize': 12,
        'scalesProperties.textColor': 'rgb(101,101,101)',
        'paneProperties.vertGridProperties.color':
          Theme === 'Dark' ? '#222325' : '#f5f5f5',
        'paneProperties.horzGridProperties.color':
          Theme === 'Dark' ? '#222325' : '#f5f5f5',
        'mainSeriesProperties.style': 1,
        volumePaneSize: 'large',
      },
      favorites: {
        intervals: [],
        chartTypes: ['candle'],
      },
      disabled_features: [
        'header_screenshot',
        'header_settings',
        'header_saveload',
        'timeframes_toolbar',
        'header_compare',
        'header_symbol_search',
        //'create_volume_indicator_by_default',
      ],
      enabled_features: ['hide_left_toolbar_by_default', 'items_favoriting'],
      custom_css_url: './custome.css',
    };
  }, [tInterval.current]);
  useEffect(() => {
    //@ts-ignore
    UBWidget.current = new widget(options);
    const Subscription = Subscriber.subscribe((message: any) => {
      if (message.name === MessageNames.MAIN_CHART_SUMMARY) {
        if (message.payload.payload) {
          newBarInfo = prepareNewBarInfo(message.payload.payload);
          if (onRealtimeCallbackFunc.current) {
            onRealtimeCallbackFunc.current(newBarInfo);
            //historyCallBack(...barsInfo, newBarInfo);
          }
        }
      }
      if (message.name === MessageNames.CHANGE_THEME) {
        //UBWidget.changeTheme()
        if (message.payload === Themes.DARK) {
          UBWidget.current.changeTheme('Dark');
          UBWidget.current.applyOverrides({
            'paneProperties.background': 'rgb(28, 28, 33)',
            'paneProperties.vertGridProperties.color': '#222325',
            'paneProperties.horzGridProperties.color': '#222325',
          });
          return;
        }
        UBWidget.current.changeTheme('Light');
        UBWidget.current.applyOverrides({
          'paneProperties.background': 'rgb(255, 255, 255)',
          'paneProperties.vertGridProperties.color': '#f5f5f5',
          'paneProperties.horzGridProperties.color': '#f5f5f5',
        });
      }
    });

    const SideSubscription = SideSubscriber.subscribe((message: any) => {
      if (message.name === MessageNames.SET_TRADE_PAGE_CURRENCY_PAIR) {
        UBWidget.current.setSymbol(
          message.payload.name.replace('-', '/'),
          tInterval.current,
          null,
        );
      }
    });
    const TradeChartSubscription = TradeChartSubscriber.subscribe(
      (message: any) => {
        if (
          //!document.hidden &&
          message.name.includes(
            `${CentrifugoChannels.TradeChartPrefix}${timeFrame()}:${
              selectedPair.current
            }`,
          )
        ) {
          MessageService.send({
            name: MessageNames.MAIN_CHART_SUMMARY,
            payload: message,
          });
        }
      },
    );
    UBWidget.current.onChartReady(() => {
      //@ts-ignore
      //  UBWidget.current._options.preset = 'mobile';
      chart.current = UBWidget.current.chart();
      //  const button = UBWidget.current.createButton();
      //  button.setAttribute('title', 'Click to show a notification popup');
      //  button.classList.add('apply-common-tooltip');
      //  button.addEventListener('click', () =>
      //    UBWidget.current.showNoticeDialog({
      //      title: 'Notification',
      //      body: 'TradingView Charting Library API works correctly',
      //      callback: () => {
      //        console.log('Noticed!');
      //      },
      //    }),
      //  );
      //  button.innerHTML = 'Check API';
      //
      chart.current
        .onIntervalChanged()
        .subscribe(null, function (interval, obj) {
          mqtt2.current.DisconnectFromSubject({
            subject: `${CentrifugoChannels.TradeChartPrefix}${timeFrame()}:${
              selectedPair.current
            }`,
          });

          tInterval.current = interval;

          localStorage[LocalStorageKeys.TIME_FRAME] = interval;
          mqtt2.current.ConnectToSubject({
            subject: `${CentrifugoChannels.TradeChartPrefix}${timeFrame()}:${
              selectedPair.current
            }`,
          });
        });

      chart.current
        .onSymbolChanged()
        .subscribe(null, function (newSymbol, obj) {
          const symbol = newSymbol.name;
          const [symbolDependent, symbolBase] = prepareSymbolName(symbol);
          mqtt2.current.DisconnectFromSubject({
            subject: `${CentrifugoChannels.TradeChartPrefix}${timeFrame()}:${
              selectedPair.current
            }`,
          });
          selectedPair.current = symbolDependent + '-' + symbolBase;
          mqtt2.current.ConnectToSubject({
            subject: `${CentrifugoChannels.TradeChartPrefix}${timeFrame()}:${
              selectedPair.current
            }`,
          });
        });
    });

    return () => {
      mqtt2.current.DisconnectFromSubject({
        subject: `${CentrifugoChannels.TradeChartPrefix}${timeFrame()}:${
          selectedPair.current
        }`,
      });
      Subscription.unsubscribe();
      SideSubscription.unsubscribe();
      TradeChartSubscription.unsubscribe();
      if (UBWidget.current) {
        UBWidget.current.remove();
      }
      historyCallBack = undefined;
      barsInfo = [];
    };
  }, []);

  return (
    <MainWrapper>
      <Wrapper id='mainChartWrapper'>
        <div id='ubChart' />
      </Wrapper>
    </MainWrapper>
  );
};

export default memo(TVChartContainer, () => true);
const MainWrapper = styled.div`
  height: 100%;
  position: absolute;
  width: 100%;
  overflow: hidden;
`;
const Wrapper = styled.div`
  width: 100%;
  height: 100%;
  border-bottom-right-radius: var(--cardBorderRadius);
  border-bottom-left-radius: var(--cardBorderRadius);
  overflow: hidden;

  #ubChart {
    width: calc(100% + 4px);
    height: calc(100% + 2px);
    border-radius: 7px;
    overflow: hidden;
  }
`;
const prepareNewBarInfo = msg => {
  //  msg = msg[0];
  const lastBar = {
    time: prepareTimestamp(msg.startTime),
    open: msg.openPrice,
    high: msg.highPrice,
    low: msg.lowPrice,
    close: msg.closePrice,
    volume: +msg.baseVolume,
  };
  return lastBar;
};

const prepareTimestamp = (strDate: string) => {
  let offset = new Date().getTimezoneOffset();
  offset = offset * 60 * 1000;

  const datum = new Date(strDate);
  return datum.getTime() - offset;
};
