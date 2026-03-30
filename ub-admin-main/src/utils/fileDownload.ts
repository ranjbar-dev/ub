import { LocalStorageKeys } from 'services/constants';

/**
 * Triggers a browser file download from a Blob.
 * Creates a temporary anchor element, assigns the blob URL, clicks it, then cleans up.
 */
const fileDownloader = (data: BlobPart, filename: string): void => {
  const blob = new Blob([data], { type: 'application/octet-stream' });
  const blobURL = URL.createObjectURL(blob);
  const tempLink = document.createElement('a');

  tempLink.style.display = 'none';
  tempLink.href = blobURL;
  tempLink.setAttribute('download', filename);

  // Safari without HTML5 download attribute support
  if (typeof tempLink.download === 'undefined') {
    tempLink.setAttribute('target', '_blank');
  }

  document.body.appendChild(tempLink);
  tempLink.click();

  // Cleanup — fixes "webkit blob resource error 1"
  setTimeout(() => {
    document.body.removeChild(tempLink);
    URL.revokeObjectURL(blobURL);
  }, 200);
};

/**
 * Fetches a file from a URL and triggers browser download.
 *
 * @param params.url - URL to fetch the file from
 * @param params.filename - Filename for the downloaded file
 */
const downloadFile = async (params: { url: string; filename: string }): Promise<void> => {
  try {
    const token = localStorage.getItem(LocalStorageKeys.ACCESS_TOKEN) || '';
    const response = await fetch(params.url, {
      headers: {
        Authorization: `Bearer ${token}`,
      },
    });
    if (!response.ok) {
      throw new Error(`Download failed with status ${response.status}`);
    }
    const blob = await response.blob();
    fileDownloader(blob, params.filename);
  } catch (error) {
    console.error('File download failed:', error);
  }
};

export default downloadFile;
