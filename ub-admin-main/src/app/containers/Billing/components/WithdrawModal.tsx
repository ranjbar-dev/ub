import { Button, Checkbox, FormControlLabel } from '@material-ui/core';
import IsLoadingWithTextAuto from 'app/components/isLoadingWithText/isLoadingWithTextAuto';
import PopupModal from 'app/components/materialModal/modal';
import UbCheckbox from 'app/components/UbCheckBox/UbCheckbox';
import UBInput from 'app/components/UBInput/UBInput';
import SnackBar from 'app/components/UserDetailsWindow/SnackBar';
import { Buttons } from 'app/constants';
import { translations } from 'locales/i18n';
import React, { memo, useCallback, useEffect, useRef, useState } from 'react';
import { useTranslation } from 'react-i18next';
import { useDispatch, useSelector } from 'react-redux';
import { MessageNames, MessageService, Subscriber, BroadcastMessage } from 'services/messageService';
import { CurrencyFormater } from 'utils/formatters';
import { useInjectReducer, useInjectSaga } from 'utils/redux-injectors';

import { billingSaga } from '../saga';
import AdminCommentPopup from './AdminCommentPopup';
import ConfirmPopup from './confirmPopup';
import RejectPopup from './RejectPopup';
import { BillingActions } from '../slice';
import { BillingReducer, sliceKey } from '../slice';
import { selectBillingCommissionsData } from '../selectors';
import { Payment, PaymentDetails } from '../types';

interface Props {
	initialModalData: {
		rowData: Payment;
		details: PaymentDetails;
	};
}

