import { GridLoading } from 'app/components/grid_loading/gridLoading';
import { InitialUserDetails } from 'app/containers/UserAccounts/types';
import React, { memo, useEffect, useState, useCallback, useRef } from 'react';
import { useDispatch, useSelector } from 'react-redux';
import styled from 'styled-components/macro';

import ImagesGrid from './ImagesGrid';
import ImageWrapper from './ImageWrapper';
import LeftInfo from './LeftInfo';
import { IdentityTypes } from '../constants';
import { selectUserImages } from '../selectors';
import { VerificationWindowActions } from '../slice';
import { ProfileImageData } from '../types';


interface Props {
  initialData: InitialUserDetails;
  type: IdentityTypes;
}

function Identity(props: Props) {
  const { initialData, type } = props;

  const dispatch = useDispatch();
  const [IsLoading, setIsLoading] = useState(true);
  const [Images, setImages] = useState<ProfileImageData[]>([]);
  const metaData = useRef<{ identity: { name: string }[]; address: { name: string }[] } | null>(null);
  const [SelectedImage, setSelectedImage] = useState<ProfileImageData | null>(null);
  const userImagesState = useSelector(selectUserImages);

  useEffect(() => {
    setIsLoading(true);
    dispatch(
      VerificationWindowActions.GetUserImagesAction({
        user_id: initialData.id,
        type,
      }),
    );
  }, [type]);

  useEffect(() => {
    if (
      userImagesState.userId === initialData.id &&
      userImagesState.data !== null
    ) {
      metaData.current = {
        address: userImagesState.data.userProfileImagesMetaData.types[0].subTypes,
        identity: userImagesState.data.userProfileImagesMetaData.types[1].subTypes,
      };
      setIsLoading(false);
      setImages(userImagesState.data.profileImages);
      setSelectedImage(userImagesState.data.profileImages[0]);
    }
  }, [userImagesState]);
  const handleGridSelect = useCallback(
    (image: ProfileImageData) => {
      setSelectedImage({ ...image });
    },
    [Images, SelectedImage],
  );
  return (
    <>
      {IsLoading === true && <GridLoading />}
      <Wrapper>
        <ImageAndInfoWrapper style={{ marginBottom: '10px' }}>
          <div className="info">
            <LeftInfo
              subTypes={metaData.current!}
              selectedImage={SelectedImage ?? Images[0]}
              type={type}
              data={initialData}
            />
          </div>
          <div className="imageWrapper">
            {SelectedImage && (
              <ImageWrapper
                type={type}
                selectedImage={SelectedImage}
                data={initialData}
              />
            )}
            {!SelectedImage && <div className="imageWrapperPlaceHolder"></div>}
          </div>
        </ImageAndInfoWrapper>
        <ImagesGrid
          images={Images}
          selectedImage={SelectedImage}
          onImageSelect={handleGridSelect}
        />
      </Wrapper>
    </>
  );
}

export default memo(Identity);
const ImageAndInfoWrapper = styled.div`
  display: flex;
  .info {
    flex: 3;
    padding: 7px;
  }
  .imageWrapper {
    flex: 5;
    max-width: 66%;
    min-width: 66%;
    padding: 0 20px;
  }
`;
const Wrapper = styled.div`
  display: flex;
  flex-direction: column;
`;
