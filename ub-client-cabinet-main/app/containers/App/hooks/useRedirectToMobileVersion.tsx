import { useEffect } from 'react';
import { minWidthToRedirectToMobile, mobileUrl } from 'services/constants';
import { AppPages } from '../constants';

const ignoreUrls = [AppPages.VerifyEmail, AppPages.UpdatePassword];

export const useRedirectToMobile = ({ path }: { path: string }) => {
  useEffect(() => {
    if (window.innerWidth < minWidthToRedirectToMobile) {
      let found = false;
      ignoreUrls.forEach(item => {
        if (path.includes(item)) {
          found = true;
        }
      });
      if (found === false) {
        window.location.href = mobileUrl;
        return;
      }
    }
    return () => {};
  }, [path]);
};
