/*
 *
 * TradeChart constants
 *
 */
import {ChartingLibraryWidgetOptions} from './charting_library/charting_library.min';
import {getLanguageFromURL} from './utils/methods';

enum ActionTypes {
	DEFAULT_ACTION='app/TradeChart/DEFAULT_ACTION',
	GET_CHART_CONFIG='app/TradeChart/GET_CHART_CONFIG',
	SET_CHART_CONFIG='app/TradeChart/SET_CHART_CONFIG',
}
export default ActionTypes;
//@ts-ignore
export const widgetOptions: ChartingLibraryWidgetOptions={
	container_id: 'ubChart',
	library_path: '/charting_library/',
	locale: getLanguageFromURL()||'en',
	disabled_features: ['use_localstorage_for_settings'],
	charts_storage_url: 'https://saveload.tradingview.com',
	charts_storage_api_version: '1.1',
	client_id: 'tradingview.com',
	user_id: 'public_user_id',
	fullscreen: false,
	autosize: true,
	studies_overrides: {},
};
