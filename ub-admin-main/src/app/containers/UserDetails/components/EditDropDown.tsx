import { Select, MenuItem } from '@material-ui/core';
import React, { memo, useState, useCallback } from 'react';
import styled from 'styled-components/macro';

interface Props {
  initialValue?: { name: string; id: string };
  options: {
    name: string;
    id: string;
  }[];
  className?: string;
  onSelect: Function;
}

function EditDropDown(props: Props) {
  const { options, onSelect, initialValue, className } = props;
  const [SelectedIndex, setSelectedIndex] = useState(
    initialValue ? initialValue.id + '' : '',
  );
  const handleChange = useCallback(e => {
    setSelectedIndex(e.target.value + '');
    onSelect(e.target.value);
  }, []);
  return (
    <Wrapper>
      <Select
        MenuProps={{
          getContentAnchorEl: null,
          anchorOrigin: {
            vertical: 'bottom',
            horizontal: 'left',
          },
        }}
        className={`ddown ${className ?? ''} ${initialValue?.name}`}
        variant="outlined"
        margin="dense"
        value={SelectedIndex}
        onChange={handleChange}
      >
        {initialValue && !initialValue.name.includes('+') && (
          <MenuItem className={initialValue?.name} value={initialValue.id}>
            {initialValue.name}
          </MenuItem>
        )}
        {options.map((item, index) => {
          if (
            !initialValue ||
            (initialValue &&
              initialValue.name != item.name &&
              !item.name.includes('+'))
          )
            return (
              <MenuItem className={item.name} key={item.id} value={item.id}>
                {item.name}
              </MenuItem>
            );
          return null;
        })}
      </Select>
    </Wrapper>
  );
}
const Wrapper = styled.div`
  position: relative;
  max-height: 20px;
  .ddown .MuiSelect-outlined.MuiSelect-outlined {
    padding: 4px !important;
    font-size: 13px !important;
    background: white;
    min-width: 63px;
  }

  .MuiOutlinedInput-root.Mui-focused .MuiOutlinedInput-notchedOutline {
    border-width: 1px;
    border-color: rgb(230, 227, 227) !important;
  }
  fieldSet {
    min-width: 110px;
  }
  .MuiInputBase-root {
    min-width: 125px;
  }
`;
export default memo(EditDropDown);
