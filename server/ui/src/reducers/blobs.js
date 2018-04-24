import {
  TRACK_BLOB
} from '../constants.js';

const reducer = (state = {}, action) => {
  switch (action.type) {
    case TRACK_BLOB:
      let blobs = state.blobs || [];
      blobs.push({ filename: action.filename, blob: action.blob });
      return { ...state, blobs: blobs };
    default:
      return { blobs: [] };
  };
};

export default reducer;
