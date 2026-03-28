import AppBar from '@material-ui/core/AppBar';
import Tab from '@material-ui/core/Tab';
import Tabs from '@material-ui/core/Tabs';
import { InitialUserDetails } from 'app/containers/UserAccounts/types';
import { translations } from 'locales/i18n';
import React, { useState, memo, useEffect, useRef } from 'react';
import { useTranslation } from 'react-i18next';
import { Subscriber, MessageNames, BroadcastMessage } from 'services/messageService';
import useDimensions from 'utils/hooks/useDimensions/useDimensions';

import SegmentSelector from './components/segmentSelector';
import SnackBar from './SnackBar';
import { MainTabsWrapper } from '../wrappers/MainTabsWrapper';

/**
 * Tabbed detail panel for a single user, opened in a floating window.
 * Tabs: User Details, Billings, Reports, Orders.
 *
 * @example
 * ```tsx
 * <UserDetailsWindow initialData={{ id: 1, userEmail: 'user@example.com' }} />
 * ```
 */
function UserDetailsWindow(props: { initialData: InitialUserDetails }) {
  const { initialData }: { initialData: InitialUserDetails } = props;

  const [ActiveMainTabIndex, setActiveMainTabIndex] = useState(0);
  const [UserData, setUserData] = useState(initialData);
  const handleChange = (_e: React.ChangeEvent<{}>, value: number) => {
    setActiveMainTabIndex(value);
  };

  useEffect(() => {
    const Subscription = Subscriber.subscribe((message: BroadcastMessage) => {
      if (message.name === MessageNames.SET_USER_DATA) {
        setUserData(message.payload as InitialUserDetails);
      }
    });
    return () => {
      Subscription.unsubscribe();
    };
  }, [UserData]);

  const { t } = useTranslation();

  return (
    <MainTabsWrapper id="UserDetailsWindow">
      <SnackBar userId={UserData.id} />
      <AppBar position="static">
        <Tabs value={ActiveMainTabIndex} onChange={handleChange}>
          <Tab
            className="mainTab"
            label={t(translations.UserAccounts.UserDetails())}
          />
          <Tab
            className="mainTab"
            label={t(translations.UserAccounts.Billings())}
          />
          <Tab
            className="mainTab"
            label={t(translations.UserAccounts.Reports())}
          />
          <Tab
            className="mainTab"
            label={t(translations.UserAccounts.Orders())}
          />
        </Tabs>
      </AppBar>
      {
        <SegmentSelector
          activeIndex={ActiveMainTabIndex}
          initialData={UserData}
        />
      }
    </MainTabsWrapper>
  );
}

export default memo(UserDetailsWindow);
