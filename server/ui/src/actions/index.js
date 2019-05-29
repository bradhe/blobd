import uuid from 'uuid/v5';

import { uploadWithProgress } from '../api.js';

import {
  UPLOAD_START,
  UPLOAD_PROGRESS,
  UPLOAD_COMPLETE,
  UPLOAD_FAIL,
  NEW_BLOB,
} from '../constants.js';

const NAMESPACE = '9cf62dae-39db-4df7-8408-8eadab53dc01'

const getHost = () => 'localhost:5001';

const getScheme = () => 'http';

export function upload (file) {
  return dispatch => {
    const id = uuid(file.name, NAMESPACE);

    // Indicate that the upload has started.
    dispatch({ id, type: UPLOAD_START, file: file });

    const withProgress = (percent) => {
      dispatch({ id, type: UPLOAD_PROGRESS, progress: percent });
    };

    return uploadWithProgress(file, withProgress).then((res) => {
      dispatch({ id, type: UPLOAD_COMPLETE, file: file, blob: { ...res } })
      return res;
    }).then((res) => {
      dispatch({ id, type: NEW_BLOB, filename: file.name, ...res });
    }).catch((res) => {
      dispatch({ id, type: UPLOAD_FAIL, file: file, error: { ...res } });
    });
  };
};

export function copyLink(blob) {
  return dispatch => {
    const url = `${getScheme()}://${getHost()}/${blob.id}?token=${blob.read_jwt}&dl=1`;
    console.log(url);
  };
};
