import React, { memo, useMemo } from 'react';
import Anime, { anime } from 'react-anime';
import { GridLoading } from 'components/grid_loading/gridLoading';
interface Props {
  children: React.ReactNode;
  isLoading: boolean;
  memoize?: boolean;
}
const AnimateChildren = (props: Props) => {
  const { isLoading, memoize = true } = props;
  const dependencies = memoize ? [isLoading] : undefined;
  return useMemo(
    () =>
      isLoading === false ? (
        <Anime
          className='animm'
          duration={400}
          delay={anime.stagger(100)}
          easing='easeOutCirc'
          scale={[0.5, 1]}
          opacity={[0, 1]}
        >
          {props && props.children}
        </Anime>
      ) : (
        <GridLoading />
      ),
    dependencies,
  );
};
export default memo(AnimateChildren);
