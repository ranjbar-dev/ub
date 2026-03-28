import React, { memo } from 'react';
import styled from 'styles/styled-components';
import { MenuItem, Select } from '@material-ui/core';

import Translate from 'containers/TradePage/messages';
import { FormattedMessage } from 'react-intl';

const PrecisionSelect = (props: { onPrecisionChange: Function }) => {
  const [Precision, setPrecision] = React.useState<string | number>(8);
  const handleChange = (event: React.ChangeEvent<{ value: unknown }>) => {
    props.onPrecisionChange(event.target.value);
    setPrecision(event.target.value as number);
  };

  return (
    <Wrapper>
      <div className='label'>
        <FormattedMessage {...Translate.Precision} />
      </div>
      <Select
        //variant="outlined"
        margin='dense'
        value={Precision}
        onChange={handleChange}
      >
        <MenuItem value={1}>1</MenuItem>
        <MenuItem value={2}>2</MenuItem>
        <MenuItem value={3}>3</MenuItem>
        <MenuItem value={4}>4</MenuItem>
        <MenuItem value={5}>5</MenuItem>
        <MenuItem value={6}>6</MenuItem>
        <MenuItem value={7}>7</MenuItem>
        <MenuItem value={8}>8</MenuItem>
      </Select>
    </Wrapper>
  );
};
export default memo(PrecisionSelect);
const Wrapper = styled.div`
  position: absolute;
  z-index: 1;
  top: 1px;
  right: 1px;
  fieldset {
    padding-left: 8px;
    padding: 0;
    max-height: 25px;
  }
  .MuiSelect-select.MuiSelect-select {
    padding: 0 15px;
    padding-right: 25px;
    font-size: 11px !important;
  }
  .label {
    span {
      line-height: 1px;
      font-size: 12px;
      font-weight: 600;
      color: var(--textGrey);
      position: absolute;
      left: -65px;
      top: 10px;
    }
  }
`;
