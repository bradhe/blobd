const reducer = (state = {}, action) => {
  switch(action.type) {
    case 'UPLOAD_START':
      return { ...state, upload: 'start', file: action.file, };
    case 'UPLOAD_PROGRESS':
      return { ...state, upload: 'progress', progress: action.progress };
    case 'UPLOAD_COMPLETE':
      return { ...state, upload: 'complete', file: action.file, };
    case 'UPLOAD_FAIL':
      return { upload: 'fail', file: action.file };
    default:
      return state;
  }
};

export default reducer;
