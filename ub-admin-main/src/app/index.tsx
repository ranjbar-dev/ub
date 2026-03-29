/**
 *
 * App
 *
 * This component is the skeleton around the actual pages, and should only
 * contain code that should be seen on all pages. (e.g. navigation bar)
 */

import * as React from 'react';
import { Helmet } from 'react-helmet-async';
import { useTranslation } from 'react-i18next';
import { useSelector, useDispatch } from 'react-redux';
import { Switch, Route } from 'react-router-dom';
import PrivateRoute from './components/PrivateRoute';
import { LocalStorageKeys } from 'services/constants';
import { MessageService, MessageNames, Subscriber, BroadcastMessage } from 'services/messageService';
import downloadFile from 'utils/fileDownload';
import { history } from 'utils/history';
import {useInjectReducer, useInjectSaga} from 'utils/redux-injectors';

import { selectRouter } from './appSelectors';
import { AppPages } from './constants';
import { GlobalStyle } from '../styles/global-styles';
import 'styles/ag-grid.min.css';
import 'styles/ag-theme-balham.min.css';
import './components/Customized/react-toastify/dist/ReactToastify.css';

import { SideNav } from './components/sideNav';
import { Admins } from './containers/Admins';
import { Balances } from './containers/Balances';
import { CurrencyPairs } from './containers/CurrencyPairs';
import { Deposits } from './containers/Deposits';
import { ExternalExchange } from './containers/ExternalExchange';
import { ExternalOrders } from './containers/ExternalOrders';
import { FilledOrders } from './containers/FilledOrders';
import { FinanceMethods } from './containers/FinanceMethods';
import { HomePage } from './containers/HomePage';
import { LiquidityOrders } from './containers/LiquidityOrders';
import { LoginHistory } from './containers/LoginHistory';
import { LoginPage } from './containers/LoginPage';
import { MarketTicks } from './containers/MarketTicks';
import { NavBar } from './containers/NavBar';
import { NotFoundPage } from './containers/NotFoundPage';

import { ConnectedRouter, replace } from 'connected-react-router';
import { globalActions } from 'store/slice';
import { translations } from 'locales/i18n';

import { OpenOrders } from './containers/OpenOrders';
import { ScanBlock } from './containers/ScanBlock';
import { UserAccounts } from './containers/UserAccounts';

import UBWindow from 'app/components/newWindow/NewWindow';

import {userAccountsSaga} from './containers/UserAccounts/saga';
import {sliceKey, UserAccountsReducer} from './containers/UserAccounts/slice';
import { Withdrawals } from './containers/Withdrawals';

import { useEffect, useMemo } from 'react';

import ForceStyles from './ForceStyles';

import PersonOutlineIcon from '@material-ui/icons/PersonOutline';
import AssignmentIcon from '@material-ui/icons/Assignment';
import SupervisorAccountIcon from '@material-ui/icons/SupervisorAccount';
import SettingsIcon from '@material-ui/icons/Settings';
import VisibilityIcon from '@material-ui/icons/Visibility';

import NewWindowContainer from './NewWindowContainer';


let timeOut: ReturnType<typeof setTimeout> | undefined;

