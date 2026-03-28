import React from 'react';
import translate from './messages';
import { MessageService, MessageNames } from 'services/message_service';
import { FormattedMessage } from 'react-intl';
import { strongRegex } from './regex';

export const PasswordValidator = (properties: {
  uniqueInputId: string;
  value: string;
  compareWithInputWithId?: string;
  isMainPassword?: boolean;
  onMainPasswordChange?: (notEqual: boolean) => void;
}) => {
  const value = properties.value;
  const inputId = properties.uniqueInputId;
  const { onMainPasswordChange } = properties;
  if (properties.compareWithInputWithId && !properties.isMainPassword) {
    const inp: any = document.getElementById(properties.compareWithInputWithId);
    if (inp && inp.value !== '' && inp.value !== value) {
      MessageService.send({
        name: MessageNames.SET_INPUT_ERROR,
        value: inputId,
        payload: <FormattedMessage {...translate.oldAndNewPassword} />,
      });
      return false;
    }
  }
  if (properties.compareWithInputWithId && properties.isMainPassword === true) {
    const inp: any = document.getElementById(properties.compareWithInputWithId);
    if (inp && inp.value !== '' && inp.value !== value) {
      MessageService.send({
        name: MessageNames.SET_INPUT_ERROR,
        value: properties.compareWithInputWithId,
        payload: <FormattedMessage {...translate.oldAndNewPassword} />,
      });
      //@ts-ignore
      onMainPasswordChange(false);
      return true;
    }
  }
  if (!strongRegex.test(value)) {
    MessageService.send({
      name: MessageNames.SET_INPUT_ERROR,
      value: inputId,
      payload: <FormattedMessage {...translate.strongPasswordError} />,
      additional: true,
    });
    return false;
  }
  if (value == null || value.length === 0) {
    MessageService.send({
      name: MessageNames.SET_INPUT_ERROR,
      value: inputId,
      payload: <FormattedMessage {...translate.passwordIsRequired} />,
    });
    return false;
  } else if (value.length < 8) {
    MessageService.send({
      name: MessageNames.SET_INPUT_ERROR,
      value: inputId,
      payload: <FormattedMessage {...translate.minimum8Character} />,
    });
    return false;
  } else {
    MessageService.send({
      name: MessageNames.SET_INPUT_ERROR,
      value: inputId,
      payload: null,
    });
    return true;
  }
};
