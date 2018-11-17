import {
  NEW_BLOB
} from '../constants.js';

const addBlob = (state, attrs) => {
  return state.concat(attrs);
}

const reducer = (state = [], action) => {
  switch (action.type) {
    case NEW_BLOB:
      return addBlob(state, action);
    default:
      return state;
  };

  return state;
};

export default reducer;
