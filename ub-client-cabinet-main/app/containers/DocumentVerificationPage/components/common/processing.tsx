import React, { memo, useState, useMemo } from 'react';
import styled from 'styles/styled-components';
import translate from '../../messages';
import { FormattedMessage } from 'react-intl';
import PopupModal from 'components/materialModal/modal';
import { IconButton } from '@material-ui/core';
import ExpandIcon from 'images/themedIcons/expandIcon';
import Anime from 'react-anime';

interface Props {
  image: string;
}

function Processing (props: Props) {
  const [IsImageOpen, setIsImageOpen] = useState(false);

  const { image } = props;
  return (
    <Wrapper>
      <PopupModal
        onClose={() => {
          setIsImageOpen(false);
        }}
        isOpen={IsImageOpen}
      >
        <div
          style={{
            display: 'flex',
            flexDirection: 'column',
            alignItems: 'center',
            padding: '48px',
          }}
        >
          <img src={image} style={{ maxWidth: '610px' }} />
        </div>
      </PopupModal>

      <div className='proccessing'>
        {useMemo(
          () => (
            <Anime
              className='animated'
              duration={500}
              easing='easeInCirc'
              scale={[0, 1]}
              opacity={[0, 1]}
            >
              <FormattedMessage
                {...translate.YourDocumentHasBeenReceivedItWillBeReviewedSoon}
              />
            </Anime>
          ),
          [],
        )}
      </div>
      <div className='image'>
        <img src={image} alt='' />
      </div>
      <div className='expandImage'>
        <IconButton
          onClick={() => {
            setIsImageOpen(true);
          }}
        >
          <ExpandIcon />
        </IconButton>
      </div>
    </Wrapper>
  );
}

export default memo(Processing);
const Wrapper = styled.div`
  width: 100%;
  position: relative;
  height: 100%;
  border-radius: 12px;
  border: 1px dashed var(--textBlue);
  .proccessing {
    position: absolute;
    width: 100%;
    background: var(--textBlue);
    border-top-right-radius: 8px;
    border-top-left-radius: 8px;
    height: 65px;
    display: flex;
    align-items: center;
    justify-content: center;
    z-index: 1;
    text-align: center;
    span {
      color: white;
      padding: 0 12px;
      font-size: 13px;
      font-weight: 600;
    }
  }
  .image {
    max-width: 100%;
    display: flex;
    justify-content: center;
    height: 100%;
    align-items: center;
    img {
      max-width: 95%;
      max-height: 100%;
    }
  }
  .expandImage {
    position: absolute;
    left: calc(50% - 22px);
    z-index: 5;
    top: calc(50% - 23px);
  }
`;
