import { toast } from 'components/Customized/react-toastify';

export const ToastMessages = (messages: any) => {
  for (const key in messages) {
    const fixedKey = key
      .replace('_', ' ')
      .replace('[', '')
      .replace(']', '');
    if (
      messages[key].includes('This value') ||
      messages[key].includes('This field')
    ) {
      const fixedMessage = messages[key]
        .replace('This value', fixedKey)
        .replace('This field', fixedKey);
      toast.error(fixedMessage);
      continue;
    }
    toast.error(messages[key]);
  }
};
