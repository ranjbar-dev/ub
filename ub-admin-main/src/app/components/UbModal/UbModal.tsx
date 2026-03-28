import React, { memo, useState, useEffect } from 'react';
import Anime, { anime } from 'react-anime';
import styled from 'styled-components/macro';

interface Props {
  children: React.ReactNode;
  onClose: () => void;
  title?: string;
}

/**
 * Animated overlay modal with a backdrop click-to-close behaviour.
 * Uses anime.js for scale-in/scale-out transitions.
 *
 * @example
 * ```tsx
 * <UbModal onClose={() => setOpen(false)}>
 *   <MyPopupContent />
 * </UbModal>
 * ```
 */
function UbModal(props: Props) {
  const { children, title } = props;
  const [AnimationStart, setAnimationStart] = useState(false);
  useEffect(() => {
    setAnimationStart(true);
    return () => {};
  }, []);
  useEffect(() => {
    const handleKeyDown = (e: KeyboardEvent) => {
      if (e.key === 'Escape') {
        handleClose();
      }
    };
    document.addEventListener('keydown', handleKeyDown);
    return () => {
      document.removeEventListener('keydown', handleKeyDown);
    };
  }, []);
  const handleClose = () => {
    setAnimationStart(false);
    setTimeout(() => {
      props.onClose();
    }, 200);
  };
  return (
    <Overlay>
      <BackDrop
        className="backDrop"
        visible={AnimationStart}
        onClick={handleClose}
      />
      <ChildrenWrapper
        className="children"
        role="dialog"
        aria-modal="true"
        aria-labelledby={title ? 'ub-modal-title' : undefined}
        visible={AnimationStart}
      >
        {title && (
          <h2 id="ub-modal-title" style={{ margin: 0 }}>
            {title}
          </h2>
        )}
        {children}
      </ChildrenWrapper>
    </Overlay>
  );
}

export default memo(UbModal);

const Overlay = styled.div`
  position: fixed;
  width: 100%;
  height: 100%;
  display: flex;
  align-items: center;
  z-index: 1;
  top: 0;
  justify-content: center;
`;

const BackDrop = styled.div<{ visible: boolean }>`
  position: absolute;
  opacity: ${p => (p.visible ? '1' : '0')};
  transition: opacity 0.2s;
  width: 100%;
  height: 100%;
  background: rgba(124, 126, 130, 0.38); /* TODO: add backdrop overlay color to theme */
  z-index: 1;
`;

const ChildrenWrapper = styled.div<{ visible: boolean }>`
  z-index: 2;
  position: absolute;
  transition: transform 0.2s;
  transform: ${p => (p.visible ? 'scale(1)' : 'scale(0)')};
`;
