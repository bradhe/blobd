import {
  UPLOAD_START,
  UPLOAD_PROGRESS,
  UPLOAD_COMPLETE,
  UPLOAD_FAIL,
  TRACK_BLOB,
} from '../constants.js';

import { uploadWithProgress } from '../api.js';

export function upload (file) {
  return dispatch => {
    // Indicate that the upload has started.
    dispatch({ type: UPLOAD_START, file: file });

    const withProgress = (percent) => {
      dispatch({ type: UPLOAD_PROGRESS, progress: percent });
    };

    return uploadWithProgress(file, withProgress).then((res) => {
      dispatch({ type: UPLOAD_COMPLETE, file: file, blob: { ...res } });
      dispatch({ type: TRACK_BLOB, filename: file.name, blob: { ...res } });
    }).catch((res) => {
      dispatch({ type: UPLOAD_FAIL, file: file });
    });
  };
};

export function addFile(filename, blob) {
  return dispatch => {
    dispatch({ type: TRACK_BLOB, filename, blob });
  };
}
