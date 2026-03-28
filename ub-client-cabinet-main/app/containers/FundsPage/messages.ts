/*
 * FundsPage Messages
 *
 * This contains all the text for the FundsPage container.
 */

import {defineMessages} from 'react-intl';
import {GlobalTranslateScope} from 'containers/App/constants';

export const scope='containers.FundsPage';
const globalScope=GlobalTranslateScope;

export default defineMessages({
	header: {
		id: `${scope}.header`,
		defaultMessage: 'This is the FundsPage container!',
	},

	balance: {
		id: `${scope}.balance`,
		defaultMessage: 'ET_balance',
	},
	deposite: {
		id: `${scope}.deposite`,
		defaultMessage: 'ET_deposite',
	},
	withdrawals: {
		id: `${scope}.withdrawals`,
		defaultMessage: 'ET_withdrawals',
	},
	transactionHistory: {
		id: `${scope}.transactionHistory`,
		defaultMessage: 'ET_transactionHistory',
	},
	EstimatedBalance: {
		id: `${scope}.EstimatedBalance`,
		defaultMessage: 'ET_EstimatedBalance',
	},
	AvailableBalance: {
		id: `${scope}.AvailableBalance`,
		defaultMessage: 'ET_AvailableBalance',
	},
	InOrders: {
		id: `${scope}.InOrders`,
		defaultMessage: 'ET_InOrders',
	},
	show: {
		id: `${scope}.show`,
		defaultMessage: 'ET_show',
	},
	hide: {
		id: `${scope}.hide`,
		defaultMessage: 'ET_hide',
	},
	Smallbalances: {
		id: `${scope}.Smallbalances`,
		defaultMessage: 'ET_Smallbalances',
	},
	Exchangeaccoount: {
		id: `${scope}.Exchangeaccoount`,
		defaultMessage: 'ET_Exchangeaccoount',
	},
	DEPOSITSADDRESS: {
		id: `${scope}.DEPOSITSADDRESS`,
		defaultMessage: 'ET_DEPOSITSADDRESS',
	},
	DEPOSITSHISTORY: {
		id: `${scope}.DEPOSITSHISTORY`,
		defaultMessage: 'ET_DEPOSITSHISTORY',
	},
	sendyour: {
		id: `${scope}.sendyour`,
		defaultMessage: 'ET_sendyour',
	},
	tothisaddress: {
		id: `${scope}.tothisaddress`,
		defaultMessage: 'ET_tothisaddress',
	},
	CopyAddress: {
		id: `${scope}.CopyAddress`,
		defaultMessage: 'ET_CopyAddress',
	},
	ShowQRcode: {
		id: `${scope}.ShowQRcode`,
		defaultMessage: 'ET_ShowQRcode',
	},
	Getnewaddress: {
		id: `${scope}.Getnewaddress`,
		defaultMessage: 'ET_Getnewaddress',
	},
	Pleaseselectanycointogetdepositaddress: {
		id: `${scope}.Pleaseselectanycointogetdepositaddress`,
		defaultMessage: 'ET_Pleaseselectanycointogetdepositaddress',
	},
	Pleaseselectanycointowithdraw: {
		id: `${scope}.Pleaseselectanycointowithdraw`,
		defaultMessage: 'ET_Pleaseselectanycointowithdraw',
	},
	addressCopiedToClipboard: {
		id: `${scope}.addressCopiedToClipboard`,
		defaultMessage: 'ET_addressCopiedToClipboard',
	},
	withdrawalAddress: {
		id: `${scope}.withdrawalAddress`,
		defaultMessage: 'ET_withdrawalAddress',
	},
	withdrawalHistory: {
		id: `${scope}.withdrawalHistory`,
		defaultMessage: 'ET_withdrawalHistory',
	},
	Totalbalance: {
		id: `${scope}.Totalbalance`,
		defaultMessage: 'ET_Totalbalance',
	},
	Inorder: {
		id: `${scope}.Inorder`,
		defaultMessage: 'ET_Inorder',
	},
	dontWithdraw: {
		id: `${scope}.dontWithdraw`,
		defaultMessage: 'ET_dontWithdraw',
	},
	Minimumwithdrawal: {
		id: `${scope}.Minimumwithdrawal`,
		defaultMessage: 'ET_Minimumwithdrawal',
	},
	TransactionFee: {
		id: `${scope}.TransactionFee`,
		defaultMessage: 'ET_TransactionFee',
	},
	YouWillGet: {
		id: `${scope}.YouWillGet`,
		defaultMessage: 'ET_YouWillGet',
	},
	Available: {
		id: `${scope}.Available`,
		defaultMessage: 'ET_Available',
	},
	maximumWithdrawAmountIs: {
		id: `${scope}.maximumWithdrawAmountIs`,
		defaultMessage: 'ET.maximumWithdrawAmountIs',
	},
	minimumWithdrawAmountIs: {
		id: `${scope}.minimumWithdrawAmountIs`,
		defaultMessage: 'ET.minimumWithdrawAmountIs',
	},
	transferNetwork: {
		id: `${scope}.transferNetwork`,
		defaultMessage: 'ET.transferNetwork',
	},
	/////////
	coin: {
		id: `${globalScope}.coin`,
		defaultMessage: 'ET_coin',
	},
	submit: {
		id: `${globalScope}.submit`,
		defaultMessage: 'ET_submit',
	},
	next: {
		id: `${globalScope}.next`,
		defaultMessage: 'ET.next',
	},
	all: {
		id: `${globalScope}.all`,
		defaultMessage: 'ET_all',
	},
	EnterCodeHere: {
		id: `${globalScope}.EnterCodeHere`,
		defaultMessage: 'ET.EnterCodeHere',
	},
	PleaseopenGoogleAuthenticatorappinyourphoneandenter2FaCode: {
		id: `${globalScope}.PleaseopenGoogleAuthenticatorappinyourphoneandenter2FaCode`,
		defaultMessage:
			'ET.PleaseopenGoogleAuthenticatorappinyourphoneandenter2FaCode',
	},
	EnterVerificationEmail: {
		id: `${globalScope}.EnterVerificationEmail`,
		defaultMessage:
			'ET.EnterVerificationEmail',
	},
});
