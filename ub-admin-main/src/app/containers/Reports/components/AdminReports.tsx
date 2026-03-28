import { GridLoading } from 'app/components/grid_loading/gridLoading';
import { InitialUserDetails } from 'app/containers/UserAccounts/types';
import React, { memo, useEffect } from 'react';
import { useDispatch, useSelector } from 'react-redux';

import { ReportsActions } from '../slice';
import { selectAdminReportsData } from '../selectors';
import AdminReportsInput from './AdminReportsInput';
import Comments from './Comments';

interface Props {
  data: InitialUserDetails;
}

function AdminReports(props: Props) {
  const dispatch = useDispatch();
  const { data } = props;
  const adminReports = useSelector(selectAdminReportsData);
  const IsLoading = adminReports === null;

  useEffect(() => {
    dispatch(ReportsActions.GetAdminReportsAction({ id: data.id }));
  }, []);

  return (
    <div className="AdminReports__Wrapper">
      {IsLoading === true && <GridLoading />}
      <AdminReportsInput userId={data.id} />
      <Comments comments={adminReports || []} user_id={data.id} />
    </div>
  );
}

export default memo(AdminReports);
