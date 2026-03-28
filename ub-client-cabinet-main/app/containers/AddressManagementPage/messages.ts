/*
 * AddressManagementPage Messages
 *
 * This contains all the text for the AddressManagementPage container.
 */

import {defineMessages} from 'react-intl';

import {GlobalTranslateScope} from 'containers/App/constants';
export const scope='containers.AddressManagementPage';

export default defineMessages({
	Withdrawaddress: {
		id: `${scope}.Withdrawaddress`,
		defaultMessage: 'ET_Withdrawaddress',
	},
	Createaddress: {
		id: `${scope}.Createaddress`,
		defaultMessage: 'ET_Createaddress',
	},
	address: {
		id: `${scope}.address`,
		defaultMessage: 'ET_address',
	},
	label: {
		id: `${scope}.label`,
		defaultMessage: 'ET_label',
	},
	Yourhavenowithdrawaddress: {
		id: `${scope}.Yourhavenowithdrawaddress`,
		defaultMessage: 'You have no withdraw address',
	},
	Pleasecreateaddressandwithdrawcoins: {
		id: `${scope}.Pleasecreateaddressandwithdrawcoins`,
		defaultMessage: 'Please create address and withdraw coins',
	},
	//////////////
	create: {
		id: `${GlobalTranslateScope}.create`,
		defaultMessage: 'ET_create',
	},
	coin: {
		id: `${GlobalTranslateScope}.coin`,
		defaultMessage: 'ET_coin',
	},
	network: {
		id: `${GlobalTranslateScope}.network`,
		defaultMessage: 'et.network',
	},
	all: {
		id: `${GlobalTranslateScope}.all`,
		defaultMessage: 'ET_all',
	},
});
