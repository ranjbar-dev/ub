import { IconButton } from '@material-ui/core';
import EditIcon from 'images/themedIcons/edit';
import React, { memo, useRef, useState, useCallback } from 'react';

interface Props {
  initialValue: string;
}

function EditableValue(props: Props) {
  const { initialValue } = props;
  const edditingValue = useRef(initialValue);
  const [IsEditting, setIsEditting] = useState(false);
  const handleEditIconClick = useCallback(() => {
    setIsEditting(IsEditting => !IsEditting);
  }, [IsEditting]);
  return (
    <>
      <span>
        {IsEditting === true ? (
          <input
            type="text"
            onChange={e => {
              edditingValue.current = e.target.value;
            }}
          />
        ) : (
          initialValue
        )}
      </span>
      <span onClick={handleEditIconClick}>{<EditIcon />}</span>
    </>
  );
}

export default memo(EditableValue);
