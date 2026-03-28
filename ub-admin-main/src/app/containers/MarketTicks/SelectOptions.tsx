import { Select, MenuItem } from '@material-ui/core';
import React, { memo, useState } from 'react';

interface SelectOption {
  name: string;
  value: string;
  obj?: unknown;
}
interface Props {
  options: SelectOption[];
  initialValue: string;
  styles?: React.CSSProperties;
  onSelect?: (value: string) => void;
  onSelectObj?: (value: SelectOption | undefined) => void;
}

function SelectOptions(props: Props) {
  const { options, initialValue, onSelect,onSelectObj,styles } = props;
  const [SelectedIndex, setSelectedIndex] = useState(initialValue);

  return (
    <Select
      variant="outlined"
      style={styles}
      margin="dense"
      MenuProps={{
        getContentAnchorEl: null,
        anchorOrigin: {
          vertical: 'bottom',
          horizontal: 'left',
        },
      }}
      value={SelectedIndex}
      onChange={(e: React.ChangeEvent<{ value: unknown }>) => {
        setSelectedIndex(e.target.value as string);
        if(onSelect){onSelect(e.target.value as string);}
        if(onSelectObj){
        const selected=options.find((el)=>el.value===e.target.value as string)
         onSelectObj(selected)
        }
        //console.log(e);
      }}
    >
      {options.map((item, index) => {
        return (
          <MenuItem key={item.name} value={item.value}>
            {item.name}
          </MenuItem>
        );
      })}
    </Select>
  );
}

export default memo(SelectOptions);
