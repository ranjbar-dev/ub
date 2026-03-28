import { InitialUserDetails } from 'app/containers/UserAccounts/types';
import React, { memo, useRef, useCallback, useState, useEffect } from 'react';
import Cropper from 'react-cropper';
import NewWindow from 'react-new-window';
import { useDispatch } from 'react-redux';
import { MessageNames, MessageService, Subscriber, BroadcastMessage } from 'services/messageService';
import downloadFile from 'utils/fileDownload';

import { IdentityTypes } from '../constants';
import { VerificationWindowActions } from '../slice';
import { ProfileImageData } from '../types';

import './imageWrapperStyle.scss';
import 'cropperjs/dist/cropper.css';
import { Buttons } from 'app/constants';
import IsLoadingWithTextAuto from 'app/components/isLoadingWithText/isLoadingWithTextAuto';
import PopupModal from 'app/components/materialModal/modal';

import RejectModal from './RejectModal';

import {
	RotateRight,
	RotateLeft,
	ZoomIn,
	CenterFocusStrong,
	ZoomOut,
	GetApp,
} from '@material-ui/icons';
interface Props {
	data: InitialUserDetails;
	selectedImage: ProfileImageData;
	type: IdentityTypes;
}

function ImageWrapper(props: Props) {
	const { data, selectedImage, type } = props;
	const dispatch = useDispatch();
	const [IsRejectModalOpen, setIsRejectModalOpen] = useState(false);
	const [IsDownloadOpen, setIsDownloadOpen] = useState(false);
	const additionalData = useRef<Record<string, unknown>>({});
	// eslint-disable-next-line @typescript-eslint/no-explicit-any
	const containerRef = useRef<any>();

	const rotate = (degree: number) => {
		containerRef.current.cropper.rotate(degree);
	};

	const center = () => {
		containerRef.current.cropper.reset();
		//.getContainerData();
		//containerRef.current.cropper
		//  .zoomTo(1, {
		//    x: containerData.width,
		//    y: containerData.height,
		//  })
		//  .scaleX(1)
		//  .scaleY(1)
		//  .setCanvasData({
		//    left: containerData.width / 3.5,
		//    top: containerData.height / 3.5,
		//  });
	};

	const zoom = (type: string) => {
		containerRef.current.cropper.zoom(type === 'in' ? 0.2 : -0.2);
	};
	const handleUpdateClick = (loadingButtonId: string, confirmation_status: string) => {
		dispatch(
			VerificationWindowActions.UpdateProfileImageStatusAction({
				confirmation_status,
				id: selectedImage.id,
				loadingButtonId,
				type,
				user_id: data.id,
				...additionalData.current,
			}),
		);
	};
	const handleDownload = () => {
		// window.open(selectedImage.imagePath)
		// MessageService.send({
		// 	name: MessageNames.DOWNLOAD_FILE,
		// 	payload: { url: selectedImage.imagePath, filename: selectedImage.originalFileName },
		// });

		setIsDownloadOpen(true);
	};
	useEffect(() => {
		const Subscription = Subscriber.subscribe((message: BroadcastMessage) => {
			const payload = message.payload as Record<string, unknown>;
			if (
				message.name === MessageNames.DATASEND &&
				payload.userId === data.id
			) {
				if (payload.documentId) {
					additionalData.current.id_card_code = payload.documentId;
				}
				if (payload.type) {
					additionalData.current.newType = payload.type;
				}
				if (payload.sub_type) {
					additionalData.current.sub_type = payload.sub_type;
				}
			}
		});
		return () => {
			Subscription.unsubscribe();
		};
	}, []);
	return (
		<div className="mainWrapper">
			{IsDownloadOpen && (
				<NewWindow
					//  center="parent"
					url={selectedImage.imagePath}
					title={selectedImage.imagePath}
					features={{ height: 600, width: 800 }}
					onUnload={() => setIsDownloadOpen(false)}
				>
					{/* <img
						height="400"
						style={{ position: 'absolute', left: 'calc(50% - 200px)' }}
						src={selectedImage.imagePath}
						alt=""
					/> */}
				</NewWindow>
			)}
			{/*<div ref={containerRef} style={{ height: '300px' }}></div>*/}
			{/*{RRef && RRef.current && (
        <Viewer
          visible
          container={RRef.current}
          images={[{ src: selectedImage.imagePath, alt: 'ima' }]}
        />
      )}*/}
			<PopupModal
				onClose={() => {
					setIsRejectModalOpen(false);
				}}
				isOpen={IsRejectModalOpen}
			>
				<RejectModal
					type={type}
					imageId={selectedImage.id}
					user_id={data.id}
					onCancel={() => {
						setIsRejectModalOpen(false);
					}}
				/>
			</PopupModal>
			<div className="cropperWrapper">
				<Cropper
					ref={containerRef}
					src={selectedImage.imagePath}
					style={{
						height: 390,
						width: '100%',
						background: '#b8b8b8 0 0 no-repeat padding-box',
						borderRadius: '5px',
					}}

					guides={false}
					highlight={true}
					movable={true}
					zoomable={true}
					zoomOnTouch={true}
					zoomOnWheel={true}
					rotatable={true}
					autoCrop={false}
					background={false}
					checkCrossOrigin={false}
					dragMode={'move'}
				/>
			</div>
			<div className="imageButtonsWrapper">
				<div className="btn" onClick={center}>
					<CenterFocusStrong />
				</div>
				<div className="btn" onClick={() => rotate(-90)}>
					<RotateLeft />
				</div>
				<div className="btn" onClick={() => rotate(90)}>
					<RotateRight />
				</div>
				<div className="btn" onClick={() => zoom('in')}>
					<ZoomIn />
				</div>
				<div className="btn" onClick={() => zoom('out')}>
					<ZoomOut />
				</div>
				<div className="btn downloadButton" onClick={() => handleDownload()}>
					<GetApp />
				</div>
			</div>
			<div className="actionButtonsWrapper">
				<IsLoadingWithTextAuto
					text="Pending"
					className={Buttons.BlackButton + ' t6'}
					loadingId={'PendingImageButton' + data.id}
					onClick={() => {
						handleUpdateClick('PendingImageButton' + data.id, 'processing');
					}}
				/>
				<IsLoadingWithTextAuto
					text="Reject"
					className={Buttons.RedButton + ' t6'}
					loadingId={'RejectImageButton' + data.id}
					onClick={() => {
						setIsRejectModalOpen(true);
					}}
				/>
				<IsLoadingWithTextAuto
					text="Temp Verify"
					className={Buttons.SkyBlueButton + ' t6'}
					loadingId={'TempVerify' + data.id}
					onClick={() => {
						handleUpdateClick('TempVerify' + data.id, 'partially_confirmed');
					}}
				/>
				<IsLoadingWithTextAuto
					text="Verify"
					className={Buttons.GreenButton + ' t6'}
					loadingId={'VerifyImageButton' + data.id}
					onClick={() => {
						handleUpdateClick('VerifyImageButton' + data.id, 'confirmed');
					}}
				/>
			</div>
		</div>
	);
}

export default memo(ImageWrapper);
