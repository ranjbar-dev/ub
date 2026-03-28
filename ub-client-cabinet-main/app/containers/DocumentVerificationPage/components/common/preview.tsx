import React, { useState, useMemo } from 'react';
import styled from 'styles/styled-components';
import Anime from 'react-anime';
import { IconButton } from '@material-ui/core';
import PopupModal from 'components/materialModal/modal';
import ExpandIcon from 'images/themedIcons/expandIcon';
import DeleteIcon from 'images/themedIcons/deleteIcon';

export default (props: { previewFile: File; deletePreview: () => void }) => {
  const { previewFile, deletePreview } = props;
  const [IsImageOpen, setIsImageOpen] = useState(false);
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
          <img
            src={URL.createObjectURL(previewFile)}
            style={{ maxWidth: '610px' }}
          />
        </div>
      </PopupModal>

      {useMemo(
        () => (
          <Anime
            duration={300}
            easing='linear'
            translateY={[40, 0]}
            opacity={[0, 1]}
            className='animated'
          >
            <PreviewImage src={URL.createObjectURL(previewFile)} alt='' />
          </Anime>
        ),
        [previewFile],
      )}
      <div className='expandImage'>
        {useMemo(
          () => (
            <Anime
              duration={200}
              delay={200}
              easing='easeOutElastic'
              scale={[0.5, 1]}
              opacity={[0, 1]}
              className='animated'
            >
              <IconButton
                onClick={() => {
                  setIsImageOpen(true);
                }}
              >
                <ExpandIcon />
              </IconButton>
            </Anime>
          ),
          [previewFile],
        )}
      </div>
      <div className='deleteButton'>
        {useMemo(
          () => (
            <Anime
              duration={200}
              delay={400}
              easing='easeOutElastic'
              scale={[0.5, 1]}
              opacity={[0, 1]}
            >
              <IconButton
                onClick={() => {
                  deletePreview();
                }}
              >
                <DeleteIcon />
              </IconButton>
            </Anime>
          ),
          [previewFile],
        )}
      </div>
    </Wrapper>
  );
};

const PreviewImage = styled.img`
  pointer-events: 'none';
  -webkit-user-drag: none;
  -khtml-user-drag: none;
  -moz-user-drag: none;
  -o-user-drag: none;
  user-drag: none;
`;

const Wrapper = styled.div`
  display: flex;
  justify-content: center;
  height: 100%;
  align-items: center;
  position: relative;
  img {
    max-width: 95%;
    max-height: 100%;
  }
  .deleteButton {
    position: absolute;
    left: 24px;
    bottom: 24px;
    .MuiIconButton-root {
      background: var(--darkTrans);
      &:hover {
        background: var(--darkTrans) !important;
        opacity: 0.9;
      }
    }
  }
  .animated0 {
    display: flex;
    justify-content: center;
    height: 100%;
    align-items: center;
    position: relative;
  }
`;
