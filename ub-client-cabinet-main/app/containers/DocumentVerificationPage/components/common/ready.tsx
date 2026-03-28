import React from 'react';
import uploadIcon from 'images/uploadIcon.svg';
import { FormattedMessage } from 'react-intl';
import translate from '../../messages';
import { Button } from '@material-ui/core';
import Anime, { anime } from 'react-anime';
import styled from 'styles/styled-components';

export default function Ready (props: {
  readyToUploadDescs: { fileSubtype: any; boldTitle: any };
  onBrowsClick: () => void;
}) {
  const { readyToUploadDescs, onBrowsClick } = props;
  return (
    <div className='contentWrapper'>
      <Anime
        className='animated'
        duration={1000}
        delay={anime.stagger(100)}
        easing='easeOutElastic'
        // scale={[0, 1]}
        translateY={[50, 0]}
        opacity={[0, 1]}
      >
        <div className='uploadIcon' style={{ marginBottom: '16px' }}>
          <img src={uploadIcon} alt='' />
        </div>
        <div className='decripts' style={{ marginBottom: '16px' }}>
          <ReadyTextWrapper>
            {' '}
            <FormattedMessage {...translate.Dropfiletouploadorbrowse} />{' '}
            <span className='blueItalicBold'>
              {readyToUploadDescs.boldTitle}{' '}
            </span>
            <FormattedMessage {...translate.Of} />{' '}
            <FormattedMessage {...translate.your} />{' '}
            <span>{readyToUploadDescs.fileSubtype}</span>
          </ReadyTextWrapper>
        </div>
        <div className='browsButtonWrapper'>
          <Button
            onClick={() => {
              onBrowsClick();
            }}
            variant='outlined'
            className='browsButton'
            color='primary'
          >
            <FormattedMessage {...translate.browse} />
          </Button>
        </div>
      </Anime>
      {/* </Anime> */}
    </div>
  );
}
const ReadyTextWrapper = styled.div`
  padding: 0 12px;
  text-align: center;
`;
