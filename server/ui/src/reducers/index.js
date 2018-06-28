import { combineReducers } from 'redux';
import uploads from './uploads';
import blobs from './blobs';

export default combineReducers({
  uploads,
  blobs,
});
