import { uploadFile } from '../api.js';

export function upload (file) {
  return dispatch => {
    // Indicate that the upload has started.
    dispatch({ type: 'UPLOAD_START', file: file });

    const withProgress = (percent) => {
      dispatch({ type: 'UPLOAD_PROGRESS', progress: percent });
    };

    return uploadFile(file, withProgress).then((res) => {
      console.log('complete', res);
      dispatch({ type: 'UPLOAD_COMPLETE', file: file });
    }).catch((res) => {
      console.log('failed', res);
      dispatch({ type: 'UPLOAD_FAIL', file: file });
    });
  };
};
