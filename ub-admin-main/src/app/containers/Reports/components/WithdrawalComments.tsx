import PostAddOutlinedIcon from '@material-ui/icons/PostAddOutlined';
import { ColDef } from 'ag-grid-community';
import IsLoadingWithTextAuto from 'app/components/isLoadingWithText/isLoadingWithTextAuto';
import { SimpleGrid } from 'app/components/SimpleGrid/SimpleGrid';
import UbTextAreaAutoHeight from 'app/components/UbTextAreaAutoHeight/UbTextAreaAutoHeight';
import { billingSaga } from 'app/containers/Billing/saga';
import { BillingActions, BillingReducer, sliceKey as billingSliceKey } from 'app/containers/Billing/slice';
import { InitialUserDetails } from 'app/containers/UserAccounts/types';
import { translations } from 'locales/i18n';
import React, { memo, useMemo, useRef } from 'react';
import { useTranslation } from 'react-i18next';
import { useDispatch, useSelector } from 'react-redux';
import { useInjectReducer, useInjectSaga } from 'utils/redux-injectors';

import { reportsSaga } from '../saga';
import { ReportsActions, ReportsReducer, sliceKey } from '../slice';
import { selectWithdrawalComments } from '../selectors';



interface Props {
  data: InitialUserDetails;
}
function WithdrawalComments(props: Props) {

  const dispatch = useDispatch()

  useInjectReducer({ key: sliceKey, reducer: ReportsReducer });
  useInjectSaga({ key: sliceKey, saga: reportsSaga });

  useInjectReducer({ key: billingSliceKey, reducer: BillingReducer });
  useInjectSaga({ key: billingSliceKey, saga: billingSaga });

  const withdrawalComments = useSelector(selectWithdrawalComments);
  const newComment = useRef<string>('')

  const { data } = props;
  const { t } = useTranslation();
  const staticRows: ColDef[] = useMemo(
    () => [
      {
        headerName: t(translations.CommonTitles.IDNO()),
        field: 'id',
        maxWidth: 80,
      },
      {
        headerName: t(translations.CommonTitles.Admin()),
        field: 'adminName',
        maxWidth: 140,
      },
      {
        headerName: t(translations.CommonTitles.Date()),
        field: 'date',
        maxWidth: 180,
      },
      {
        headerName: t(translations.CommonTitles.Message()),
        field: 'comment',
      },
    ],
    [],
  );
  const handleCommentChange = (e: string) => {
    newComment.current = e
  }
  const handleSubmitClick = () => {
    if (newComment.current === '') {
      return;
    }
    const sendingData = {
      payment_id: data.id,
      comment: newComment.current,
    }
    dispatch(
      BillingActions.AddPaymentCommentAction(sendingData),
    );

  }
  return (
    <div style={{ width: '100%' }}>

      {/* <div className="adminReportsInput" style={{ padding: '12px' }}>
        <UbTextAreaAutoHeight
          placeHolder="Write Your comment ... "
          initialValue=""
          style={{
            marginBottom: '12px',
            height: '50px',
            background: '#FBFBFB',
            borderRadius: '5px',
            border: '1px solid #C5C5C5',
          }}
          onChange={handleCommentChange}
        />
        <IsLoadingWithTextAuto
          icon={<PostAddOutlinedIcon />}
          style={{ background: '#649FFF' }}
          text="Add Comment"
          loadingId={'addWithdrawComment' + data.id}
          onClick={handleSubmitClick}
        />
      </div> */}

      <SimpleGrid
        containerId="UserDetailsWindow"
        additionalInitialParams={{ user_id: data.id }}
        arrayFieldName="comments"
        immutableId="id"
        userId={data.id}
        //onRowClick={handleRowClick}
        initialAction={ReportsActions.GetWithdrawalCommentsAction}
        externalData={withdrawalComments}
        staticRows={staticRows}
      />
    </div>
  );
}

export default memo(WithdrawalComments);
