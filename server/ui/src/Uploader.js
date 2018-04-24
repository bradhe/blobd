import React, { Component } from 'react';
import { connect } from 'react-redux';
import { upload } from './actions';

import {
  UPLOAD_START,
  UPLOAD_PROGRESS,
  UPLOAD_COMPLETE,
  UPLOAD_FAIL
} from './constants.js';

class Uploader extends Component {
  constructor(props, context) {
    super(props, context);

    this.onFileChange = this.onFileChange.bind(this);
  }

  onFileChange(e) {
    const { dispatch } = this.props;

    // TODO: Should there be support for multiple files?
    let file = e.target.files[0];
    dispatch(upload(file));
  }

  renderUploader() {
    const disabled = this.props.status == 'complete' || !this.props.status ? false : true;

    return (
      <input type="file" onChange={this.onFileChange} disabled={disabled} />
    );
  }

  renderProgress() {
    return (
      <span className="blobd-file-progress">{this.props.progress}%</span>
    );
  }

  render() {
    let jsx;

    switch (this.props.status) {
      default:
      case 'completed':
      case 'started':
        jsx = this.renderUploader();
        break;
      case 'progress':
        jsx = this.renderProgress();
        break;
      case 'failed':
        jsx = <span className="blobd-uploader-failed">Upload failed.</span>;
        break;
    }

    return (
      <div className="blobd-file-uploade">
        {jsx}
      </div>
    );
  }
}

const mapStateToProps = (state, ownProps) => {
  return {}
};

export default connect(mapStateToProps)(Uploader);
