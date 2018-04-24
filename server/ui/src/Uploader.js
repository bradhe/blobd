import React, { Component } from 'react';
import { connect } from 'react-redux';
import { upload } from './actions';
import { PropTypes } from 'prop-types';

import FloatingActionButton from 'material-ui/FloatingActionButton';
import ContentAdd from 'material-ui/svg-icons/content/add';

import {
  UPLOAD_READY
} from './constants.js';

const propTypes = {
  status: PropTypes.string,
  file: PropTypes.object,
  progress: PropTypes.number,
};

const uploadButtonDisabled = (status) => {
  return status !== UPLOAD_READY
}

class Uploader extends Component {
  constructor(props, context) {
    super(props, context);

    this.onFileChange = this.onFileChange.bind(this);
    this.onUploadButtonClick = this.onUploadButtonClick.bind(this);

    this.fileInput = React.createRef();
  }

  onUploadButtonClick(e) {
    this.fileInput.current.click();
  }

  onFileChange(e) {
    const { dispatch } = this.props;

    // TODO: Should there be support for multiple files?
    let file = e.target.files[0];
    dispatch(upload(file));
  }

  renderUploader() {
    return (
      <div className="blobd-uploader-button">
        <input ref={this.fileInput} type="file" id="blobd-uploader-file" style={{display: 'none'}} onChange={this.onFileChange}/>

        <FloatingActionButton disabled={uploadButtonDisabled(this.props.status)} onClick={this.onUploadButtonClick}>
          <ContentAdd />
        </FloatingActionButton>
      </div>
    );
  }

  render() {
    return (
      <div className="blobd-uploader">
        {this.renderUploader()}
      </div>
    );
  }
}

Uploader.propTypes = propTypes;

const mapStateToProps = (state, ownProps) => {
  return {
    ...state.upload,
  }
};

export default connect(mapStateToProps)(Uploader);
