import { TextField } from '@material-ui/core';
import { toast } from 'app/components/Customized/react-toastify';
import { GridLoading } from 'app/components/grid_loading/gridLoading';
import IsLoadingWithTextAuto from 'app/components/isLoadingWithText/isLoadingWithTextAuto';
import { Buttons } from 'app/constants';
import React, { memo, useState, useEffect, useRef } from 'react';
import { useDispatch } from 'react-redux';
import { LocalStorageKeys } from 'services/constants';
import { Subscriber, BroadcastMessage } from 'services/messageService';
import styled from 'styled-components/macro';
import { LoadingIds } from 'utils/loading';

import { ScanBlockActions } from './slice';
import { Currency } from './types';
import SelectOptions from '../MarketTicks/SelectOptions';

interface Props { }

function ScanPage(props: Props) {
    const networks = useRef<{ name: string; code: string }[]>([])
    const currencies: Currency[] = JSON.parse(localStorage[LocalStorageKeys.CURRENCIES]).currencies
    const [hasNetwork, setHasNetwork] = useState(false)
    const dispatch = useDispatch();
    let sendingData = useRef({
        network: '',
        block_number: -1
    });
    useEffect(() => {
        const Subscription = Subscriber.subscribe((message: BroadcastMessage) => {

        });
        return () => {
            Subscription.unsubscribe()
        };
    }, []);




    const handleScanClick = () => {
        if (sendingData.current.network === '' || sendingData.current.block_number === -1) {
            toast.warn('All fields are required')
            return
        }
        let data = sendingData.current;
        dispatch(ScanBlockActions.Scan(data));
    };
    const handleBlockNumberChange = (v: string) => {
        if (v == '') {
            sendingData.current.block_number = -1
            return
        }
        sendingData.current.block_number = Number(v)
    }
    const handleCoinSelect = (e: Currency) => {
        if (e.otherBlockChainNetworks.length > 0) {
            networks.current = [{ code: e.code, name: e.name }, ...(e.otherBlockChainNetworks)]
            sendingData.current.network = ''
            setHasNetwork(true)
        }
        else {
            networks.current = []
            setHasNetwork(false)
            if(e.mainNetwork==''){
                e.mainNetwork=e.code
            }
            sendingData.current.network = e.mainNetwork
        }
    }
    const handleNetworkSelect = (network: string) => {
        sendingData.current.network = network
    }
    return (
        <>
            {currencies.length === 0 ? (
                <GridLoading />
            ) : (
                <Wrapper>
                    <div className="elementWrapper">
                        <div className="elementTitle">Coin : </div>
                        <SelectOptions

                            onSelectObj={e => {
                                // console.log(e);
                                if (e?.obj) handleCoinSelect(e.obj as Currency)
                            }}
                            initialValue={-1 + ''}
                            options={currencies.map((item, index) => {
                                return { name: item.name, value: item.id + '', obj: item };
                            })}
                        />
                    </div>
                    {hasNetwork && <div className="elementWrapper">
                        <div className="elementTitle">network : </div>
                        <SelectOptions
                            styles={{ minWidth: '150px' }}
                            onSelectObj={e => {
                                // console.log(e);
                                const obj = e?.obj as { code: string } | undefined;
                                if (obj) handleNetworkSelect(obj.code)
                            }}
                            initialValue={-1 + ''}
                            options={networks.current.map((item: { name: string; code: string }, index: number) => {
                                return { name: item.name, value: item.code + '', obj: item };
                            })}
                        />
                    </div>}
                    <div className="elementWrapper last">
                        <div className="elementTitle">Block Number : </div>
                        <TextField variant="outlined" margin="dense" onChange={(e) => handleBlockNumberChange(e.target.value)}

                        />
                    </div>

                    <IsLoadingWithTextAuto
                        text="Scan"
                        loadingId={LoadingIds.ScanButton}
                        style={{ position: 'absolute', left: '356px', marginTop: '20px' }}
                        className={Buttons.SkyBlueButton}
                        onClick={() => {
                            handleScanClick();
                        }}
                    />
                </Wrapper>
            )}
        </>
    );
}

export default memo(ScanPage);
const Wrapper = styled.div`
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
