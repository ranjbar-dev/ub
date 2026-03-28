import { Select, MenuItem } from '@material-ui/core';
import React, { memo, useState, useCallback } from 'react';
import styled from 'styled-components/macro';

interface Props {
  options: {
    name: string;
    value: string | boolean;
  }[];
  initialValue: string | boolean;
  label?: string;
  /** Called with the selected value when the user changes the dropdown. */
  onSelect: (value: string) => void;
  style?: React.CSSProperties;
  'aria-label'?: string;
}

/**
 * Labeled dropdown (Select) for form fields inside modals and detail panels.
 * Displays an optional label prefix before the select control.
 *
 * @example
 * ```tsx
 * <UbDropDown
 *   label="Status"
 *   initialValue="active"
 *   options={[{ name: 'Active', value: 'active' }, { name: 'Inactive', value: 'inactive' }]}
 *   onSelect={(val) => setStatus(val)}
 * />
 * ```
 */
function DropDown(props: Props) {
  const { options, onSelect, initialValue, style } = props;
  const ariaLabel = props['aria-label'];
  const [SelectedValue, setSelectedValue] = useState(initialValue);
  const handleChange = useCallback(e => {
    setSelectedValue(e.target.value);
    onSelect(e.target.value);
  }, []);
  return (
    <Wrapper style={style ?? {}}>
      {props.label && (
        <div className="label">
          {props.label}
          {' : '}
        </div>
      )}
      <Select
        className="select"
        variant="outlined"
        margin="dense"
        inputProps={ariaLabel ? { 'aria-label': ariaLabel } : undefined}
        MenuProps={{
          getContentAnchorEl: null,
          anchorOrigin: {
            vertical: 'bottom',
            horizontal: 'left',
          },
        }}
        value={SelectedValue}
        onChange={handleChange}
      >
        {options.map((item, index) => {
          return (
            <MenuItem key={item.value + ''} value={item.value + ''}>
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
  width: 100%;
  display: flex;
  align-items: center;
  margin-bottom: 12px;
  .label {
    min-width: 120px;
    font-size: 13px;
    color: #535353;
  }
  .MuiOutlinedInput-root {
    flex: 1;
  }
  .MuiSelect-selectMenu {
    padding: 7px 10px !important;
  }
`;
export default memo(DropDown);
