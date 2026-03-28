import {Button} from '@material-ui/core';
import IsLoadingWithTextAuto from 'app/components/isLoadingWithText/isLoadingWithTextAuto';
import UbTextAreaAutoHeight from 'app/components/UbTextAreaAutoHeight/UbTextAreaAutoHeight';
import {Buttons} from 'app/constants';
import {translations} from 'locales/i18n';
import React,{memo,useRef} from 'react';
import {useTranslation} from 'react-i18next';

interface Props {
	onRejectReason: (value: string) => void;
	onCancel: () => void;
	initialValue: string;
	id: string
}

function RejectPopup(props: Props) {
	const {t}=useTranslation();
	const {initialValue,onCancel,onRejectReason,id}=props;
	const rejReason=useRef(initialValue);
	return (
		<div className="rejectPopup">
			<div className="title">{t(translations.CommonTitles.RejectReason())}</div>
			<div className="inp" style={{marginBottom: '24px'}}>
				<UbTextAreaAutoHeight
					initialValue={initialValue}
					onChange={e => {
						rejReason.current=e;
					}}
				/>
			</div>
			<div className="buttons">
				<IsLoadingWithTextAuto
					style={{marginRight: '12px'}}
					text={t(translations.CommonTitles.Reject())}
					loadingId={id}
					className={Buttons.RedButton}
					onClick={() => onRejectReason(rejReason.current)}
				/>

				<Button
					onClick={onCancel}
					variant="contained"
					className={Buttons.BlackButton}
				>
					{t(translations.CommonTitles.Cancel())}
				</Button>
			</div>
		</div>
	);
}

export default memo(RejectPopup);
