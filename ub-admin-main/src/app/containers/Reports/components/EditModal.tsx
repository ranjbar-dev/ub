import {
  DialogTitle,
  DialogContent,
  DialogActions,
  Button,
} from '@material-ui/core';
import IsLoadingWithTextAuto from 'app/components/isLoadingWithText/isLoadingWithTextAuto';
import UbTextAreaAutoHeight from 'app/components/UbTextAreaAutoHeight/UbTextAreaAutoHeight';
import { Buttons } from 'app/constants';
import React, { memo, useCallback, useEffect, useRef } from 'react';
import { useDispatch } from 'react-redux';
import { Subscriber, MessageNames, BroadcastMessage } from 'services/messageService';

import { ReportsActions } from '../slice';
import { Report } from '../types';

interface Props {
  onCancel: () => void;
  commentData: Report;
  user_id: number;
}

function EditModal(props: Props) {
  const { onCancel, commentData, user_id } = props;
  const dispatch = useDispatch();
  const value = useRef(commentData.comment);
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
      ReportsActions.EditAdminCommentAction({
        id: commentData.id,
        user_id: user_id,
        comment: value.current,
      }),
    );
  }, []);
  return (
    <div style={{ minWidth: '490px' }}>
      <DialogTitle id="alert-dialog-title">Edit Admin Comment</DialogTitle>
      <DialogContent>
        <UbTextAreaAutoHeight
          initialValue={commentData.comment}
          onChange={handleChange}
        />
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
          className={Buttons.SkyBlueButton}
          text="Update"
          loadingId={'editComentButton' + user_id}
          onClick={handleEdit}
        />
      </DialogActions>
    </div>
  );
}

export default memo(EditModal);