export function App() {
//this part is written here, because we will need to access user account reducer and saga from any place in the app
useInjectReducer({ key: sliceKey, reducer: UserAccountsReducer });
useInjectSaga({ key: sliceKey, saga: userAccountsSaga });

	const dispatch = useDispatch();

	useEffect(() => {
		const subscription = Subscriber.subscribe((message: BroadcastMessage) => {
			if (message.name === MessageNames.DOWNLOAD_FILE) {
				downloadFile(message.payload as { url: string; filename: string })
			}
			if (message.name === MessageNames.AUTH_ERROR_EVENT) {
				dispatch(globalActions.setIsLoggedIn(false));
				dispatch(replace(AppPages.RootPage));
			}
		})
		return () => {
			subscription.unsubscribe()
		}
	}, [dispatch])

	useEffect(() => {
		const handleResize = (e: UIEvent) => {
			clearTimeout(timeOut);
			timeOut = setTimeout(() => {
				MessageService.send({
					name: MessageNames.RESIZE,
					payload: (e.target as Window).innerWidth,
				});
			}, 150);
		};
		window.addEventListener('resize', handleResize);
		return () => {
			window.removeEventListener('resize', handleResize);
			clearTimeout(timeOut);
		};
	}, []);
	const router = useSelector(selectRouter);
	const { t } = useTranslation();

	const Categories = useMemo(
		() => [
			{
				name: t(translations.PageNames.UserManagement()),
				icon: <PersonOutlineIcon />,
				childs: [
					{
						name: t(translations.PageNames.UserAcounts()),
						page: AppPages.UserAccounts,
					},

					{
						name: t(translations.PageNames.Verification()),
						page: AppPages.UserAccounts + 'verification',
					},
					{
						name: t(translations.PageNames.Groups()),
					},
					{
						name: t(translations.PageNames.LoginHistory()),
						page: AppPages.LoginHistory,
					},
					//  {
					//    name: 'placeholder',
					//    page: AppPages.PlaceHolder,
					//  },
				],
			},
			{
				name: t(translations.PageNames.OrderManagement()),
				icon: <AssignmentIcon />,
				childs: [
					{
						name: t(translations.PageNames.OpenOrders()),
						page: AppPages.OpenOrders,
					},
					{
						name: t(translations.PageNames.FilledOrders()),
						page: AppPages.FilledOrders,
					},
					{
						name: t(translations.PageNames.ExternalOrders()),
						page: AppPages.ExternalOrders,
					},
					{
						name: t(translations.CommonTitles.LiquidityOrders()),
						page: AppPages.LiquidityOrders,
					},
				],
			},
			{
				name: t(translations.PageNames.Accounting()),
				icon: <SupervisorAccountIcon />,
				childs: [
					{
						name: t(translations.PageNames.Deposits()),
						page: AppPages.Deposits,
					},
					{
						name: t(translations.PageNames.Withdrawals()),
						page: AppPages.Withdrawals,
					},
					{
						name: t(translations.PageNames.Balances()),
						page: AppPages.Balances,
					},
					{
						name: t(translations.PageNames.ScanBlock()),
						page: AppPages.ScanBlock,
					},
				],
			},
			{
				name: t(translations.PageNames.Configuration()),
				icon: <SettingsIcon />,
				childs: [
					{
						name: t(translations.PageNames.FinanceMethods()),
						page: AppPages.FinanceMethods,
					},
					{
						name: t(translations.PageNames.CurrencyPairs()),
						page: AppPages.CurrencyPairs,
					},
					{
						name: t(translations.PageNames.ExternalExchange()),
						page: AppPages.ExternalExchange,
					},
					{
						name: t(translations.PageNames.MarketTicks()),
						page: AppPages.MarketTicks,
					},
				],
			},
			{
				name: t(translations.PageNames.Administration()),
				icon: <VisibilityIcon />,
				childs: [
					{
						name: t(translations.PageNames.Admins()),
						page: AppPages.Admins,
					},
					{
						name: t(translations.PageNames.AdminRules()),
					},
					{
						name: t(translations.PageNames.Logs()),
					},
				],
			},
		],

		[]
	);
	const ShowSideNav = (): boolean => {
		if (
			router.location.pathname !== AppPages.RootPage &&
			localStorage[LocalStorageKeys.ACCESS_TOKEN]
		) {
			return true;
		}
		return false;
	};

	return (
		<ConnectedRouter history={history}>
			{useMemo(
				() => (
					<ForceStyles />
				),
				[],
			)}
			{useMemo(
				() => (
					<NewWindowContainer />
				),
				[],
			)}
			<Helmet titleTemplate="%s - unitedBit admin" defaultTitle="UB-Admin">
				<meta name="description" content="unitedBit" />
			</Helmet>
			<UBWindow />
			{ShowSideNav() === true ? <SideNav mainCategories={Categories} /> : ''}
			{ShowSideNav() === true ? <NavBar /> : ''}
			<Switch>
				<Route exact path={AppPages.RootPage} component={LoginPage} />
				<PrivateRoute path={AppPages.HomePage} component={HomePage} />
				<PrivateRoute path={AppPages.UserAccounts} component={UserAccounts} />
				<PrivateRoute path={AppPages.LoginHistory} component={LoginHistory} />
				<PrivateRoute path={AppPages.OpenOrders} component={OpenOrders} />
				<PrivateRoute path={AppPages.FilledOrders} component={FilledOrders} />
				<PrivateRoute path={AppPages.ExternalOrders} component={ExternalOrders} />
				<PrivateRoute path={AppPages.Deposits} component={Deposits} />
				<PrivateRoute path={AppPages.Withdrawals} component={Withdrawals} />
				<PrivateRoute path={AppPages.FinanceMethods} component={FinanceMethods} />
				<PrivateRoute path={AppPages.CurrencyPairs} component={CurrencyPairs} />
				<PrivateRoute path={AppPages.ExternalExchange} component={ExternalExchange} />
				<PrivateRoute path={AppPages.MarketTicks} component={MarketTicks} />
				<PrivateRoute path={AppPages.Balances} component={Balances} />
				<PrivateRoute path={AppPages.ScanBlock} component={ScanBlock} />
				<PrivateRoute path={AppPages.LiquidityOrders} component={LiquidityOrders} />
				{/*<PrivateRoute path={AppPages.PlaceHolder} component={PlaceHolder} />*/}
				<PrivateRoute path={AppPages.Admins} component={Admins} allowedRoles={['superadmin']} />
				<Route component={NotFoundPage} />
			</Switch>
			<GlobalStyle />
		</ConnectedRouter>
	);
}
