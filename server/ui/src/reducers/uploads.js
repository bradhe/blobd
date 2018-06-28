import {
  UPLOAD_START,
  UPLOAD_PROGRESS,
  UPLOAD_COMPLETE,
  UPLOAD_FAIL
} from '../constants.js';

const addUpload = (state, id, file) => {
  return state.concat({ id, file, status: UPLOAD_START });
}

const updateProgress = (state, id, progress) => {
  return state.map((upload) => {
    if (upload.id === id) {
      upload.status = UPLOAD_PROGRESS;
      upload.progress = progress;
    }

    return upload;
  });
};

const removeUpload = (state, id) => {
  return state.filter((upload) => upload.id !== id);
};

const failUpload = (state, id, error) => {
  return state.map((upload) => {
    if (upload.id === id) {
      upload.status = UPLOAD_FAIL;
      upload.error = error;
      delete(upload.progress);
    }

    return upload;
  });
};

const reducer = (state = [], action) => {
  switch(action.type) {
    case UPLOAD_START:
      return addUpload(state, action.id, action.file)
    case UPLOAD_PROGRESS:
      return updateProgress(state, action.id, action.progress);
    case UPLOAD_COMPLETE:
      return removeUpload(state, action.id);
    case UPLOAD_FAIL:
      return failUpload(state, action.id, action.error);
    default:
      return state;
  }
};

export default reducer;
