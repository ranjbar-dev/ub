import axios,{AxiosRequestConfig} from 'axios';
import {BaseUrl,UploadUrls} from './constants';
import {MessageService,MessageNames} from './message_service';
import {ISetUploaderStateMessage} from 'containers/DocumentVerificationPage/components/common/fileUploader';
import {UploadState} from 'containers/DocumentVerificationPage/constants';
import {IUploadPercentMessage} from 'containers/DocumentVerificationPage/components/common/uploading';
import {cookies,CookieKeys} from './cookie';
export interface BaseUploadModel { }
export interface UploadModel {
	file: any;
	uploadUrl?: UploadUrls;
	type: string;
	isBack: boolean;
	headerType?: string;
	subtype: string;
	id?: number;
	mainImageId?: number;
}
export const UploadFile=(data: UploadModel) => {
	const {subtype,type,isBack,file}=data;
	const setUploaderStateMessage: ISetUploaderStateMessage={
		name: MessageNames.SET_UPLOADER_STATE,
		payload: {
			uploaderId: type+subtype+isBack,
			setTo: UploadState.UPLOADING,
		},
	};
	MessageService.send(setUploaderStateMessage);
	const config={
		onUploadProgress: function(progressEvent) {
			const percentCompleted=Math.round(
				(progressEvent.loaded*100)/progressEvent.total,
			);
			const uploadPercentageMessage: IUploadPercentMessage={
				name: MessageNames.UPLOAD_PERCENTAGE,
				payload: {
					uploaderId: type+subtype+isBack,
					percent: percentCompleted,
				},
			};
			MessageService.send(uploadPercentageMessage);
		},
		headers: {
			Authorization: `Bearer ${cookies.get(CookieKeys.Token)}`,
			Type: data.headerType,
		},
	};

	const sendingData=new FormData();
	sendingData.append('file',file);
	sendingData.append('type',type);
	sendingData.append('sub_type',subtype??'');

	if(data.isBack===true) {
		sendingData.append(
			'main_image_id',
			(data.mainImageId??data.id??'undefined').toString(),
		);
		sendingData.append('is_back','1');
	}

	return axios.post(
		BaseUrl+
		(data.uploadUrl
			? data.uploadUrl
			:'user-profile-image/upload?need_id=true'),
		sendingData,
		config,
	);
};
export const UploadMultiFile=(data: {
	frontImage: File;
	backImage: File;
	type: string;
	subtype: string;
	front_image_id?: number|string;
	back_image_id?: number|string;
}) => {
	const sendingData=new FormData();
	if(data.frontImage) {
		sendingData.append('front_image',data.frontImage);
	}
	const {
		subtype,
		type,
		frontImage,
		backImage,
		front_image_id,
		back_image_id,
	}=data;
	if(backImage) {
		sendingData.append('back_image',backImage);
	}
	sendingData.append('sub_type',subtype);
	sendingData.append('type',type);

	if(front_image_id) {
		sendingData.append('front_image_id',front_image_id+'');
	}
	if(back_image_id) {
		sendingData.append('back_image_id',back_image_id+'');
	}
	if(frontImage) {
		const setFrontUploaderStateMessage: ISetUploaderStateMessage={
			name: MessageNames.SET_UPLOADER_STATE,
			payload: {
				uploaderId: type+subtype+false,
				setTo: UploadState.UPLOADING,
			},
		};
		MessageService.send(setFrontUploaderStateMessage);
	}
	if(backImage) {
		const setbackUploaderStateMessage: ISetUploaderStateMessage={
			name: MessageNames.SET_UPLOADER_STATE,
			payload: {
				uploaderId: type+subtype+true,
				setTo: UploadState.UPLOADING,
			},
		};
		MessageService.send(setbackUploaderStateMessage);
	}
	// mocUpload(data);
	const config: AxiosRequestConfig={
		timeout: 30000,
		onUploadProgress: function(progressEvent) {
			const percentCompleted=Math.round(
				(progressEvent.loaded*100)/progressEvent.total,
			);
			if(frontImage) {
				const frontUploadPercentageMessage: IUploadPercentMessage={
					name: MessageNames.UPLOAD_PERCENTAGE,
					payload: {
						uploaderId: type+subtype+false,
						percent: percentCompleted,
					},
				};
				MessageService.send(frontUploadPercentageMessage);
			}
			if(backImage) {
				const backUploadPercentageMessage: IUploadPercentMessage={
					name: MessageNames.UPLOAD_PERCENTAGE,
					payload: {
						uploaderId: type+subtype+true,
						percent: percentCompleted,
					},
				};
				MessageService.send(backUploadPercentageMessage);
			}
		},
		headers: {
			Authorization: `Bearer ${cookies.get(CookieKeys.Token)}`,
			// Type: data.headerType,
		},
	};

	return axios.post(
		BaseUrl+'user-profile-image/multiple-upload',
		sendingData,
		config,
	);
};
export const mocUpload=(data: {
	frontImage: File;
	backImage: File;
	type: string;
	subtype: string;
}) => {
	const {subtype,type,frontImage,backImage}=data;
	let percent=0;
	setInterval(() => {
		percent+=50;
		if(percent<100) {
			if(frontImage) {
				const frontUploadPercentageMessage: IUploadPercentMessage={
					name: MessageNames.UPLOAD_PERCENTAGE,
					payload: {
						uploaderId: type+subtype+false,
						percent,
					},
				};
				MessageService.send(frontUploadPercentageMessage);
			}
			if(backImage) {
				const backUploadPercentageMessage: IUploadPercentMessage={
					name: MessageNames.UPLOAD_PERCENTAGE,
					payload: {
						uploaderId: type+subtype+true,
						percent,
					},
				};
				MessageService.send(backUploadPercentageMessage);
			}
		} else if(percent===100) {
			if(frontImage) {
				MessageService.send({
					name: MessageNames.SET_UPLOADED_IMAGE,
					payload: {
						uploaderId: type+subtype+false,
						disable: type,
						image:
							'https://i.pinimg.com/originals/16/4a/3a/164a3ab1c8fe20267e1b58e71cf1869c.png',
					},
				});
			}
			if(backImage) {
				MessageService.send({
					name: MessageNames.SET_UPLOADED_IMAGE,
					payload: {
						uploaderId: type+subtype+true,
						disable: type,
						image:
							'https://i.pinimg.com/originals/0e/47/b1/0e47b1ab2e12a7c5aa81a275a8cbb2ed.png',
					},
				});
			}
		}
	},1000);
};
