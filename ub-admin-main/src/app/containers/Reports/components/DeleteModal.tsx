import {
  DialogTitle,
  DialogContent,
  DialogContentText,
  DialogActions,
  Button,
} from '@material-ui/core';
import IsLoadingWithTextAuto from 'app/components/isLoadingWithText/isLoadingWithTextAuto';
import { Buttons } from 'app/constants';
import React, { memo, useCallback, useEffect } from 'react';
import { useDispatch } from 'react-redux';
import { Subscriber, MessageNames, BroadcastMessage } from 'services/messageService';

import { ReportsActions } from '../slice';
import { Report } from '../types';

interface Props {
  onCancel: () => void;
  commentData: Report;
  user_id: number;
}

function DeleteModal(props: Props) {
  const { onCancel, commentData, user_id } = props;
  const dispatch = useDispatch();
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
  const handleDelete = useCallback(() => {
    dispatch(
      ReportsActions.DeleteAdminCommentAction({
        id: commentData.id,
        user_id: user_id,
      }),
    );
  }, []);
  return (
    <div>
      <DialogTitle id="alert-dialog-title">delete comment?</DialogTitle>
      <DialogContent>
        <DialogContentText id="alert-dialog-description">
          Are You Sure You Want To Delete This Comment?
        </DialogContentText>
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
          text="Delete"
          className={Buttons.RedButton}
          loadingId={'deleteComentButton' + user_id}
          onClick={handleDelete}
        />
      </DialogActions>
    </div>
  );
}

export default memo(DeleteModal);
