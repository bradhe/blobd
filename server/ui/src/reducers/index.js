import { combineReducers } from 'redux';
import upload from './upload';
import blobs from './blobs';

export default combineReducers({
  upload,
  blobs,
});
