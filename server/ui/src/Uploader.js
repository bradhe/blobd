import React, { Component } from 'react';
import { connect } from 'react-redux';
import { upload } from './actions';
import { PropTypes } from 'prop-types';

import FloatingActionButton from 'material-ui/FloatingActionButton';
import ContentAdd from 'material-ui/svg-icons/content/add';
import ContentCopy from 'material-ui/svg-icons/content/content-copy';
import FileDownload from 'material-ui/svg-icons/file/file-download';
import LinearProgress from 'material-ui/LinearProgress';
import Paper from 'material-ui/Paper';
import {List, ListItem} from 'material-ui/List';
import IconButton from 'material-ui/IconButton';

import {
  UPLOAD_START,
  UPLOAD_PROGRESS,
  UPLOAD_COMPLETE,
  UPLOAD_FAIL
} from './constants.js';

const propTypes = {
  status: PropTypes.string,
  file: PropTypes.object,
  progress: PropTypes.number,
};

const uploadButtonDisabled = (status) => {
  return status === UPLOAD_START || status === UPLOAD_PROGRESS;
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

  renderProgress() {
    switch (this.props.status) {
    case UPLOAD_START:
        return <LinearProgress />;
    case UPLOAD_PROGRESS:
        return <LinearProgress mode="determinate" min={0} max={100} value={this.props.progress} />;
    }
  }

  renderUploadingFile() {
    let buttons;

    if (this.props.status == UPLOAD_COMPLETE) {
      buttons = (
        <div>
          <IconButton>
            <ContentCopy/>
          </IconButton>
          <IconButton>
            <FileDownload/>
          </IconButton>
        </div>
      )
    }
      
    return (
      <Paper className="blobd-uploader-status">
        <List>
          <ListItem
            disabled={true}
            primaryText={this.props.file.name}
            secondaryText={this.renderProgress()}
            rightIconButton={buttons} />
        </List>
      </Paper>
    );
  }

  renderFailedMessage() {
    return (
      <span className="blobd-uploader-failed">Upload failed.</span>
    );
  }

  renderContent() {
    switch (this.props.status) {
      case UPLOAD_START:
      case UPLOAD_PROGRESS:
      case UPLOAD_COMPLETE:
        return this.renderUploadingFile()
      case UPLOAD_FAIL:
        return this.renderFailedMessage()
    }
  }

  render() {
    return (
      <div className="blobd-uploader">
        {this.renderContent()}
        {this.renderUploader()}
      </div>
    );
  }
}

Uploader.propTypes = propTypes;

const mapStateToProps = (state, ownProps) => {
  return {}
};

export default connect(mapStateToProps)(Uploader);
