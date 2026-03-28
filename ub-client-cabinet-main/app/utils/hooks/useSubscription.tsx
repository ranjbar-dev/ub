import { useEffect } from 'react';
import { Subject } from 'rxjs';
import { MessageNames } from 'services/message_service';

interface Props {
  onRecieve: Function;
  subscriber: Subject<unknown>;
  messageName: MessageNames;
  dependencies?: any[];
  id?: string;
}
const useSubscription = (props: Props) => {
  const { onRecieve, dependencies, id, messageName, subscriber } = props;
  useEffect(() => {
    const Subscription = subscriber.subscribe((message: any) => {
      if (message.name === messageName && (!id || id === message.id)) {
        onRecieve(message.payload);
      }
    });
    return () => {
      Subscription.unsubscribe();
    };
  }, dependencies ?? []);
  return;
};
export default useSubscription;
