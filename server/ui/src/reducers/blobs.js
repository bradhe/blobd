import {
  NEW_BLOB
} from '../constants.js';

const addBlob = (state, attrs) => {
  return state.concat(attrs);
}

const reducer = (state = [], action) => {
  if (action.type === NEW_BLOB) {
    return addBlob(state, action);
  }

  return state;
};

export default reducer;
