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
      <input className="blobd-uploader-input" type="file" onChange={this.onFileChange} disabled={disabled} />
    );
  }

  renderProgress() {
    return (
      <span className="blobd-uploader-progress">{this.props.progress}%</span>
    );
  }

  renderFailedMessage() {
    return (
      <span className="blobd-uploader-failed">Upload failed.</span>
    );
  }

  renderContent() {
    switch (this.props.status) {
      case UPLOAD_PROGRESS:
        return this.renderProgress();
      case UPLOAD_FAIL:
        return this.renderFailedMessage()
      default:
      case UPLOAD_COMPLETE:
      case UPLOAD_START:
        return this.renderUploader();
    }
  }

  render() {
    return (
      <div className="blobd-uploader">
        {this.renderContent()}
      </div>
    );
  }
}

const mapStateToProps = (state, ownProps) => {
  return {}
};

export default connect(mapStateToProps)(Uploader);
