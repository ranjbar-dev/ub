import { Button } from '@material-ui/core';
import { Buttons } from 'app/constants';
import React, { memo, useRef } from 'react';


interface Props {
    onConfirm: () => void
    content: string
}

function ConfirmPopup(props: Props) {
    return (
        <div className="withdrawWindowConfirmPopupWrapper">
            <p>{props.content}</p>
            <Button
                onClick={() => {
                    props.onConfirm()
                }}
                className={Buttons.SkyBlueButton}
            >
                OK
            </Button>
        </div>
    );
}

export default memo(ConfirmPopup);
