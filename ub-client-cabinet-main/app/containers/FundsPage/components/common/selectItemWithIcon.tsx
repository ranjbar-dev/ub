import React from 'react';
import { ListItemText } from '@material-ui/core';

export const SelectItemWithIcon = props => (
  <ListItemText
    primary={
      <div className="coinWrapper">
        <img
          src={props.item.image}
          style={{
            height: '25px',
            width: '25px',
            border: '1px solid #c1c1c1',
            borderRadius: '50px',
            padding: '1px',
          }}
        />
        <span
          style={{
            margin: '0 8px',
            fontWeight: 600,
            color: `var(--blackText)`,
          }}
        >
          {props.item.code} - {props.item.name}
        </span>
      </div>
    }
  />
);
