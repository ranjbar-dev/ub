import {
  DialogTitle,
  DialogContent,
  DialogContentText,
  DialogActions,
  Button,
} from '@material-ui/core';
import IsLoadingWithTextAuto from 'app/components/isLoadingWithText/isLoadingWithTextAuto';
import UbTextAreaAutoHeight from 'app/components/UbTextAreaAutoHeight/UbTextAreaAutoHeight';
import { Buttons } from 'app/constants';
import React, { memo, useCallback, useEffect, useRef } from 'react';
import { useDispatch } from 'react-redux';
import { Subscriber, MessageNames, BroadcastMessage } from 'services/messageService';

import { VerificationWindowActions } from '../slice';

interface Props {
  onCancel: () => void;
  imageId: number;
  type: string;
  user_id: number;
}

function RejectModal(props: Props) {
  const { onCancel, user_id, imageId, type } = props;
  const dispatch = useDispatch();
  const value = useRef('');
  const handleChange = useCallback(e => {
    value.current = e;
  }, []);
  useEffect(() => {
    const Subscription = Subscriber.subscribe((message: BroadcastMessage) => {
      if (message.name === MessageNames.CLOSE_POPUP) {
        onCancel();
      }
    });
    return () => {
      Subscription.unsubscribe();
    };
  }, []);
  const handleEdit = useCallback(() => {
    dispatch(
      VerificationWindowActions.UpdateProfileImageStatusAction({
        confirmation_status: 'rejected',
        id: imageId,
        user_id,
        type,
        loadingButtonId: 'RejectProfileImage' + user_id,
        rejection_reason: value.current,
      }),
    );
  }, []);
  return (
    <div style={{ minWidth: '490px' }}>
      <DialogTitle id="alert-dialog-title">Reject Reason</DialogTitle>
      <DialogContent>
        <UbTextAreaAutoHeight initialValue={''} onChange={handleChange} />
      </DialogContent>
      <DialogActions>
        <Button
          onClick={onCancel}
          color="primary"
          className={Buttons.BlackButton}
        >
          Cancel
        </Button>
        <IsLoadingWithTextAuto
          className={Buttons.RedButton}
          text="Reject"
          loadingId={'RejectProfileImage' + user_id}
          onClick={handleEdit}
        />
      </DialogActions>
    </div>
  );
}

export default memo(RejectModal);
