import { Permission } from 'app/containers/UserDetails/types';

/* --- STATE --- */
export interface UserImagesData {
  profileImages: ProfileImageData[];
  userProfileImagesMetaData: {
    types: Array<{ subTypes: Array<{ id: number; name: string; type?: string }> }>;
  };
}

export interface VerificationWindowState {
  userImages: { userId: number | null; data: UserImagesData | null };
  permissionsData: { userId: number | null; data: Permission[] | null };
}

export type ContainerState = VerificationWindowState;
export interface ProfileImageData {
  confirmationStatus: string;
  createdAt: string;
  id: number;
  idCardCode: string;
  imagePath: string;
  isBack: boolean;
  mainImageId: number | null;
  originalFileName: string;
  rejectionReason: string | null;
  subType: string;
  type: string;
  updatedAt: string;
}
export enum ImageStatusStrings {
  Confirmed = 'CONFIRMED',
  Rejected = 'REJECTED',
}