const WithdrawModal = (props: Props) => {
	const { t } = useTranslation();
	const dispatch = useDispatch();
	useInjectReducer({ key: sliceKey, reducer: BillingReducer });
	useInjectSaga({ key: sliceKey, saga: billingSaga });
	const { initialModalData } = props;
	const { rowData, details } = initialModalData;
	const { status } = rowData
	const [RejectPopupOpen, setRejectPopupOpen] = useState(false);
	const [AdminCommentPopupOpen, setAdminCommentPopupOpen] = useState(false);
	const [confirmPopupData, setConfirmPopupData] = useState({ isOpen: false, message: '' })


	const pressedAdminStatusButtonId = useRef('');
	const [Comments, setComments] = useState(details.adminComments);
	const Commitions = (useSelector(selectBillingCommissionsData) ?? {}) as {
		depositAndWithdrawData?: {
			btcTotalDeposit: number;
			btcTotalWithdraw: number;
		};
		commissionData?: {
			btcWithdrawFee: number;
		};
		balanceData?: {
			btcTotalSum: number;
			btcInOrderSum: number;
		};
	};
	useEffect(() => {
		dispatch(BillingActions.GetCommitionsAction({ userId: rowData.userId }));
		const Subscription = Subscriber.subscribe((message: BroadcastMessage) => {
			if (message.name === MessageNames.CLOSE_REJECT_POPUP) {
				setRejectPopupOpen(false);
			}
			if (message.name === MessageNames.SHOW_WITHDRAW_CONFIRM_POPUP && message.userId === rowData.userId + '' + details.id) {
				setConfirmPopupData({ isOpen: true, message: message.payload as string })
			}
		});
		return () => {
			Subscription.unsubscribe();
		};
	}, [details]);

	const DetailRow = useCallback(propes => {
		let { title, value, last, small } = propes;
		if (!value) {
			value = '-';
		}
		return (
			<div
				className={`detailRow ${last === true ? 'last' : ''} ${small === true ? 'small' : ''
					}`}
			>
				<div className="title">
					{title}
					{' : '}
				</div>
				<div className="value">{value}</div>
			</div>
		);
	}, []);
	const changedData = useRef<Record<string, unknown>>({
		user_id: rowData.userId,
		id: details.id,
		fee: details.fee,
		auto_transfer: details.autoTransfer,
	});
	const handleAmountChange = (e: string) => {
		changedData.current.amount = e;
	};
	const handleFeeChange = (e: string) => {
		changedData.current.fee = e;
	};

	const handleStatusChange = useCallback((e: string) => {
		delete changedData.current.admin_status;
		changedData.current.status = e;
		sendData(e);
	}, []);
	const handleRejectClick = useCallback(() => {
		setRejectPopupOpen(true);
	}, []);
	const handleReject = useCallback((e: string) => {
		delete changedData.current.admin_status;
		changedData.current.status = 'rejected';
		changedData.current.rejection_reason = e;
		sendData('rejected');
	}, []);
	const handleAdminStatusChange = useCallback((e: string) => {
		delete changedData.current.status;
		changedData.current.admin_status = e;
		if (e === 'pending') {
			sendData(e);
			return;
		}

		pressedAdminStatusButtonId.current = e;
		setAdminCommentPopupOpen(true);
	}, []);
	const handleAutotransferChange = useCallback((e: boolean) => {
		changedData.current.auto_transfer = e;
	}, []);
	const sendData = (buttonLabel: string) => {

		dispatch(
			BillingActions.UpdateBillingWithdrawAction({
				...changedData.current,
				buttonId: buttonLabel + rowData.userId + '' + details.id,
			}),
		);
	};
	const setAdminComment = (e: string) => {
		dispatch(
			BillingActions.AddPaymentCommentAction({
				payment_id: details.id,
				comment: e,
			}),
		);
		sendData(pressedAdminStatusButtonId.current);
		setAdminCommentPopupOpen(false);
	};

	const handleCloseWindow = () => {
		const messageData = {
			name: MessageNames.CLOSE_WINDOW,
			payload: {
				id: 'WD' + details.id,
			},
		};
		MessageService.send(messageData);
	}



	const handleTestPopup = () => {
		MessageService.send({
			name: MessageNames.SHOW_WITHDRAW_CONFIRM_POPUP,
			payload: 'withdraw updated ',
			type: 'success',
			userId: rowData.userId + '' + details.id,
		});
	}
	const sendMoneyTitle = (status: Props["initialModalData"]["rowData"]["status"]) => {
		return (status === 'rejected' || status === 'canceled') ? 'ReActivate' : status === 'completed' ? "Completed" : t(translations.CommonTitles.SendMoney())
	}
	const handleSendMoneyClick = (status: Props["initialModalData"]["rowData"]["status"]) => {
		const toPrevent = ['completed', 'in_progress', 'user_canceled']
		if (toPrevent.indexOf(status) !== -1) {

			MessageService.send({
				name: MessageNames.TOAST,
				payload: `Can't do this action for ${status.replace(/_/g, ' ')} withdraws`,
				type: 'error',
				userId: rowData.userId + '' + details.id,
			});

			return
		}


		if (status !== 'rejected' && status !== 'completed' && status !== 'canceled') {
			handleStatusChange('in_progress')
		}
		else if (status === 'rejected' || status === 'canceled') {
			handleStatusChange('created')
		}
	}
	return (
		<div className="withdrawModal">
			<SnackBar userId={rowData.userId + '' + details.id} />
			<SnackBar userId={'withdrawPopup' + details.id} priority={'secound'} />
			<PopupModal
				isOpen={confirmPopupData.isOpen}
				onClose={handleCloseWindow}
			>
				<ConfirmPopup content={confirmPopupData.message} onConfirm={handleCloseWindow} />
			</PopupModal>
			<PopupModal
				isOpen={RejectPopupOpen}
				onClose={() => {
					setRejectPopupOpen(false);
				}}
			>
				<RejectPopup
					initialValue={''}
					onCancel={() => {
						setRejectPopupOpen(false);
					}}
					id={'rejected' + rowData.userId + '' + details.id}
					onRejectReason={handleReject}
				/>
			</PopupModal>
			<PopupModal
				isOpen={AdminCommentPopupOpen}
				onClose={() => {
					setAdminCommentPopupOpen(false);
				}}
			>
				<AdminCommentPopup
					initialValue={''}
					onCancel={() => {
						setAdminCommentPopupOpen(false);
					}}
					onAdminCommentSubmitted={setAdminComment}
				/>
			</PopupModal>
			<div className="wrapperTitle">
				{t(translations.CommonTitles.WithdrawalRequest())}
			</div>
			<div className="content">
				<div className="detailsContainer">
					<div className="detailRowsContainer">
						<DetailRow
							title={t(translations.CommonTitles.WithdrawalID())}
							value={details.id}
						/>
						<DetailRow
							title={t(translations.CommonTitles.Name())}
							value={details.name}
						/>
						<DetailRow
							title={t(translations.CommonTitles.Email())}
							value={details.userEmail + ` (${rowData.userId})`}
						/>
						<DetailRow
							title={t(translations.CommonTitles.Country())}
							value={details.country}
						/>
						<DetailRow
							title={t(translations.CommonTitles.TrustLevel())}
							value={details.level}
						/>
						<DetailRow
							title={t(translations.CommonTitles.RequestDate())}
							value={details.createdAt}
						/>
						<DetailRow
							title={t(translations.CommonTitles.LastResponseDate())}
							value={Comments[0]?.updatedAt ?? ''}
						/>
						<DetailRow
							title={t(translations.CommonTitles.IPAddress())}
							value={details.ip}
						/>
						<DetailRow
							title={t(translations.CommonTitles.Method())}
							value={details.currencyName}
						/>
						<DetailRow
							title={t(translations.CommonTitles.Amount())}
							value={CurrencyFormater(details.amount)}
						/>
						<DetailRow
							title={t(translations.CommonTitles.WalletAddress())}
							value={details.toAddress}
							last={!details.txId}
						/>
						{details.txId && <DetailRow
							title={t(translations.CommonTitles.TransactionId())}
							value={details.txId}
							last={true}
						/>}
						{details.rejectionReason && <DetailRow
							title={t(translations.CommonTitles.RejectReason())}
							value={details.rejectionReason}
							last={true}
						/>}
					</div>
					<div className="commentsContainer">
						<div className={'adminPaymentComment'}>
							{Comments.length > 0 && (
								<>
									{' '}
									<div className="adminName">
										{Comments[0].updatedAt ?? ''}{' '}
										{`(${Comments[0].adminName ?? ''})`}
									</div>
									{/*<div className="commentDate"></div>*/}
									<div className="adminComment">
										{Comments[0].comment ?? ''}
									</div>
								</>
							)}
						</div>
					</div>
				</div>
				<div className="actionsContainer">
					<div className="inp1">
						<UBInput
							style={{ fontSize: '16px', fontWeight: 700 }}
							onChange={e => {
								handleAmountChange(e);
							}}
							initialValue={details.amount + ''}
							properties={{
								fullWidth: true,
								margin: 'dense',
							}}
							endText={details.currencyCode}
						/>
					</div>
					<div className="inp2">
						<div className="fee">
							{t(translations.CommonTitles.Fee())}
							{' : '}
						</div>
						<UBInput
							onChange={e => {
								handleFeeChange(e);
							}}
							style={{ fontSize: '14px', fontWeight: 700, color: '#414141' }}
							initialValue={details.fee + ''}
							properties={{
								fullWidth: true,
								margin: 'dense',
							}}
							endText={details.currencyCode}
						/>
					</div>
					<div className="autoTransfer">
						<UbCheckbox
							title={t(translations.CommonTitles.AutomaticTransfer())}
							initialValue={changedData.current.auto_transfer as boolean}
							onChange={handleAutotransferChange}
						/>
					</div>

					<div className="sendMoney">
						<IsLoadingWithTextAuto
							text={sendMoneyTitle(status)}
							loadingId={'in_progress' + rowData.userId + '' + details.id}
							//  className={'sendMoneyButton'}
							className={status === 'rejected' || status === 'completed' ? Buttons.SkyBlueButton : Buttons.LightGreenButton}
							onClick={() => handleSendMoneyClick(status)}
						/>
						{/*<Button
              onClick={() => {
                handleStatusChange('in_progress');
              }}
              color="primary"
              className="sendMoneyButton"
              variant="contained"
            >
              {t(translations.CommonTitles.SendMoney())}
            </Button>*/}
					</div>
					<div className="rejectCancel" style={{ filter: status === 'rejected' || status === 'completed' ? 'grayscale(1)' : 'none' }} >
						<IsLoadingWithTextAuto
							disabled={status === 'rejected' || status === 'completed'}
							text={t(translations.CommonTitles.Reject())}
							loadingId={'rejected' + rowData.userId + '' + details.id}
							className={Buttons.RedButton}
							onClick={handleRejectClick}
						/>
						<IsLoadingWithTextAuto
							disabled={status === 'rejected' || status === 'completed'}
							text={t(translations.CommonTitles.Cancel())}
							loadingId={'canceled' + rowData.userId + '' + details.id}
							className={Buttons.BlackButton}
							onClick={() => {
								handleStatusChange('canceled');
							}}
						/>
						{/*<Button
              onClick={() => {
                handleStatusChange('cancel');
              }}
              className={Buttons.BlackButton}
            >
              {t(translations.CommonTitles.Cancel())}
            </Button>*/}
					</div>
					<div className="actionDetailRows">
						<DetailRow
							title={t(translations.CommonTitles.TotalDeposit())}
							value={
								Commitions.depositAndWithdrawData
									? CurrencyFormater(
										Commitions.depositAndWithdrawData.btcTotalDeposit +
										' BTC',
									)
									: '-'
							}
							small={true}
						/>
						<DetailRow
							title={t(translations.CommonTitles.TotalWithdraw())}
							value={
								Commitions.depositAndWithdrawData
									? CurrencyFormater(
										Commitions.depositAndWithdrawData.btcTotalWithdraw +
										' BTC',
									)
									: '-'
							}
							small={true}
						/>
						<DetailRow
							title={t(translations.CommonTitles.TotalCommissions())}
							value={
								Commitions.commissionData
									? Commitions.commissionData.btcWithdrawFee + ' BTC'
									: '-'
							}
							small={true}
						/>
						<DetailRow
							title={t(translations.CommonTitles.TotalBalances())}
							value={
								Commitions.balanceData
									? CurrencyFormater(
										Commitions.balanceData.btcTotalSum + ' BTC',
									)
									: '-'
							}
							small={true}
						/>
						<DetailRow
							title={t(translations.CommonTitles.TotalOntrade())}
							value={
								Commitions.balanceData
									? CurrencyFormater(
										Commitions.balanceData.btcInOrderSum + ' BTC',
									)
									: '-'
							}
							small={true}
							last={true}
						/>
					</div>
					<div className="bottomActions">
						<IsLoadingWithTextAuto
							lightBackground={true}
							text={t(translations.CommonTitles.SetAsUncheck())}
							loadingId={'pending' + rowData.userId + '' + details.id}
							className={Buttons.VeryLightBlueButton}
							onClick={() => {
								handleAdminStatusChange('pending');
							}}
						/>
						<IsLoadingWithTextAuto
							lightBackground={true}
							text={t(translations.CommonTitles.NeedsRecheck())}
							loadingId={'recheck' + rowData.userId + '' + details.id}
							className={Buttons.LightYellowButton}
							onClick={() => {
								handleAdminStatusChange('recheck');
							}}
						/>
						<IsLoadingWithTextAuto
							lightBackground={true}
							text={t(translations.CommonTitles.MoveToAccepts())}
							loadingId={'approved' + rowData.userId + '' + details.id}
							className={Buttons.VeryLightGreenButton}
							onClick={() => {
								handleAdminStatusChange('approved');
							}}
						/>
					</div>
				</div>
			</div>
		</div>
	);
};
export default memo(WithdrawModal);


