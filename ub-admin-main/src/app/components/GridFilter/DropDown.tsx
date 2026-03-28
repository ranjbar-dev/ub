import { Select, MenuItem } from '@material-ui/core';
import React, { memo, useState, useCallback } from 'react';
import styled from 'styled-components/macro';

interface Props {
  options: {
    name: string;
    value: string;
  }[];
  title: string;
  /** Called with the selected value, or empty string when reset to the title option. */
  onSelect: (value: string) => void;
}

/**
 * Dropdown filter input used inside GridFilter for enum-type columns.
 * Renders a Material-UI Select with a placeholder title option.
 *
 * @example
 * ```tsx
 * <DropDown
 *   title="Status"
 *   options={[{ name: 'Pending', value: 'pending' }]}
 *   onSelect={(value) => applyFilter(value)}
 * />
 * ```
 */
function DropDown(props: Props) {
  const { options, onSelect, title } = props;
  const [SelectedIndex, setSelectedIndex] = useState(title);
  const handleChange = useCallback(e => {
    setSelectedIndex(e.target.value);
    if (e.target.value === title) {
      onSelect('');
      return;
    }
    onSelect(e.target.value);
  }, []);
  return (
    <Wrapper style={{}}>
      <Select
        MenuProps={{
          getContentAnchorEl: null,
          anchorOrigin: {
            vertical: 'bottom',
            horizontal: 'left',
          },
        }}
        className="select"
        variant="outlined"
        margin="dense"
        value={SelectedIndex}
        onChange={handleChange}
      >
        <MenuItem value={title}>{title}</MenuItem>
        {options.map((item, index) => {
          return (
            <MenuItem key={item.value} value={item.value}>
              {item.name}
            </MenuItem>
          );
        })}
      </Select>
    </Wrapper>
  );
}
const Wrapper = styled.div`
  position: relative;
  .MuiSelect-selectMenu {
    padding: 7px 10px !important;
  }
`;
export default memo(DropDown);
