/*
 * validator Messages
 *
 * This contains all the text for the LoginPage container.
 */

import {defineMessages} from 'react-intl';

export const scope='validationMessages';

export default defineMessages({
	emailIsRequired: {
		id: `${scope}.emailIsRequired`,
		defaultMessage: 'ET.emailIsRequired',
	},
	emailIsNotValid: {
		id: `${scope}.emailIsNotValid`,
		defaultMessage: 'ET.emailIsNotValid',
	},
	passwordIsRequired: {
		id: `${scope}.passwordIsRequired`,
		defaultMessage: 'ET.passwordIsRequired',
	},
	minimum8Character: {
		id: `${scope}.minimum8Character`,
		defaultMessage: 'ET.minimum8Character',
	},
	upper: {
		id: `${scope}.upper`,
		defaultMessage: 'ET.upper',
	},
	number: {
		id: `${scope}.number`,
		defaultMessage: 'ET.number',
	},
	special: {
		id: `${scope}.special`,
		defaultMessage: 'ET.special',
	},
	oldAndNewPassword: {
		id: `${scope}.oldAndNewPassword`,
		defaultMessage: 'ET.oldAndNewPassword',
	},
	strongPasswordError: {
		id: `${scope}.strongPasswordError`,
		defaultMessage: 'ET.strongPasswordError',
	},
});
