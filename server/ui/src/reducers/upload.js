import {
  UPLOAD_READY,
  UPLOAD_START,
  UPLOAD_PROGRESS,
  UPLOAD_COMPLETE,
  UPLOAD_FAIL
} from '../constants.js';

const reducer = (state = {}, action) => {
  switch(action.type) {
    case UPLOAD_START:
      return { ...state, status: UPLOAD_START, file: action.file };
    case UPLOAD_PROGRESS:
      return { ...state, status: UPLOAD_PROGRESS, progress: action.progress };
    case UPLOAD_COMPLETE:
      return { ...state, status: UPLOAD_READY };
    case UPLOAD_FAIL:
      return { status: UPLOAD_FAIL, file: action.file };
    default:
      return { status: UPLOAD_READY };
  }
};

export default reducer;
