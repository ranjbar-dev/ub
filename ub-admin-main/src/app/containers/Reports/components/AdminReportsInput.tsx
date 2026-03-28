import PostAddOutlinedIcon from '@material-ui/icons/PostAddOutlined';
import IsLoadingWithTextAuto from 'app/components/isLoadingWithText/isLoadingWithTextAuto';
import UbTextAreaAutoHeight from 'app/components/UbTextAreaAutoHeight/UbTextAreaAutoHeight';
import React, { memo, useRef, useCallback } from 'react';
import { useDispatch } from 'react-redux';

import { ReportsActions } from '../slice';

interface Props {
  userId: number;
}

function AdminReportsInput(props: Props) {
  const value = useRef('');
  const { userId } = props;
  const dispatch = useDispatch();
  const handleChange = useCallback((e: string) => {
    value.current = e;
  }, []);
  const handleSubmitClick = useCallback(() => {
    if (value.current != '') {
      dispatch(
        ReportsActions.AddAdmiCommentAction({
          comment: value.current,
          user_id: userId,
        }),
      );
    }
  }, []);
  return (
    <div className="adminReportsInput">
      <UbTextAreaAutoHeight
        placeHolder="Write Your Message ... "
        initialValue=""
        style={{
          marginBottom: '12px',
          height: '80px',
          background: '#FBFBFB',
          borderRadius: '5px',
          border: '1px solid #C5C5C5',
        }}
        onChange={handleChange}
      />
      <IsLoadingWithTextAuto
        icon={<PostAddOutlinedIcon />}
        style={{ background: '#649FFF' }}
        text="Send Message"
        loadingId={'userReportsSubmitButton' + userId}
        onClick={handleSubmitClick}
      />
    </div>
  );
}

export default memo(AdminReportsInput);
