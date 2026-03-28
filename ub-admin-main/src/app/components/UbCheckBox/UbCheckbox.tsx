import { FormControlLabel, Checkbox } from '@material-ui/core';
import tick from 'images/staticImages/tick.png';
import React, { memo, useState, useCallback } from 'react';
interface Props {
  onChange: (value: boolean) => void;
  initialValue: boolean;
  title?: string;
  titlePlacement?:"start"|"end"
}

/**
 * Controlled checkbox with a custom tick image, wrapped in a FormControlLabel.
 *
 * @example
 * ```tsx
 * <UbCheckbox initialValue={true} title="Active" onChange={(checked) => setActive(checked)} />
 * ```
 */
function UbCheckbox(props: Props) {
  const [Checked, setChecked] = useState(props.initialValue);
  const { onChange,titlePlacement } = props;
  const handleChange = useCallback((e: boolean) => {
    setChecked(e);
    onChange(e);
  }, []);
  return (
    <FormControlLabel
      control={
        <Checkbox
          checked={Checked}
          checkedIcon={<img style={{ width: '24px' }} src={tick} alt="" aria-hidden="true" />}
          onChange={(event, checked) =>
            //handleChange(index, checked)
            {
              handleChange(checked);
            }
          }
        />
      }
      label={props.title ?? ''}
      labelPlacement={titlePlacement??"start"}
    />
  );
}

export default memo(UbCheckbox);
