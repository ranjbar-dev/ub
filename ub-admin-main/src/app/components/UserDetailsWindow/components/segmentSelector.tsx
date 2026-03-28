import { Billing } from 'app/containers/Billing';
import { Orders } from 'app/containers/Orders';
import { Reports } from 'app/containers/Reports';
import { InitialUserDetails } from 'app/containers/UserAccounts/types';
import { UserDetails } from 'app/containers/UserDetails';
import React, { memo } from 'react';

/**
 * Renders the correct sub-page component for the active tab in UserDetailsWindow.
 *
 * @example
 * ```tsx
 * <SegmentSelector activeIndex={0} initialData={userData} />
 * ```
 */
function SegmentSelector(props: { activeIndex: number; initialData: InitialUserDetails }) {
  const { initialData, activeIndex } = props;
  switch (activeIndex) {
    case 0:
      return <UserDetails initialData={initialData} />;
    case 1:
      return <Billing initialData={initialData} />;
    case 2:
      return <Reports initialData={initialData} />;
    case 3:
      return <Orders initialData={initialData} />;

    default:
      return <div></div>;
  }
}

export default memo(SegmentSelector);
