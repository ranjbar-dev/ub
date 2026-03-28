import {Button} from '@material-ui/core';
import React,{useState,useEffect} from 'react';
import {Subscriber,MessageNames,BroadcastMessage} from 'services/messageService';

import LoadingInButton from '../loadingInButton/loadingInButton';

/**
 * Button that shows a spinner while an async operation is in progress.
 * Loading state is driven by SET_BUTTON_LOADING messages on the message bus.
 *
 * @example
 * ```tsx
 * <IsLoadingWithTextAuto text="Submit" loadingId="submitBtn" onClick={handleSubmit} />
 * ```
 */
export default function IsLoadingWithTextAuto(props: {
	text: React.ReactNode;
	loadingId: string;
	onClick: () => void;
	style?: React.CSSProperties;
	disabled?: boolean;
	className?: string;
	small?: boolean;
	icon?: React.ReactNode;
	lightBackground?: boolean;
}) {
	const {style,lightBackground,disabled}=props;
	const [IsLoading,setIsLoading]=useState(false);
	useEffect(() => {
		const Subscription=Subscriber.subscribe((message: BroadcastMessage) => {
			if(
				message.name===MessageNames.SET_BUTTON_LOADING&&
				message.loadingId===props.loadingId
			) {
				setIsLoading(message.payload as boolean);
			}
		});
		return () => {
			Subscription.unsubscribe();
		};
	},[props.loadingId]);
	let loadingOpacity=IsLoading===true? 1:0;
	let textOpacity=IsLoading===false? 1:0;
	return (
		<Button
			disabled={disabled}
			color="primary"
			variant="contained"
			startIcon={props.icon??null}
			style={props.style? {...style}:{}}
			className={props.className??''}
			onClick={IsLoading===false? () => props.onClick():() => { }}
		>
			<span
				className="loadingCircle"
				style={{
					opacity: loadingOpacity,
					position: 'absolute',
					left: 'calc(50% - 8px)',
					top: props.small? '5px !important':'11px',
				}}
			>
				<LoadingInButton lightBackground={lightBackground} />
			</span>
			<span style={{opacity: textOpacity}}>{props.text}</span>
		</Button>
	);
}
