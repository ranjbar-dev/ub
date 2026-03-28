import TextField from '@material-ui/core/TextField';
import React, { memo } from 'react';
//import 'date-fns';
//import DateFnsUtils from '@date-io/date-fns';
//import {
//  MuiPickersUtilsProvider,
//  KeyboardTimePicker,
//  KeyboardDatePicker,
//} from '@material-ui/pickers';
interface Props {
  title: string;
  onChange: (date: string) => void;
}

/**
 * Backup date/time picker using a plain HTML datetime-local input.
 * Used as a fallback when the Vaadin date picker is unavailable.
 *
 * @example
 * ```tsx
 * <DatePick title="Select Date" onChange={(date) => setDate(date)} />
 * ```
 */
function DatePick(props: Props) {
  const { onChange, title } = props;
  return (
    <div>
      <TextField
        onChange={e => {
          onChange(e.target.value);
        }}
        id="datetime-local"
        label={title}
        type="datetime-local"
        variant="outlined"
        margin="dense"
        //defaultValue="2017-05-24T10:30"
        InputLabelProps={{
          shrink: true,
        }}
      />
    </div>
  );
}

export default memo(DatePick);
