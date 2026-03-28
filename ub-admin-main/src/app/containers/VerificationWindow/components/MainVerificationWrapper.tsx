import AppBar from '@material-ui/core/AppBar';
import Tab from '@material-ui/core/Tab';
import Tabs from '@material-ui/core/Tabs';
import SnackBar from 'app/components/UserDetailsWindow/SnackBar';
import { MainTabsWrapper } from 'app/components/wrappers/MainTabsWrapper';
import { InitialUserDetails } from 'app/containers/UserAccounts/types';
import { translations } from 'locales/i18n';
import React, { useState, memo, useLayoutEffect } from 'react';
import { useTranslation } from 'react-i18next';
import { LocalStorageKeys } from 'services/constants';

import SegmentSelector from './segmentSelector';


function MainVerificationWrapper(props: { initialData: InitialUserDetails }) {
  const { initialData }: { initialData: InitialUserDetails } = props;

  const [ActiveMainTabIndex, setActiveMainTabIndex] = useState(
    localStorage[LocalStorageKeys.VERIFICATION_WINDOW_TYPE] === 'Address'
      ? 1
      : 0,
  );
  useLayoutEffect(() => {
    localStorage.removeItem(LocalStorageKeys.VERIFICATION_WINDOW_TYPE);
    return () => {};
  }, []);
  const handleChange = (_e: React.ChangeEvent<{}>, value: number) => {
    setActiveMainTabIndex(value);
  };
  const { t } = useTranslation();
  return (
    <MainTabsWrapper className="verificateWindow">
      <SnackBar userId={initialData.id} />
      <AppBar position="static">
        <Tabs value={ActiveMainTabIndex} onChange={handleChange}>
          <Tab
            className="mainTab"
            label={t(translations.CommonTitles.Identity())}
          />
          <Tab
            className="mainTab"
            label={t(translations.CommonTitles.Address())}
          />
          <Tab
            className="mainTab"
            label={t(translations.CommonTitles.Permissions())}
          />
        </Tabs>
      </AppBar>

      <SegmentSelector
        activeIndex={ActiveMainTabIndex}
        initialData={initialData}
      />
    </MainTabsWrapper>
  );
}

export default memo(MainVerificationWrapper);
