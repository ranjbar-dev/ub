# Task: Convert fileDownload.js to TypeScript

**ID:** p1-filedownload  
**Phase:** 1 — Type Safety Foundation  
**Severity:** 🟠 MEDIUM  
**Dependencies:** None  

## Problem

`src/utils/fileDownload.js` is the only JavaScript file in the `src/` directory. It has zero type safety and uses the deprecated `msSaveBlob` IE11 API.

## File to Modify

**Rename:** `src/utils/fileDownload.js` → `src/utils/fileDownload.ts`

### Current Content (51 lines)
```javascript
const FileDownloader = (data, filename) => {
  if (window) {
    const blobData = [data]
    const blob = new window.Blob(blobData, { type: 'application/octet-stream' })
    if (typeof window.navigator.msSaveBlob !== 'undefined') {
      window.navigator.msSaveBlob(blob, filename)
    } else {
      const blobURL =
        window.URL && window.URL.createObjectURL
          ? window.URL.createObjectURL(blob)
          : window.webkitURL.createObjectURL(blob)
      const tempLink = document.createElement('a')
      tempLink.style.display = 'none'
      tempLink.href = blobURL
      tempLink.setAttribute('download', filename)

      // Safari thinks _blank anchor are pop ups. We only want to set _blank
      // target if the browser does not support the HTML5 download attribute.
      // This allows us to download files in desktop safari if pop up blocking
      // is enabled.

      if (typeof tempLink.download === 'undefined') {
        tempLink.setAttribute('target', '_blank')
      }

      document.body.appendChild(tempLink)
      tempLink.click()

      // Fixes "webkit blob resource error 1"
      setTimeout(function () {
        document.body.removeChild(tempLink)
        window.URL.revokeObjectURL(blobURL)
      }, 200)
    }
  }
}

const downloadFile = ({ url, filename }) => {
  try {
    fetch(url)
      .then(res => res.blob())
      .then(blob => {
        FileDownloader(blob, filename)
      })
  } catch (e) {
    console.log(e)
  }
}

export default downloadFile
```

### Target Content
```typescript
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
    const response = await fetch(params.url);
    const blob = await response.blob();
    fileDownloader(blob, params.filename);
  } catch (error) {
    console.error('File download failed:', error);
  }
};

export default downloadFile;
```

### Changes Made
1. Added TypeScript types to all parameters and return values
2. Removed deprecated `msSaveBlob` IE11 code path
3. Removed `webkitURL` fallback (obsolete in modern browsers)
4. Converted `downloadFile` to async/await
5. Added JSDoc comments
6. Renamed `FileDownloader` to `fileDownloader` (camelCase convention)

## Import Consumers

The only import is in `src/app/index.tsx`:
```typescript
import downloadFile from 'utils/fileDownload';
```
This import does NOT include the file extension, so renaming `.js` → `.ts` requires no import changes.

## Validation

```bash
npm run checkTs   # Must pass
npm test          # Must pass
npm run build     # Must build
```
