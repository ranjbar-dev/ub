import { InitialUserDetails } from 'app/containers/UserAccounts/types';
import React, { memo } from 'react';

interface Props {
  data: InitialUserDetails;
}

function UserLogs(props: Props) {
  const { data } = props;

  return <>UserLogs</>;
}

export default memo(UserLogs);
