import { toast } from 'app/components/Customized/react-toastify';
export const ToastMessages = (messages: Record<string, string>) => {
  for (const key in messages) {
    let fixedKey = key.replace('_', ' ').replace('[', '').replace(']', '');
    if (
      messages[key].includes('This value') ||
      messages[key].includes('This field')
    ) {
      let fixedMessage = messages[key]
        .replace('This value', fixedKey)
        .replace('This field', fixedKey);
      toast.error(fixedMessage);
      continue;
    }
    toast.error(messages[key]);
  }
};
