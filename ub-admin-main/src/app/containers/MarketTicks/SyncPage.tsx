import SyncIcon from '@material-ui/icons/Sync';
import DatePick from 'app/components/BackupDatePicker/datePick';
import {toast} from 'app/components/Customized/react-toastify';
import {GridLoading} from 'app/components/grid_loading/gridLoading';
import IsLoadingWithTextAuto from 'app/components/isLoadingWithText/isLoadingWithTextAuto';
import UbCheckbox from 'app/components/UbCheckBox/UbCheckbox';
import {Buttons} from 'app/constants';
import React,{memo,useState,useEffect,useRef} from 'react';
import {useDispatch} from 'react-redux';
import {Subscriber,MessageNames,BroadcastMessage} from 'services/messageService';
import styled from 'styled-components/macro';

import {MarketTicksTimeFrames} from './constants';
import SelectOptions from './SelectOptions';
import {MarketTicksActions} from './slice';
interface Props { }

function SyncPage(props: Props) {
	const [CurrencyPairs,setCurrencyPairs] = useState<{ name: string; id: number }[]>([]);
	const dispatch=useDispatch();
	let sendingData=useRef({
		end_time: '',
		pair_currency_id: '1',
		start_time: '',
		time_frame: '1hour',
		update: 0
	});

	//t(translations.CommonTitles.ID())
	useEffect(() => {
		const Subscription=Subscriber.subscribe((message: BroadcastMessage) => {
			if(message.name===MessageNames.SET_CURRENCY_PAIRS) {
				setCurrencyPairs(message.payload as { name: string; id: number }[]);
			}
		});
		dispatch(MarketTicksActions.GetCurrencyPairsAction({}));
		return () => {
			Subscription.unsubscribe();
		};
	},[]);
	const handleSyncClick=() => {
		let data=sendingData.current;
		if(!data.end_time||!data.start_time) {
			toast.warn('all fields are required');
			return;
		}
		dispatch(MarketTicksActions.SyncTicksAction(sendingData.current));
	};
	const handleUpdateCheckboxChange = (e: boolean) => {
		sendingData.current.update=e? 1:0
	}
	return (
		<>
			{CurrencyPairs.length===0? (
				<GridLoading />
			):(
					<Wrapper>
						<div className="elementWrapper">
							<div className="elementTitle">PairName : </div>
							<SelectOptions
								onSelect={(e: string) => {
									//console.log(e);
									sendingData.current.pair_currency_id=e+'';
								}}
								initialValue={1+''}
								options={CurrencyPairs.map((item: { name: string; id: number },index: number) => {
									return {name: item.name,value: item.id+''};
								})}
							/>
						</div>
						<div className="elementWrapper">
							<div className="elementTitle">Time Frame : </div>
							<SelectOptions
								onSelect={(e: string) => {
									//console.log(e);
									sendingData.current.time_frame=e;
								}}
								initialValue={'1hour'}
								options={MarketTicksTimeFrames.map((item,index) => {
									return {name: item.name,value: item.value+''};
								})}
							/>
						</div>
						<div className="elementWrapper">
							<div className="elementTitle">Start Time : </div>
							<DatePick
								title="from"
								onChange={e => {
									sendingData.current.start_time=e.replace('T',' ')+':00';
								}}
							/>
						</div>
						<div className="elementWrapper last">
							<div className="elementTitle">End Time : </div>
							<DatePick
								title="to"
								onChange={e => {
									sendingData.current.end_time=e.replace('T',' ')+':00';
								}}
							/>
						</div>
						<div className="updateCheckbox">
							<UbCheckbox initialValue={false} title='Update' titlePlacement='start' onChange={handleUpdateCheckboxChange} />
						</div>

						<IsLoadingWithTextAuto
							icon={<SyncIcon />}
							text="Start Sync"
							loadingId="syncButton"
							style={{position: 'absolute',left: '356px',marginTop: '20px'}}
							className={Buttons.SkyBlueButton}
							onClick={() => {
								handleSyncClick();
							}}
						/>
					</Wrapper>
				)}
		</>
	);
}

export default memo(SyncPage);
const Wrapper=styled.div`
  padding-top: 40px;
  background: ${p => p.theme.white};
  height: 765px;
  margin-top: -10px;
  .updateCheckbox{
	  margin-left:12px;
	.MuiTypography-body1{
		font-size: 13px;
        font-weight: 600;
        color: #787878;
	}
  }
  .elementWrapper {
    display: flex;
    max-width: 473px;
    padding: 8px 16px;
    border-bottom: 2px solid #c2c0c0;
    align-items: center;
    &.last {
      margin-bottom: 12px;
    }
    .elementTitle {
      margin: 0 10px;
      min-width: 105px;
      font-size: 13px;
      font-weight: 600;
      color: #787878;
    }
  }
  .MuiSelect-outlined.MuiSelect-outlined {
    background: #f5f6f8;
  }
  .MuiOutlinedInput-inputMarginDense {
    background: #f5f6f8;
  }
  .MuiSelect-selectMenu {
    padding: 7px 10px !important;
    min-width: 100px;
  }
  .MuiOutlinedInput-inputMarginDense {
    padding-top: 7.5px;
    padding-bottom: 6.5px;
  }
`;
