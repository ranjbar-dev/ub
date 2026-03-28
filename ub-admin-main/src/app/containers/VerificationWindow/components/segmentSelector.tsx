import { InitialUserDetails } from 'app/containers/UserAccounts/types';
import React, { memo } from 'react';

import ImageAndInfo from './ImageAndInfo';
import { IdentityTypes } from '../constants';
import PermissionsSegment from './PermissionsSegment';

function SegmentSelector(props: {
  activeIndex: number;
  initialData: InitialUserDetails;
}) {
  const { initialData, activeIndex } = props;
  switch (activeIndex) {
    case 0:
      return (
        <ImageAndInfo type={IdentityTypes.Identity} initialData={initialData} />
      );
    case 1:
      return (
        <ImageAndInfo type={IdentityTypes.Address} initialData={initialData} />
      );
    case 2:
      return <PermissionsSegment data={initialData} />;

    default:
      return <div></div>;
  }
}

export default memo(SegmentSelector);
