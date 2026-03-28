import { MessageService, MessageNames } from 'services/messageService';

const emailtext = new RegExp(
  '^[a-zA-Z0-9.!#$%&’*+/=?^_`{|}~-]+@[a-zA-Z0-9-]+(?:\\.[a-zA-Z0-9-]+)*$',
);
export const EmailValidator = (properties: {
  uniqueInputId: string;
  value: string;
  errors: string[];
}) => {
  let value = properties.value;
  let inputId = properties.uniqueInputId;
  const sendError = (error: string | null) => {
    MessageService.send({
      name: MessageNames.SET_INPUT_ERROR,
      value: inputId,
      payload: error,
    });
  };
  if (value == null || value.length === 0) {
    sendError(properties.errors[0]);
    return false;
  } else if (!emailtext.test(value)) {
    sendError(properties.errors[1]);
    return false;
  } else {
    sendError(null);
    return true;
  }
};
