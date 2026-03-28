import {Zoom,Dialog,IconButton,DialogTitle} from '@material-ui/core';
import {TransitionProps} from '@material-ui/core/transitions/transition';
import CloseIcon from 'images/themedIcons/CloseIcon';
import React from 'react';
import styled from 'styled-components/macro';

const Transition=React.forwardRef<unknown,TransitionProps>(
	function Transition(props,ref) {
		// @ts-expect-error — Zoom ref type incompatibility with forwardRef generic
		return <Zoom ref={ref} {...props} {...{timeout: 300}} />;
	},
);
/**
 * Material-UI Dialog wrapper with a close button and smooth Zoom transition.
 * Used as the base modal for all admin panel popups.
 *
 * @example
 * ```tsx
 * <PopupModal isOpen={open} onClose={() => setOpen(false)} title="Confirm">
 *   <p>Are you sure?</p>
 * </PopupModal>
 * ```
 */
export default function PopupModal(props: {
	isOpen: boolean;
	children: React.ReactNode;
	onClose: () => void;
	fullScreen?: boolean;
	wrapperClassName?: string;
	minWidth?: string;
	title?: string;
}) {
	return (
		<StyledDialog
			TransitionComponent={Transition}
			fullScreen={props.fullScreen??false}
			disablePortal
			open={props.isOpen}
			className={props.wrapperClassName}
			onClose={() => {
				props.onClose();
			}}
		>
			<CloseWrapper>
				<IconButton
					style={{float: 'right'}}
					onClick={() => {
						props.onClose();
					}}
					size="small"
				>
					<CloseIcon />
					{/*<img src={crossIcon} alt="" />*/}
				</IconButton>
			</CloseWrapper>
			{props.title&&<DialogTitle>{props.title}</DialogTitle>}

			<ChildsWrapper>{props.children}</ChildsWrapper>
		</StyledDialog>
	);
}

const StyledDialog=styled(Dialog)`
&.min1200{
	.MuiDialog-paper{
		min-width:1200px;
	}
}

`


const CloseWrapper=styled.div`
  position: absolute;
  right: 5px;
  top: 5px;
  z-index: 2;
  img {
    width: 35px;
  }
`;
const ChildsWrapper=styled.div`
  .alertWrapper {
    padding: 72px 102px 24px 102px;
    min-height: 220px;
    display: flex;
    flex-direction: column;
    align-items: center;
    label {
      color: var(--blackText);
    }
    p {
      margin-top: 0;
      color: var(--blackText);
    }
  }
  .addAddress {
    .MuiInputLabel-outlined.MuiInputLabel-marginDense {
      margin-top: -0.5px !important;
    }
    .loadingCircle {
      top: 8px !important;
    }
    .MuiInputBase-root legend {
      zoom: 0.72;
      margin-left: 5px;
    }
  }
`;
