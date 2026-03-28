import { takeLatest, put, call, takeEvery } from 'redux-saga/effects';
import ActionTypes, { ProfileImageStatus, UploadState } from './constants';
import {
  setIsLoadingUserProfileDataAction,
  setUserProfileAction,
  setUploadedFileAction,
  deleteUserImageAction,
  getUserProfileAction,
} from './actions';
import { toast } from 'components/Customized/react-toastify';
import { StandardResponse } from 'services/constants';
import {
  getUserProfileAPI,
  deleteUserImageAPI,
} from 'services/user_acount_service';
import {
  UploadModel,
  UploadFile,
  UploadMultiFile,
} from 'services/upload_service';
import { MessageService, MessageNames } from 'services/message_service';
import { replace } from 'redux-first-history';
import { AppPages } from 'containers/App/constants';
import { KycStatus } from 'containers/AcountPage/constants';
import { IUserProfileImage, IUserProfileMetaData } from './types';

function * getUserProfile (action: {
  type: string;
  payload: { silent?: boolean };
}) {
  if (!action.payload.silent) {
    yield put(setIsLoadingUserProfileDataAction(true));
  }

  try {
    const response: StandardResponse = yield call(getUserProfileAPI);
    if (response.status === false) {
      toast.error('error while getting user profile');
      if (!action.payload.silent) {
        yield put(setIsLoadingUserProfileDataAction(false));
      }
      return;
    }
    if (response.data.status === KycStatus.CONFIRMED) {
      toast.info('your acount has been verified');
      yield put(replace(AppPages.AcountPage));
    }

    const {
      userProfileImages,
      userProfileImagesMetaData,
    }: {
      userProfileImages: IUserProfileImage[];
      userProfileImagesMetaData: IUserProfileMetaData;
    } = response.data;

    yield put(setUserProfileAction(response.data));

    MessageService.send({
      name: MessageNames.SET_DOCUMENT_IMAGES,
      payload: {
        userProfileImages: userProfileImages ?? [],
        userProfileImagesMetaData,
      },
    });
  } catch (error) {
    toast.error('error getting user profile');
    if (!action.payload.silent) {
      yield put(setIsLoadingUserProfileDataAction(false));
    }
  }
}

