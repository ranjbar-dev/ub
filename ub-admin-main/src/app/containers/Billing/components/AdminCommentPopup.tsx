import { Button } from '@material-ui/core';
import IsLoadingWithTextAuto from 'app/components/isLoadingWithText/isLoadingWithTextAuto';
import UbTextAreaAutoHeight from 'app/components/UbTextAreaAutoHeight/UbTextAreaAutoHeight';
import { Buttons } from 'app/constants';
import { translations } from 'locales/i18n';
import React, { memo, useRef } from 'react';
import { useTranslation } from 'react-i18next';

interface Props {
  onAdminCommentSubmitted: (value: string) => void;
  onCancel: () => void;
  initialValue: string;
}

function AdminCommentPopup(props: Props) {
  const { t } = useTranslation();
  const { initialValue, onCancel, onAdminCommentSubmitted } = props;
  const adminComment = useRef(initialValue);
  return (
    <div className="rejectPopup">
      <div className="title">{t(translations.CommonTitles.AdminReports())}</div>
      <div className="inp" style={{ marginBottom: '24px' }}>
        <UbTextAreaAutoHeight
          initialValue={initialValue}
          onChange={e => {
            adminComment.current = e;
          }}
        />
      </div>
      <div className="buttons">
        <IsLoadingWithTextAuto
          style={{ marginRight: '12px' }}
          text={t(translations.CommonTitles.Submit())}
          loadingId="reject"
          className={Buttons.SkyBlueButton}
          onClick={() => onAdminCommentSubmitted(adminComment.current)}
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

export default memo(AdminCommentPopup);
