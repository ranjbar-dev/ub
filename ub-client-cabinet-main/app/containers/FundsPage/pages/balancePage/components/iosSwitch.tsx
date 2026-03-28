import React from 'react';
import {
  withStyles,
  Theme,
  createStyles,
  FormGroup,
  FormControlLabel,
  Switch,
  SwitchProps,
  SwitchClassKey,
} from '@material-ui/core';
import styled from 'styles/styled-components';

interface Styles extends Partial<Record<SwitchClassKey, string>> {
  focusVisible?: string;
}

interface Props extends SwitchProps {
  classes: Styles;
}

const IOSSwitch = withStyles((theme: Theme) =>
  createStyles({
    root: {
      width: 35,
      height: 19,
      padding: 0,
      margin: theme.spacing(1),
    },
    switchBase: {
      padding: 2,
      color: 'var(--white)',
      '&$checked': {
        transform: 'translate(16px, 0px)',
        color: 'var(--white)',
        '& + $track': {
          //   backgroundColor: '#52d869',
          opacity: 1,
          border: 'none',
        },
      },
      '&$focusVisible $thumb': {
        // color: '#52d869',
        border: '6px solid #fff',
      },
    },
    thumb: {
      width: 15,
      height: 15,
    },
    track: {
      borderRadius: 26 / 2,
      border: `none`,
      backgroundColor: 'var(--swithBackBround)',
      opacity: 1,
      transition: theme.transitions.create(['background-color', 'border']),
    },
    checked: {},
    focusVisible: {},
  }),
)(({ classes, ...props }: Props) => {
  return (
    <Switch
      color='primary'
      focusVisibleClassName={classes.focusVisible}
      disableRipple
      classes={{
        root: classes.root,
        switchBase: classes.switchBase,
        thumb: classes.thumb,
        track: classes.track,
        checked: classes.checked,
      }}
      {...props}
    />
  );
});

export default function IosSwitch (props: { onChange: Function; title: any }) {
  const [state, setState] = React.useState({
    checked: true,
  });

  const handleChange = (name: string) => (
    event: React.ChangeEvent<HTMLInputElement>,
  ) => {
    setState({ ...state, [name]: event.target.checked });
    props.onChange(event.target.checked);
  };

  return (
    <StyledFormGroup>
      <FormControlLabel
        control={
          <IOSSwitch
            checked={state.checked}
            onChange={handleChange('checked')}
            value='checked'
          />
        }
        label={props.title}
      />
    </StyledFormGroup>
  );
}
const StyledFormGroup = styled(FormGroup)`
  .MuiTypography-body1 {
    line-height: 0 !important;
  }
`;
