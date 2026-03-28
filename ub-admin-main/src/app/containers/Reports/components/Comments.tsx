import { IconButton } from '@material-ui/core';
import DeleteIcon from '@material-ui/icons/Delete';
import EditIcon from '@material-ui/icons/Edit';
import PopupModal from 'app/components/materialModal/modal';
import React, { useState } from 'react';

import { Report } from '../types';
import DeleteModal from './DeleteModal';
import EditModal from './EditModal';

interface Props {
  comments: Report[];
  user_id: number;
}

function Comments(props: Props) {
  const { comments, user_id } = props;
  const [IsEditPopUpOpen, setIsEditPopUpOpen] = useState(false);
  const [IsDeletePopupOpen, setIsDeletePopupOpen] = useState(false);
  const [PopupData, setPopupData] = useState<Report>({} as Report);
  return (
    <div className="adminReportsComments">
      <PopupModal
        isOpen={IsEditPopUpOpen}
        onClose={() => setIsEditPopUpOpen(false)}
      >
        <div className="edit">
          <EditModal
            onCancel={() => setIsEditPopUpOpen(false)}
            commentData={PopupData}
            user_id={user_id}
          />
        </div>
      </PopupModal>
      <PopupModal
        isOpen={IsDeletePopupOpen}
        onClose={() => setIsDeletePopupOpen(false)}
      >
        <div className="delete">
          <DeleteModal
            commentData={PopupData}
            user_id={user_id}
            onCancel={() => {
              setIsDeletePopupOpen(false);
            }}
          />
        </div>
      </PopupModal>
      {comments.map((item, index) => {
        return (
          <div key={'comment' + item.id} className="mainCommentWrapper">
            <div className="top">
              <span className="bold grey">
                {item.adminFullName + ' : '}
                {item.createdAt}
              </span>
            </div>
            <div className="commentWrapper">
              <div className="comment">{item.comment}</div>
              <div className="actions">
                <IconButton
                  onClick={() => {
                    setPopupData(item);
                    setIsDeletePopupOpen(true);
                  }}
                  size="small"
                >
                  <DeleteIcon />
                </IconButton>
                <IconButton
                  onClick={() => {
                    setPopupData(item);
                    setIsEditPopUpOpen(true);
                  }}
                  size="small"
                >
                  <EditIcon />
                </IconButton>
              </div>
            </div>
          </div>
        );
      })}
    </div>
  );
}

export default Comments;