function * uploadFile (action: { type: string; payload: UploadModel }) {
  const response: StandardResponse = yield call(UploadFile, action.payload);
  try {
    const res: StandardResponse = response.data;
    if (res.status === false) {
      MessageService.send({
        name: MessageNames.SET_UPLOADER_STATE,
        payload: {
          type: action.payload.type,
          uploadState: UploadState.READY,
        },
      });
      if (res.message) {
        toast.error(res.message);
        return;
      }
      toast.error('error while uploading file');
      return;
    } else if (res.status === true) {
      yield put(
        setUploadedFileAction({
          type: action.payload.type,
          image: res.data[0].image,
          id: res.data[0].id,
          isBack: action.payload.isBack ?? undefined,
        }),
      );

      MessageService.send({
        name: MessageNames.SET_UPLOADED_IMAGE,
        payload: {
          type: action.payload.type,
          image: res.data[0].image,
          id: res.data[0].id,
          fromSaga: true,
          status: ProfileImageStatus.PROCCESSING,
        },
      });
    }
  } catch (error) {
    toast.error('error while uplading image,maximum image size is 5mb');
    MessageService.send({
      name: MessageNames.SET_ERROR,
      errorId: 'upload' + action.type,
    });
  }
}
function * uploadProfileImage (action: { type: string; payload: UploadModel }) {
  const response: StandardResponse = yield call(UploadFile, action.payload);
  try {
    const res: StandardResponse = response.data;
    if (res.status === false) {
      MessageService.send({
        name: MessageNames.RESET_IMAGES,
        payload: {
          uploaderId:
            action.payload.type +
            action.payload.subtype +
            action.payload.isBack,
          newState: UploadState.READY,
        },
      });
      if (res.message) {
        toast.error(res.message);
        return;
      }
      toast.error('error while uploading file');
      return;
    } else if (res.status === true) {
      // yield put(
      //   setUploadedFileAction({
      //     type: action.payload.type,
      //     image: res.data[0].image,
      //     id: res.data[0].id,
      //     isBack: action.payload.isBack ?? undefined,
      //   }),
      // );
      // MessageService.send({
      //   name: MessageNames.SET_UPLOADED_IMAGE,
      //   payload: {
      //     type: action.payload.type,
      //     image: res.data[0].image,
      //     id: res.data[0].id,
      //     fromSaga: true,
      //     status: ProfileImageStatus.PROCCESSING,
      //   },
      // });
      yield put(getUserProfileAction({ silent: true }));
      // MessageService.send({
      //   name: MessageNames.RESET_IMAGES,
      //   payload: {
      //     uploaderId:
      //       action.payload.type +
      //       action.payload.subtype +
      //       action.payload.isBack,
      //     newState: UploadState.PROCESSING,
      //   },
      // });
      MessageService.send({
        name: MessageNames.TOGGLE_SEND_IMAGE_BUTTON,
        payload: false,
        additional: action.payload.type,
      });
      // MessageService.send({
      //   name: MessageNames.RESET_IMAGES,
      //   payload: {
      //     uploaderId: action.payload.type + action.payload.subtype + true,
      //     newState: UploadState.READY,
      //   },
      // });
    }
  } catch (error) {
    // toast.error('error while uplading image,maximum image size is 5mb');
    // MessageService.send({
    //   name: MessageNames.SET_ERROR,
    //   errorId: 'upload' + action.type,
    // });
  }
}
function * uploadMultiProfileImage (action: {
  type: string;
  payload: {
    frontImage: File;
    backImage: File;
    type: string;
    subtype: string;
    front_image_id?: number | string;
    back_image_id?: number | string;
  };
}) {
  const { type, subtype } = action.payload;
  try {
    const response: StandardResponse = yield call(
      UploadMultiFile,
      action.payload,
    );
    const res: StandardResponse = response.data;
    if (res.status === false) {
      if (action.payload.frontImage) {
        MessageService.send({
          name: MessageNames.RESET_IMAGES,
          payload: {
            uploaderId: action.payload.type + action.payload.subtype + 'false',
            newState: UploadState.READY,
          },
        });
      }
      if (action.payload.backImage) {
        MessageService.send({
          name: MessageNames.RESET_IMAGES,
          payload: {
            uploaderId: action.payload.type + action.payload.subtype + 'true',
            newState: UploadState.READY,
          },
        });
      }
      if (res.message) {
        toast.error(res.message);
        return;
      }
      toast.error('error while uploading file');
      return;
    } else if (res.status === true) {
      const images = res.data;
      if (images.frontImage) {
        MessageService.send({
          name: MessageNames.SET_UPLOADED_IMAGE,
          payload: {
            uploaderId: type + subtype + false,
            image: images.frontImage.path,
            disable: type,
            imageId: images.frontImage.id,
          },
        });
      }
      if (images.backImage) {
        MessageService.send({
          name: MessageNames.SET_UPLOADED_IMAGE,
          payload: {
            uploaderId: type + subtype + true,
            image: images.backImage.path,
            imageId: images.backImage.id,
            disable: type,
          },
        });
      }

      // yield put(
      //   setUploadedFileAction({
      //     type: action.payload.type,
      //     image: res.data[0].image,
      //     id: res.data[0].id,
      //     isBack: action.payload.isBack ?? undefined,
      //   }),
      // );
      // MessageService.send({
      //   name: MessageNames.SET_UPLOADED_IMAGE,
      //   payload: {
      //     type: action.payload.type,
      //     image: res.data[0].image,
      //     id: res.data[0].id,
      //     fromSaga: true,
      //     status: ProfileImageStatus.PROCCESSING,
      //   },
      // });
      // yield put(getUserProfileAction({ silent: true }));
    }
  } catch (error) {
    console.log(error);
    // if (action.payload.frontImage) {
    //   MessageService.send({
    //     name: MessageNames.RESET_IMAGES,
    //     payload: {
    //       uploaderId: action.payload.type + action.payload.subtype + 'false',
    //       newState: UploadState.READY,
    //     },
    //   });
    // }
    // if (action.payload.backImage) {
    //   MessageService.send({
    //     name: MessageNames.RESET_IMAGES,
    //     payload: {
    //       uploaderId: action.payload.type + action.payload.subtype + 'true',
    //       newState: UploadState.READY,
    //     },
    //   });
    // }
    // toast.error('error uploading file');
    // return;
  }
}
function * deleteFile (action: {
  type: string;
  payload: { id: number; type: string };
}) {
  const response: StandardResponse = yield call(deleteUserImageAPI, {
    id: action.payload.id,
  });
  const res: StandardResponse = response.data;
  if (res.status === false) {
    toast.error('error while uploading file');
    return;
  }
  toast.success('image successfuly deleted');
  yield put(deleteUserImageAction(action.payload));
  MessageService.send({
    name: MessageNames.DELETE_UPLOADED_IMAGE,
    payload: action.payload,
  });
}

export default function * documentVerificationPageSaga () {
  yield takeLatest(ActionTypes.GET_USER_PROFILE, getUserProfile);
  yield takeEvery(ActionTypes.UPLOAD_FILE, uploadProfileImage);
  yield takeEvery(ActionTypes.UPLOAD_MULTI_FILE, uploadMultiProfileImage);
  yield takeLatest(ActionTypes.DELETE_FILE, deleteFile);
}
